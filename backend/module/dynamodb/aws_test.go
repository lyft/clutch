package dynamodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
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
	assert.NoError(t, r.HasAPI("clutch.aws.dynamodb.v1.DDBAPI"))
	assert.True(t, r.JSONRegistered())
}

func TestDDBAPIDescribeTable(t *testing.T) {
	c := awsmock.New()
	api := newDDBAPI(c)
	resp, err := api.DescribeTable(context.Background(), &dynamodbv1.DescribeTableRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
