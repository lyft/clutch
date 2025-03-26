package aws

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/timestamppb"

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
	KeySchema: []types.KeySchemaElement{{
		AttributeName: aws.String("ID"),
		KeyType:       types.KeyTypeHash,
	},
		{
			AttributeName: aws.String("Status"),
			KeyType:       types.KeyTypeRange,
		},
	},
	AttributeDefinitions: []types.AttributeDefinition{{
		AttributeName: aws.String("OrderName"),
		AttributeType: types.ScalarAttributeTypeS,
	},
		{
			AttributeName: aws.String("OrderWeight"),
			AttributeType: types.ScalarAttributeTypeN,
		},
	},
}

var testTableOutput = &dynamodbv1.Table{
	Name:    "test-table",
	Account: "default",
	Region:  "us-east-1",
	ProvisionedThroughput: &dynamodbv1.Throughput{
		ReadCapacityUnits:  100,
		WriteCapacityUnits: 200,
	},
	GlobalSecondaryIndexes: []*dynamodbv1.GlobalSecondaryIndex{},
	Status:                 dynamodbv1.Table_Status(5),
	BillingMode:            dynamodbv1.Table_BillingMode(2),
	KeySchemas: []*dynamodbv1.KeySchema{
		{AttributeName: "ID",
			Type: dynamodbv1.KeySchema_HASH,
		},
		{AttributeName: "Status",
			Type: dynamodbv1.KeySchema_RANGE,
		},
	},
	AttributeDefinitions: []*dynamodbv1.AttributeDefinition{
		{
			AttributeName: "OrderName",
			AttributeType: "S",
		},
		{
			AttributeName: "OrderWeight",
			AttributeType: "N",
		},
	},
}

var testMaxDynamodbTable = &types.TableDescription{
	TableName: aws.String("test-max-table"),
	ProvisionedThroughput: &types.ProvisionedThroughputDescription{
		ReadCapacityUnits:  aws.Int64(45000),
		WriteCapacityUnits: aws.Int64(45000),
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndexDescription{
		{IndexName: aws.String("test-gsi"),
			KeySchema: []types.KeySchemaElement{},
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				ReadCapacityUnits:  aws.Int64(45000),
				WriteCapacityUnits: aws.Int64(45000),
			},
			IndexStatus: "ACTIVE",
		},
	},
	TableStatus: "ACTIVE",
	BillingModeSummary: &types.BillingModeSummary{
		BillingMode: "PROVISIONED",
	},
}

var testMaxDynamodbTableResponse = &types.TableDescription{
	TableName: aws.String("test-max-table"),
	ProvisionedThroughput: &types.ProvisionedThroughputDescription{
		ReadCapacityUnits:  aws.Int64(45000),
		WriteCapacityUnits: aws.Int64(45000),
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndexDescription{
		{IndexName: aws.String("test-gsi"),
			KeySchema: []types.KeySchemaElement{},
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				ReadCapacityUnits:  aws.Int64(45000),
				WriteCapacityUnits: aws.Int64(45000),
			},
			IndexStatus: "UPDATING",
		},
	},
	TableStatus: "UPDATING",
	BillingModeSummary: &types.BillingModeSummary{
		BillingMode: "PROVISIONED",
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
			IndexStatus: "ACTIVE",
		},
	},
	TableStatus: "ACTIVE",
	BillingModeSummary: &types.BillingModeSummary{
		BillingMode: "PROVISIONED",
	},
}

var testDynamodbTableWithGSIResponse = &types.TableDescription{
	TableName: aws.String("test-max-table"),
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

var testTableWithGSIOutput = &dynamodbv1.Table{
	Name:    "test-gsi-table",
	Account: "default",
	Region:  "us-east-1",
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
	Status:               dynamodbv1.Table_Status(5),
	BillingMode:          dynamodbv1.Table_BillingMode(2),
	KeySchemas:           []*dynamodbv1.KeySchema{},
	AttributeDefinitions: []*dynamodbv1.AttributeDefinition{},
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

var testContinuousBackupsDescription = &types.ContinuousBackupsDescription{
	ContinuousBackupsStatus: "ENABLED",
	PointInTimeRecoveryDescription: &types.PointInTimeRecoveryDescription{
		PointInTimeRecoveryStatus:  "DISABLED",
		LatestRestorableDateTime:   aws.Time(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
		EarliestRestorableDateTime: aws.Time(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
}

var testContinuousBackupsOutput = &dynamodbv1.ContinuousBackups{
	ContinuousBackupsStatus:    dynamodbv1.ContinuousBackups_Status(2),
	PointInTimeRecoveryStatus:  dynamodbv1.ContinuousBackups_Status(3),
	LatestRestorableDateTime:   timestamppb.New(*aws.Time(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))),
	EarliestRestorableDateTime: timestamppb.New(*aws.Time(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))),
}

var testContinuousBackupsDescriptionNoPitr = &types.ContinuousBackupsDescription{
	ContinuousBackupsStatus: "DISABLED",
}

var testContinuousBackupsNoPitrOutput = &dynamodbv1.ContinuousBackups{
	ContinuousBackupsStatus:    dynamodbv1.ContinuousBackups_Status(3),
	PointInTimeRecoveryStatus:  dynamodbv1.ContinuousBackups_Status(0),
	LatestRestorableDateTime:   nil,
	EarliestRestorableDateTime: nil,
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

	result, err := c.DescribeTable(context.Background(), "default", "us-east-1", "test-table")
	assert.NoError(t, err)
	assert.Equal(t, testTableOutput, result)

	m.tableErr = errors.New("error")
	_, err1 := c.DescribeTable(context.Background(), "default", "us-east-1", "test-table")
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
	_, err := c.DescribeTable(context.Background(), "default", "us-east-1", "nonexistent-table")
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

	result, err := c.DescribeTable(context.Background(), "default", "us-east-1", "test-gsi-table")
	assert.NoError(t, err)
	assert.Equal(t, testTableWithGSIOutput, result)
}

func TestDescribeContinuousBackups(t *testing.T) {
	m := &mockDynamodb{
		backups: testContinuousBackupsDescription,
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

	result, err := c.DescribeContinuousBackups(context.Background(), "default", "us-east-1", "test-table")
	assert.NoError(t, err)
	assert.Equal(t, testContinuousBackupsOutput, result)
}

func TestDescribeContinuousBackupsError(t *testing.T) {
	m := &mockDynamodb{
		backups:    testContinuousBackupsDescription,
		backupsErr: fmt.Errorf("describe continuous backups error"),
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

	_, err := c.DescribeContinuousBackups(context.Background(), "default", "us-east-1", "test-table")
	assert.Error(t, err)
}

func TestDescribeContinuousBackupsNoPitr(t *testing.T) {
	m := &mockDynamodb{
		backups: testContinuousBackupsDescriptionNoPitr,
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

	result, err := c.DescribeContinuousBackups(context.Background(), "default", "us-east-1", "test-table")
	assert.NoError(t, err)
	assert.Equal(t, testContinuousBackupsNoPitrOutput, result)
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

// case: using custom scaling factor
// These tests confirm the scaling factor limit is correctly applied
// by attempting to provision capacities that represent the following edge conditions:
// condition 1 (success) : target is exactly x times the scaling factor
// condition 2 (failure) : target is 1 above x times the scaling factor
func TestCustomScaleFactor(t *testing.T) {
	// scale factor to be used in this test
	// Change the number to represent a different limit
	testScaleFactor := float32(3.5)

	// for readability in the tests
	// casting capacities to float32 because dynamic test values are being generated by
	// multiplying them with the scale factor
	tableRCU := float32(*testDynamodbTableWithGSI.ProvisionedThroughput.ReadCapacityUnits)
	tableWCU := float32(*testDynamodbTableWithGSI.ProvisionedThroughput.WriteCapacityUnits)
	indexRCU := float32(*testDynamodbTableWithGSI.GlobalSecondaryIndexes[0].ProvisionedThroughput.ReadCapacityUnits)
	indexWCU := float32(*testDynamodbTableWithGSI.GlobalSecondaryIndexes[0].ProvisionedThroughput.WriteCapacityUnits)

	successTests := []struct {
		name           string
		targetTableRCU float32
		targetTableWCU float32
		targetIndexRCU float32
		targetIndexWCU float32
		status         int64
	}{
		{
			"table: valid target rcu within limits",
			(tableRCU * testScaleFactor),
			tableWCU,
			indexRCU,
			indexWCU,
			3,
		},
		{
			"table: valid target wcu within limits",
			tableRCU,
			tableWCU * testScaleFactor,
			indexRCU,
			indexWCU,
			3,
		},
		{
			"index: valid target rcu within limits",
			tableRCU,
			tableWCU,
			(indexRCU * testScaleFactor),
			indexWCU,
			3,
		},
		{
			"index: valid target wcu within limits",
			tableRCU,
			tableWCU,
			indexRCU,
			(indexWCU * testScaleFactor),
			3,
		},
	}

	errorTests := []struct {
		name           string
		targetTableRCU float32
		targetTableWCU float32
		targetIndexRCU float32
		targetIndexWCU float32
		want           string
	}{
		{
			"table: invalid target rcu above scaling limit",
			(tableRCU*testScaleFactor + 1),
			tableWCU,
			indexRCU,
			indexWCU,
			fmt.Sprintf("rpc error: code = FailedPrecondition desc = Target read capacity exceeds the scale limit of [%1.1f]x current capacity", testScaleFactor),
		},
		{
			"table: invalid target wcu above scaling limit",
			tableRCU,
			(tableWCU*testScaleFactor + 1),
			indexRCU,
			indexWCU,
			fmt.Sprintf("rpc error: code = FailedPrecondition desc = Target write capacity exceeds the scale limit of [%1.1f]x current capacity", testScaleFactor),
		},
		{
			"index: invalid target rcu above scaling limit",
			tableRCU,
			tableWCU,
			(indexRCU*testScaleFactor + 1),
			indexWCU,
			fmt.Sprintf("rpc error: code = FailedPrecondition desc = Target read capacity exceeds the scale limit of [%1.1f]x current capacity", testScaleFactor),
		},
		{
			"index: invalid target wcu above scaling limit",
			tableRCU,
			tableWCU,
			indexRCU,
			(indexWCU*testScaleFactor + 1),
			fmt.Sprintf("rpc error: code = FailedPrecondition desc = Target write capacity exceeds the scale limit of [%1.1f]x current capacity", testScaleFactor),
		},
	}

	m := &mockDynamodb{
		table:  testDynamodbTableWithGSI,
		update: testDynamodbTableWithGSIResponse,
	}

	cfg := &awsv1.Config{
		Regions: regions,
		ClientConfig: &awsv1.ClientConfig{
			Retries: 10,
		},
		DynamodbConfig: &awsv1.DynamodbConfig{
			ScalingLimits: &awsv1.ScalingLimits{
				MaxScaleFactor: testScaleFactor, // custom scale factor set here
				EnableOverride: false,
			},
		},
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

	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			targetTableCapacity := &dynamodbv1.Throughput{
				ReadCapacityUnits:  int64(tt.targetTableRCU),
				WriteCapacityUnits: int64(tt.targetTableWCU),
			}

			targetIndexCapacity := &dynamodbv1.Throughput{
				ReadCapacityUnits:  int64(tt.targetIndexRCU),
				WriteCapacityUnits: int64(tt.targetIndexWCU),
			}

			// generate and append index updates
			gsiUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)
			update := dynamodbv1.IndexUpdateAction{
				Name:            "test-gsi",
				IndexThroughput: targetIndexCapacity,
			}
			gsiUpdates = append(gsiUpdates, &update)

			result, err := c.UpdateCapacity(context.Background(), "default", "us-east-1", "test-max-table", targetTableCapacity, gsiUpdates, false)
			assert.Nil(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, result.Status, dynamodbv1.Table_Status(tt.status))
		})
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			targetTableCapacity := &dynamodbv1.Throughput{
				ReadCapacityUnits:  int64(tt.targetTableRCU),
				WriteCapacityUnits: int64(tt.targetTableWCU),
			}

			targetIndexCapacity := &dynamodbv1.Throughput{
				ReadCapacityUnits:  int64(tt.targetIndexRCU),
				WriteCapacityUnits: int64(tt.targetIndexWCU),
			}

			// generate and append index updates
			gsiUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)
			update := dynamodbv1.IndexUpdateAction{
				Name:            "test-gsi",
				IndexThroughput: targetIndexCapacity,
			}
			gsiUpdates = append(gsiUpdates, &update)

			result, err := c.UpdateCapacity(context.Background(), "default", "us-east-1", "test-max-table", targetTableCapacity, gsiUpdates, false)

			if err == nil {
				t.Error("Received nil instead of expected error")
			}

			if err.Error() != tt.want {
				t.Errorf("\nWant error msg: %s\nGot error msg: %s", tt.want, err)
			}

			assert.Nil(t, result)
		})
	}
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
		{"table rcu change scale too high", 201, 200, "rpc error: code = FailedPrecondition desc = Target read capacity exceeds the scale limit of [2.0]x current capacity"},
		{"table wcu change scale too high", 100, 401, "rpc error: code = FailedPrecondition desc = Target write capacity exceeds the scale limit of [2.0]x current capacity"},
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
		// capture range variable
		t.Run(tt.name, func(t *testing.T) {
			targetTableCapacity := &dynamodbv1.Throughput{
				ReadCapacityUnits:  tt.targetRCU,
				WriteCapacityUnits: tt.targetWCU,
			}

			gsiUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)

			result, err := c.UpdateCapacity(context.Background(), "default", "us-east-1", "test-table", targetTableCapacity, gsiUpdates, false)
			if err.Error() != tt.want {
				t.Errorf("\nWant error msg: %s\nGot error msg: %s", tt.want, err)
			}
			assert.Nil(t, result)
		})
	}
}

// Case: table already above the maximum
// These tests confirm the correct error messages are displayed if
// attempting to provision capacities above the set scale factor (default 2x).
// Capacity increases for tables already above the max limit
// (but within the scale factor limit) should pass.
func TestMaxTableUpdates(t *testing.T) {
	errorTests := []struct {
		name      string
		targetRCU int64
		targetWCU int64
		want      string
	}{
		{"target table rcu lower than current", 44000, 45000, "rpc error: code = FailedPrecondition desc = Target read capacity [44000] is lower than current capacity [45000]"},
		{"target table wcu lower than current", 45000, 44000, "rpc error: code = FailedPrecondition desc = Target write capacity [44000] is lower than current capacity [45000]"},
		{"table rcu change scale too high", 90001, 45000, "rpc error: code = FailedPrecondition desc = Target read capacity exceeds the scale limit of [2.0]x current capacity"},
		{"table wcu change scale too high", 45000, 90001, "rpc error: code = FailedPrecondition desc = Target write capacity exceeds the scale limit of [2.0]x current capacity"},
	}

	successTests := []struct {
		name      string
		targetRCU int64
		targetWCU int64
		status    int64
	}{
		{"target table rcu above limit", 46000, 45000, 3},
		{"target table wcu above limit", 45000, 46000, 3},
	}

	m := &mockDynamodb{
		table:  testMaxDynamodbTable,
		update: testMaxDynamodbTableResponse,
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

	emptyIndexUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)

	for _, tt := range errorTests {
		// capture range variable
		t.Run(tt.name, func(t *testing.T) {
			targetTableCapacity := &dynamodbv1.Throughput{
				ReadCapacityUnits:  tt.targetRCU,
				WriteCapacityUnits: tt.targetWCU,
			}

			result, err := c.UpdateCapacity(context.Background(), "default", "us-east-1", "test-max-table", targetTableCapacity, emptyIndexUpdates, false)
			if err.Error() != tt.want {
				t.Errorf("\nWant error msg: %s\nGot error msg: %s", tt.want, err)
			}
			assert.Nil(t, result)
		})
	}

	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			targetTableCapacity := &dynamodbv1.Throughput{
				ReadCapacityUnits:  tt.targetRCU,
				WriteCapacityUnits: tt.targetWCU,
			}

			targetIndexCapacity := &dynamodbv1.Throughput{
				ReadCapacityUnits:  45000,
				WriteCapacityUnits: 45000,
			}

			gsiUpdates := make([]*dynamodbv1.IndexUpdateAction, 0)
			update := dynamodbv1.IndexUpdateAction{
				Name:            "test-gsi",
				IndexThroughput: targetIndexCapacity,
			}
			gsiUpdates = append(gsiUpdates, &update)

			result, err := c.UpdateCapacity(context.Background(), "default", "us-east-1", "test-max-table", targetTableCapacity, gsiUpdates, false)
			assert.Nil(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, result.Status, dynamodbv1.Table_Status(tt.status))
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
		{"gsi rcu change scale too high", 21, 25, "rpc error: code = FailedPrecondition desc = Target read capacity exceeds the scale limit of [2.0]x current capacity"},
		{"gsi wcu change scale too high", 15, 41, "rpc error: code = FailedPrecondition desc = Target write capacity exceeds the scale limit of [2.0]x current capacity"},
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
		// capture range variable
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

			result, err := c.UpdateCapacity(context.Background(), "default", "us-east-1", "test-table", targetTableCapacity, gsiUpdates, false)
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

	got, err := c.UpdateCapacity(context.Background(), "default", "us-east-1", "test-gsi-table", targetTableCapacity, gsiUpdates, false)
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

	got, err := c.UpdateCapacity(context.Background(), "default", "us-east-1", "test-gsi-table", targetTableCapacity, gsiUpdates, true)
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

	result, err := c.UpdateCapacity(context.Background(), "default", "us-east-1", "test-ondemand-table", targetTableCapacity, gsiUpdates, false)
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

	result, err := c.UpdateCapacity(context.Background(), "default", "us-east-1", "test-ondemand-table", targetTableCapacity, gsiUpdates, false)
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

	got, err := c.UpdateCapacity(context.Background(), "default", "us-east-1", "test-table-no-billing", targetTableCapacity, gsiUpdates, false)
	assert.NotNil(t, got)
	assert.Nil(t, err)
	assert.Equal(t, got.Status, dynamodbv1.Table_Status(3))
}

func CreateBatchGetInput(tableName string) dynamodb.BatchGetItemInput {
	return dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			tableName: {
				Keys: []map[string]types.AttributeValue{
					{
						"Artist": &types.AttributeValueMemberS{
							Value: "No One You Know",
						},
						"SongTitle": &types.AttributeValueMemberS{
							Value: "Call Me Today",
						},
					},
					{
						"Artist": &types.AttributeValueMemberS{
							Value: "Acme Band",
						},
						"SongTitle": &types.AttributeValueMemberS{
							Value: "Happy Day",
						},
					},
				},
				ConsistentRead:       aws.Bool(false),
				ProjectionExpression: aws.String("AlbumTitle,SongTitle"),
			},
		},
	}
}

func TestBatchGetItemsValid(t *testing.T) {
	testTableName := "MusicTable"
	mockDDB := &mockDynamodb{
		batchGetOutput: dynamodb.BatchGetItemOutput{
			ConsumedCapacity: nil,
			Responses: map[string][]map[string]types.AttributeValue{
				testTableName: {
					{
						"AlbumTitle": &types.AttributeValueMemberS{
							Value: "AlbumOne",
						},
						"SongTitle": &types.AttributeValueMemberS{
							Value: "Call Me Today",
						},
					},
					{
						"AlbumTitle": &types.AttributeValueMemberS{
							Value: "AlbumTwo",
						},
						"SongTitle": &types.AttributeValueMemberS{
							Value: "Happy Day",
						},
					},
				},
			},
			UnprocessedKeys: nil,
		},
	}

	client := &client{
		log:                 zaptest.NewLogger(t),
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodb: mockDDB},
				},
			},
		},
	}

	input := CreateBatchGetInput(testTableName)
	result, err := client.BatchGetItem(context.Background(), "default", "us-east-1", &input)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Responses[testTableName]))
}

func TestBatchGetItemsError(t *testing.T) {
	mockDDB := &mockDynamodb{
		batchGetOutput: dynamodb.BatchGetItemOutput{},
		batchGetErr:    fmt.Errorf("test batch get items error"),
	}
	client := &client{
		log:                 zaptest.NewLogger(t),
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", dynamodb: mockDDB},
				},
			},
		},
	}

	input := CreateBatchGetInput("MusicTable")
	_, err := client.BatchGetItem(context.Background(), "default", "us-east-1", &input)
	assert.Error(t, err)
}

type mockDynamodb struct {
	dynamodbClient

	tableErr error
	table    *types.TableDescription

	updateErr error
	update    *types.TableDescription

	batchGetErr    error
	batchGetOutput dynamodb.BatchGetItemOutput

	backupsErr error
	backups    *types.ContinuousBackupsDescription
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

func (m *mockDynamodb) BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error) {
	if m.batchGetErr != nil {
		return nil, m.batchGetErr
	}

	return &m.batchGetOutput, nil
}

func (m *mockDynamodb) DescribeContinuousBackups(ctx context.Context, params *dynamodb.DescribeContinuousBackupsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeContinuousBackupsOutput, error) {
	if m.backupsErr != nil {
		return nil, m.backupsErr
	}

	ret := &dynamodb.DescribeContinuousBackupsOutput{
		ContinuousBackupsDescription: m.backups,
	}

	return ret, nil
}
