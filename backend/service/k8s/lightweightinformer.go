package k8s

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

type lightweightCacheObject struct {
	metav1.Object
	Name      string
	Namespace string
}

func (lw *lightweightCacheObject) GetName() string      { return lw.Name }
func (lw *lightweightCacheObject) GetNamespace() string { return lw.Namespace }

// The LightweightInformer is an informer thats optimized for memory usage with drawbacks.
//
// The reduction in memory consumption does come at a cost, to achieve this we store small objects
// in the informers cache store. We do this by utilizing storing `lightweightCacheObject` instead
// of the full Kubernetes object.
// `lightweightCacheObject` has just enough metadata for the cache store and DeltaFIFO components to operate normally.
//
// There are drawbacks too using a LightweightInformer and its does not fit all use cases.
// For the Topology Caching this type of solution helped to reduce memory footprint significantly
// for large scale Kubernetes deployments.
//
// Also to note the memory footprint of the cache store is only part of the story.
// While the informers controller is receiving Kubernetes objects it stores that full object in the DeltaFIFO queue.
// This queue while processed quickly does store a vast amount of objects at any given time and contributes to memory usage greatly.
//
// Drawbacks
// - Update resource event handler does not function as expected, old objects will always return nil.
//   This is because we dont cache the full k8s object to compute deltas as we are using lightweightCacheObjects instead.
// - Resync does not work as expected becuase the cache is filled with lightweightCacheObjects,
//   for this reason Resync is disabled.

func NewLightweightInformer(
	lw cache.ListerWatcher,
	objType runtime.Object,
	h cache.ResourceEventHandler,
	recieveUpdates bool,
	clusterName string,
) cache.Controller {
	cacheStore := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, cache.Indexers{})
	fifo := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KnownObjects:          cacheStore,
		EmitDeltaTypeReplaced: true,
	})

	return cache.New(&cache.Config{
		Queue:            fifo,
		ListerWatcher:    lw,
		ObjectType:       objType,
		FullResyncPeriod: 0,
		RetryOnError:     false,
		Process: func(obj interface{}) error {
			for _, d := range obj.(cache.Deltas) {
				incomingObjectMeta, err := meta.Accessor(d.Object)
				if err != nil {
					return err
				}

				lightweightObj := &lightweightCacheObject{
					Name:      incomingObjectMeta.GetName(),
					Namespace: incomingObjectMeta.GetNamespace(),
				}

				// ClusterName is still not set in kube v1.20 so we are setting this manually.
				// https://github.com/kubernetes/apimachinery/blob/2456ebdaba229616fab2161a615148884b46644b/pkg/apis/meta/v1/types.go#L266-L270
				incomingObjectMeta.SetClusterName(clusterName)

				switch d.Type {
				case cache.Sync, cache.Replaced, cache.Added, cache.Updated:
					if _, exists, err := cacheStore.Get(lightweightObj); err == nil && exists {
						// Not all use-cases of this informer require updates to Kubernetes objects
						// For this reason you can disable updates completely by setting `recieveUpdates` to false
						// This both disables the cache update and the OnUpdate handler
						if recieveUpdates {
							if err := cacheStore.Update(lightweightObj); err != nil {
								return err
							}
							h.OnUpdate(nil, d.Object)
						}
					} else {
						if err := cacheStore.Add(lightweightObj); err != nil {
							return err
						}
						h.OnAdd(d.Object)
					}
				case cache.Deleted:
					if err := cacheStore.Delete(lightweightObj); err != nil {
						return err
					}
					h.OnDelete(d.Object)
				default:
					return fmt.Errorf("Cache type not supported: %s", d.Type)
				}
			}
			return nil
		},
	})
}
