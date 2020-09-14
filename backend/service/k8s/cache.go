package k8s

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/lyft/clutch/backend/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func (s *svc) CacheEnabled() bool {
	return true
}

func (s *svc) GetTopologyObjectChannel() chan types.TopologyObject {
	stop := make(chan struct{})

	for name, cs := range s.manager.Clientsets() {
		log.Printf("starting informer for cluster: %s", name)
		s.startInformers(cs, stop)
	}

	return s.TopologyObjectChan
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
	podObj, _ := obj.(*corev1.Pod)
	clutchPod := podDescription(podObj, "")
	protoPod, _ := ptypes.MarshalAny(clutchPod)

	s.TopologyObjectChan <- types.TopologyObject{
		ResolverTypeURL: protoPod.GetTypeUrl(),
		Pb:              protoPod,
		Metadata:        clutchPod.GetLabels(),
		Action:          types.CREATE,
	}
}

func (s *svc) informerUpdateHandler(oldObj, newObj interface{}) {
	log.Print("Update Handler")
	podObj, _ := newObj.(*corev1.Pod)
	clutchPod := podDescription(podObj, "")
	protoPod, _ := ptypes.MarshalAny(clutchPod)

	s.TopologyObjectChan <- types.TopologyObject{
		ResolverTypeURL: protoPod.GetTypeUrl(),
		Pb:              protoPod,
		Metadata:        clutchPod.GetLabels(),
		Action:          types.UPDATE,
	}
}

func (s *svc) informerDeleteHandler(obj interface{}) {
	log.Print("Delete handler")
	podObj, _ := obj.(*corev1.Pod)
	clutchPod := podDescription(podObj, "")
	protoPod, _ := ptypes.MarshalAny(clutchPod)

	s.TopologyObjectChan <- types.TopologyObject{
		ResolverTypeURL: protoPod.GetTypeUrl(),
		Pb:              protoPod,
		Metadata:        clutchPod.GetLabels(),
		Action:          types.DELETE,
	}
}
