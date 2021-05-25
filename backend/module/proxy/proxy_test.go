package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"

	proxyv1cfg "github.com/lyft/clutch/backend/api/config/module/proxy/v1"
	"github.com/lyft/clutch/backend/module/moduletest"
)

func TestModule(t *testing.T) {
	config, _ := anypb.New(&proxyv1cfg.Config{})
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	m, err := New(config, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.True(t, r.JSONRegistered())
}
