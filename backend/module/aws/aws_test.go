package aws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	"github.com/lyft/clutch/backend/mock/service/awsmock"
	"github.com/lyft/clutch/backend/module/moduletest"
	"github.com/lyft/clutch/backend/service"
)

func TestModule(t *testing.T) {
	service.Registry["clutch.service.aws"] = awsmock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	m, err := New(nil, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("clutch.aws.ec2.v1.EC2API"))
	assert.True(t, r.JSONRegistered())
}

func TestEC2APIGetInstance(t *testing.T) {
	c := awsmock.New()
	api := newEC2API(c)
	resp, err := api.GetInstance(context.Background(), &ec2v1.GetInstanceRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
