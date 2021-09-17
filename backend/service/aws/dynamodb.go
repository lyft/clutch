package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
	awsv1 "github.com/lyft/clutch/backend/api/config/service/aws/v1"
)

// defaults for the dynamodb settings config
const (
	AwsMaxRCU       = 40000
	AwsMaxWCU       = 40000
	SafeScaleFactor = 2.0
)

// get or set defaults for dynamodb scaling
func getScalingLimits(cfg *awsv1.Config) *awsv1.ScalingLimits {
	if cfg.GetDynamodbConfig() == nil && cfg.DynamodbConfig.GetScalingLimits() == nil {
		ds := &awsv1.ScalingLimits{
			MaxReadCapacityUnits:  AwsMaxRCU,
			MaxWriteCapacityUnits: AwsMaxWCU,
			MaxScaleFactor:        SafeScaleFactor,
			EnableOverride:        false,
		}
		return ds
	}
	return cfg.DynamodbConfig.ScalingLimits
}

func (c *client) DescribeTable(ctx context.Context, region string, tableName string) (*dynamodbv1.Table, error) {
	cl, err := c.getRegionalClient(region)
	if err != nil {
		c.log.Error("unable to get regional client", zap.Error(err))
		return nil, err
	}

	result, err := getTable(ctx, cl, tableName)

	if err != nil {
		c.log.Warn("unable to find table in region: "+region, zap.Error(err))
		return nil, err
	}

	c.log.Info("generating new proto")
	ret := newProtoForTable(result.Table, region)

	return ret, nil
}

func getTable(ctx context.Context, client *regionalClient, tableName string) (*dynamodb.DescribeTableOutput, error) {
	input := &dynamodb.DescribeTableInput{TableName: aws.String(tableName)}
	return client.dynamodb.DescribeTable(ctx, input)
}

// generator for the list of GSI protos
func getGlobalSecondaryIndexes(indexes []types.GlobalSecondaryIndexDescription) []*dynamodbv1.GlobalSecondaryIndex {
	gsis := make([]*dynamodbv1.GlobalSecondaryIndex, len(indexes))
	for idx, i := range indexes {
		gsis[idx] = newProtoForGlobalSecondaryIndex(i)
	}
	return gsis
}

// retrieve one GSI from table description
func getGlobalSecondaryIndex(indexes []types.GlobalSecondaryIndexDescription, targetIndexName string) (*types.GlobalSecondaryIndexDescription, error) {
	for _, i := range indexes {
		if *i.IndexName == targetIndexName {
			return &i, nil
		}
	}
	return nil, status.Error(codes.NotFound, "Global secondary index not found.")
}

func newProtoForTable(t *types.TableDescription, region string) *dynamodbv1.Table {
	currentCapacity := &dynamodbv1.Throughput{
		WriteCapacityUnits: aws.ToInt64(t.ProvisionedThroughput.WriteCapacityUnits),
		ReadCapacityUnits:  aws.ToInt64(t.ProvisionedThroughput.ReadCapacityUnits),
	}

	globalSecondaryIndexes := getGlobalSecondaryIndexes(t.GlobalSecondaryIndexes)

	tableStatus := newProtoForTableStatus(t.TableStatus)

	billingMode := newProtoForBillingMode(t.BillingModeSummary)

	ret := &dynamodbv1.Table{
		Name:                   aws.ToString(t.TableName),
		Region:                 region,
		GlobalSecondaryIndexes: globalSecondaryIndexes,
		ProvisionedThroughput:  currentCapacity,
		Status:                 tableStatus,
		BillingMode:            billingMode,
	}

	return ret
}

func newProtoForTableStatus(s types.TableStatus) dynamodbv1.Table_Status {
	value, ok := dynamodbv1.Table_Status_value[string(s)]
	if !ok {
		return dynamodbv1.Table_UNKNOWN
	}
	return dynamodbv1.Table_Status(value)
}

func newProtoForIndexStatus(s types.IndexStatus) dynamodbv1.GlobalSecondaryIndex_Status {
	value, ok := dynamodbv1.GlobalSecondaryIndex_Status_value[string(s)]
	if !ok {
		return dynamodbv1.GlobalSecondaryIndex_UNKNOWN
	}
	return dynamodbv1.GlobalSecondaryIndex_Status(value)
}

func newProtoForBillingMode(s *types.BillingModeSummary) dynamodbv1.Table_BillingMode {
	value, ok := dynamodbv1.Table_BillingMode_value[string(s.BillingMode)]
	if !ok {
		return dynamodbv1.Table_BILLING_UNKNOWN
	}
	return dynamodbv1.Table_BillingMode(value)
}

func newProtoForGlobalSecondaryIndex(index types.GlobalSecondaryIndexDescription) *dynamodbv1.GlobalSecondaryIndex {
	currentCapacity := &dynamodbv1.Throughput{
		ReadCapacityUnits:  aws.ToInt64(index.ProvisionedThroughput.ReadCapacityUnits),
		WriteCapacityUnits: aws.ToInt64(index.ProvisionedThroughput.WriteCapacityUnits),
	}

	indexStatus := newProtoForIndexStatus(index.IndexStatus)

	return &dynamodbv1.GlobalSecondaryIndex{
		Name:                  aws.ToString(index.IndexName),
		ProvisionedThroughput: currentCapacity,
		Status:                indexStatus,
	}
}

func isValidIncrease(client *regionalClient, current *types.ProvisionedThroughputDescription, target types.ProvisionedThroughput, ignore_maximums bool) error {
	// check for targets that are lower than current (can't scale down)
	if *current.ReadCapacityUnits > *target.ReadCapacityUnits {
		return status.Errorf(codes.FailedPrecondition, "Target read capacity [%d] is lower than current capacity [%d]", *target.ReadCapacityUnits, *current.ReadCapacityUnits)
	}
	if *current.WriteCapacityUnits > *target.WriteCapacityUnits {
		return status.Errorf(codes.FailedPrecondition, "Target write capacity [%d] is lower than current capacity [%d]", *target.WriteCapacityUnits, *current.WriteCapacityUnits)
	}

	// override not enabled in config or override not set to true in args
	if !client.dynamodbCfg.ScalingLimits.EnableOverride || !ignore_maximums {
		// check for targets that exceed max limits
		if *target.ReadCapacityUnits > client.dynamodbCfg.ScalingLimits.MaxReadCapacityUnits {
			return status.Errorf(codes.FailedPrecondition, "Target read capacity exceeds maximum allowed limits [%d]", client.dynamodbCfg.ScalingLimits.MaxReadCapacityUnits)
		}
		if *target.WriteCapacityUnits > client.dynamodbCfg.ScalingLimits.MaxWriteCapacityUnits {
			return status.Errorf(codes.FailedPrecondition, "Target write capacity exceeds maximum allowed limits [%d]", client.dynamodbCfg.ScalingLimits.MaxWriteCapacityUnits)
		}

		// check for increases that exceed max increase scale
		if (float32(*target.ReadCapacityUnits / *current.ReadCapacityUnits)) > client.dynamodbCfg.ScalingLimits.MaxScaleFactor {
			return status.Errorf(codes.FailedPrecondition, "Target read capacity exceeds the scale limit of [%.1f]x current capacity", client.dynamodbCfg.ScalingLimits.MaxScaleFactor)
		}
		if (float32(*target.WriteCapacityUnits / *current.WriteCapacityUnits)) > client.dynamodbCfg.ScalingLimits.MaxScaleFactor {
			return status.Errorf(codes.FailedPrecondition, "Target write capacity exceeds the scale limit of [%.1f]x current capacity", client.dynamodbCfg.ScalingLimits.MaxScaleFactor)
		}
	}
	return nil
}

func isProvisioned(t *dynamodb.DescribeTableOutput) bool {
	return t.Table.BillingModeSummary.BillingMode == "PROVISIONED"
}

func (c *client) UpdateCapacity(ctx context.Context, region string, tableName string, targetTableCapacity *dynamodbv1.Throughput, indexUpdates []*dynamodbv1.IndexUpdateAction, ignore_maximums bool) (*dynamodbv1.Table, error) {
	cl, err := c.getRegionalClient(region)
	if err != nil {
		c.log.Error("unable to get regional client", zap.Error(err))
		return nil, err
	}

	currentTable, err := getTable(ctx, cl, tableName)
	if err != nil {
		c.log.Error("unable to find table", zap.Error(err))
		return nil, err
	}

	if !(isProvisioned(currentTable)) {
		return nil, status.Error(codes.FailedPrecondition, "Table billing mode is not set to PROVISIONED, cannot scale capacities.")
	}

	// initialize the input we'll pass to AWS
	input := &dynamodb.UpdateTableInput{
		TableName: aws.String(tableName),
	}

	// add target table capacity to input if valid
	if targetTableCapacity != nil { // received new vals for table
		t := types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(targetTableCapacity.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(targetTableCapacity.WriteCapacityUnits),
		}

		err = isValidIncrease(cl, currentTable.Table.ProvisionedThroughput, t, ignore_maximums)
		if err != nil {
			return nil, err
		}

		input.ProvisionedThroughput = &t
	}

	// add target index capacities to input if valid
	if len(indexUpdates) > 0 { // received at least one index update
		updates, err := generateIndexUpdates(cl, currentTable, indexUpdates, ignore_maximums)
		if err != nil {
			return nil, err
		}
		input.GlobalSecondaryIndexUpdates = append(input.GlobalSecondaryIndexUpdates, updates...)
	}

	result, err := cl.dynamodb.UpdateTable(ctx, input)
	if err != nil {
		return nil, err
	}

	ret := newProtoForTable(result.TableDescription, region)

	return ret, nil
}

// given a list of ddbv1.updates, generate the AWS types for updates
// TO DO: this function currently terminates upon encountering one invalid index update,
// preserving the AWS SDK design of "all-or-nothing" UpdateTable functionality.
// confirm if that is still the desired behavior. Alternatively, we can selectively make updates.
func generateIndexUpdates(cl *regionalClient, t *dynamodb.DescribeTableOutput, indexUpdates []*dynamodbv1.IndexUpdateAction, ignore_maximums bool) ([]types.GlobalSecondaryIndexUpdate, error) {
	currentIndexes := t.Table.GlobalSecondaryIndexes

	var updates []types.GlobalSecondaryIndexUpdate

	for _, i := range indexUpdates {
		currentIndex, err := getGlobalSecondaryIndex(currentIndexes, i.Name)
		if err != nil {
			return nil, err
		}

		targetIndexCapacity := types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(i.IndexThroughput.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(i.IndexThroughput.WriteCapacityUnits),
		}

		err = isValidIncrease(cl, currentIndex.ProvisionedThroughput, targetIndexCapacity, ignore_maximums)
		if err != nil {
			return nil, err
		}

		newIndexUpdate := types.GlobalSecondaryIndexUpdate{
			Update: &types.UpdateGlobalSecondaryIndexAction{
				IndexName: aws.String(i.Name),
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(i.IndexThroughput.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(i.IndexThroughput.WriteCapacityUnits),
				},
			},
		}

		updates = append(updates, newIndexUpdate)
	}

	return updates, nil
}
