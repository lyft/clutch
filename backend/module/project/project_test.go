package project

import (
	"testing"

	"github.com/lyft/clutch/backend/mock/service/projectmock"
	"github.com/lyft/clutch/backend/module/moduletest"
	"github.com/lyft/clutch/backend/service"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestModule(t *testing.T) {
	service.Registry["clutch.service.project"] = projectmock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	m, err := New(&anypb.Any{}, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("clutch.project.v1.ProjectAPI"))
	assert.True(t, r.JSONRegistered())
}
