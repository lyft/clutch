package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func testListJobClientset() *fake.Clientset {
	testJobs := []runtime.Object{
		&v1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testing-job-name",
				Namespace: "testing-namespace",
				Labels:    map[string]string{"foo": "bar"},
			},
		},
		&v1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testing-job-name-1",
				Namespace: "testing-namespace",
				Labels:    map[string]string{"foo": "bar"},
			},
		},
		&v1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testing-job-name-2",
				Namespace: "testing-namespace",
				Labels:    map[string]string{"foo1": "bar"},
			},
		},
	}

	return fake.NewSimpleClientset(testJobs...)
}

func testJobClientset() k8s.Interface {
	job := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testing-job-name",
			Namespace: "testing-namespace",
			Labels:    map[string]string{"foo": "bar"},
		},
	}

	return fake.NewSimpleClientset(job)
}

func TestListJobs(t *testing.T) {
	t.Parallel()

	cs := testListJobClientset()

	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"testing-clientset": &ctxClientsetImpl{
				Interface: cs,
				namespace: "testing-namespace",
				cluster:   "testing-cluster",
			}},
		},
	}

	// No matching Jobs
	result, err := s.ListJobs(
		context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		&k8sv1.ListOptions{Labels: map[string]string{"unknown-label": "bar"}},
	)
	assert.NoError(t, err)
	assert.Empty(t, result)

	// Two matching Jobs
	result, err = s.ListJobs(
		context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		&k8sv1.ListOptions{Labels: map[string]string{"foo": "bar"}},
	)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// All Jobs in the namespace
	result, err = s.ListJobs(
		context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		&k8sv1.ListOptions{},
	)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestDeleteJob(t *testing.T) {
	t.Parallel()

	cs := testJobClientset()

	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"testing-clientset": &ctxClientsetImpl{
				Interface: cs,
				namespace: "testing-namespace",
				cluster:   "testing-cluster",
			}},
		},
	}

	// Not found
	err := s.DeleteJob(context.Background(), "testing-clientset", "testing-cluster", "testing-namespace", "abc")
	assert.Error(t, err)

	err = s.DeleteJob(context.Background(), "testing-clientset", "testing-cluster", "testing-namespace", "testing-job-name")
	assert.NoError(t, err)
}

var batchConfig = `
apiVersion: batch/v1
kind: Job
metadata:
  name: test-job
  labels:
    environment: staging
    facet-type: batch
`
var invalidBatchConfig = `
apiVersion: abc/v1
kind: foo
metadata:
  name: test-job
`

func TestCreateJob(t *testing.T) {
	cs := testJobClientset()

	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"testing-clientset": &ctxClientsetImpl{
				Interface: cs,
				namespace: "testing-namespace",
				cluster:   "testing-cluster",
			}},
		},
	}

	value := structpb.NewStringValue(batchConfig)
	job, err := s.CreateJob(context.Background(), "testing-clientset", "testing-cluster", "testing-namespace", value)
	assert.NoError(t, err)
	assert.NotNil(t, job)
	assert.Equal(t, "test-job", job.Name)

	invalidValue := structpb.NewStringValue(invalidBatchConfig)
	_, err = s.CreateJob(context.Background(), "testing-clientset", "testing-cluster", "testing-namespace1", invalidValue)
	assert.Error(t, err)
}
func TestProtoForJob(t *testing.T) {
	t.Parallel()

	var jobTestCases = []struct {
		id                  string
		inputClusterName    string
		expectedClusterName string
		job                 *v1.Job
	}{
		{
			id:                  "clustername already set",
			inputClusterName:    "notprod",
			expectedClusterName: "production",
			job: &v1.Job{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "production",
				},
			},
		},
		{
			id:                  "custername is not set",
			inputClusterName:    "staging",
			expectedClusterName: "staging",
			job: &v1.Job{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "",
				},
			},
		},
	}

	for _, tt := range jobTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			job := protoForJob(tt.inputClusterName, tt.job)
			assert.Equal(t, tt.expectedClusterName, job.Cluster)
		})
	}
}
