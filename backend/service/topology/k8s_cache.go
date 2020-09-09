package topology

import (
	"encoding/json"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	k8sservice "github.com/lyft/clutch/backend/service/k8s"
)

func (c *client) startInformers(cs k8sservice.ContextClientset, stop chan struct{}) {
	factory := informers.NewSharedInformerFactoryWithOptions(cs, time.Minute*1)

	podInformer := factory.Core().V1().Pods().Informer()
	deploymentInformer := factory.Apps().V1().Deployments().Informer()

	informerHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.informerAddHandler,
		UpdateFunc: c.informerUpdateHandler,
		DeleteFunc: c.informerDeleteHandler,
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
func (c *client) informerAddHandler(obj interface{}) {
	log.Print("Add Handler")
	// log.Printf("%v", obj.(runtime.Object))

	k8sObj := obj.(*corev1.Pod)
	// todo: make a switch for this choiceniss
	// k8sObj.GetObjectKind()
	b, _ := json.Marshal(k8sObj)

	c.SetCache(k8sObj.Name, "pod", b)
}

func (c *client) informerUpdateHandler(oldObj, newObj interface{}) {
	log.Print("Update Handler")

	k8sObj := newObj.(*corev1.Pod)
	b, _ := json.Marshal(k8sObj)
	c.SetCache(k8sObj.Name, "pod", b)
}

func (c *client) informerDeleteHandler(obj interface{}) {
	log.Print("Delete handler")
	// log.Printf("%v", obj.(runtime.Object))
	k8sObj := obj.(*corev1.Pod)

	c.DeleteCache(k8sObj.Name)
}
