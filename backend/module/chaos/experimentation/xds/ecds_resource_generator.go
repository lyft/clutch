package xds

import (
	gcpCoreV3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

var ECDSGeneratorsByTypeUrl = map[string]ECDSResourceGenerator{}

type ECDSResource struct {
	Cluster         string
	ExtensionConfig *gcpCoreV3.TypedExtensionConfig
}

func NewECDSResource(cluster string, config *gcpCoreV3.TypedExtensionConfig) (*ECDSResource, error) {
	return &ECDSResource{
		Cluster:         cluster,
		ExtensionConfig: config,
	}, nil
}

func NewEmptyECDSResource() *ECDSResource {
	return &ECDSResource{
		Cluster:         "",
		ExtensionConfig: nil,
	}
}

func (r *ECDSResource) Empty() bool {
	return r.ExtensionConfig == nil
}

type ECDSResourceGenerator interface {
	// Generates an ECDS resource for a given experiment. The implementation should return
	// a resource created with a `NewEmptyECDSResource` method call if it's not interested
	// in creating any resources for a given experiment.
	GenerateResource(experiment *experimentstore.Experiment) (*ECDSResource, error)
	// Generates an ECDS resource for a given experiment and a given resource name that
	// moves ECDS state back to its initial state - as if the receiver has not generated any ECDS
	// resource. The implementation of the method should return a resource created with
	// a `NewEmptyECDSResource` method call if it receives a `resourceName` that it doesn't recognize.
	GenerateDefaultResource(cluster string, resourceName string) (*ECDSResource, error)
}
