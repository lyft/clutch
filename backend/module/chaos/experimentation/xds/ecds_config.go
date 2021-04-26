package xds

import "sync"

type SafeEcdsResourceMap struct {
	mu                    sync.Mutex
	requestedResourcesMap map[string]map[string]struct{}
}

func (safeResource *SafeEcdsResourceMap) getResourcesFromCluster(cluster string) []string {
	var resourceNames []string

	safeResource.mu.Lock()
	if _, exists := safeResource.requestedResourcesMap[cluster]; exists {
		for resourceName := range safeResource.requestedResourcesMap[cluster] {
			resourceNames = append(resourceNames, resourceName)
		}
	}
	safeResource.mu.Unlock()

	return resourceNames
}

func (safeResource *SafeEcdsResourceMap) setResourcesForCluster(cluster string, newResources []string) {
	safeResource.mu.Lock()
	defer safeResource.mu.Unlock()

	resources := make(map[string]struct{})
	if _, exists := safeResource.requestedResourcesMap[cluster]; exists {
		resources = safeResource.requestedResourcesMap[cluster]
	}

	for _, newResource := range newResources {
		resources[newResource] = struct{}{}
	}

	safeResource.requestedResourcesMap[cluster] = resources
}

type ECDSConfig struct {
	enabledClusters map[string]struct{}

	ecdsResourceMap *SafeEcdsResourceMap
}
