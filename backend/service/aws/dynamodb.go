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
		c.log.Error("unable to find table", zap.Error(err))
		return nil, err
	}

	currentCapacity := &dynamodbv1.ProvisionedThroughput{
		WriteCapacityUnits: aws.ToInt64(result.Table.ProvisionedThroughput.WriteCapacityUnits),
		ReadCapacityUnits:  aws.ToInt64(result.Table.ProvisionedThroughput.ReadCapacityUnits),
	}

	globalSecondaryIndexes := getGlobalSecondaryIndexes(result.Table.GlobalSecondaryIndexes)

	status := newProtoForTableStatus(result.Table.TableStatus)

	ret := &dynamodbv1.Table{
		Name:                   aws.ToString(result.Table.TableName),
		Region:                 region,
		GlobalSecondaryIndexes: globalSecondaryIndexes,
		ProvisionedThroughput:  currentCapacity,
		Status:                 status,
	}
	return ret, nil
}

func getTable(ctx context.Context, client *regionalClient, tableName string) (*dynamodb.DescribeTableOutput, error) {
	input := &dynamodb.DescribeTableInput{TableName: aws.String(tableName)}
	return client.dynamodb.DescribeTable(ctx, input)
}

func getGlobalSecondaryIndexes(indexes []types.GlobalSecondaryIndexDescription) []*dynamodbv1.GlobalSecondaryIndex {
	gsis := make([]*dynamodbv1.GlobalSecondaryIndex, len(indexes))
	for idx, i := range indexes {
		gsis[idx] = newProtoForGlobalSecondaryIndex(i)
	}
	return gsis
}

func newProtoForTableStatus(s types.TableStatus) dynamodbv1.Status {
	value, ok := dynamodbv1.Status_value[string(s)]
	if !ok {
		return dynamodbv1.Status_UNSPECIFIED
	}
	return dynamodbv1.Status(value)
}

func newProtoForGlobalSecondaryIndex(index types.GlobalSecondaryIndexDescription) *dynamodbv1.GlobalSecondaryIndex {
	currentCapacity := &dynamodbv1.ProvisionedThroughput{
		ReadCapacityUnits:  aws.ToInt64(index.ProvisionedThroughput.ReadCapacityUnits),
		WriteCapacityUnits: aws.ToInt64(index.ProvisionedThroughput.WriteCapacityUnits),
	}

	return &dynamodbv1.GlobalSecondaryIndex{
		Name:                  aws.ToString(index.IndexName),
		ProvisionedThroughput: currentCapacity,
	}
}

func isValidIncrease(client *regionalClient, current *types.ProvisionedThroughputDescription, target types.ProvisionedThroughput) error {
	// check for targets that are lower than current (can't scale down)
	if *current.ReadCapacityUnits > *target.ReadCapacityUnits {
		return status.Errorf(codes.FailedPrecondition, "Target read capacity [%d] is lower than current capacity [%d]", *target.ReadCapacityUnits, *current.ReadCapacityUnits)
	}
	if *current.WriteCapacityUnits > *target.WriteCapacityUnits {
		return status.Errorf(codes.FailedPrecondition, "Target write capacity [%d] is lower than current capacity [%d]", *target.WriteCapacityUnits, *current.WriteCapacityUnits)
	}

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
	return nil
}

func (c *client) UpdateTableCapacity(ctx context.Context, region string, tableName string, targetTableRcu int64, targetTableWcu int64) (*dynamodbv1.Status, error) {
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

	targetCapacity := types.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(targetTableRcu),
		WriteCapacityUnits: aws.Int64(targetTableWcu),
	}

	err = isValidIncrease(cl, currentTable.Table.ProvisionedThroughput, targetCapacity)
	if err != nil {
		c.log.Error("invalid requested amount for capacity increase", zap.Error(err))
		return nil, err
	}

	input := &dynamodb.UpdateTableInput{
		TableName:             aws.String(tableName),
		ProvisionedThroughput: &targetCapacity,
	}

	result, err := cl.dynamodb.UpdateTable(ctx, input)

	if err != nil {
		c.log.Error("update table failed", zap.Error(err))
		return nil, err
	}

	tableStatus := newProtoForTableStatus(result.TableDescription.TableStatus)

	return &tableStatus, nil
}
