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
	clusterName := deployment.ClusterName
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

		LivenessProbeObject := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_Probe{}
		ReadinessProbeObject := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_Probe{}
		if container.LivenessProbe != nil {
			LivenessProbeHTTPObject := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_HTTPGetAction{}
			LivenessProbeHTTPHeaders := make([]*k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_HTTPHeader, 0)
			if container.LivenessProbe.ProbeHandler.HTTPGet != nil {
				if container.LivenessProbe.ProbeHandler.HTTPGet.HTTPHeaders != nil {
					LivenessProbeHTTPHeaders := make([]*k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_HTTPHeader, 0, len(container.LivenessProbe.ProbeHandler.HTTPGet.HTTPHeaders))
					for _, value := range container.LivenessProbe.ProbeHandler.HTTPGet.HTTPHeaders {
						UniqueLivenessHeader := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_HTTPHeader{
							Name:  &value.Name,
							Value: &value.Value,
						}
						LivenessProbeHTTPHeaders = append(LivenessProbeHTTPHeaders, UniqueLivenessHeader)
					}
				}
				LivenessProbeHTTPObject = &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_HTTPGetAction{
					Path:        &container.LivenessProbe.ProbeHandler.HTTPGet.Path,
					Port:        &container.LivenessProbe.ProbeHandler.HTTPGet.Port.IntVal,
					Host:        &container.LivenessProbe.ProbeHandler.HTTPGet.Host,
					Scheme:      (*string)(&container.LivenessProbe.ProbeHandler.HTTPGet.Scheme),
					HttpHeaders: LivenessProbeHTTPHeaders,
				}
			}
			LivenessProbeExec := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_ExecAction{}
			if container.LivenessProbe.ProbeHandler.Exec != nil {
				LivenessProbeExec = &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_ExecAction{
					Command: container.LivenessProbe.ProbeHandler.Exec.Command,
				}
			}

			LivenessProbeTCPSocket := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_TCPSocketAction{}
			if container.LivenessProbe.ProbeHandler.TCPSocket != nil {
				LivenessProbeTCPSocket = &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_TCPSocketAction{
					Port: &container.LivenessProbe.ProbeHandler.TCPSocket.Port.IntVal,
					Host: &container.LivenessProbe.ProbeHandler.TCPSocket.Host,
				}
			}

			LivenessProbeGRPC := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_GRPCAction{}
			if container.LivenessProbe.ProbeHandler.GRPC != nil {
				LivenessProbeGRPC = &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_GRPCAction{
					Port:    &container.LivenessProbe.ProbeHandler.GRPC.Port,
					Service: container.LivenessProbe.ProbeHandler.GRPC.Service,
				}
			}

			LivenessProbeObject = &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_Probe{
				InitialDelaySeconds:           &container.LivenessProbe.InitialDelaySeconds,
				TimeoutSeconds:                &container.LivenessProbe.TimeoutSeconds,
				PeriodSeconds:                 &container.LivenessProbe.PeriodSeconds,
				SuccessThreshold:              &container.LivenessProbe.SuccessThreshold,
				FailureThreshold:              &container.LivenessProbe.FailureThreshold,
				TerminationGracePeriodSeconds: container.LivenessProbe.TerminationGracePeriodSeconds,
				Handler: &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_ProbeHandler{
					Exec:      LivenessProbeExec,
					HttpGet:   LivenessProbeHTTPObject,
					TcpSocket: LivenessProbeTCPSocket,
					Grpc:      LivenessProbeGRPC,
				},
			}
		}

		if container.ReadinessProbe != nil {
			ReadinessProbeHTTPObject := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_HTTPGetAction{}
			ReadinessProbeHTTPHeaders := make([]*k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_HTTPHeader, 0)
			if container.ReadinessProbe.ProbeHandler.HTTPGet != nil {
				if container.ReadinessProbe.ProbeHandler.HTTPGet.HTTPHeaders != nil {
					ReadinessProbeHTTPHeaders := make([]*k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_HTTPHeader, 0, len(container.ReadinessProbe.ProbeHandler.HTTPGet.HTTPHeaders))
					for _, value := range container.ReadinessProbe.ProbeHandler.HTTPGet.HTTPHeaders {
						UniqueReadnessHeader := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_HTTPHeader{
							Name:  &value.Name,
							Value: &value.Value,
						}
						ReadinessProbeHTTPHeaders = append(ReadinessProbeHTTPHeaders, UniqueReadnessHeader)
					}
				}
				ReadinessProbeHTTPObject = &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_HTTPGetAction{
					Path:        &container.ReadinessProbe.ProbeHandler.HTTPGet.Path,
					Port:        &container.ReadinessProbe.ProbeHandler.HTTPGet.Port.IntVal,
					Host:        &container.ReadinessProbe.ProbeHandler.HTTPGet.Host,
					Scheme:      (*string)(&container.ReadinessProbe.ProbeHandler.HTTPGet.Scheme),
					HttpHeaders: ReadinessProbeHTTPHeaders,
				}
			}
			ReadinessProbeExec := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_ExecAction{}
			if container.ReadinessProbe.ProbeHandler.Exec != nil {
				ReadinessProbeExec = &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_ExecAction{
					Command: container.ReadinessProbe.ProbeHandler.Exec.Command,
				}
			}

			ReadinessProbeTCPSocket := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_TCPSocketAction{}
			if container.ReadinessProbe.ProbeHandler.TCPSocket != nil {
				ReadinessProbeTCPSocket = &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_TCPSocketAction{
					Port: &container.ReadinessProbe.ProbeHandler.TCPSocket.Port.IntVal,
					Host: &container.ReadinessProbe.ProbeHandler.TCPSocket.Host,
				}
			}

			ReadinessProbeGRPC := &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_GRPCAction{}
			if container.ReadinessProbe.ProbeHandler.GRPC != nil {
				ReadinessProbeGRPC = &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_GRPCAction{
					Port:    &container.ReadinessProbe.ProbeHandler.GRPC.Port,
					Service: container.ReadinessProbe.ProbeHandler.GRPC.Service,
				}
			}

			ReadinessProbeObject = &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_Probe{
				InitialDelaySeconds:           &container.ReadinessProbe.InitialDelaySeconds,
				TimeoutSeconds:                &container.ReadinessProbe.TimeoutSeconds,
				PeriodSeconds:                 &container.ReadinessProbe.PeriodSeconds,
				SuccessThreshold:              &container.ReadinessProbe.SuccessThreshold,
				FailureThreshold:              &container.ReadinessProbe.FailureThreshold,
				TerminationGracePeriodSeconds: container.ReadinessProbe.TerminationGracePeriodSeconds,
				Handler: &k8sapiv1.Deployment_DeploymentSpec_PodTemplateSpec_PodSpec_Container_ProbeHandler{
					Exec:      ReadinessProbeExec,
					HttpGet:   ReadinessProbeHTTPObject,
					TcpSocket: ReadinessProbeTCPSocket,
					Grpc:      ReadinessProbeGRPC,
				},
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
		Replicas:             uint32(deploymentStatus.Replicas),
		UpdatedReplicas:      uint32(deploymentStatus.UpdatedReplicas),
		ReadyReplicas:        uint32(deploymentStatus.ReadyReplicas),
		AvailableReplicas:    uint32(deploymentStatus.AvailableReplicas),
		UnavailableReplicas:  uint32(deploymentStatus.UnavailableReplicas),
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

func updateContainerProbes(deployment *appsv1.Deployment, fields *k8sapiv1.UpdateDeploymentRequest_Fields) error {
	for _, containerProbes := range fields.ContainerProbes {
		for _, container := range deployment.Spec.Template.Spec.Containers {
			if container.Name == containerProbes.ContainerName {
				if containerProbes.LivenessProbe != nil {
					resourceProbe := containerProbes.LivenessProbe
					if resourceProbe.InitialDelaySeconds != nil {
						container.LivenessProbe.InitialDelaySeconds = *resourceProbe.InitialDelaySeconds
					}
					if resourceProbe.PeriodSeconds != nil {
						container.LivenessProbe.PeriodSeconds = *resourceProbe.PeriodSeconds
					}
					if resourceProbe.TimeoutSeconds != nil {
						container.LivenessProbe.TimeoutSeconds = *resourceProbe.TimeoutSeconds
					}
					if resourceProbe.SuccessThreshold != nil {
						container.LivenessProbe.SuccessThreshold = *resourceProbe.SuccessThreshold
					}
					if resourceProbe.FailureThreshold != nil {
						container.LivenessProbe.FailureThreshold = *resourceProbe.FailureThreshold
					}
					if resourceProbe.TerminationGracePeriodSeconds != nil {
						container.LivenessProbe.TerminationGracePeriodSeconds = resourceProbe.TerminationGracePeriodSeconds
					}
					if resourceProbe.Handler != nil {
						if resourceProbe.Handler.Exec != nil {
							container.LivenessProbe.ProbeHandler.Exec.Command = resourceProbe.Handler.Exec.Command
						}
						if resourceProbe.Handler.HttpGet != nil {
							if resourceProbe.Handler.HttpGet.Host != nil {
								container.LivenessProbe.ProbeHandler.HTTPGet.Host = *resourceProbe.Handler.HttpGet.Host
							}
							if resourceProbe.Handler.HttpGet.Path != nil {
								container.LivenessProbe.ProbeHandler.HTTPGet.Path = *resourceProbe.Handler.HttpGet.Path
							}
							if resourceProbe.Handler.HttpGet.Port != nil {
								container.LivenessProbe.ProbeHandler.HTTPGet.Port.IntVal = *resourceProbe.Handler.HttpGet.Port
							}
							if resourceProbe.Handler.HttpGet.Scheme != nil {
								container.LivenessProbe.ProbeHandler.HTTPGet.Scheme = (v1.URIScheme)(*resourceProbe.Handler.HttpGet.Scheme)
							}
							if resourceProbe.Handler.HttpGet.HttpHeaders != nil {
								LivenessProbeHTTPHeaders := make([]v1.HTTPHeader, 0, len(resourceProbe.Handler.HttpGet.HttpHeaders))
								for _, value := range resourceProbe.Handler.HttpGet.HttpHeaders {
									UniqueLivenessHeader := v1.HTTPHeader{
										Name:  *value.Name,
										Value: *value.Value,
									}
									LivenessProbeHTTPHeaders = append(LivenessProbeHTTPHeaders, UniqueLivenessHeader)
								}
								container.LivenessProbe.ProbeHandler.HTTPGet.HTTPHeaders = LivenessProbeHTTPHeaders
							}
						}
						if resourceProbe.Handler.TcpSocket != nil {
							container.LivenessProbe.ProbeHandler.TCPSocket.Port.IntVal = *resourceProbe.Handler.TcpSocket.Port
							if resourceProbe.Handler.TcpSocket.Host != nil {
								container.LivenessProbe.ProbeHandler.TCPSocket.Host = *resourceProbe.Handler.TcpSocket.Host
							}
						}
						if resourceProbe.Handler.Grpc != nil {
							container.LivenessProbe.ProbeHandler.GRPC.Port = *resourceProbe.Handler.Grpc.Port
							container.LivenessProbe.ProbeHandler.GRPC.Service = resourceProbe.Handler.Grpc.Service
						}
					}
				}
				if containerProbes.ReadinessProbe == nil {
					return nil
				}
				resourceReadinessProbe := containerProbes.ReadinessProbe
				if resourceReadinessProbe.InitialDelaySeconds != nil {
					container.ReadinessProbe.InitialDelaySeconds = *resourceReadinessProbe.InitialDelaySeconds
				}
				if resourceReadinessProbe.PeriodSeconds != nil {
					container.ReadinessProbe.PeriodSeconds = *resourceReadinessProbe.PeriodSeconds
				}
				if resourceReadinessProbe.TimeoutSeconds != nil {
					container.ReadinessProbe.TimeoutSeconds = *resourceReadinessProbe.TimeoutSeconds
				}
				if resourceReadinessProbe.SuccessThreshold != nil {
					container.ReadinessProbe.SuccessThreshold = *resourceReadinessProbe.SuccessThreshold
				}
				if resourceReadinessProbe.FailureThreshold != nil {
					container.ReadinessProbe.FailureThreshold = *resourceReadinessProbe.FailureThreshold
				}
				if resourceReadinessProbe.TerminationGracePeriodSeconds != nil {
					container.ReadinessProbe.TerminationGracePeriodSeconds = resourceReadinessProbe.TerminationGracePeriodSeconds
				}
				if resourceReadinessProbe.Handler != nil {
					if resourceReadinessProbe.Handler.Exec != nil {
						container.ReadinessProbe.ProbeHandler.Exec.Command = resourceReadinessProbe.Handler.Exec.Command
					}
					if resourceReadinessProbe.Handler.HttpGet != nil {
						if resourceReadinessProbe.Handler.HttpGet.Host != nil {
							container.ReadinessProbe.ProbeHandler.HTTPGet.Host = *resourceReadinessProbe.Handler.HttpGet.Host
						}
						if resourceReadinessProbe.Handler.HttpGet.Path != nil {
							container.ReadinessProbe.ProbeHandler.HTTPGet.Path = *resourceReadinessProbe.Handler.HttpGet.Path
						}
						if resourceReadinessProbe.Handler.HttpGet.Port != nil {
							container.ReadinessProbe.ProbeHandler.HTTPGet.Port.IntVal = *resourceReadinessProbe.Handler.HttpGet.Port
						}
						if resourceReadinessProbe.Handler.HttpGet.Scheme != nil {
							container.ReadinessProbe.ProbeHandler.HTTPGet.Scheme = (v1.URIScheme)(*resourceReadinessProbe.Handler.HttpGet.Scheme)
						}
						if resourceReadinessProbe.Handler.HttpGet.HttpHeaders != nil {
							ReadinessProbeHTTPHeaders := make([]v1.HTTPHeader, 0, len(resourceReadinessProbe.Handler.HttpGet.HttpHeaders))
							for _, value := range resourceReadinessProbe.Handler.HttpGet.HttpHeaders {
								UniqueReadinessHeader := v1.HTTPHeader{
									Name:  *value.Name,
									Value: *value.Value,
								}
								ReadinessProbeHTTPHeaders = append(ReadinessProbeHTTPHeaders, UniqueReadinessHeader)
							}
							container.ReadinessProbe.ProbeHandler.HTTPGet.HTTPHeaders = ReadinessProbeHTTPHeaders
						}
					}
					if resourceReadinessProbe.Handler.TcpSocket != nil {
						container.ReadinessProbe.ProbeHandler.TCPSocket.Port.IntVal = *resourceReadinessProbe.Handler.TcpSocket.Port
						if resourceReadinessProbe.Handler.TcpSocket.Host != nil {
							container.ReadinessProbe.ProbeHandler.TCPSocket.Host = *resourceReadinessProbe.Handler.TcpSocket.Host
						}
					}
					if resourceReadinessProbe.Handler.Grpc != nil {
						container.ReadinessProbe.ProbeHandler.GRPC.Port = *resourceReadinessProbe.Handler.Grpc.Port
						container.ReadinessProbe.ProbeHandler.GRPC.Service = resourceReadinessProbe.Handler.Grpc.Service
					}
				}
			}
		}
	}
	return nil
}
