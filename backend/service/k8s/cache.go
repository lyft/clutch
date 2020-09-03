package k8s

import (
	"database/sql"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// setup a channel to stop the informers
func (s *svc) PopulateCache(db *sql.DB) {
	for _, cs := range s.manager.Clientsets() {
		startInformers(cs)
	}
}

func startInformers(cs ContextClientset) {
	factory := informers.NewSharedInformerFactoryWithOptions(cs, time.Second*10)

	podInformer := factory.Core().V1().Pods().Informer()
	deploymentInformer := factory.Apps().V1().Deployments().Informer()

	informerHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    informerAddHandler,
		UpdateFunc: informerUpdateHandler,
		DeleteFunc: informerDeleteHandler,
	}

	podInformer.AddEventHandler(informerHandlers)
	deploymentInformer.AddEventHandler(informerHandlers)
}

func informerAddHandler(obj interface{}) {
	log.Print("Add Handler")
	log.Printf("%v", obj.(runtime.Object))
}

func informerUpdateHandler(oldObj, newObj interface{}) {
	log.Print("Update Handler")
	log.Printf("%v", newObj.(runtime.Object))
}

func informerDeleteHandler(obj interface{}) {
	log.Print("Delete handler")
	log.Printf("%v", obj.(runtime.Object))
}
