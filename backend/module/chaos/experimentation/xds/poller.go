package xds

import (
	"context"
	"fmt"
	"time"

	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcpCacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

type Poller struct {
	ctx           context.Context
	storer        experimentstore.Storer
	snapshotCache gcpCacheV3.SnapshotCache

	resourceTtl          time.Duration
	cacheRefreshInterval time.Duration

	rtdsConfig *RTDSConfig
	ecdsConfig *ECDSConfig

	rtdsGeneratorsByTypeUrl map[string]RTDSResourceGenerator
	ecdsGeneratorsByTypeUrl map[string]ECDSResourceGenerator

	rtdsResourceGenerationFailureCount        tally.Counter
	ecdsResourceGenerationFailureCount        tally.Counter
	ecdsDefaultResourceGenerationFailureCount tally.Counter
	setCacheSnapshotSuccessCount              tally.Counter
	setCacheSnapshotFailureCount              tally.Counter
	activeFaultsGauge                         tally.Gauge

	logger *zap.SugaredLogger
	ticker *time.Ticker
}

func (p *Poller) Start() {
	p.ticker = time.NewTicker(p.cacheRefreshInterval)
	go func() {
		for range p.ticker.C {
			p.logger.Info("Refreshing xDS cache")
			p.refreshCache()
		}
	}()
}

func (p *Poller) Stop() {
	p.ticker.Stop()
}

func (p *Poller) refreshCache() {
	allRunningExperiments, err := p.storer.GetExperiments(p.ctx, "", experimentation.GetExperimentsRequest_STATUS_RUNNING)

	if err != nil {
		p.logger.Errorw("Failed to get data from experiments store", "error", err)

		// If failed to get data from DB, stop all ongoing experiments.
		allRunningExperiments = []*experimentstore.Experiment{}
	}

	clustersWithResources := make(map[string]bool)
	rtdsResourcesByCluster := make(map[string][]*RTDSResource)
	ecdsResourcesByCluster := make(map[string][]*ECDSResource)

	for _, experiment := range allRunningExperiments {
		if g, ok := p.rtdsGeneratorsByTypeUrl[TypeUrl(experiment.Config.Message)]; ok {
			r, err := g.GenerateResource(experiment)
			if err != nil {
				p.rtdsResourceGenerationFailureCount.Inc(1)
				p.logger.Errorw("Error when generating RTDS resource for experiment", "error", err, "experimentRunID", experiment.Run.Id)
				continue
			}
			if r.Empty() {
				// Nothing to do if resource is empty
				continue
			}
			if _, exists := p.ecdsConfig.enabledClusters[r.Cluster]; exists {
				// Allow ECDS to control resource for a given cluster instead
				continue
			}

			rtdsResourcesByCluster[r.Cluster] = append(rtdsResourcesByCluster[r.Cluster], r)
			clustersWithResources[r.Cluster] = true
		}

		if g, ok := p.ecdsGeneratorsByTypeUrl[TypeUrl(experiment.Config.Message)]; ok {
			r, err := g.GenerateResource(experiment)
			if err != nil {
				p.ecdsResourceGenerationFailureCount.Inc(1)
				p.logger.Errorw("Error when generating ECDS resource for experiment", "error", err, "experimentRunID", experiment.Run.Id)
				continue
			}
			if r.Empty() {
				// Nothing to do if resource is empty
				continue
			}
			if _, exists := p.ecdsConfig.enabledClusters[r.Cluster]; !exists {
				// RTDS controls resource for a given cluster
				continue
			}

			ecdsResourcesByCluster[r.Cluster] = append(ecdsResourcesByCluster[r.Cluster], r)
			clustersWithResources[r.Cluster] = true
		}
	}

	// Settings snapshot with empty experiments to remove the experiments
	for _, cluster := range p.snapshotCache.GetStatusKeys() {
		if _, ok := clustersWithResources[cluster]; !ok {
			// in order to remove fault, we need to set the snapshot with default ecds config and default runtime resource
			emptyResources := make(map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl)

			p.logger.Debugw("Removing RTDS resource for cluster", "cluster", cluster)
			runtime, _ := p.createRuntime([][]*RuntimeKeyValue{}, nil)
			emptyResources[gcpTypes.Runtime] = runtime

			if _, exists := p.ecdsConfig.ecdsResourceMap.requestedResourcesMap[cluster]; exists {
				p.logger.Debugw("Removing ECDS resources for cluster", "cluster", cluster)
				resourceNames := p.ecdsConfig.ecdsResourceMap.getResourcesFromCluster(cluster)

				resourcesWithTTL := make([]gcpTypes.ResourceWithTtl, 0)
				// iterate over all of the ECDS resource generators and all of the resource names
				for _, g := range p.ecdsGeneratorsByTypeUrl {
					for _, n := range resourceNames {
						r, err := g.GenerateDefaultResource(cluster, n)
						if err != nil {
							p.ecdsDefaultResourceGenerationFailureCount.Inc(1)
							p.logger.Errorw("Cannot generate empty ECDS resource for cluster", "cluster", cluster, "resource", n)
						}
						if r.Empty() {
							continue
						}

						resourceWithTTL := p.createResourceWithTTLForECDSResource(r, nil)
						resourcesWithTTL = append(resourcesWithTTL, resourceWithTTL...)
					}
				}

				emptyResources[gcpTypes.ExtensionConfig] = resourcesWithTTL
			}

			err := p.setSnapshot(emptyResources, cluster)
			if err != nil {
				p.setCacheSnapshotFailureCount.Inc(1)
				p.logger.Errorw("Unable to unset the fault for cluster", "cluster", cluster, "error", err)
			} else {
				p.setCacheSnapshotSuccessCount.Inc(1)
			}
		}
	}

	activeFaults := 0
	for cluster := range clustersWithResources {
		resources := make(map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl)

		if rs, ok := rtdsResourcesByCluster[cluster]; ok {
			p.logger.Debugw("Injecting fault for cluster using RTDS", "cluster", cluster)
			runtimeKeyValues := make([][]*RuntimeKeyValue, 0)
			for _, r := range rs {
				runtimeKeyValues = append(runtimeKeyValues, r.RuntimeKeyValues)
			}

			runtime, faultsCount := p.createRuntime(runtimeKeyValues, &p.resourceTtl)
			resources[gcpTypes.Runtime] = runtime
			activeFaults += faultsCount
		}

		if rs, ok := ecdsResourcesByCluster[cluster]; ok {
			p.logger.Debugw("Injecting fault for cluster using ECDS", "cluster", cluster)
			for _, r := range rs {
				resources[gcpTypes.ExtensionConfig] = p.createResourceWithTTLForECDSResource(r, &p.resourceTtl)
				activeFaults += 1
				// we apply only the first ECDS resource and ignore the rest
				break
			}
		}

		err := p.setSnapshot(resources, cluster)
		if err != nil {
			p.logger.Errorw("Unable to set the fault for cluster using ECDS", "cluster", cluster, "error", err)
		}
	}

	p.activeFaultsGauge.Update(float64(activeFaults))
}

// Only one experiment can end up controlling a runtime key-value pair
// for a given cluster but it's  still possible for multiple different experiments to
// control runtime values for one cluster as long as each of them controls different
// key-value pairs.
func (p *Poller) createRuntime(runtimeKeyValues [][]*RuntimeKeyValue, ttl *time.Duration) ([]gcpTypes.ResourceWithTtl, int) {
	isUsed := func(usedKeys map[string]bool, keyValues []*RuntimeKeyValue) bool {
		for _, kv := range keyValues {
			if val, ok := usedKeys[kv.Key]; ok && val {
				return true
			}
		}
		return false
	}

	updateIsUsed := func(usedKeys map[string]bool, keyValues []*RuntimeKeyValue) {
		for _, kv := range keyValues {
			usedKeys[kv.Key] = true
		}
	}

	appliedKeValueSetsCount := 0
	usedKeyValues := make([]*RuntimeKeyValue, 0)
	usedKeys := make(map[string]bool)
	for _, kvs := range runtimeKeyValues {
		if isUsed(usedKeys, kvs) {
			// Skip if there is an overlap between already processed keys and
			// keys that are being currently evaluated.
			continue
		}

		appliedKeValueSetsCount += 1
		updateIsUsed(usedKeys, kvs)
		usedKeyValues = append(usedKeyValues, kvs...)
	}

	var fieldsMap = map[string]*pstruct.Value{}
	for _, kv := range usedKeyValues {
		fieldsMap[kv.Key] = &pstruct.Value{
			Kind: &pstruct.Value_NumberValue{
				NumberValue: float64(kv.Value),
			},
		}
	}

	// We only want to set a TTL for non-empty layers. This ensures that we don't spam clients with heartbeat responses
	// unless they have an active runtime override.
	var resourceTTL *time.Duration
	if len(usedKeyValues) > 0 {
		resourceTTL = ttl
	}

	return []gcpTypes.ResourceWithTtl{{
		Resource: &gcpRuntimeServiceV3.Runtime{
			Name: p.rtdsConfig.layerName,
			Layer: &pstruct.Struct{
				Fields: fieldsMap,
			},
		},
		Ttl: resourceTTL,
	}}, appliedKeValueSetsCount
}

func (p *Poller) createResourceWithTTLForECDSResource(resource *ECDSResource, ttl *time.Duration) []gcpTypes.ResourceWithTtl {
	return []gcpTypes.ResourceWithTtl{{
		Resource: resource.ExtensionConfig,
		Ttl:      ttl,
	}}
}

func (p *Poller) setSnapshot(resourceMap map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl, cluster string) error {
	currentSnapshot, _ := p.snapshotCache.GetSnapshot(cluster)

	snapshot := gcpCacheV3.Snapshot{}
	isNewSnapshotSame := true
	for resourceType, resources := range resourceMap {
		newComputedVersion, err := computeChecksum(resources)
		if err != nil {
			p.logger.Errorw("Unable to compute checksum", "cluster", cluster, "resourceType", resourceType, "resources", resources)
			return nil
		}

		if newComputedVersion != currentSnapshot.Resources[resourceType].Version {
			isNewSnapshotSame = false
		}

		snapshot.Resources[resourceType] = gcpCacheV3.NewResourcesWithTtl(newComputedVersion, resources)
	}

	if isNewSnapshotSame {
		// No change in snapshot of this cluster
		p.logger.Debugw("Not setting snapshot as its same as old one", "cluster", cluster, "snapshot", currentSnapshot, "resourceMap", resourceMap)
		return nil
	}

	p.logger.Infow("Setting snapshot", "cluster", cluster, "resources", resourceMap)
	err := p.snapshotCache.SetSnapshot(cluster, snapshot)
	if err != nil {
		return err
	}

	return nil
}

func computeChecksum(item interface{}) (string, error) {
	hash, err := hashstructure.Hash(item, hashstructure.FormatV1, &hashstructure.HashOptions{TagName: "json"})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", hash), nil
}
