package k8s

// <!-- START clutchdoc -->
// description: Multi-clientset Kubernetes interface.
// <!-- END clutchdoc -->

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"

	k8sconfigv1 "github.com/lyft/clutch/backend/api/config/service/k8s/v1"
	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.k8s"

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	k8sConfig := &k8sconfigv1.Config{}

	// Use the default kubeconfig (environment or well-known path) if kubeconfigs are not passed in.
	// https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/
	if cfg != nil {
		if err := ptypes.UnmarshalAny(cfg, k8sConfig); err != nil {
			return nil, err
		}

		if k8sConfig.Kubeconfigs != nil {
			loadingRules = &clientcmd.ClientConfigLoadingRules{
				Precedence: k8sConfig.Kubeconfigs,
			}
		}
	}

	c, err := newClientsetManager(loadingRules, k8sConfig.RestClientConfig, logger)
	if err != nil {
		return nil, err
	}

	return NewWithClientsetManager(c, logger, scope)
}

type Service interface {
	// All names of clientsets.
	Clientsets() []string

	// Pod management functions.
	DescribePod(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.Pod, error)
	DeletePod(ctx context.Context, clientset, cluster, namespace, name string) error
	ListPods(ctx context.Context, clientset, cluster, namespace string, listOptions *k8sapiv1.ListOptions) ([]*k8sapiv1.Pod, error)
	UpdatePod(ctx context.Context, clientset, cluster, namespace, name string, expected_object_meta_fields *k8sapiv1.ExpectedObjectMetaFields, object_meta_fields *k8sapiv1.ObjectMetaFields, remove_object_meta_fields *k8sapiv1.RemoveObjectMetaFields) error

	// HPA management functions.
	DescribeHPA(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.HPA, error)
	ResizeHPA(ctx context.Context, clientset, cluster, namespace, name string, sizing *k8sapiv1.ResizeHPARequest_Sizing) error
	DeleteHPA(ctx context.Context, clientset, cluster, namespace, name string) error

	// Deployment management functions.
	DescribeDeployment(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.Deployment, error)
	UpdateDeployment(ctx context.Context, clientset, cluster, namespace, name string, fields *k8sapiv1.UpdateDeploymentRequest_Fields) error
	DeleteDeployment(ctx context.Context, clientset, cluster, namespace, name string) error

	// Service management functions.
	DescribeService(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.Service, error)
	DeleteService(ctx context.Context, clientset, cluster, namespace, name string) error

	// StatefulSet management functions.
	DescribeStatefulSet(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.StatefulSet, error)
	UpdateStatefulSet(ctx context.Context, clientset, cluster, namespace, name string, fields *k8sapiv1.UpdateStatefulSetRequest_Fields) error
}

type svc struct {
	manager ClientsetManager

	topologyObjectChan   chan *topologyv1.UpdateCacheRequest
	topologyInformerLock *semaphore.Weighted
	log                  *zap.Logger
	scope                tally.Scope
}

func NewWithClientsetManager(manager ClientsetManager, logger *zap.Logger, scope tally.Scope) (Service, error) {
	return &svc{
		manager:              manager,
		topologyObjectChan:   make(chan *topologyv1.UpdateCacheRequest, topologyObjectChanBufferSize),
		topologyInformerLock: semaphore.NewWeighted(1),
		log:                  logger,
		scope:                scope,
	}, nil
}

func (s *svc) Clientsets() []string {
	ret := make([]string, 0, len(s.manager.Clientsets()))
	for name := range s.manager.Clientsets() {
		ret = append(ret, name)
	}
	return ret
}
