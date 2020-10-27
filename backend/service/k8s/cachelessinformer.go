package k8s

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

func NewCachelessInformer(
	cs ContextClientset,
	lw cache.ListerWatcher,
	objType runtime.Object,
	resync time.Duration,
	h cache.ResourceEventHandler,
) cache.Controller {
	fifo := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		// just to satisify the interface were not going to use this
		KnownObjects:          cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, cache.Indexers{}),
		EmitDeltaTypeReplaced: true,
	})

	return cache.New(&cache.Config{
		Queue:            fifo,
		ListerWatcher:    lw,
		ObjectType:       objType,
		FullResyncPeriod: resync,
		RetryOnError:     false,
		Process: func(obj interface{}) error {
			for _, d := range obj.(cache.Deltas) {
				switch d.Type {
				case cache.Sync, cache.Replaced, cache.Added, cache.Updated:
					h.OnAdd(d.Object)
				case cache.Deleted:
					h.OnDelete(d.Object)
				}
			}
			return nil
		},
	})
}
