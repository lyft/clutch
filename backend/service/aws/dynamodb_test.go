package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
)

var testDynamodbTable = &types.TableDescription{
	TableName: aws.String("test-table"),
	ProvisionedThroughput: &types.ProvisionedThroughputDescription{
		ReadCapacityUnits:  aws.Int64(100),
		WriteCapacityUnits: aws.Int64(200),
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndexDescription{},
}

var testTableOutput = &dynamodbv1.Table{
	Name:   "test-table",
	Region: "us-east-1",
	ProvisionedThroughput: &dynamodbv1.ProvisionedThroughput{
		ReadCapacityUnits:  100,
		WriteCapacityUnits: 200,
	},
	GlobalSecondaryIndexes: []*dynamodbv1.GlobalSecondaryIndex{},
}

var testGlobalSecondaryIndex = &types.GlobalSecondaryIndexDescription{
	IndexName: aws.String("test-gsi"),
	KeySchema: []types.KeySchemaElement{},
	ProvisionedThroughput: &types.ProvisionedThroughputDescription{
		ReadCapacityUnits:  aws.Int64(10),
		WriteCapacityUnits: aws.Int64(20),
	},
}

var testDynamodbTableWithGSI = &types.TableDescription{
	TableName: aws.String("test-gsi-table"),
	ProvisionedThroughput: &types.ProvisionedThroughputDescription{
		ReadCapacityUnits:  aws.Int64(100),
		WriteCapacityUnits: aws.Int64(200),
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndexDescription{
		{IndexName: aws.String("test-gsi"),
			KeySchema: []types.KeySchemaElement{},
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(20),
			},
		},
	},
}

var testTableWithGSIOutput = &dynamodbv1.Table{
	Name:   "test-gsi-table",
	Region: "us-east-1",
	ProvisionedThroughput: &dynamodbv1.ProvisionedThroughput{
		ReadCapacityUnits:  100,
		WriteCapacityUnits: 200,
	},
	GlobalSecondaryIndexes: []*dynamodbv1.GlobalSecondaryIndex{
		{
			Name: "test-gsi",
			ProvisionedThroughput: &dynamodbv1.ProvisionedThroughput{
				ReadCapacityUnits:  10,
				WriteCapacityUnits: 20,
			},
		},
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
	assert.Equal(t, testTableOutput, result)

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

	m.tableErr = errors.New("resource not found")
	_, err := c.DescribeTable(context.Background(), "us-east-1", "nonexistent-table")
	assert.EqualError(t, err, "resource not found")
}

func TestDescribeTableWithGsiValid(t *testing.T) {
	m := &mockDynamodb{
		table: testDynamodbTableWithGSI,
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", dynamodb: m}},
	}

	result, err := c.DescribeTable(context.Background(), "us-east-1", "test-gsi-table")
	assert.NoError(t, err)
	assert.Equal(t, testTableWithGSIOutput, result)
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
