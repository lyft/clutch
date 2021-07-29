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

var testDynamodbTable = &types.TableDescription{
	TableName: "test-table",
	ProvisionedThroughput: &types.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(100),
		WriteCapacityUnits: aws.Int64(200),
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{},
}

var testGlobalSecondaryIndex = &types.GlobalSecondaryIndex{
	IndexName:  "test-gsi",
	KeySchema:  []types.KeySchema{},
	Projection: &types.Projection,
	ProvisionedThroughput: &types.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(10),
		WriteCapacityUnits: aws.Int64(20),
	},
}

var testDynamodbTableWithGSI = &types.TableDescription{
	TableName: "test-gsi-table",
	ProvisionedThroughput: &types.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(100),
		WriteCapacityUnits: aws.Int64(200),
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{testGlobalSecondaryIndex},
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
	assert.Equal(t, testDynamodbTable, result)

	m.tableErr = errors.New("error")
	_, err1 := c.DescribeTable(context.Background(), "us-east-1", "test-table")
	assert.EqualError(t, err1, "error")
}

func TestDescribeTableNotValid(t *testing.T) {
	m := &mockDynamodb{
		table: testDynamodbTable,
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", dynamodb: m}},
	}

	_, err := c.DescribeTable(context.Background(), "us-east-1", "nonexistent-table")
	assert.EqualError(t, err, "error")
}

func TestDescribeGsiTableValid(t *testing.T) {
	m := &mockDynamodb{
		table: testDynamodbTableWithGSI,
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", dynamodb: m}},
	}

	result, err := c.DescribeTable(context.Background(), "us-east-1", "test-gsi-table")
	assert.NoError(t, err)
	assert.Equal(t, testDynamodbTableWithGSI, result)

	m.tableErr = errors.New("error")
	_, err1 := c.DescribeTable(context.Background(), "us-east-1", "test-gsi-table")
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
