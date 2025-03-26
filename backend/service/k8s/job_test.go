package k8s

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

func testJobService(n int) *svc {
	jobs := make([]runtime.Object, n)
	for i := 0; i < n; i++ {
		job := &v1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("testing-job-name-%b", i),
				Namespace: "testing-namespace",
				Labels:    map[string]string{"foo": "bar"},
			},
		}
		jobs[i] = job
	}

	return &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"testing-clientset": {
				Interface: fake.NewSimpleClientset(jobs...),
				namespace: "testing-namespace",
				cluster:   "testing-cluster",
			}},
		},
	}
}

func TestDescribeJob(t *testing.T) {
	// 1 existing job
	s := testJobService(1)
	job, err := s.DescribeJob(context.Background(), "testing-clientset", "testing-cluster", "testing-namespace", "testing-job-name")
	assert.NoError(t, err)
	assert.NotNil(t, job)

	// No existing job
	s = testJobService(0)
	job, err = s.DescribeJob(context.Background(), "testing-clientset", "testing-cluster", "testing-namespace", "testing-job-name")
	assert.Error(t, err)
	assert.Nil(t, job)

	// 2 existing jobs
	s = testJobService(2)
	job, err = s.DescribeJob(context.Background(), "testing-clientset", "testing-cluster", "testing-namespace", "testing-job-name")
	assert.Error(t, err)
	assert.Nil(t, job)
}

func TestListJobs(t *testing.T) {
	t.Parallel()

	cs := testListJobClientset()

	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"testing-clientset": {
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

	s := testJobService(1)

	// Not found
	err := s.DeleteJob(context.Background(), "testing-clientset", "testing-cluster", "testing-namespace", "abc")
	assert.Error(t, err)

	err = s.DeleteJob(context.Background(), "testing-clientset", "testing-cluster", "testing-namespace", "testing-job-name-0")
	assert.NoError(t, err)
}

func TestCreateJob(t *testing.T) {
	s := testJobService(1)

	jobConfig := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-job",
			Namespace: "testing-namespace",
			Labels:    map[string]string{"foo": "bar"},
		},
	}

	job, err := s.CreateJob(context.Background(), "testing-clientset", "testing-cluster", "testing-namespace", jobConfig)
	assert.NoError(t, err)
	assert.NotNil(t, job)
	assert.Equal(t, "test-job", job.Name)

	// invalid job config
	_, err = s.CreateJob(context.Background(), "testing-clientset", "testing-cluster", "testing-namespace1", nil)
	assert.Error(t, err)
}

func TestProtoForJob(t *testing.T) {
	t.Parallel()

	jobTestCases := []struct {
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
			job: &v1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clutchLabelClusterName: "staging",
					},
				},
			},
		},
	}

	for _, tt := range jobTestCases {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			job := protoForJob(tt.inputClusterName, tt.job)
			assert.Equal(t, tt.expectedClusterName, job.Cluster)
		})
	}
}
