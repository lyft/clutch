package k8s

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	v1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func testCronService() *svc {
	cron := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "testing-cron-name",
			Namespace:   "testing-namespace",
			Labels:      map[string]string{"test": "foo"},
			Annotations: map[string]string{"test": "bar"},
		},
	}

	cs := fake.NewSimpleClientset(cron)
	return &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": &ctxClientsetImpl{
				Interface: cs,
				namespace: "default",
				cluster:   "core-testing",
			}},
		},
	}
}

func TestDescribeCron(t *testing.T) {
	s := testCronService()
	cron, err := s.DescribeCronJob(context.Background(), "foo", "core-testing", "testing-namespace", "testing-cron-name")
	assert.NoError(t, err)
	assert.NotNil(t, cron)
}

func TestListCron(t *testing.T) {
	s := testCronService()
	opts := &k8sapiv1.ListOptions{Labels: map[string]string{"test": "foo"}}
	list, err := s.ListCronJobs(context.Background(), "foo", "core-testing", "testing-namespace", opts)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	// Not Found
	opts = &k8sapiv1.ListOptions{Labels: map[string]string{"unknown": "bar"}}
	list, err = s.ListCronJobs(context.Background(), "foo", "core-testing", "testing-namespace", opts)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(list))
}

func TestDeleteCron(t *testing.T) {
	s := testCronService()
	// Not found.
	err := s.DeleteCronJob(context.Background(), "foo", "core-testing", "testing-namespace", "abc")
	assert.Error(t, err)

	err = s.DeleteCronJob(context.Background(), "foo", "core-testing", "testing-namespace", "testing-cron-name")
	assert.NoError(t, err)

	// Not found.
	_, err = s.DescribeCronJob(context.Background(), "foo", "core-testing", "testing-namespace", "testing-cron-name")
	assert.Error(t, err)
}

func TestProtoForCron(t *testing.T) {
	t.Parallel()

	var cronTestCases = []struct {
		id                  string
		inputClusterName    string
		expectedClusterName string
		expectedName        string
		cron                *v1beta1.CronJob
	}{
		{
			id:                  "clustername already set",
			inputClusterName:    "abc",
			expectedClusterName: "production",
			expectedName:        "test1",
			cron: &v1beta1.CronJob{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "production",
					Name:        "test1",
				},
			},
		},
		{
			id:                  "clustername is not set",
			inputClusterName:    "staging",
			expectedClusterName: "staging",
			expectedName:        "test2",
			cron: &v1beta1.CronJob{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "",
					Name:        "test2",
				},
				Spec: v1beta1.CronJobSpec{
					ConcurrencyPolicy:       v1beta1.AllowConcurrent,
					Schedule:                "5 4 * * *",
					Suspend:                 &[]bool{true}[0],
					StartingDeadlineSeconds: &[]int64{69}[0],
				},
				Status: v1beta1.CronJobStatus{
					Active: []v1.ObjectReference{{}, {}},
				},
			},
		},
	}

	for _, tt := range cronTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			cron := ProtoForCronJob(tt.inputClusterName, tt.cron)
			assert.Equal(t, tt.expectedClusterName, cron.Cluster)
			assert.Equal(t, tt.expectedName, cron.Name)

			assert.Equal(t, tt.cron.Spec.Schedule, cron.Schedule)

			if tt.cron.Spec.ConcurrencyPolicy != "" {
				assert.Equal(t, strings.ToUpper(string(tt.cron.Spec.ConcurrencyPolicy)), cron.ConcurrencyPolicy.String())
			}
			if tt.cron.Spec.Suspend != nil {
				assert.Equal(t, *tt.cron.Spec.Suspend, cron.Suspend)
			}
			if tt.cron.Spec.StartingDeadlineSeconds != nil {
				assert.Equal(t, *tt.cron.Spec.StartingDeadlineSeconds, cron.StartingDeadlineSeconds)
			}
			if tt.cron.Status.Active != nil {
				assert.Equal(t, int32(len(tt.cron.Status.Active)), cron.ActiveJobs)
			}
		})
	}
}
