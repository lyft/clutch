package k8s

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

var t1 = time.Date(2021, time.May, 4, 0, 0, 0, 0, time.UTC)

func testDeploymentClientset() k8s.Interface {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "testing-deployment-name",
			Namespace:         "testing-namespace",
			Labels:            map[string]string{"foo": "bar"},
			Annotations:       map[string]string{"baz": "quuz"},
			CreationTimestamp: metav1.NewTime(t1),
		},
	}

	return fake.NewSimpleClientset(deployment)
}

func TestDescribeDeployment(t *testing.T) {
	cs := testDeploymentClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": {
				Interface: cs,
				namespace: "default",
				cluster:   "core-testing",
			}},
		},
	}

	d, err := s.DescribeDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "testing-deployment-name")
	assert.NoError(t, err)
	assert.Equal(t, "testing-deployment-name", d.Name)
	assert.NotNil(t, d.CreationTimeMillis)
	// Not found.
	_, err = s.DescribeDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "abc")
	assert.NoError(t, err)
}

func TestListDeployments(t *testing.T) {
	cs := testDeploymentClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": {
				Interface: cs,
				namespace: "default",
				cluster:   "core-testing",
			}},
		},
	}

	opts := &k8sapiv1.ListOptions{Labels: map[string]string{"foo": "bar"}}
	list, err := s.ListDeployments(context.Background(), "foo", "core-testing", "testing-namespace", opts)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	// Not Found
	opts = &k8sapiv1.ListOptions{Labels: map[string]string{"unknown": "bar"}}
	list, err = s.ListDeployments(context.Background(), "foo", "core-testing", "testing-namespace", opts)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(list))
}

func TestUpdateDeployment(t *testing.T) {
	cs := testDeploymentClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": {
				Interface: cs,
				namespace: "default",
				cluster:   "core-testing",
			}},
		},
	}

	// Not found.
	err := s.UpdateDeployment(context.Background(), "", "", "", "", nil)
	assert.Error(t, err)

	err = s.UpdateDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "testing-deployment-name", &k8sapiv1.UpdateDeploymentRequest_Fields{})
	assert.NoError(t, err)
}

func TestDeleteDeployment(t *testing.T) {
	cs := testDeploymentClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": {
				Interface: cs,
				namespace: "default",
				cluster:   "core-testing",
			}},
		},
	}

	// Not found.
	err := s.DeleteDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "abc")
	assert.Error(t, err)

	err = s.DeleteDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "testing-deployment-name")
	assert.NoError(t, err)

	// Not found.
	_, err = s.DescribeDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "testing-deployment-name")
	assert.Error(t, err)
}

func TestMergeDeploymentLabelsAndAnnotations(t *testing.T) {
	t.Parallel()

	mergeDeploymentLabelAnnotationsTestCases := []struct {
		id     string
		fields *k8sapiv1.UpdateDeploymentRequest_Fields
		expect *appsv1.Deployment
	}{
		{
			id:     "no changes",
			fields: &k8sapiv1.UpdateDeploymentRequest_Fields{},
			expect: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"foo": "bar"},
					Annotations: map[string]string{"baz": "quuz"},
				},
			},
		},
		{
			id: "Adding labels",
			fields: &k8sapiv1.UpdateDeploymentRequest_Fields{
				Labels: map[string]string{"field1": "field1val"},
			},
			expect: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"foo":    "bar",
						"field1": "field1val",
					},
					Annotations: map[string]string{
						"baz": "quuz",
					},
				},
			},
		},
		{
			id: "Adding annotations",
			fields: &k8sapiv1.UpdateDeploymentRequest_Fields{
				Annotations: map[string]string{"field1": "field1val"},
			},
			expect: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"foo": "bar",
					},
					Annotations: map[string]string{
						"baz":    "quuz",
						"field1": "field1val",
					},
				},
			},
		},
		{
			id: "Adding both labels and annotations",
			fields: &k8sapiv1.UpdateDeploymentRequest_Fields{
				Labels:      map[string]string{"label1": "label1val"},
				Annotations: map[string]string{"annotation1": "annotation1val"},
			},
			expect: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"foo":    "bar",
						"label1": "label1val",
					},
					Annotations: map[string]string{
						"baz":         "quuz",
						"annotation1": "annotation1val",
					},
				},
			},
		},
	}

	for _, tt := range mergeDeploymentLabelAnnotationsTestCases {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"foo": "bar"},
					Annotations: map[string]string{"baz": "quuz"},
				},
			}

			mergeDeploymentLabelsAndAnnotations(deployment, tt.fields)
			assert.Equal(t, tt.expect.ObjectMeta, deployment.ObjectMeta)
		})
	}
}

func TestProtoForDeploymentClusterName(t *testing.T) {
	t.Parallel()

	deploymentTestCases := []struct {
		id                  string
		inputClusterName    string
		expectedClusterName string
		deployment          *appsv1.Deployment
	}{
		{
			id:                  "clustername already set",
			inputClusterName:    "notprod",
			expectedClusterName: "production",
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clutchLabelClusterName: "production",
					},
				},
			},
		},
		{
			id:                  "custername is not set",
			inputClusterName:    "staging",
			expectedClusterName: "staging",
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clutchLabelClusterName: "",
					},
				},
			},
		},
	}

	for _, tt := range deploymentTestCases {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			deployment := ProtoForDeployment(tt.inputClusterName, tt.deployment)
			assert.Equal(t, tt.expectedClusterName, deployment.Cluster)
		})
	}
}

func TestProtoForDeploymentSpec(t *testing.T) {
	t.Parallel()

	deploymentTestCases := []struct {
		id         string
		deployment *appsv1.Deployment
	}{
		{
			id: "foo",
			deployment: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Resources: v1.ResourceRequirements{
										Limits: v1.ResourceList{
											"cpu":    resource.MustParse("500m"),
											"memory": resource.MustParse("128Mi"),
										},
										Requests: v1.ResourceList{
											"cpu":    resource.MustParse("250m"),
											"memory": resource.MustParse("64Mi"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range deploymentTestCases {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			deployment := ProtoForDeployment("", tt.deployment)
			assert.Equal(t, tt.deployment.Spec.Template.Spec.Containers[0].Resources.Limits["cpu"], resource.MustParse(deployment.DeploymentSpec.Template.Spec.Containers[0].Resources.Limits["cpu"]))
			assert.Equal(t, tt.deployment.Spec.Template.Spec.Containers[0].Resources.Limits["memory"], resource.MustParse(deployment.DeploymentSpec.Template.Spec.Containers[0].Resources.Limits["memory"]))
			assert.Equal(t, tt.deployment.Spec.Template.Spec.Containers[0].Resources.Requests["cpu"], resource.MustParse(deployment.DeploymentSpec.Template.Spec.Containers[0].Resources.Requests["cpu"]))
			assert.Equal(t, tt.deployment.Spec.Template.Spec.Containers[0].Resources.Requests["memory"], resource.MustParse(deployment.DeploymentSpec.Template.Spec.Containers[0].Resources.Requests["memory"]))
		})
	}
}

func TestProtoForDeploymentStatus(t *testing.T) {
	t.Parallel()

	deploymentTestCases := []struct {
		id         string
		deployment *appsv1.Deployment
	}{
		{
			id: "foo",
			deployment: &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					Replicas:            60,
					UpdatedReplicas:     60,
					AvailableReplicas:   10,
					UnavailableReplicas: 20,
					Conditions: []appsv1.DeploymentCondition{
						{
							Type:    appsv1.DeploymentAvailable,
							Status:  v1.ConditionTrue,
							Reason:  "reason",
							Message: "message",
						},
					},
				},
			},
		},
		{
			id: "oof",
			deployment: &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					Replicas:            40,
					UpdatedReplicas:     31,
					AvailableReplicas:   3,
					UnavailableReplicas: 8,
					Conditions: []appsv1.DeploymentCondition{
						{
							Type:    appsv1.DeploymentReplicaFailure,
							Status:  v1.ConditionFalse,
							Reason:  "space",
							Message: "moon",
						},
					},
				},
			},
		},
	}

	for _, tt := range deploymentTestCases {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			deployment := ProtoForDeployment("", tt.deployment)
			assert.Equal(t, tt.deployment.Status.Replicas, int32(deployment.DeploymentStatus.Replicas))
			assert.Equal(t, tt.deployment.Status.UpdatedReplicas, int32(deployment.DeploymentStatus.UpdatedReplicas))
			assert.Equal(t, tt.deployment.Status.ReadyReplicas, int32(deployment.DeploymentStatus.ReadyReplicas))
			assert.Equal(t, tt.deployment.Status.AvailableReplicas, int32(deployment.DeploymentStatus.AvailableReplicas))
			assert.Equal(t, len(tt.deployment.Status.Conditions), len(deployment.DeploymentStatus.DeploymentConditions))
		})
	}
}

func TestProtoForDeploymentSpecWithProbesLivenessHTTPGet(t *testing.T) {
	t.Parallel()
	var terminationVar int64 = 30
	var deploymentTestCase = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							LivenessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									HTTPGet: &v1.HTTPGetAction{
										Path: "/",
										Port: intstr.IntOrString{
											IntVal: 8080,
										},
									},
								},
								InitialDelaySeconds:           10,
								PeriodSeconds:                 30,
								TimeoutSeconds:                1,
								SuccessThreshold:              1,
								FailureThreshold:              3,
								TerminationGracePeriodSeconds: &terminationVar,
							},
						},
					},
				},
			},
		},
	}

	t.Run("LivenessProbeHttpGetAction", func(t *testing.T) {
		t.Parallel()
		deployment := ProtoForDeployment("", deploymentTestCase)
		assert.Equal(t, deploymentTestCase.Spec.Template.Spec.Containers[0].LivenessProbe.InitialDelaySeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].LivenessProbe.InitialDelaySeconds)
		assert.Equal(t, deploymentTestCase.Spec.Template.Spec.Containers[0].LivenessProbe.PeriodSeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].LivenessProbe.PeriodSeconds)

		err := updateContainerProbes(deploymentTestCase, &k8sapiv1.UpdateDeploymentRequest_Fields{})
		assert.NoError(t, err)
		err = updateContainerProbes(deploymentTestCase, &k8sapiv1.UpdateDeploymentRequest_Fields{
			ContainerProbes: []*k8sapiv1.UpdateDeploymentRequest_Fields_ContainerProbes{
				{
					LivenessProbe: &k8sapiv1.Probe{
						Handler: &k8sapiv1.Probe_HttpGet{
							HttpGet: &k8sapiv1.HTTPGetAction{
								Port: 8081,
							},
						},
						InitialDelaySeconds: newInt32(20),
						PeriodSeconds:       newInt32(15),
					},
				},
			},
		})
		assert.NoError(t, err)
	})
}

func TestProtoForDeploymentSpecWithProbesLivenessExec(t *testing.T) {
	t.Parallel()
	var terminationVar int64 = 30
	var deploymentTestCases = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							LivenessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									Exec: &v1.ExecAction{
										Command: []string{"ls", "-l"},
									},
								},
								InitialDelaySeconds:           5,
								PeriodSeconds:                 25,
								TimeoutSeconds:                5,
								SuccessThreshold:              4,
								FailureThreshold:              8,
								TerminationGracePeriodSeconds: &terminationVar,
							},
						},
					},
				},
			},
		},
	}

	t.Run("LivenessProbeExecAction", func(t *testing.T) {
		t.Parallel()
		deployment := ProtoForDeployment("", deploymentTestCases)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].LivenessProbe.InitialDelaySeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].LivenessProbe.InitialDelaySeconds)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].LivenessProbe.PeriodSeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].LivenessProbe.PeriodSeconds)

		err := updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{})
		assert.NoError(t, err)
		err = updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{
			ContainerProbes: []*k8sapiv1.UpdateDeploymentRequest_Fields_ContainerProbes{
				{
					LivenessProbe: &k8sapiv1.Probe{
						Handler: &k8sapiv1.Probe_Exec{
							Exec: &k8sapiv1.ExecAction{
								Command: []string{"ps"},
							},
						},
						InitialDelaySeconds: newInt32(20),
						PeriodSeconds:       newInt32(15),
					},
				},
			},
		})
		assert.NoError(t, err)
	})
}

func TestProtoForDeploymentSpecWithProbesLivenessTCP(t *testing.T) {
	t.Parallel()
	var terminationVar int64 = 30
	var deploymentTestCases = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							LivenessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									TCPSocket: &v1.TCPSocketAction{
										Port: intstr.IntOrString{
											IntVal: 8080,
										},
										Host: "/",
									},
								},
								InitialDelaySeconds:           6,
								PeriodSeconds:                 26,
								TimeoutSeconds:                6,
								SuccessThreshold:              10,
								FailureThreshold:              9,
								TerminationGracePeriodSeconds: &terminationVar,
							},
						},
					},
				},
			},
		},
	}

	t.Run("LivenessProbeTCPAction", func(t *testing.T) {
		t.Parallel()
		deployment := ProtoForDeployment("", deploymentTestCases)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].LivenessProbe.InitialDelaySeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].LivenessProbe.InitialDelaySeconds)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].LivenessProbe.PeriodSeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].LivenessProbe.PeriodSeconds)

		err := updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{})
		assert.NoError(t, err)
		err = updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{
			ContainerProbes: []*k8sapiv1.UpdateDeploymentRequest_Fields_ContainerProbes{
				{
					LivenessProbe: &k8sapiv1.Probe{
						Handler: &k8sapiv1.Probe_TcpSocket{
							TcpSocket: &k8sapiv1.TCPSocketAction{
								Port: 8081,
							},
						},
						InitialDelaySeconds: newInt32(20),
						PeriodSeconds:       newInt32(15),
					},
				},
			},
		})
		assert.NoError(t, err)
	})
}

func TestProtoForDeploymentSpecWithProbesLivenessGRPC(t *testing.T) {
	t.Parallel()
	var terminationVar int64 = 30
	var portGRPC int32 = 8080
	var serviceGRPC string = "service"
	var deploymentTestCases = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							LivenessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									GRPC: &v1.GRPCAction{
										Port:    portGRPC,
										Service: &serviceGRPC,
									},
								},
								InitialDelaySeconds:           4,
								PeriodSeconds:                 24,
								TimeoutSeconds:                2,
								SuccessThreshold:              4,
								FailureThreshold:              5,
								TerminationGracePeriodSeconds: &terminationVar,
							},
						},
					},
				},
			},
		},
	}

	t.Run("LivenessProbeGRPCAction", func(t *testing.T) {
		t.Parallel()
		deployment := ProtoForDeployment("", deploymentTestCases)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].LivenessProbe.InitialDelaySeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].LivenessProbe.InitialDelaySeconds)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].LivenessProbe.PeriodSeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].LivenessProbe.PeriodSeconds)

		err := updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{})
		assert.NoError(t, err)
		err = updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{
			ContainerProbes: []*k8sapiv1.UpdateDeploymentRequest_Fields_ContainerProbes{
				{
					LivenessProbe: &k8sapiv1.Probe{
						Handler: &k8sapiv1.Probe_Grpc{
							Grpc: &k8sapiv1.GRPCAction{
								Service: "tmp",
							},
						},
						InitialDelaySeconds: newInt32(20),
						PeriodSeconds:       newInt32(15),
					},
				},
			},
		})
		assert.NoError(t, err)
	})
}

func TestProtoForDeploymentSpecWithProbesReadinessHTTPGet(t *testing.T) {
	t.Parallel()
	var terminationVar int64 = 30
	var deploymentTestCases = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							ReadinessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									HTTPGet: &v1.HTTPGetAction{
										Path: "/",
										Port: intstr.IntOrString{
											IntVal: 8080,
										},
									},
								},
								InitialDelaySeconds:           10,
								PeriodSeconds:                 30,
								TimeoutSeconds:                1,
								SuccessThreshold:              1,
								FailureThreshold:              3,
								TerminationGracePeriodSeconds: &terminationVar,
							},
						},
					},
				},
			},
		},
	}

	t.Run("ReadinessProbeHttpGetAction", func(t *testing.T) {
		t.Parallel()
		deployment := ProtoForDeployment("", deploymentTestCases)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.InitialDelaySeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].ReadinessProbe.InitialDelaySeconds)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.PeriodSeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].ReadinessProbe.PeriodSeconds)

		probeDeployment := processObjProbe(deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.HTTPGet.Host, probeDeployment.GetHttpGet().Host)

		err := updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{})
		assert.NoError(t, err)
		err = updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{
			ContainerProbes: []*k8sapiv1.UpdateDeploymentRequest_Fields_ContainerProbes{
				{
					ReadinessProbe: &k8sapiv1.Probe{
						Handler: &k8sapiv1.Probe_HttpGet{
							HttpGet: &k8sapiv1.HTTPGetAction{
								Host: "/test",
							},
						},
						InitialDelaySeconds: newInt32(20),
						PeriodSeconds:       newInt32(15),
					},
				},
			},
		})
		assert.NoError(t, err)
	})
}

func TestProtoForDeploymentSpecWithProbesReadinessExec(t *testing.T) {
	t.Parallel()
	var terminationVar int64 = 30
	var deploymentTestCases = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							ReadinessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									Exec: &v1.ExecAction{
										Command: []string{"ls", "-l"},
									},
								},
								InitialDelaySeconds:           5,
								PeriodSeconds:                 25,
								TimeoutSeconds:                5,
								SuccessThreshold:              4,
								FailureThreshold:              8,
								TerminationGracePeriodSeconds: &terminationVar,
							},
						},
					},
				},
			},
		},
	}

	t.Run("ReadinessProbeExecAction", func(t *testing.T) {
		t.Parallel()
		deployment := ProtoForDeployment("", deploymentTestCases)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.InitialDelaySeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].ReadinessProbe.InitialDelaySeconds)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.PeriodSeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].ReadinessProbe.PeriodSeconds)

		probeDeployment := processObjProbe(deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.Exec.Command, probeDeployment.GetExec().Command)

		err := updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{})
		assert.NoError(t, err)
		err = updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{
			ContainerProbes: []*k8sapiv1.UpdateDeploymentRequest_Fields_ContainerProbes{
				{
					ReadinessProbe: &k8sapiv1.Probe{
						Handler: &k8sapiv1.Probe_Exec{
							Exec: &k8sapiv1.ExecAction{
								Command: []string{"pwd"},
							},
						},
						InitialDelaySeconds: newInt32(20),
						PeriodSeconds:       newInt32(15),
					},
				},
			},
		})
		assert.NoError(t, err)
	})
}

func TestProtoForDeploymentSpecWithProbesReadinessTCP(t *testing.T) {
	t.Parallel()
	var terminationVar int64 = 30
	var deploymentTestCases = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							ReadinessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									TCPSocket: &v1.TCPSocketAction{
										Port: intstr.IntOrString{
											IntVal: 8080,
										},
										Host: "/",
									},
								},
								InitialDelaySeconds:           6,
								PeriodSeconds:                 26,
								TimeoutSeconds:                6,
								SuccessThreshold:              10,
								FailureThreshold:              9,
								TerminationGracePeriodSeconds: &terminationVar,
							},
						},
					},
				},
			},
		},
	}

	t.Run("ReadinessProbeTCPAction", func(t *testing.T) {
		t.Parallel()
		deployment := ProtoForDeployment("", deploymentTestCases)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.InitialDelaySeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].ReadinessProbe.InitialDelaySeconds)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.PeriodSeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].ReadinessProbe.PeriodSeconds)

		probeDeployment := processObjProbe(deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.TCPSocket.Host, probeDeployment.GetTcpSocket().Host)

		err := updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{})
		assert.NoError(t, err)
		err = updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{
			ContainerProbes: []*k8sapiv1.UpdateDeploymentRequest_Fields_ContainerProbes{
				{
					ReadinessProbe: &k8sapiv1.Probe{
						Handler: &k8sapiv1.Probe_TcpSocket{
							TcpSocket: &k8sapiv1.TCPSocketAction{
								Host: "/test",
							},
						},
						InitialDelaySeconds: newInt32(20),
						PeriodSeconds:       newInt32(15),
					},
				},
			},
		})
		assert.NoError(t, err)
	})
}

func TestProtoForDeploymentSpecWithProbesReadinessGRPC(t *testing.T) {
	t.Parallel()
	var terminationVar int64 = 30
	var portGRPC int32 = 8080
	var serviceGRPC string = "service"
	var deploymentTestCases = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							ReadinessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									GRPC: &v1.GRPCAction{
										Port:    portGRPC,
										Service: &serviceGRPC,
									},
								},
								InitialDelaySeconds:           4,
								PeriodSeconds:                 24,
								TimeoutSeconds:                2,
								SuccessThreshold:              4,
								FailureThreshold:              5,
								TerminationGracePeriodSeconds: &terminationVar,
							},
						},
					},
				},
			},
		},
	}

	t.Run("ReadinessProbeGRPCAction", func(t *testing.T) {
		t.Parallel()
		deployment := ProtoForDeployment("", deploymentTestCases)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.InitialDelaySeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].ReadinessProbe.InitialDelaySeconds)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.PeriodSeconds, *deployment.DeploymentSpec.Template.Spec.Containers[0].ReadinessProbe.PeriodSeconds)

		probeDeployment := processObjProbe(deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe)
		assert.Equal(t, deploymentTestCases.Spec.Template.Spec.Containers[0].ReadinessProbe.GRPC.Port, probeDeployment.GetGrpc().Port)

		err := updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{})
		assert.NoError(t, err)
		err = updateContainerProbes(deploymentTestCases, &k8sapiv1.UpdateDeploymentRequest_Fields{
			ContainerProbes: []*k8sapiv1.UpdateDeploymentRequest_Fields_ContainerProbes{
				{
					ReadinessProbe: &k8sapiv1.Probe{
						Handler: &k8sapiv1.Probe_Grpc{
							Grpc: &k8sapiv1.GRPCAction{
								Service: "tmp",
							},
						},
						InitialDelaySeconds: newInt32(20),
						PeriodSeconds:       newInt32(15),
					},
				},
			},
		})
		assert.NoError(t, err)
	})
}
