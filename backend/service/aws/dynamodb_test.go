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
	TableStatus:            "ACTIVE",
	BillingModeSummary: &types.BillingModeSummary{
		BillingMode: "PROVISIONED",
	},
}

var testTableOutput = &dynamodbv1.Table{
	Name:   "test-table",
	Region: "us-east-1",
	ProvisionedThroughput: &dynamodbv1.Throughput{
		ReadCapacityUnits:  100,
		WriteCapacityUnits: 200,
	},
	GlobalSecondaryIndexes: []*dynamodbv1.GlobalSecondaryIndex{},
	Status:                 dynamodbv1.Table_Status(5),
	BillingMode:            dynamodbv1.Table_BillingMode(2),
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
			IndexStatus: "ACTIVE",
		},
	},
	TableStatus: "ACTIVE",
	BillingModeSummary: &types.BillingModeSummary{
		BillingMode: "PROVISIONED",
	},
}

var testTableWithGSIOutput = &dynamodbv1.Table{
	Name:   "test-gsi-table",
	Region: "us-east-1",
	ProvisionedThroughput: &dynamodbv1.Throughput{
		ReadCapacityUnits:  100,
		WriteCapacityUnits: 200,
	},
	GlobalSecondaryIndexes: []*dynamodbv1.GlobalSecondaryIndex{
		{
			Name: "test-gsi",
			ProvisionedThroughput: &dynamodbv1.Throughput{
				ReadCapacityUnits:  10,
				WriteCapacityUnits: 20,
			},
			Status: dynamodbv1.GlobalSecondaryIndex_Status(5),
		},
	},
	Status:      dynamodbv1.Table_Status(5),
	BillingMode: dynamodbv1.Table_BillingMode(2),
}

var testOnDemandTable = &types.TableDescription{
	TableName: aws.String("test-ondemand-table"),
	ProvisionedThroughput: &types.ProvisionedThroughputDescription{
		ReadCapacityUnits:  aws.Int64(0),
		WriteCapacityUnits: aws.Int64(0),
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndexDescription{
		{IndexName: aws.String("test-gsi"),
			KeySchema: []types.KeySchemaElement{},
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				ReadCapacityUnits:  aws.Int64(0),
				WriteCapacityUnits: aws.Int64(0),
			},
			IndexStatus: "ACTIVE",
		},
	},
	TableStatus: "ACTIVE",
	BillingModeSummary: &types.BillingModeSummary{
		BillingMode: "PAY_PER_REQUEST",
	},
}

var testOnDemandInferredTable = &types.TableDescription{
	TableName: aws.String("test-ondemand-table"),
	ProvisionedThroughput: &types.ProvisionedThroughputDescription{
		ReadCapacityUnits:  aws.Int64(0),
		WriteCapacityUnits: aws.Int64(0),
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndexDescription{
		{IndexName: aws.String("test-gsi"),
			KeySchema: []types.KeySchemaElement{},
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				ReadCapacityUnits:  aws.Int64(0),
				WriteCapacityUnits: aws.Int64(0),
			},
			IndexStatus: "ACTIVE",
		},
	},
	TableStatus: "ACTIVE",
}

var testUpdateResponse = &types.TableDescription{
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
			IndexStatus: "UPDATING",
		},
	},
	TableStatus: "UPDATING",
	BillingModeSummary: &types.BillingModeSummary{
		BillingMode: "PROVISIONED",
	},
}

var testTableNoBillingMode = &types.TableDescription{
	TableName: aws.String("test-table-no-billing"),
	ProvisionedThroughput: &types.ProvisionedThroughputDescription{
		ReadCapacityUnits:  aws.Int64(100),
		WriteCapacityUnits: aws.Int64(200),
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndexDescription{},
	TableStatus:            "ACTIVE",
}

var testTableNoBillingModeResponse = &types.TableDescription{
	TableName: aws.String("test-table-no-billing"),
	ProvisionedThroughput: &types.ProvisionedThroughputDescription{
		ReadCapacityUnits:  aws.Int64(100),
		WriteCapacityUnits: aws.Int64(200),
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndexDescription{},
	TableStatus:            "UPDATING",
}

func TestDescribeTableValid(t *testing.T) {
	m := &mockDynamodb{
		table: testDynamodbTable,
	}
	c := &client{
		log:                 zaptest.NewLogger(t),
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodb: m},
				},
			},
		},
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
		log:                 zaptest.NewLogger(t),
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodb: m},
				},
			},
		},
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
		log:                 zaptest.NewLogger(t),
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodb: m},
				},
			},
		},
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

func TestIncreaseTableCapacityErrors(t *testing.T) {
	tests := []struct {
		name      string
		targetRCU int64
		targetWCU int64
		want      string
	}{
		{"target table rcu above max", 100000, 250, "rpc error: code = FailedPrecondition desc = Target read capacity exceeds maximum allowed limits [40000]"},
		{"target table wcu above max", 100, 100000, "rpc error: code = FailedPrecondition desc = Target write capacity exceeds maximum allowed limits [40000]"},
		{"target table rcu lower than current", 1, 1500, "rpc error: code = FailedPrecondition desc = Target read capacity [1] is lower than current capacity [100]"},
		{"target table wcu lower than current", 1500, 1, "rpc error: code = FailedPrecondition desc = Target write capacity [1] is lower than current capacity [200]"},
		{"table rcu change scale too high", 400, 200, "rpc error: code = FailedPrecondition desc = Target read capacity exceeds the scale limit of [2.0]x current capacity"},
		{"table wcu change scale too high", 100, 600, "rpc error: code = FailedPrecondition desc = Target write capacity exceeds the scale limit of [2.0]x current capacity"},
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
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodbCfg: d, dynamodb: m},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			targetTableCapacity := &dynamodbv1.Throughput{
				ReadCapacityUnits:  tt.targetRCU,
				WriteCapacityUnits: tt.targetWCU,
			}

			gsiUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)

			result, err := c.UpdateCapacity(context.Background(), "us-east-1", "test-table", targetTableCapacity, gsiUpdates, false)
			if err.Error() != tt.want {
				t.Errorf("\nWant error msg: %s\nGot error msg: %s", tt.want, err)
			}
			assert.Nil(t, result)
		})
	}
}

func TestUpdateGSICapacityErrors(t *testing.T) {
	tests := []struct {
		name      string
		targetRCU int64
		targetWCU int64
		want      string
	}{
		{"target gsi rcu above max", 100000, 25, "rpc error: code = FailedPrecondition desc = Target read capacity exceeds maximum allowed limits [40000]"},
		{"target gsi wcu above max", 25, 100000, "rpc error: code = FailedPrecondition desc = Target write capacity exceeds maximum allowed limits [40000]"},
		{"target gsi rcu lower than current", 1, 25, "rpc error: code = FailedPrecondition desc = Target read capacity [1] is lower than current capacity [10]"},
		{"target gsi wcu lower than current", 25, 1, "rpc error: code = FailedPrecondition desc = Target write capacity [1] is lower than current capacity [20]"},
		{"gsi rcu change scale too high", 400, 25, "rpc error: code = FailedPrecondition desc = Target read capacity exceeds the scale limit of [2.0]x current capacity"},
		{"gsi wcu change scale too high", 25, 600, "rpc error: code = FailedPrecondition desc = Target write capacity exceeds the scale limit of [2.0]x current capacity"},
	}

	m := &mockDynamodb{
		table: testDynamodbTableWithGSI,
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
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodbCfg: d, dynamodb: m},
				},
			},
		},
	}

	targetTableCapacity := &dynamodbv1.Throughput{
		ReadCapacityUnits:  101,
		WriteCapacityUnits: 202,
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			targetIndexCapacity := &dynamodbv1.Throughput{
				ReadCapacityUnits:  tt.targetRCU,
				WriteCapacityUnits: tt.targetWCU,
			}

			gsiUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)
			update := dynamodbv1.IndexUpdateAction{
				Name:            "test-gsi",
				IndexThroughput: targetIndexCapacity,
			}
			gsiUpdates = append(gsiUpdates, &update)

			result, err := c.UpdateCapacity(context.Background(), "us-east-1", "test-table", targetTableCapacity, gsiUpdates, false)
			if err.Error() != tt.want {
				t.Errorf("\nWant error msg: %s\nGot error msg: %s", tt.want, err)
			}
			assert.Nil(t, result)
		})
	}
}

func TestUpdateCapacitySuccess(t *testing.T) {
	m := &mockDynamodb{
		table:  testDynamodbTableWithGSI,
		update: testUpdateResponse,
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
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodbCfg: d, dynamodb: m},
				},
			},
		},
	}

	targetTableCapacity := &dynamodbv1.Throughput{
		ReadCapacityUnits:  101,
		WriteCapacityUnits: 202,
	}

	targetIndexCapacity := &dynamodbv1.Throughput{
		ReadCapacityUnits:  11,
		WriteCapacityUnits: 21,
	}

	gsiUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)
	update := dynamodbv1.IndexUpdateAction{
		Name:            "test-gsi",
		IndexThroughput: targetIndexCapacity,
	}
	gsiUpdates = append(gsiUpdates, &update)

	got, err := c.UpdateCapacity(context.Background(), "us-east-1", "test-gsi-table", targetTableCapacity, gsiUpdates, false)
	assert.NotNil(t, got)
	assert.Nil(t, err)
	assert.Equal(t, got.Status, dynamodbv1.Table_Status(3))
	assert.Equal(t, got.GlobalSecondaryIndexes[0].Status, dynamodbv1.GlobalSecondaryIndex_Status(3))
}

func TestIgnoreMaximums(t *testing.T) {
	m := &mockDynamodb{
		table:  testDynamodbTableWithGSI,
		update: testUpdateResponse,
	}

	ds := getScalingLimits(cfg)

	d := &awsv1.DynamodbConfig{
		ScalingLimits: &awsv1.ScalingLimits{
			MaxReadCapacityUnits:  ds.MaxReadCapacityUnits,
			MaxWriteCapacityUnits: ds.MaxWriteCapacityUnits,
			MaxScaleFactor:        ds.MaxScaleFactor,
			EnableOverride:        true,
		},
	}

	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodbCfg: d, dynamodb: m},
				},
			},
		},
	}

	targetTableCapacity := &dynamodbv1.Throughput{
		ReadCapacityUnits:  100000,
		WriteCapacityUnits: 100000,
	}

	targetIndexCapacity := &dynamodbv1.Throughput{
		ReadCapacityUnits:  100000,
		WriteCapacityUnits: 100000,
	}

	gsiUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)
	update := dynamodbv1.IndexUpdateAction{
		Name:            "test-gsi",
		IndexThroughput: targetIndexCapacity,
	}
	gsiUpdates = append(gsiUpdates, &update)

	got, err := c.UpdateCapacity(context.Background(), "us-east-1", "test-gsi-table", targetTableCapacity, gsiUpdates, true)
	assert.NotNil(t, got)
	assert.Nil(t, err)
	assert.Equal(t, got.Status, dynamodbv1.Table_Status(3))
	assert.Equal(t, got.GlobalSecondaryIndexes[0].Status, dynamodbv1.GlobalSecondaryIndex_Status(3))
}

func TestOnDemandCheck(t *testing.T) {
	m := &mockDynamodb{
		table: testOnDemandTable,
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
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodbCfg: d, dynamodb: m},
				},
			},
		},
	}

	targetTableCapacity := &dynamodbv1.Throughput{
		ReadCapacityUnits:  10,
		WriteCapacityUnits: 10,
	}

	gsiUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)

	result, err := c.UpdateCapacity(context.Background(), "us-east-1", "test-ondemand-table", targetTableCapacity, gsiUpdates, false)
	assert.Nil(t, result)
	assert.Equal(t, "rpc error: code = FailedPrecondition desc = Table billing mode is not set to PROVISIONED, cannot scale capacities.", err.Error())
}

func TestOnDemandCheckInferred(t *testing.T) {
	m := &mockDynamodb{
		table: testOnDemandInferredTable,
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
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodbCfg: d, dynamodb: m},
				},
			},
		},
	}

	targetTableCapacity := &dynamodbv1.Throughput{
		ReadCapacityUnits:  10,
		WriteCapacityUnits: 10,
	}

	gsiUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)

	result, err := c.UpdateCapacity(context.Background(), "us-east-1", "test-ondemand-table", targetTableCapacity, gsiUpdates, false)
	assert.Nil(t, result)
	assert.Equal(t, "rpc error: code = FailedPrecondition desc = Table billing mode is not set to PROVISIONED, cannot scale capacities.", err.Error())
}

func TestNoBillingMode(t *testing.T) {
	m := &mockDynamodb{
		table:  testTableNoBillingMode,
		update: testTableNoBillingModeResponse,
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
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodbCfg: d, dynamodb: m},
				},
			},
		},
	}

	targetTableCapacity := &dynamodbv1.Throughput{
		ReadCapacityUnits:  101,
		WriteCapacityUnits: 202,
	}

	gsiUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)

	got, err := c.UpdateCapacity(context.Background(), "us-east-1", "test-table-no-billing", targetTableCapacity, gsiUpdates, false)
	assert.NotNil(t, got)
	assert.Nil(t, err)
	assert.Equal(t, got.Status, dynamodbv1.Table_Status(3))
}

type mockDynamodb struct {
	dynamodbClient

	tableErr error
	table    *types.TableDescription

	updateErr error
	update    *types.TableDescription
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

	ret := &dynamodb.UpdateTableOutput{
		TableDescription: m.update,
	}

	return ret, nil
}
