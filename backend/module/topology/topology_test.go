package topology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	"github.com/lyft/clutch/backend/mock/service/topologymock"
	"github.com/lyft/clutch/backend/module/moduletest"
	"github.com/lyft/clutch/backend/service"
)

func TestModule(t *testing.T) {
	service.Registry["clutch.service.topology"] = topologymock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	m, err := New(nil, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.True(t, r.JSONRegistered())
}
