package featureflag

import (
	"context"
	"testing"

	featureflagv1 "github.com/lyft/clutch/backend/api/config/module/featureflag/v1"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	"github.com/lyft/clutch/backend/module/moduletest"
)

func TestModule(t *testing.T) {
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	m, err := New(nil, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("clutch.featureflag.v1.FeatureFlagAPI"))
	assert.True(t, r.JSONRegistered())
}

func TestAPI(t *testing.T) {
	m, err := newModuleImpl(
		&featureflagv1.Simple{
			Flags: map[string]bool{
				"foo": true,
				"bar": false,
			},
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, m)

	resp, err := m.GetFlags(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Flags["foo"].GetBooleanValue())
	assert.Equal(t, false, resp.Flags["bar"].GetBooleanValue())
	assert.Len(t, resp.Flags, 2)
}
