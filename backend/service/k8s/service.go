package k8s

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeService(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.Service, error) {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	services, err := cs.CoreV1().Services(cs.Namespace()).List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	if len(services.Items) == 1 {
		return ProtoForService(cs.Cluster(), &services.Items[0]), nil
	} else if len(services.Items) > 1 {
		return nil, fmt.Errorf("Located multiple Services")
	}
	return nil, fmt.Errorf("Unable to locate service")
}

func (s *svc) DeleteService(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return cs.CoreV1().Services(cs.Namespace()).Delete(ctx, name, opts)
}

func ProtoForService(cluster string, k8sservice *corev1.Service) *k8sapiv1.Service {
	clusterName := k8sservice.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	return &k8sapiv1.Service{
		Cluster:     clusterName,
		Namespace:   k8sservice.Namespace,
		Name:        k8sservice.Name,
		Type:        protoForServiceType(k8sservice.Spec.Type),
		Labels:      k8sservice.Labels,
		Annotations: k8sservice.Annotations,
	}
}

func protoForServiceType(serviceType corev1.ServiceType) k8sapiv1.Service_Type {
	// Look up value in generated enum map after ensuring consistent case with generated code.
	val, ok := k8sapiv1.Service_Type_value[strings.ToUpper(string(serviceType))]
	if !ok {
		return k8sapiv1.Service_UNKNOWN
	}

	return k8sapiv1.Service_Type(val)
}
