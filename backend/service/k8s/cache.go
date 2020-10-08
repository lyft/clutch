package k8s

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
)

func (s *svc) CacheEnabled() bool {
	return true
}

func (s *svc) GetTopologyObjectChannel(ctx context.Context) chan *topologyv1.UpdateCacheRequest {
	for name, cs := range s.manager.Clientsets() {
		log.Printf("starting informer for cluster: %s", name)
		s.startInformers(ctx, cs)
	}

	return s.topologyObjectChan
}

func (s *svc) startInformers(ctx context.Context, cs ContextClientset) {
	stop := make(chan struct{})
	// TODO (mcutalo): Make this configurable or keep this at a high value
	factory := informers.NewSharedInformerFactoryWithOptions(cs, time.Minute*60)

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

	select {
	case <-ctx.Done():
		s.log.Info("Shutting down the kubernetes cache informers")
		close(stop)
		close(s.topologyObjectChan)
	default:
	}
}

func (s *svc) informerAddHandler(obj interface{}) {
	s.processInformerEvent(obj, topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE)
}

func (s *svc) informerUpdateHandler(oldObj, newObj interface{}) {
	s.processInformerEvent(newObj, topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE)
}

func (s *svc) informerDeleteHandler(obj interface{}) {
	s.processInformerEvent(obj, topologyv1.UpdateCacheRequest_DELETE)
}

func (s *svc) processInformerEvent(obj interface{}, action topologyv1.UpdateCacheRequest_Action) {
	switch objType := obj.(type) {
	case *corev1.Pod:
		pod := podDescription(obj.(*corev1.Pod), "")
		protoPod, err := ptypes.MarshalAny(pod)
		if err != nil {
			s.log.Error("unable to marshal pod", zap.Error(err))
			return
		}
		s.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
			Resource: &topologyv1.Resource{
				Id:       pod.Name,
				Pb:       protoPod,
				Metadata: pod.GetLabels(),
			},
			Action: action,
		}
	case *appsv1.Deployment:
		deployment := ProtoForDeployment("", obj.(*appsv1.Deployment))
		protoDeployment, err := ptypes.MarshalAny(deployment)
		if err != nil {
			s.log.Error("unable to marshal deployment", zap.Error(err))
			return
		}
		s.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
			Resource: &topologyv1.Resource{
				Id:       deployment.Name,
				Pb:       protoDeployment,
				Metadata: deployment.GetLabels(),
			},
			Action: action,
		}
	case *autoscalingv1.HorizontalPodAutoscaler:
		hpa := ProtoForHPA("", obj.(*autoscalingv1.HorizontalPodAutoscaler))
		protoHpa, err := ptypes.MarshalAny(hpa)
		if err != nil {
			s.log.Error("unable to marshal hpa", zap.Error(err))
			return
		}
		s.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
			Resource: &topologyv1.Resource{
				Id:       hpa.Name,
				Pb:       protoHpa,
				Metadata: hpa.GetLabels(),
			},
			Action: action,
		}
	default:
		s.log.Warn("unable to determin topology object type", zap.Any("object type", objType))
	}
}
