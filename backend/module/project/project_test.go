package project

import (
	"testing"

	projectv1cfg "github.com/lyft/clutch/backend/api/config/module/project/v1"
	"github.com/lyft/clutch/backend/mock/service/projectmock"
	"github.com/lyft/clutch/backend/module/moduletest"
	"github.com/lyft/clutch/backend/service"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestModule(t *testing.T) {
	service.Registry["clutch.service.project"] = projectmock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	config, _ := anypb.New(&projectv1cfg.Config{})

	m, err := New(config, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("clutch.project.v1.ProjectAPI"))
	assert.True(t, r.JSONRegistered())
}

func TestModuleServiceOverride(t *testing.T) {
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	config, _ := anypb.New(&projectv1cfg.Config{
		ProjectServiceOverride: "custom.service.project",
	})

	// Register custom service
	service.Registry["custom.service.project"] = projectmock.New()

	// Test override was valid
	_, err := New(config, log, scope)
	assert.NoError(t, err)

	// Register custom service with wrong type, should error.
	service.Registry["custom.service.project"] = "meow"
	_, err = New(config, log, scope)
	assert.Error(t, err)
}
