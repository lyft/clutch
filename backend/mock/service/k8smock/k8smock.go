package k8smock

import (
	"context"
	"math/rand"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"

	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
	"github.com/lyft/clutch/backend/service"
	k8sservice "github.com/lyft/clutch/backend/service/k8s"
)

type svc struct{}

func (s *svc) DescribeHPA(ctx context.Context, clientset, cluster, namespace, name string) (*k8sv1.HPA, error) {
	hpa := &k8sv1.HPA{
		Cluster:   "fake-cluster-name",
		Namespace: namespace,
		Name:      name,
		Sizing: &k8sv1.HPA_Sizing{
			MinReplicas:     1,
			MaxReplicas:     100,
			CurrentReplicas: uint32(rand.Int31n(100)),
			DesiredReplicas: uint32(rand.Int31n(100)),
		},
		Labels:      map[string]string{"Label key": "Value"},
		Annotations: map[string]string{"Annotation key": "Value"},
	}
	return hpa, nil
}

func (*svc) ResizeHPA(ctx context.Context, clientset, cluster, namespace, name string, sizing *k8sv1.ResizeHPARequest_Sizing) error {
	return nil
}

func (*svc) Manager() k8sservice.ClientsetManager {
	return nil
}

func (s *svc) DescribePod(_ context.Context, clientset, cluster, namespace, name string) (*k8sv1.Pod, error) {
	pod := &k8sv1.Pod{
		Cluster:   "fake-cluster-name",
		Namespace: namespace,
		Name:      name,
		NodeIp:    "10.0.0.1",
		PodIp:     "8.1.1.8",
		State:     k8sv1.Pod_State(rand.Intn(len(k8sv1.Pod_State_value))),
		//StartTime:   ptypes.TimestampNow(),
		Labels:      map[string]string{"Key": "value"},
		Annotations: map[string]string{"Key": "value"},
	}
	return pod, nil
}

func (s *svc) ListPods(_ context.Context, clientset, cluster, namespace string, listOptions *k8sv1.ListOptions) ([]*k8sv1.Pod, error) {
	pods := []*k8sv1.Pod{
		&k8sv1.Pod{
			Cluster:     cluster,
			Namespace:   namespace,
			Name:        "name1",
			NodeIp:      "10.0.0.1",
			PodIp:       "8.1.1.8",
			State:       k8sv1.Pod_State(rand.Intn(len(k8sv1.Pod_State_value))),
			Labels:      listOptions.Labels,
			Annotations: map[string]string{"Key": "value"},
		},
		&k8sv1.Pod{
			Cluster:     cluster,
			Namespace:   namespace,
			Name:        "name2",
			NodeIp:      "10.0.0.2",
			PodIp:       "8.1.1.9",
			State:       k8sv1.Pod_State(rand.Intn(len(k8sv1.Pod_State_value))),
			Labels:      listOptions.Labels,
			Annotations: map[string]string{"Key": "value"},
		},
	}
	return pods, nil
}

func (*svc) DescribeDeployment(ctx context.Context, clientset, cluster, namespace, name string) (*k8sv1.Deployment, error) {
	return &k8sv1.Deployment{
		Cluster:     cluster,
		Namespace:   namespace,
		Name:        "deployment1",
		Labels:      map[string]string{"Key": "value"},
		Annotations: map[string]string{"Key": "value"},
	}, nil
}

func (*svc) UpdateDeployment(ctx context.Context, clientset, cluster, namespace, name string, fields *k8sv1.UpdateDeploymentRequest_Fields) error {
	return nil
}

func (*svc) DeletePod(ctx context.Context, clientset, cluster, namespace, name string) error {
	return nil
}

func (s *svc) UpdatePod(ctx context.Context, clientset, cluster, namespace, name string, expectedObjectMetaFields *k8sv1.ExpectedObjectMetaFields, objectMetaFields *k8sv1.ObjectMetaFields, removeObjectMetaFields *k8sv1.RemoveObjectMetaFields) error {
	return nil
}

func (*svc) Clientsets() []string {
	return []string{"fake-user@fake-cluster"}
}

func (s *svc) GetClientSets() map[string]k8sservice.ContextClientset {
	return map[string]k8sservice.ContextClientset{
		"fake-cluster": k8sservice.NewContextClientset("ns", "cluster", &kubernetes.Clientset{}),
	}
}

func New() k8sservice.Service {
	return &svc{}
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}
