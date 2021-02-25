package k8s

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeConfigMap(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.ConfigMap, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	configMapList, err := cs.CoreV1().ConfigMaps(cs.Namespace()).List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})

	if err != nil {
		return nil, err
	}

	if len(configMapList.Items) == 1 {
		return protoForConfigMap(cs.Cluster(), &configMapList.Items[0]), nil
	} else if len(configMapList.Items) > 1 {
		return nil, status.Errorf(codes.FailedPrecondition, "located multiple config maps with name '%s'", name)
	}
	return nil, status.Error(codes.NotFound, "unable to locate specified config map")
}

func (s *svc) DeleteConfigMap(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return cs.CoreV1().ConfigMaps(cs.Namespace()).Delete(ctx, name, opts)
}

func (s *svc) ListConfigMaps(ctx context.Context, clientset, cluster, namespace string, listOptions *k8sapiv1.ListOptions) ([]*k8sapiv1.ConfigMap, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	opts := ApplyListOptions(listOptions)

	configMapList, err := cs.CoreV1().ConfigMaps(cs.Namespace()).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	var configMaps []*k8sapiv1.ConfigMap
	for _, cm := range configMapList.Items {
		configMap := cm
		configMaps = append(configMaps, protoForConfigMap(cs.Cluster(), &configMap))
	}

	return configMaps, nil
}

func protoForConfigMap(cluster string, k8sconfigMap *v1.ConfigMap) *k8sapiv1.ConfigMap {
	clusterName := k8sconfigMap.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}

	return &k8sapiv1.ConfigMap{
		Cluster:     clusterName,
		Namespace:   k8sconfigMap.Namespace,
		Name:        k8sconfigMap.Name,
		Labels:      k8sconfigMap.Labels,
		Annotations: k8sconfigMap.Annotations,
	}
}
