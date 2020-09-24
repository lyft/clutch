package k8s

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeService(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.Service, error) {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	svc, err := cs.CoreV1().Services(cs.Namespace()).List(metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	if len(svc.Items) == 1 {
		return svcDescription(&svc.Items[0], cs.Cluster()), nil
	} else if len(svc.Items) > 1 {
		return nil, fmt.Errorf("Located multiple Pods")
	}
	return nil, fmt.Errorf("Unable to locate pod")
}

func (s *svc) ListServices(ctx context.Context, clientset, cluster, namespace string, listOpts *k8sapiv1.ListOptions) ([]*k8sapiv1.Service, error) {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	opts := ApplyListOptions(listOpts)

	svcs, err := cs.CoreV1().Services(cs.Namespace()).List(opts)
	if err != nil {
		return nil, err
	}

	var svcList []*k8sapiv1.Service
	for _, s := range svcs.Items {
		svc := s
		svcList = append(svcList, svcDescription(&svc, cs.Cluster()))
	}

	return svcList, nil
}

func svcDescription(k8sSvc *corev1.Service, cluster string) *k8sapiv1.Service {
	// TODO: There's a mismatch between the serialization of the timestamp here and what's expected on the frontend.
	// TODO: (cpuri) Currently I use a string to serialize the timestamp from k8s instead of the protobuf std timestamp...
	//var launch *timestamp.Timestamp
	//if converted, err := ptypes.TimestampProto(k8spod.Status.StartTime.Time); err == nil {
	//	launch = converted
	//}
	clusterName := k8sSvc.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	return &k8sapiv1.Service{
		Cluster:    clusterName,
		Namespace:  k8sSvc.Namespace,
		Name:       k8sSvc.Name,
		ClusterIp:  k8sSvc.Spec.ClusterIP,
		//ExternalIp: k8sSvc.Spec.ExternalIPs, // TODO ExternalIPs is a string array
		Type:      protoForServiceType(k8sSvc.Spec.Type),
		StartTime:   k8sSvc.CreationTimestamp.String(),
		Labels:      k8sSvc.Labels,
		Annotations: k8sSvc.Annotations,
	}
}

func protoForServiceType(svcType corev1.ServiceType) k8sapiv1.Service_Type {
	// Look up value in generated enum map after ensuring consistent case with generated code.
	val, ok := k8sapiv1.Service_Type_value[string(svcType)]
	if !ok {
		return k8sapiv1.Service_Unknown
	}

	return k8sapiv1.Service_Type(val)
}