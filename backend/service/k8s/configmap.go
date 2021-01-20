package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) ListConfigMaps(ctx context.Context, clientset, cluster, namespace string, listOptions *k8sapiv1.ListOptions) ([]*k8sapiv1.ConfigMap, error) {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
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
