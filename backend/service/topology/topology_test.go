package topology

import (
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	topologyv1 "github.com/lyft/clutch/backend/api/config/service/topology/v1"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
)

func TestNew(t *testing.T) {
	cfg, _ := ptypes.MarshalAny(&topologyv1.Config{})

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	// Should throw error looking for postgres
	_, err := New(cfg, log, scope)
	assert.Error(t, err)
}

func TestNewWithWrongConfig(t *testing.T) {
	_, err := New(&any.Any{TypeUrl: "foobar"}, nil, nil)
	assert.EqualError(t, err, `mismatched message type: got "foobar" want "clutch.config.service.topology.v1.Config"`)
}
