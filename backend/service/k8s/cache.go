package k8s

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
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

	lwPod := cache.NewListWatchFromClient(cs.CoreV1().RESTClient(), "pods", corev1.NamespaceAll, fields.Everything())
	podInformer := NewLightweightInformer(
		lwPod,
		&corev1.Pod{},
		informerHandlers,
		false,
		clusterName,
	)

	lwDeployment := cache.NewListWatchFromClient(cs.AppsV1().RESTClient(), "deployments", corev1.NamespaceAll, fields.Everything())
	deploymentInformer := NewLightweightInformer(
		lwDeployment,
		&appsv1.Deployment{},
		informerHandlers,
		true,
		clusterName,
	)

	lwHPA := cache.NewListWatchFromClient(cs.AutoscalingV1().RESTClient(), "horizontalpodautoscalers", corev1.NamespaceAll, fields.Everything())
	hpaInformer := NewLightweightInformer(
		lwHPA,
		&autoscalingv1.HorizontalPodAutoscaler{},
		informerHandlers,
		true,
		clusterName,
	)

	stop := make(chan struct{})
	go podInformer.Run(stop)
	go deploymentInformer.Run(stop)
	go hpaInformer.Run(stop)
	go s.cacheFullRelist(ctx, lwPod, lwDeployment, lwHPA)

	<-ctx.Done()
	s.log.Info("Shutting down the kubernetes cache informers")
	close(stop)
	close(s.topologyObjectChan)
	s.topologyInformerLock.Release(topologyInformerLockId)
}

// cacheFullRelist will list all resources and push them to the topology cache for processing.
// The importance of doing a semi frequent full re-list give us better cache accuracy,
// while also keeping resources that are infrequently updated from being cleaned up by the cache TTL.
// This works in tandem with the LightweightInformers above.
//
// Notably we intentionally run these in serial, not only can this cause memory pressure but
// also being mindful of the kubernetes api servers to reduce burst load.
func (s *svc) cacheFullRelist(ctx context.Context, lwPods, lwDeployments, lwHPA *cache.ListWatch) {
	ticker := time.NewTicker(time.Hour * 1)
	for {
		// The informers will only ever do a full list once on boot
		// we will wait the hour before doing another full list again
		select {
		case <-ticker.C:
			pods, err := lwPods.List(metav1.ListOptions{})
			if err != nil {
				s.log.Warn("Unable to list pods to populate Kubernetes cache", zap.Error(err))
				return
			}

			podItems := pods.(*corev1.PodList).Items
			for i := range podItems {
				s.processInformerEvent(&podItems[i], topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE)
			}

			deployments, err := lwDeployments.List(metav1.ListOptions{})
			if err != nil {
				s.log.Warn("Unable to list deployments to populate Kubernetes cache", zap.Error(err))
				return
			}

			deploymentItems := deployments.(*appsv1.DeploymentList).Items
			for i := range deploymentItems {
				s.processInformerEvent(&deploymentItems[i], topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE)
			}

			hpas, err := lwHPA.List(metav1.ListOptions{})
			if err != nil {
				s.log.Warn("Unable to list HPAs to populate Kubernetes cache", zap.Error(err))
				return
			}

			hpaItems := hpas.(*autoscalingv1.HorizontalPodAutoscalerList).Items
			for i := range hpaItems {
				s.processInformerEvent(&hpaItems[i], topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE)
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
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
		pod := podDescription(objType, "")
		protoPod, err := ptypes.MarshalAny(pod)
		if err != nil {
			s.log.Error("unable to marshal pod", zap.Error(err))
			return
		}
		s.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
			Resource: &topologyv1.Resource{
				Id: pod.Name,
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
				Id: deployment.Name,
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
				Id: hpa.Name,
				Pb: protoHpa,
			},
			Action: action,
		}
	default:
		s.log.Warn("unable to determine topology object type", zap.String("object_type", fmt.Sprintf("%T", objType)))
	}
}
