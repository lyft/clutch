package kinesis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
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
	assert.NoError(t, r.HasAPI("clutch.aws.kinesis.v1.KinesisAPI"))
	assert.True(t, r.JSONRegistered())
}

func TestKinesisAPIGetStream(t *testing.T) {
	c := awsmock.New()
	api := newKinesisAPI(c)
	resp, err := api.GetStream(context.Background(), &kinesisv1.GetStreamRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
