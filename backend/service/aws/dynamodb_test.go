package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
	awsv1 "github.com/lyft/clutch/backend/api/config/service/aws/v1"
)

var regions = []string{"us-east-1", "us-west-2"}

var cfg = &awsv1.Config{
	Regions: regions,
	ClientConfig: &awsv1.ClientConfig{
		Retries: 10,
	},
}

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
		log:     zaptest.NewLogger(t),
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
		log:     zaptest.NewLogger(t),
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
		log:     zaptest.NewLogger(t),
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", dynamodb: m}},
	}

	result, err := c.DescribeTable(context.Background(), "us-east-1", "test-gsi-table")
	assert.NoError(t, err)
	assert.Equal(t, testTableWithGSIOutput, result)
}

func TestGetScalingLimitsDefault(t *testing.T) {
	ds := getScalingLimits(cfg)

	assert.Equal(t, ds.MaxReadCapacityUnits, int64(AwsMaxRCU), "Max RCU default set")
	assert.Equal(t, ds.MaxScaleFactor, float32(SafeScaleFactor), "scale factor default set")
	assert.False(t, ds.EnableOverride)
}

func TestGetScalingLimitsCustom(t *testing.T) {
	cfg := &awsv1.Config{
		Regions: regions,
		ClientConfig: &awsv1.ClientConfig{
			Retries: 10,
		},
		DynamodbConfig: &awsv1.DynamodbConfig{
			ScalingLimits: &awsv1.ScalingLimits{
				MaxReadCapacityUnits:  1000,
				MaxWriteCapacityUnits: 2000,
				MaxScaleFactor:        4.0,
				EnableOverride:        false,
			},
		},
	}

	ds := getScalingLimits(cfg)

	assert.Equal(t, ds.MaxReadCapacityUnits, cfg.DynamodbConfig.ScalingLimits.MaxReadCapacityUnits, "RCU set")
	assert.Equal(t, ds.MaxWriteCapacityUnits, cfg.DynamodbConfig.ScalingLimits.MaxWriteCapacityUnits, "WCU set")
	assert.Equal(t, ds.MaxScaleFactor, cfg.DynamodbConfig.ScalingLimits.MaxScaleFactor, "scale factor default set")
	assert.False(t, ds.EnableOverride)
}
func TestUpdateTableCapacityWithDefaultLimits(t *testing.T) {

	tests := []struct {
		name     string
		inputRCU int64
		inputWCU int64
		want     string
	}{
		{"rcu above max", 100000, 250, "rpc error: code = FailedPrecondition desc = Target read capacity exceeds maximum allowed limits [40000]"},
		{"wcu above max", 100, 100000, "rpc error: code = FailedPrecondition desc = Target write capacity exceeds maximum allowed limits [40000]"},
		{"rcu lower than current", 1, 1500, "rpc error: code = FailedPrecondition desc = Target read capacity [1] is lower than current capacity [100]"},
		{"wcu lower than current", 1500, 1, "rpc error: code = FailedPrecondition desc = Target write capacity [1] is lower than current capacity [200]"},
		{"rcu change scale too high", 400, 200, "rpc error: code = FailedPrecondition desc = Target read capacity exceeds the scale limit of [2.0]x current capacity"},
		{"wcu change scale too high", 100, 600, "rpc error: code = FailedPrecondition desc = Target write capacity exceeds the scale limit of [2.0]x current capacity"},
	}

	m := &mockDynamodb{
		table: testDynamodbTable,
	}

	ds := getScalingLimits(cfg)

	d := &awsv1.DynamodbConfig{
		ScalingLimits: &awsv1.ScalingLimits{
			MaxReadCapacityUnits:  ds.MaxReadCapacityUnits,
			MaxWriteCapacityUnits: ds.MaxWriteCapacityUnits,
			MaxScaleFactor:        ds.MaxScaleFactor,
			EnableOverride:        ds.EnableOverride,
		},
	}

	c := &client{
		log:     zaptest.NewLogger(t),
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", dynamodbCfg: d, dynamodb: m}},
	}

	err := c.UpdateTableCapacity(context.Background(), "us-east-1", "test-table", 101, 201)
	assert.NoError(t, err)

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := c.UpdateTableCapacity(context.Background(), "us-east-1", "test-table", tt.inputRCU, tt.inputWCU)
			if got.Error() != tt.want {
				t.Errorf("\nWant error msg: %s\nGot error msg: %s", tt.want, got)
			}
		})
	}
}

func TestUpdateTableCapacityWithCustomLimits(t *testing.T) {
	tests := []struct {
		name     string
		inputRCU int64
		inputWCU int64
		want     string
	}{
		{"rcu above max", 2000, 250, "rpc error: code = FailedPrecondition desc = Target read capacity exceeds maximum allowed limits [1000]"},
		{"wcu above max", 100, 3000, "rpc error: code = FailedPrecondition desc = Target write capacity exceeds maximum allowed limits [2000]"},
		{"rcu lower than current", 1, 1500, "rpc error: code = FailedPrecondition desc = Target read capacity [1] is lower than current capacity [100]"},
		{"wcu lower than current", 1500, 1, "rpc error: code = FailedPrecondition desc = Target write capacity [1] is lower than current capacity [200]"},
		{"rcu change scale too high", 900, 200, "rpc error: code = FailedPrecondition desc = Target read capacity exceeds the scale limit of [4.0]x current capacity"},
		{"wcu change scale too high", 100, 1900, "rpc error: code = FailedPrecondition desc = Target write capacity exceeds the scale limit of [4.0]x current capacity"},
	}

	m := &mockDynamodb{
		table: testDynamodbTable,
	}
	d := &awsv1.DynamodbConfig{
		ScalingLimits: &awsv1.ScalingLimits{
			MaxReadCapacityUnits:  1000,
			MaxWriteCapacityUnits: 2000,
			MaxScaleFactor:        4.0,
			EnableOverride:        false,
		},
	}
	c := &client{
		log:     zaptest.NewLogger(t),
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", dynamodbCfg: d, dynamodb: m}},
	}

	err := c.UpdateTableCapacity(context.Background(), "us-east-1", "test-table", 101, 201)
	assert.NoError(t, err)

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := c.UpdateTableCapacity(context.Background(), "us-east-1", "test-table", tt.inputRCU, tt.inputWCU)
			if got.Error() != tt.want {
				t.Errorf("\nWant error msg: %s\nGot error msg: %s", tt.want, got)
			}
		})
	}
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
