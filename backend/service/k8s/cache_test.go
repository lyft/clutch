package k8s

import (
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
)

func TestProcessInformerEvent(t *testing.T) {
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	topologyObjectChan := make(chan *topologyv1.UpdateCacheRequest, 1)
	svc := svc{
		topologyObjectChan: topologyObjectChan,
		log:                log,
		scope:              scope,
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			ClusterName: "cluster",
			Name:        "test-pod-1",
			Namespace:   "testing-namespace",
		},
	}

	expectedClutchPod := podDescription(pod, "")
	protoPod, err := ptypes.MarshalAny(expectedClutchPod)
	assert.NoError(t, err)

	expectedUpdateCacheRequest := &topologyv1.UpdateCacheRequest{
		Resource: &topologyv1.Resource{
			Id: "cluster/testing-namespace/test-pod-1",
			Pb: protoPod,
		},
		Action: topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
	}

	svc.processInformerEvent(pod, topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE)
	topologyUpdateRequest := <-svc.topologyObjectChan

	assert.Equal(t, expectedUpdateCacheRequest, topologyUpdateRequest)
}
