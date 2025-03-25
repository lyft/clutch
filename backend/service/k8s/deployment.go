package k8s

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeDeployment(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.Deployment, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	deployments, err := cs.AppsV1().Deployments(cs.Namespace()).List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	if len(deployments.Items) == 1 {
		return ProtoForDeployment(cs.Cluster(), &deployments.Items[0]), nil
	} else if len(deployments.Items) > 1 {
		return nil, status.Error(codes.FailedPrecondition, "located multiple deployments")
	}

	return nil, status.Error(codes.NotFound, "unable to locate specified deployment")
}

func (s *svc) ListDeployments(ctx context.Context, clientset, cluster, namespace string, listOptions *k8sapiv1.ListOptions) ([]*k8sapiv1.Deployment, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	opts, err := ApplyListOptions(listOptions)
	if err != nil {
		return nil, err
	}

	deploymentList, err := cs.AppsV1().Deployments(cs.Namespace()).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	var deployments []*k8sapiv1.Deployment
	for _, d := range deploymentList.Items {
		deployment := d
		deployments = append(deployments, ProtoForDeployment(cs.Cluster(), &deployment))
	}

	return deployments, nil
}

func ProtoForDeployment(cluster string, deployment *appsv1.Deployment) *k8sapiv1.Deployment {
	clusterName := GetKubeClusterName(deployment)
	if clusterName == "" {
		clusterName = cluster
	}
	k8sDeployment := &k8sapiv1.Deployment{
		Cluster:          clusterName,
		Namespace:        deployment.Namespace,
		Name:             deployment.Name,
		Labels:           deployment.Labels,
		Annotations:      deployment.Annotations,
		DeploymentSpec:   ProtoForDeploymentSpec(deployment.Spec),
		DeploymentStatus: ProtoForDeploymentStatus(deployment.Status),
	}

	if !deployment.CreationTimestamp.IsZero() {
		// Convert Unix Timestamp to milliseconds
		k8sDeployment.CreationTimeMillis = deployment.CreationTimestamp.UnixNano() / 1e6
	}

	return k8sDeployment
}

func processObjProbe(objProbe *v1.Probe) *k8sapiv1.Probe {
	HandlerObj := &k8sapiv1.Probe{}

	if objProbe.ProbeHandler.HTTPGet != nil {
		ObjProbeHTTPHeaders := make([]*k8sapiv1.HTTPHeader, 0, len(objProbe.ProbeHandler.HTTPGet.HTTPHeaders))
		for _, value := range objProbe.ProbeHandler.HTTPGet.HTTPHeaders {
			UniqueLivenessHeader := &k8sapiv1.HTTPHeader{
				Name:  &value.Name,
				Value: &value.Value,
			}
			ObjProbeHTTPHeaders = append(ObjProbeHTTPHeaders, UniqueLivenessHeader)
		}

		ObjProbeHTTPObject := &k8sapiv1.HTTPGetAction{
			Path:        objProbe.ProbeHandler.HTTPGet.Path,
			Port:        objProbe.ProbeHandler.HTTPGet.Port.IntVal,
			Host:        objProbe.ProbeHandler.HTTPGet.Host,
			Scheme:      string(objProbe.ProbeHandler.HTTPGet.Scheme),
			HttpHeaders: ObjProbeHTTPHeaders,
		}
		HandlerObj.Handler = &k8sapiv1.Probe_HttpGet{
			HttpGet: ObjProbeHTTPObject,
		}
	}

	if objProbe.ProbeHandler.Exec != nil {
		ObjProbeExec := &k8sapiv1.ExecAction{
			Command: objProbe.ProbeHandler.Exec.Command,
		}
		HandlerObj.Handler = &k8sapiv1.Probe_Exec{
			Exec: ObjProbeExec,
		}
	}

	if objProbe.ProbeHandler.TCPSocket != nil {
		ObjProbeTCPSocket := &k8sapiv1.TCPSocketAction{
			Port: objProbe.ProbeHandler.TCPSocket.Port.IntVal,
			Host: objProbe.ProbeHandler.TCPSocket.Host,
		}
		HandlerObj.Handler = &k8sapiv1.Probe_TcpSocket{
			TcpSocket: ObjProbeTCPSocket,
		}
	}

	if objProbe.ProbeHandler.GRPC != nil {
		ObjProbeGRPC := &k8sapiv1.GRPCAction{
			Port:    objProbe.ProbeHandler.GRPC.Port,
			Service: *objProbe.ProbeHandler.GRPC.Service,
		}
		HandlerObj.Handler = &k8sapiv1.Probe_Grpc{
			Grpc: ObjProbeGRPC,
		}
	}
	return HandlerObj
}

func ProtoForDeploymentSpec(deploymentSpec appsv1.DeploymentSpec) *k8sapiv1.Deployment_DeploymentSpec {
	deploymentContainers := make([]*k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container, 0, len(deploymentSpec.Template.Spec.Containers))
	for _, container := range deploymentSpec.Template.Spec.Containers {
		resourceLimits := make(map[string]string, len(container.Resources.Limits))
		resourceRequests := make(map[string]string, len(container.Resources.Requests))

		for res, quantity := range container.Resources.Limits {
			resourceLimits[string(res)] = quantity.String()
		}

		for res, quantity := range container.Resources.Requests {
			resourceRequests[string(res)] = quantity.String()
		}

		LivenessProbeObject := &k8sapiv1.Probe{}
		ReadinessProbeObject := &k8sapiv1.Probe{}
		if container.LivenessProbe != nil {
			HandlerObj := processObjProbe(container.LivenessProbe)

			LivenessProbeObject = &k8sapiv1.Probe{
				InitialDelaySeconds:           &container.LivenessProbe.InitialDelaySeconds,
				TimeoutSeconds:                &container.LivenessProbe.TimeoutSeconds,
				PeriodSeconds:                 &container.LivenessProbe.PeriodSeconds,
				SuccessThreshold:              &container.LivenessProbe.SuccessThreshold,
				FailureThreshold:              &container.LivenessProbe.FailureThreshold,
				TerminationGracePeriodSeconds: container.LivenessProbe.TerminationGracePeriodSeconds,
				Handler:                       HandlerObj.Handler,
			}
		}

		if container.ReadinessProbe != nil {
			HandlerObj := processObjProbe(container.ReadinessProbe)

			ReadinessProbeObject = &k8sapiv1.Probe{
				InitialDelaySeconds:           &container.ReadinessProbe.InitialDelaySeconds,
				TimeoutSeconds:                &container.ReadinessProbe.TimeoutSeconds,
				PeriodSeconds:                 &container.ReadinessProbe.PeriodSeconds,
				SuccessThreshold:              &container.ReadinessProbe.SuccessThreshold,
				FailureThreshold:              &container.ReadinessProbe.FailureThreshold,
				TerminationGracePeriodSeconds: container.ReadinessProbe.TerminationGracePeriodSeconds,
				Handler:                       HandlerObj.Handler,
			}
		}

		newContainer := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container{
			Name: container.Name,
			Resources: &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_ResourceRequirements{
				Limits:   resourceLimits,
				Requests: resourceRequests,
			},
			LivenessProbe:  LivenessProbeObject,
			ReadinessProbe: ReadinessProbeObject,
		}
		deploymentContainers = append(deploymentContainers, newContainer)
	}
	return &k8sapiv1.Deployment_DeploymentSpec{
		Template: &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec{
			Spec: &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec{
				Containers: deploymentContainers,
			},
		},
	}
}

func ProtoForDeploymentStatus(deploymentStatus appsv1.DeploymentStatus) *k8sapiv1.Deployment_DeploymentStatus {
	var deploymentConditions []*k8sapiv1.Deployment_DeploymentStatus_Condition
	for _, cond := range deploymentStatus.Conditions {
		var deploymentConditionType k8sapiv1.Deployment_DeploymentStatus_Condition_Type
		// TODO: Is this the preferred way of converting from one enum to another?
		if cond.Type != "" {
			deploymentConditionType = k8sapiv1.Deployment_DeploymentStatus_Condition_Type(
				k8sapiv1.Deployment_DeploymentStatus_Condition_Type_value[strings.ToUpper(string(cond.Type))])
		}
		var condStatus k8sapiv1.Deployment_DeploymentStatus_Condition_ConditionStatus
		switch cond.Status {
		case v1.ConditionTrue:
			{
				condStatus = k8sapiv1.Deployment_DeploymentStatus_Condition_CONDITION_TRUE
			}
		case v1.ConditionFalse:
			{
				condStatus = k8sapiv1.Deployment_DeploymentStatus_Condition_CONDITION_FALSE
			}
		default:
			condStatus = k8sapiv1.Deployment_DeploymentStatus_Condition_CONDITION_UNKNOWN
		}

		newCond := &k8sapiv1.Deployment_DeploymentStatus_Condition{
			Type:            deploymentConditionType,
			ConditionStatus: condStatus,
			Reason:          cond.Reason,
			Message:         cond.Message,
		}
		deploymentConditions = append(deploymentConditions, newCond)
	}

	return &k8sapiv1.Deployment_DeploymentStatus{
		Replicas:             uint32(deploymentStatus.Replicas), //nolint
		UpdatedReplicas:      uint32(deploymentStatus.UpdatedReplicas), //nolint
		ReadyReplicas:        uint32(deploymentStatus.ReadyReplicas), //nolint
		AvailableReplicas:    uint32(deploymentStatus.AvailableReplicas), //nolint
		UnavailableReplicas:  uint32(deploymentStatus.UnavailableReplicas), //nolint
		DeploymentConditions: deploymentConditions,
	}
}

func (s *svc) UpdateDeployment(ctx context.Context, clientset, cluster, namespace, name string, fields *k8sapiv1.UpdateDeploymentRequest_Fields) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	getOpts := metav1.GetOptions{}
	oldDeployment, err := cs.AppsV1().Deployments(cs.Namespace()).Get(ctx, name, getOpts)
	if err != nil {
		return err
	}

	newDeployment := oldDeployment.DeepCopy()
	mergeDeploymentLabelsAndAnnotations(newDeployment, fields)
	if err := updateContainerResources(newDeployment, fields); err != nil {
		return err
	}

	if err := updateContainerProbes(newDeployment, fields); err != nil {
		return err
	}

	patchBytes, err := GenerateStrategicPatch(oldDeployment, newDeployment, appsv1.Deployment{})
	if err != nil {
		return err
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := cs.AppsV1().Deployments(cs.Namespace()).Patch(ctx, oldDeployment.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
		return err
	})
	return retryErr
}

func (s *svc) DeleteDeployment(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return cs.AppsV1().Deployments(cs.Namespace()).Delete(ctx, name, opts)
}

func mergeDeploymentLabelsAndAnnotations(deployment *appsv1.Deployment, fields *k8sapiv1.UpdateDeploymentRequest_Fields) {
	if len(fields.Labels) > 0 {
		deployment.Labels = labels.Merge(labels.Set(deployment.Labels), labels.Set(fields.Labels))
		deployment.Spec.Template.ObjectMeta.Labels = labels.Merge(labels.Set(deployment.Spec.Template.ObjectMeta.Labels), labels.Set(fields.Labels))
	}

	if len(fields.Annotations) > 0 {
		deployment.Annotations = labels.Merge(labels.Set(deployment.Annotations), labels.Set(fields.Annotations))
		deployment.Spec.Template.ObjectMeta.Annotations = labels.Merge(labels.Set(deployment.Spec.Template.ObjectMeta.Annotations), labels.Set(fields.Annotations))
	}
}

func updateContainerResources(deployment *appsv1.Deployment, fields *k8sapiv1.UpdateDeploymentRequest_Fields) error {
	for _, containerResource := range fields.ContainerResources {
		for _, container := range deployment.Spec.Template.Spec.Containers {
			if container.Name == containerResource.ContainerName {
				for resourceName := range containerResource.Resources.Limits {
					quantity, err := resource.ParseQuantity(containerResource.Resources.Limits[resourceName])
					if err != nil {
						return err
					}
					container.Resources.Limits[v1.ResourceName(resourceName)] = quantity
				}
				for resourceName := range containerResource.Resources.Requests {
					quantity, err := resource.ParseQuantity(containerResource.Resources.Requests[resourceName])
					if err != nil {
						return err
					}
					container.Resources.Requests[v1.ResourceName(resourceName)] = quantity
				}
			}
		}
	}
	return nil
}

func setOptionalInt64Value(source *int64, target *int64) {
	if source != nil {
		*target = *source
	}
}

func setOptionalValue(source *int32, target *int32) {
	if source != nil {
		*target = *source
	}
}

func updateContainerProbes(deployment *appsv1.Deployment, fields *k8sapiv1.UpdateDeploymentRequest_Fields) error {
	for _, containerProbes := range fields.ContainerProbes {
		for _, container := range deployment.Spec.Template.Spec.Containers {
			if container.Name == containerProbes.ContainerName {
				if containerProbes.LivenessProbe != nil {
					resourceProbe := containerProbes.LivenessProbe
					setOptionalValue(resourceProbe.InitialDelaySeconds, &container.LivenessProbe.InitialDelaySeconds)
					setOptionalValue(resourceProbe.PeriodSeconds, &container.LivenessProbe.PeriodSeconds)
					setOptionalValue(resourceProbe.TimeoutSeconds, &container.LivenessProbe.TimeoutSeconds)
					setOptionalValue(resourceProbe.SuccessThreshold, &container.LivenessProbe.SuccessThreshold)
					setOptionalValue(resourceProbe.FailureThreshold, &container.LivenessProbe.FailureThreshold)
					setOptionalInt64Value(resourceProbe.TerminationGracePeriodSeconds, container.LivenessProbe.TerminationGracePeriodSeconds)
					if handler := resourceProbe.Handler; handler != nil {
						switch resourceProbe.Handler.(type) {
						case *k8sapiv1.Probe_Exec:
							container.LivenessProbe.ProbeHandler.Exec.Command = resourceProbe.GetExec().Command
						case *k8sapiv1.Probe_Grpc:
							container.LivenessProbe.ProbeHandler.GRPC.Port = resourceProbe.GetGrpc().Port
							container.LivenessProbe.ProbeHandler.GRPC.Service = &resourceProbe.GetGrpc().Service
						case *k8sapiv1.Probe_TcpSocket:
							container.LivenessProbe.ProbeHandler.TCPSocket.Port.IntVal = resourceProbe.GetTcpSocket().Port
							container.LivenessProbe.ProbeHandler.TCPSocket.Host = resourceProbe.GetTcpSocket().Host
						case *k8sapiv1.Probe_HttpGet:
							container.LivenessProbe.ProbeHandler.HTTPGet.Host = resourceProbe.GetHttpGet().Host
							container.LivenessProbe.ProbeHandler.HTTPGet.Path = resourceProbe.GetHttpGet().Path
							container.LivenessProbe.ProbeHandler.HTTPGet.Port.IntVal = resourceProbe.GetHttpGet().Port
							container.LivenessProbe.ProbeHandler.HTTPGet.Scheme = (v1.URIScheme)(resourceProbe.GetHttpGet().Scheme)
							LivenessProbeHTTPHeaders := make([]v1.HTTPHeader, 0, len(resourceProbe.GetHttpGet().HttpHeaders))
							for _, value := range resourceProbe.GetHttpGet().HttpHeaders {
								UniqueLivenessHeader := v1.HTTPHeader{
									Name:  *value.Name,
									Value: *value.Value,
								}
								LivenessProbeHTTPHeaders = append(LivenessProbeHTTPHeaders, UniqueLivenessHeader)
							}
							container.LivenessProbe.ProbeHandler.HTTPGet.HTTPHeaders = LivenessProbeHTTPHeaders
						}
					}
				}
				if containerProbes.ReadinessProbe == nil {
					return nil
				}
				resourceReadinessProbe := containerProbes.ReadinessProbe
				setOptionalValue(resourceReadinessProbe.InitialDelaySeconds, &container.ReadinessProbe.InitialDelaySeconds)
				setOptionalValue(resourceReadinessProbe.PeriodSeconds, &container.ReadinessProbe.PeriodSeconds)
				setOptionalValue(resourceReadinessProbe.TimeoutSeconds, &container.ReadinessProbe.TimeoutSeconds)
				setOptionalValue(resourceReadinessProbe.SuccessThreshold, &container.ReadinessProbe.SuccessThreshold)
				setOptionalValue(resourceReadinessProbe.FailureThreshold, &container.ReadinessProbe.FailureThreshold)
				setOptionalInt64Value(resourceReadinessProbe.TerminationGracePeriodSeconds, container.ReadinessProbe.TerminationGracePeriodSeconds)
				if handler := resourceReadinessProbe.Handler; handler != nil {
					switch resourceReadinessProbe.Handler.(type) {
					case *k8sapiv1.Probe_Exec:
						container.ReadinessProbe.ProbeHandler.Exec.Command = resourceReadinessProbe.GetExec().Command
					case *k8sapiv1.Probe_Grpc:
						container.ReadinessProbe.ProbeHandler.GRPC.Port = resourceReadinessProbe.GetGrpc().Port
						container.ReadinessProbe.ProbeHandler.GRPC.Service = &resourceReadinessProbe.GetGrpc().Service
					case *k8sapiv1.Probe_TcpSocket:
						container.ReadinessProbe.ProbeHandler.TCPSocket.Port.IntVal = resourceReadinessProbe.GetTcpSocket().Port
						container.ReadinessProbe.ProbeHandler.TCPSocket.Host = resourceReadinessProbe.GetTcpSocket().Host
					case *k8sapiv1.Probe_HttpGet:
						container.ReadinessProbe.ProbeHandler.HTTPGet.Host = resourceReadinessProbe.GetHttpGet().Host
						container.ReadinessProbe.ProbeHandler.HTTPGet.Path = resourceReadinessProbe.GetHttpGet().Path
						container.ReadinessProbe.ProbeHandler.HTTPGet.Port.IntVal = resourceReadinessProbe.GetHttpGet().Port
						container.ReadinessProbe.ProbeHandler.HTTPGet.Scheme = (v1.URIScheme)(resourceReadinessProbe.GetHttpGet().Scheme)
					}
				}
			}
		}
	}
	return nil
}
