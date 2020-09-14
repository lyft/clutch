package k8s

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/lyft/clutch/backend/types"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
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
	hpaInformer := factory.Autoscaling().V1().HorizontalPodAutoscalers().Informer()

	informerHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    s.informerAddHandler,
		UpdateFunc: s.informerUpdateHandler,
		DeleteFunc: s.informerDeleteHandler,
	}

	podInformer.AddEventHandler(informerHandlers)
	deploymentInformer.AddEventHandler(informerHandlers)
	hpaInformer.AddEventHandler(informerHandlers)

	go func() { podInformer.Run(stop) }()
	go func() { deploymentInformer.Run(stop) }()
	go func() { hpaInformer.Run(stop) }()
}

func (s *svc) informerAddHandler(obj interface{}) {
	log.Print("Add Handler")
	s.processInformerEvent(obj, types.CREATE)
}

func (s *svc) informerUpdateHandler(oldObj, newObj interface{}) {
	log.Print("Update Handler")
	s.processInformerEvent(newObj, types.UPDATE)
}

func (s *svc) informerDeleteHandler(obj interface{}) {
	log.Print("Delete handler")
	s.processInformerEvent(obj, types.DELETE)
}

func (s *svc) processInformerEvent(obj interface{}, action types.Action) {
	switch obj.(type) {
	case *corev1.Pod:
		pod := podDescription(obj.(*corev1.Pod), "")
		protoPod, _ := ptypes.MarshalAny(pod)
		s.TopologyObjectChan <- types.TopologyObject{
			Id:              pod.Name,
			ResolverTypeURL: protoPod.GetTypeUrl(),
			Pb:              protoPod,
			Metadata:        pod.GetLabels(),
			Action:          action,
		}
	case *appsv1.Deployment:
		deployment := ProtoForDeployment("", obj.(*appsv1.Deployment))
		protoDeployment, _ := ptypes.MarshalAny(deployment)
		s.TopologyObjectChan <- types.TopologyObject{
			Id:              deployment.Name,
			ResolverTypeURL: protoDeployment.GetTypeUrl(),
			Pb:              protoDeployment,
			Metadata:        deployment.GetLabels(),
			Action:          action,
		}
	case *autoscalingv1.HorizontalPodAutoscaler:
		hpa := ProtoForHPA("", obj.(*autoscalingv1.HorizontalPodAutoscaler))
		protoHpa, _ := ptypes.MarshalAny(hpa)
		s.TopologyObjectChan <- types.TopologyObject{
			Id:              hpa.Name,
			ResolverTypeURL: protoHpa.GetTypeUrl(),
			Pb:              protoHpa,
			Metadata:        hpa.GetLabels(),
			Action:          action,
		}
	default:
		s.log.Warn("unable to determin topology object type")
	}
}
