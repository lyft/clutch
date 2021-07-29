package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

var testDynamodbTable = &types.DescribeTable{
	TableName: "test-table",
	ProvisionedThroughput: &types.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(100),
		WriteCapacityUnits: aws.Int64(200),
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{},
}

var testDynamodbTableOutput = &dynamodb.DescribeTableOutput{
	Table: &types.TableDescription{
		TableName: "test",
	},
}

func TestDescribeTableValid(t *testing.T) {
	m := &mockDynamodb{
		table: testDynamodbTable,
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", dynamodb: m}},
	}

	result, err := c.DescribeTable(context.Background(), "us-east-1", "test-table")
	assert.NoError(t, err)
	assert.Equal(t, testDynamodbTableOutput, result)

	m.tableErr = errors.New("error")
	_, err1 := c.DescribeTable(context.Background(), "us-east-1", "test-table")
	assert.EqualError(t, err1, "error")
}

type mockDynamodb struct {
	dynamodbClient

	tableErr error
	table    *types.TableDescription

	updateErr error
	update    *dynamodb.UpdateTableOutput
}

func (m *mockDynamodb) DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	if m.tableErr != nil {
		return nil, m.tableErr
	}

	ret := &dynamodb.DescribeTableOutput{
		Table: m.table,
	}

	return ret, nil
}

func (m *mockDynamodb) UpdateTable(ctx context.Context, params *dynamodb.UpdateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateTableOutput, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	ret := m.update

	return ret, nil
}
