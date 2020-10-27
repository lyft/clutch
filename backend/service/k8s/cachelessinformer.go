package k8s

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
)

type lightweightObj struct {
	metav1.Object
	UID       types.UID
	Name      string
	Namespace string
}

func (cl *lightweightObj) GetUID() types.UID             { return cl.UID }
func (cl *lightweightObj) SetUID(uid types.UID)          { cl.UID = uid }
func (cl *lightweightObj) GetName() string               { return cl.Name }
func (cl *lightweightObj) SetName(name string)           { cl.Name = name }
func (cl *lightweightObj) GetNamespace() string          { return cl.Namespace }
func (cl *lightweightObj) SetNamespace(namespace string) { cl.Namespace = namespace }

func NewLightweightInformer(
	cs ContextClientset,
	lw cache.ListerWatcher,
	objType runtime.Object,
	resync time.Duration,
	h cache.ResourceEventHandler,
) cache.Controller {
	cacheStore := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, cache.Indexers{})
	fifo := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KnownObjects:          cacheStore,
		EmitDeltaTypeReplaced: true,
		KeyFunction: func(obj interface{}) (string, error) {
			theMeta, err := meta.Accessor(obj)
			if err != nil {
				return "", err
			}
			return string(theMeta.GetUID()), nil
		},
	})

	return cache.New(&cache.Config{
		Queue:            fifo,
		ListerWatcher:    lw,
		ObjectType:       objType,
		FullResyncPeriod: resync,
		RetryOnError:     false,
		Process: func(obj interface{}) error {
			for _, d := range obj.(cache.Deltas) {

				incomeingObjectMeta, _ := meta.Accessor(d.Object)
				lightweightObj := &lightweightObj{}
				lightweightObj.SetUID(incomeingObjectMeta.GetUID())
				lightweightObj.SetName(incomeingObjectMeta.GetName())
				lightweightObj.SetNamespace(incomeingObjectMeta.GetNamespace())

				switch d.Type {
				case cache.Sync, cache.Replaced, cache.Added, cache.Updated:
					if _, exists, err := cacheStore.Get(lightweightObj); err == nil && exists {
						if err := cacheStore.Update(lightweightObj); err != nil {
							log.Printf("error updating %v", err)
						}
						log.Print("trying to update")
						h.OnUpdate(nil, d.Object)
					} else {
						if err := cacheStore.Add(lightweightObj); err != nil {
							log.Printf("error adding %v", err)
						}
						log.Print("trying to add")
						h.OnAdd(d.Object)
					}
				case cache.Deleted:
					log.Print("delete event")
					if err := cacheStore.Delete(lightweightObj); err != nil {
						return err
					}
					h.OnDelete(d.Object)
				}
			}
			return nil
		},
	})
}
