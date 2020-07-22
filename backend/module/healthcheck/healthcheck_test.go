package healthcheck

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"

	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
	"github.com/lyft/clutch/backend/module/moduletest"

	"testing"

	"go.uber.org/zap/zaptest"
)

func TestModule(t *testing.T) {
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	m, err := New(nil, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("clutch.healthcheck.v1.HealthcheckAPI"))
	assert.True(t, r.JSONRegistered())
}

func TestAPI(t *testing.T) {
	api := newAPI()
	resp, err := api.Healthcheck(context.Background(), &healthcheckv1.HealthcheckRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
