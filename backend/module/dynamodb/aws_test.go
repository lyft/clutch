package dynamodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally/v4"
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

var testUpdateCapacityResponse = &dynamodbv1.UpdateCapacityResponse{
	Table: testTable,
}

var pt = &dynamodbv1.Throughput{
	ReadCapacityUnits:  100,
	WriteCapacityUnits: 200,
}
var g = []*dynamodbv1.GlobalSecondaryIndex{
	{
		Name: "test-gsi",
		ProvisionedThroughput: &dynamodbv1.Throughput{
			ReadCapacityUnits:  10,
			WriteCapacityUnits: 20,
		},
		Status: dynamodbv1.GlobalSecondaryIndex_Status(3),
	},
}

var testTable = &dynamodbv1.Table{
	Name:                   "",
	Region:                 "",
	Account:                "default",
	ProvisionedThroughput:  pt,
	GlobalSecondaryIndexes: g,
	Status:                 dynamodbv1.Table_Status(3),
}

func TestDDBAPIDescribeTable(t *testing.T) {
	c := awsmock.New()
	api := newDDBAPI(c)
	resp, err := api.DescribeTable(context.Background(), &dynamodbv1.DescribeTableRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDDBAPIUpdateCapacity(t *testing.T) {
	c := awsmock.New()
	api := newDDBAPI(c)
	resp, err := api.UpdateCapacity(context.Background(), &dynamodbv1.UpdateCapacityRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, testUpdateCapacityResponse, resp)
}
