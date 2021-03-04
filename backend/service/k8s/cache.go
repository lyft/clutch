package k8s

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
)

// Setting a large channel buffer mostly for first boot and the  resync timer,
// this really should be sized according to the size of your k8s deployment.
// However this should be a large enough buffer for the datastore to keep up with.
const topologyObjectChanBufferSize = 5000
const topologyInformerLockId = 1

func (s *svc) CacheEnabled() bool {
	return true
}

func (s *svc) StartTopologyCaching(ctx context.Context) (<-chan *topologyv1.UpdateCacheRequest, error) {
	// There should only ever be one instances of all the informers for topology caching
	// We lock here until the context is closed
	if !s.topologyInformerLock.TryAcquire(topologyInformerLockId) {
		return nil, errors.New("TopologyCaching is already in progress")
	}

	clientsets, err := s.manager.Clientsets(ctx)
	if err != nil {
		return nil, err
	}
	for name, cs := range clientsets {
		s.log.Info("starting informer for", zap.String("cluster", name))
		go s.startInformers(ctx, name, cs)
	}

	return s.topologyObjectChan, nil
}

func (s *svc) startInformers(ctx context.Context, clusterName string, cs ContextClientset) {
	informerHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    s.informerAddHandler,
		UpdateFunc: s.informerUpdateHandler,
		DeleteFunc: s.informerDeleteHandler,
	}

	podInformer := NewLightweightInformer(
		cache.NewListWatchFromClient(cs.CoreV1().RESTClient(), "pods", corev1.NamespaceAll, fields.Everything()),
		&corev1.Pod{},
		informerHandlers,
		false,
		clusterName,
	)

	deploymentInformer := NewLightweightInformer(
		cache.NewListWatchFromClient(cs.AppsV1().RESTClient(), "deployments", corev1.NamespaceAll, fields.Everything()),
		&appsv1.Deployment{},
		informerHandlers,
		true,
		clusterName,
	)

	hpaInformer := NewLightweightInformer(
		cache.NewListWatchFromClient(cs.AutoscalingV1().RESTClient(), "horizontalpodautoscalers", corev1.NamespaceAll, fields.Everything()),
		&autoscalingv1.HorizontalPodAutoscaler{},
		informerHandlers,
		true,
		clusterName,
	)

	stop := make(chan struct{})
	go podInformer.Run(stop)
	go deploymentInformer.Run(stop)
	go hpaInformer.Run(stop)

	<-ctx.Done()
	s.log.Info("Shutting down the kubernetes cache informers")
	close(stop)
	close(s.topologyObjectChan)
	s.topologyInformerLock.Release(topologyInformerLockId)
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
		pod := podDescription(objType, "")
		protoPod, err := ptypes.MarshalAny(pod)
		if err != nil {
			s.log.Error("unable to marshal pod", zap.Error(err))
			return
		}
		s.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
			Resource: &topologyv1.Resource{
				Id: meta.HydratedPatternForProto(pod),
				Pb: protoPod,
			},
			Action: action,
		}
	case *appsv1.Deployment:
		deployment := ProtoForDeployment("", objType)
		protoDeployment, err := ptypes.MarshalAny(deployment)
		if err != nil {
			s.log.Error("unable to marshal deployment", zap.Error(err))
			return
		}
		s.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
			Resource: &topologyv1.Resource{
				Id: meta.HydratedPatternForProto(deployment),
				Pb: protoDeployment,
			},
			Action: action,
		}
	case *autoscalingv1.HorizontalPodAutoscaler:
		hpa := ProtoForHPA("", objType)
		protoHpa, err := ptypes.MarshalAny(hpa)
		if err != nil {
			s.log.Error("unable to marshal hpa", zap.Error(err))
			return
		}
		s.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
			Resource: &topologyv1.Resource{
				Id: meta.HydratedPatternForProto(hpa),
				Pb: protoHpa,
			},
			Action: action,
		}
	default:
		s.log.Warn("unable to determine topology object type", zap.String("object_type", fmt.Sprintf("%T", objType)))
	}
}
