package k8s

import (
	"encoding/json"
	"log"
	"time"

	topologyservice "github.com/lyft/clutch/backend/service/topology"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func (s *svc) ManageCache(client topologyservice.Client) {
	s.topologyCache = client

	log.Print("populate cache")
	stop := make(chan struct{})

	for name, cs := range s.GetClientSets() {
		log.Printf("starting informer for cluster: %s", name)
		s.startInformers(cs, stop)
	}
}

func (s *svc) startInformers(cs ContextClientset, stop chan struct{}) {
	factory := informers.NewSharedInformerFactoryWithOptions(cs, time.Minute*1)

	podInformer := factory.Core().V1().Pods().Informer()
	deploymentInformer := factory.Apps().V1().Deployments().Informer()

	informerHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    s.informerAddHandler,
		UpdateFunc: s.informerUpdateHandler,
		DeleteFunc: s.informerDeleteHandler,
	}

	podInformer.AddEventHandler(informerHandlers)
	deploymentInformer.AddEventHandler(informerHandlers)

	go func() {
		podInformer.Run(stop)
	}()

	// go func() {
	// 	deploymentInformer.Run(stop)
	// }()
}

// switch to select type?
func (s *svc) informerAddHandler(obj interface{}) {
	log.Print("Add Handler")
	// log.Printf("%v", obj.(runtime.Object))

	k8sObj := obj.(*corev1.Pod)
	// todo: make a switch for this choiceniss
	// k8sObj.GetObjectKind()
	b, _ := json.Marshal(k8sObj)

	s.topologyCache.SetCache(k8sObj.Name, "pod", b)
}

func (s *svc) informerUpdateHandler(oldObj, newObj interface{}) {
	log.Print("Update Handler")

	k8sObj := newObj.(*corev1.Pod)
	b, _ := json.Marshal(k8sObj)
	s.topologyCache.SetCache(k8sObj.Name, "pod", b)
}

func (s *svc) informerDeleteHandler(obj interface{}) {
	log.Print("Delete handler")
	// log.Printf("%v", obj.(runtime.Object))
	k8sObj := obj.(*corev1.Pod)

	s.topologyCache.DeleteCache(k8sObj.Name)
}
