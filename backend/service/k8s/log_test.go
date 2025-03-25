package k8s

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func TestGetPodLogs(t *testing.T) {
	t.Parallel()
	testPods := []runtime.Object{
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testing-pod-name",
				Namespace: "testing-namespace",
				Labels: map[string]string{
					"foo":                  "bar",
					clutchLabelClusterName: "core-testing",
				},
				Annotations: map[string]string{"baz": "quuz"},
			},
			Status: corev1.PodStatus{
				StartTime:         &metav1.Time{},
				ContainerStatuses: []corev1.ContainerStatus{{Name: "container1"}},
			},
		},
	}

	cs := fake.NewSimpleClientset(testPods...)

	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{
				"foo": {
					Interface: cs,
					namespace: "testing-namespace",
					cluster:   "core-testing",
				},
			},
		},
	}

	resp, err := s.GetPodLogs(context.Background(), "foo", "core-testing", "testing-namespace", "testing-pod-name", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "fake logs", resp.Logs[0].S)
}

func TestBufferToResponse(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		input    string
		expected *k8sv1.GetPodLogsResponse
	}{
		{input: "", expected: &k8sv1.GetPodLogsResponse{}},
		{input: "Hello world!", expected: &k8sv1.GetPodLogsResponse{Logs: []*k8sv1.PodLogLine{{S: "Hello world!"}}}},
		{
			input: "2022-11-07T19:30:38.974187286Z Hello world!\nHello!\n",
			expected: &k8sv1.GetPodLogsResponse{
				LatestTs: "2022-11-07T19:30:38.974187286Z",
				Logs: []*k8sv1.PodLogLine{
					{
						S:  "Hello world!",
						Ts: "2022-11-07T19:30:38.974187286Z",
					},
					{
						S: "Hello!",
					},
				},
			},
		},
	}

	for idx, tc := range tcs {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			rdr := strings.NewReader(tc.input)
			resp, err := bufferToResponse(rdr)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.expected.LatestTs, resp.LatestTs)
			assert.Equal(t, tc.expected.Logs, resp.Logs)
		})
	}
}

func TestOptsConversion(t *testing.T) {
	t.Parallel()
	{
		o, err := protoOptsToK8sOpts(nil)
		assert.NoError(t, err)
		assert.NotNil(t, o)
		assert.True(t, o.Timestamps)
	}
	{
		o, err := protoOptsToK8sOpts(&k8sv1.PodLogsOptions{
			ContainerName: "foo",
			Previous:      true,
			SinceTs:       "2022-11-07T19:30:38.974187286Z",
			TailNumLines:  25,
		})
		assert.NoError(t, err)
		assert.Equal(t, "foo", o.Container)
		assert.True(t, o.Previous)
		assert.EqualValues(t, 25, *o.TailLines)
		assert.Equal(t, "2022-11-07T19:30:38.974187286Z", o.SinceTime.Format(rfc3339NanoFixed))
	}
	{
		o, err := protoOptsToK8sOpts(&k8sv1.PodLogsOptions{})

		assert.NoError(t, err)
		assert.Equal(t, "", o.Container)
		assert.False(t, o.Previous)
		assert.Nil(t, o.TailLines)
		assert.Nil(t, o.SinceTime)
	}
}

func TestBytesToLogLine(t *testing.T) {
	t.Parallel()
	{
		ll := bytesToLogLine([]byte("2022-11-07T19:30:38.974187286Z Hello world this is it!\n"))
		assert.Equal(t, "2022-11-07T19:30:38.974187286Z", ll.Ts)
		assert.Equal(t, "Hello world this is it!", ll.S)
	}
	{
		ll := bytesToLogLine([]byte("2022-11-07T19:30:38.974187287Z Hello world this is it!"))
		assert.Equal(t, "2022-11-07T19:30:38.974187287Z", ll.Ts)
		assert.Equal(t, "Hello world this is it!", ll.S)
	}
	{
		ll := bytesToLogLine([]byte("2022 Hello world"))
		assert.Equal(t, "", ll.Ts)
		assert.Equal(t, "2022 Hello world", ll.S)
	}
}
