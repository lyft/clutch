package audit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func TestEventFromResponse(t *testing.T) {
	m := &mid{}
	err := errors.New("error")
	a := (*anypb.Any)(nil)

	// case: err, passed to eventFromResponse, does not equal nil
	event, err := m.eventFromResponse(nil, err)
	assert.NoError(t, err)
	assert.NotEmpty(t, event)
	assert.Equal(t, "error", event.Status.Message)
	assert.Equal(t, 0, len(event.Resources))
	assert.Equal(t, a, event.ResponseMetadata.Body)

	resp := &k8sapiv1.Pod{Cluster: "kind-clutch", Namespace: "envoy-staging", Name: "envoy-main-579848cc64-cxnqm"}
	// case: err, passed to eventFromResponse, equals nil
	event, err = m.eventFromResponse(resp, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, event)
	assert.Equal(t, int32(0), event.Status.Code)
	assert.Equal(t, "", event.Status.Message)
	assert.Equal(t, 1, len(event.Resources))
	assert.Equal(t, "type.googleapis.com/clutch.k8s.v1.Pod", event.ResponseMetadata.Body.TypeUrl)
}
