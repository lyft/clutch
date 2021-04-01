package topology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"

	topologyv1 "github.com/lyft/clutch/backend/api/config/service/topology/v1"
)

func TestNew(t *testing.T) {
	cfg, _ := anypb.New(&topologyv1.Config{})

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	// Should throw error looking for postgres
	_, err := New(cfg, log, scope)
	assert.Error(t, err)
}

func TestNewWithWrongConfig(t *testing.T) {
	_, err := New(&anypb.Any{TypeUrl: "foobar"}, nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mismatched message")
}
