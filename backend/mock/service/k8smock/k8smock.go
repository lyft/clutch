package k8smock

import (
	"context"
	"math/rand"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

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

func (*svc) DeleteHPA(ctx context.Context, clientset, cluster, namespace, name string) error {
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

func (*svc) DeleteDeployment(ctx context.Context, clientset, cluster, namespace, name string) error {
	return nil
}

func (*svc) DescribeStatefulSet(ctx context.Context, clientset, cluster, namespace, name string) (*k8sv1.StatefulSet, error) {
	return &k8sv1.StatefulSet{
		Cluster:     cluster,
		Namespace:   namespace,
		Name:        "statefulset1",
		Labels:      map[string]string{"Key": "value"},
		Annotations: map[string]string{"Key": "value"},
	}, nil
}

func (*svc) UpdateStatefulSet(ctx context.Context, clientset, cluster, namespace, name string, fields *k8sv1.UpdateStatefulSetRequest_Fields) error {
	return nil
}

func (*svc) DeleteStatefulSet(ctx context.Context, clientset, cluster, namespace, name string) error {
	return nil
}

func (*svc) DeletePod(ctx context.Context, clientset, cluster, namespace, name string) error {
	return nil
}

func (s *svc) UpdatePod(ctx context.Context, clientset, cluster, namespace, name string, expectedObjectMetaFields *k8sv1.ExpectedObjectMetaFields, objectMetaFields *k8sv1.ObjectMetaFields, removeObjectMetaFields *k8sv1.RemoveObjectMetaFields) error {
	return nil
}

func (s *svc) DescribeService(_ context.Context, clientset, cluster, namespace, name string) (*k8sv1.Service, error) {
	return &k8sv1.Service{
		Cluster:     "fake-cluster-name",
		Namespace:   namespace,
		Name:        name,
		Type:        k8sv1.Service_Type(rand.Intn(len(k8sv1.Service_Type_value))),
		Labels:      map[string]string{"Key": "value"},
		Annotations: map[string]string{"Key": "value"},
	}, nil
}

func (*svc) DeleteService(ctx context.Context, clientset, cluster, namespace, name string) error {
	return nil
}

func (s *svc) DescribeCronJob(_ context.Context, clientset, cluster, namespace, name string) (*k8sv1.CronJob, error) {
	return &k8sv1.CronJob{
		Cluster:     "fake-cluster-name",
		Namespace:   namespace,
		Name:        name,
		Schedule:    "0 0 1 1 *",
		Labels:      map[string]string{"Key": "value"},
		Annotations: map[string]string{"Key": "value"},
	}, nil
}

func (*svc) DeleteCronJob(ctx context.Context, clientset, cluster, namespace, name string) error {
	return nil
}

func (s *svc) DescribeConfigMap(_ context.Context, clientset, cluster, namespace, name string) (*k8sv1.ConfigMap, error) {
	return &k8sv1.ConfigMap{
		Cluster:     "fake-cluster-name",
		Namespace:   namespace,
		Name:        name,
		Labels:      map[string]string{"Key": "value"},
		Annotations: map[string]string{"Key": "value"},
	}, nil
}

func (*svc) DeleteConfigMap(ctx context.Context, clientset, cluster, namespace, name string) error {
	return nil
}

func (s *svc) ListConfigMaps(_ context.Context, clientset, cluster, namespace string, listOptions *k8sv1.ListOptions) ([]*k8sv1.ConfigMap, error) {
	configMaps := []*k8sv1.ConfigMap{
		&k8sv1.ConfigMap{
			Cluster:     "fake-cluster-name",
			Namespace:   namespace,
			Name:        "name1",
			Labels:      listOptions.Labels,
			Annotations: map[string]string{"Key": "value"},
		},
		&k8sv1.ConfigMap{
			Cluster:     "fake-cluster-name",
			Namespace:   namespace,
			Name:        "name2",
			Labels:      listOptions.Labels,
			Annotations: map[string]string{"Key": "value"},
		},
	}
	return configMaps, nil
}

func (s *svc) ListJobs(_ context.Context, clientset, cluster, namespace string, listOptions *k8sv1.ListOptions) ([]*k8sv1.Job, error) {
	jobs := []*k8sv1.Job{
		&k8sv1.Job{
			Cluster:     "fake-cluster-name",
			Namespace:   namespace,
			Name:        "name1",
			Labels:      listOptions.Labels,
			Annotations: map[string]string{"Key": "value"},
		},
		&k8sv1.Job{
			Cluster:     "fake-cluster-name",
			Namespace:   namespace,
			Name:        "name2",
			Labels:      listOptions.Labels,
			Annotations: map[string]string{"Key": "value"},
		},
	}
	return jobs, nil
}

func (*svc) DeleteJob(ctx context.Context, clientset, cluster, namespace, name string) error {
	return nil
}

func (*svc) Clientsets(ctx context.Context) ([]string, error) {
	return []string{"fake-user@fake-cluster"}, nil
}

func New() k8sservice.Service {
	return &svc{}
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}
