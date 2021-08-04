package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
	awsv1 "github.com/lyft/clutch/backend/api/config/service/aws/v1"
)

<<<<<<< HEAD
// defaults for the dynamodb settings config
const (
	AwsMaxRCU       = int64(40000)
	AwsMaxWCU       = int64(40000)
	SafeScaleFactor = float32(2.0)
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
=======
const maxDynamodbReadCapacityUnits = 35000
const maxDynamodbWriteCapacityUnits = 50000
const maxDynamodbCapacityScaleFactor = 2
>>>>>>> 695a4f1671e96a397a5acff83d6e6b5232650a14

func (c *client) DescribeTable(ctx context.Context, region string, tableName string) (*dynamodbv1.Table, error) {
	cl, err := c.getRegionalClient(region)
	if err != nil {
		return nil, err
	}

	result, err := getTable(ctx, cl, tableName)
	if err != nil {
		return nil, err
	}

	currentCapacity := &dynamodbv1.ProvisionedThroughput{
		WriteCapacityUnits: aws.ToInt64(result.Table.ProvisionedThroughput.WriteCapacityUnits),
		ReadCapacityUnits:  aws.ToInt64(result.Table.ProvisionedThroughput.ReadCapacityUnits),
	}

	globalSecondaryIndexes := getGlobalSecondaryIndexes(result.Table.GlobalSecondaryIndexes)

	ret := &dynamodbv1.Table{
		Name:                   aws.ToString(result.Table.TableName),
		Region:                 region,
		GlobalSecondaryIndexes: globalSecondaryIndexes,
		ProvisionedThroughput:  currentCapacity,
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

<<<<<<< HEAD
func isValidIncrease(client *regionalClient, current *types.ProvisionedThroughputDescription, target types.ProvisionedThroughput) error {
	// check for targets that are lower than current (can't scale down)
	if *current.ReadCapacityUnits > *target.ReadCapacityUnits {
		return status.Errorf(codes.FailedPrecondition, fmt.Sprintf("Target read capacity [%d] is lower than current capacity [%d]", *target.ReadCapacityUnits, *current.ReadCapacityUnits))
	}
	if *current.WriteCapacityUnits > *target.WriteCapacityUnits {
		return status.Errorf(codes.FailedPrecondition, fmt.Sprintf("Target write capacity [%d] is lower than current capacity [%d]", *target.WriteCapacityUnits, *current.WriteCapacityUnits))
	}

	// check for targets that exceed max limits
	if *target.ReadCapacityUnits > client.dynamodbCfg.scalingLimits.maxReadCapacityUnits {
		return status.Errorf(codes.FailedPrecondition, fmt.Sprintf("Target read capacity exceeds maximum allowed limits [%d]", client.dynamodbCfg.scalingLimits.maxReadCapacityUnits))
	}
	if *target.WriteCapacityUnits > client.dynamodbCfg.scalingLimits.maxWriteCapacityUnits {
		return status.Errorf(codes.FailedPrecondition, fmt.Sprintf("Target write capacity exceeds maximum allowed limits [%d]", client.dynamodbCfg.scalingLimits.maxWriteCapacityUnits))
	}

	// check for increases that exceed max increase scale
	if (float32(*target.ReadCapacityUnits / *current.ReadCapacityUnits)) > client.dynamodbCfg.scalingLimits.maxScaleFactor {
		return status.Errorf(codes.FailedPrecondition, fmt.Sprintf("Target read capacity exceeds the scale limit of [%.1f]x current capacity", client.dynamodbCfg.scalingLimits.maxScaleFactor))
	}
	if (float32(*target.WriteCapacityUnits / *current.WriteCapacityUnits)) > client.dynamodbCfg.scalingLimits.maxScaleFactor {
		return status.Errorf(codes.FailedPrecondition, fmt.Sprintf("Target write capacity exceeds the scale limit of [%.1f]x current capacity", client.dynamodbCfg.scalingLimits.maxScaleFactor))
=======
func isValidIncrease(current *types.ProvisionedThroughputDescription, target types.ProvisionedThroughput) error {
	// check for targets that are lower than current (can't scale down)
	if *current.ReadCapacityUnits > *target.ReadCapacityUnits {
		return status.Error(codes.FailedPrecondition, "Target read capacity is lower than current capacity.")
	}
	if *current.WriteCapacityUnits > *target.WriteCapacityUnits {
		return status.Error(codes.FailedPrecondition, "Target write capacity is lower than current capacity.")
	}

	// check for targets that exceed max limits
	if *target.ReadCapacityUnits > maxDynamodbReadCapacityUnits {
		return status.Error(codes.FailedPrecondition, "Target read capacity exceeds maximum allowed limits.")
	}
	if *target.WriteCapacityUnits > maxDynamodbWriteCapacityUnits {
		return status.Error(codes.FailedPrecondition, "Target write capacity exceeds maximum allowed limits.")
	}

	// check for increases that exceed max increase scale
	if (*target.ReadCapacityUnits / *current.ReadCapacityUnits) > maxDynamodbCapacityScaleFactor {
		return status.Error(codes.FailedPrecondition, "Target read capacity exceeds scale limit.")
	}
	if (*target.WriteCapacityUnits / *current.WriteCapacityUnits) > maxDynamodbCapacityScaleFactor {
		return status.Error(codes.FailedPrecondition, "Target write capacity exceeds scale limit.")
>>>>>>> 695a4f1671e96a397a5acff83d6e6b5232650a14
	}
	return nil
}

func (c *client) UpdateTableCapacity(ctx context.Context, region string, tableName string, targetTableRcu int64, targetTableWcu int64) error {
	cl, err := c.getRegionalClient(region)
	if err != nil {
		return err
	}

<<<<<<< HEAD
<<<<<<< HEAD
=======
	// s, err := c.GetDynamodbConfig()

>>>>>>> finish merge
=======
	// s, err := c.GetDynamodbConfig()

=======
>>>>>>> 695a4f1671e96a397a5acff83d6e6b5232650a14
>>>>>>> Revert "finish merge"
	currentTable, err := getTable(ctx, cl, tableName)
	if err != nil {
		return err
	}

	targetCapacity := types.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(targetTableRcu),
		WriteCapacityUnits: aws.Int64(targetTableWcu),
	}

<<<<<<< HEAD
	err = isValidIncrease(cl, currentTable.Table.ProvisionedThroughput, targetCapacity)
=======
	err = isValidIncrease(currentTable.Table.ProvisionedThroughput, targetCapacity)
>>>>>>> 695a4f1671e96a397a5acff83d6e6b5232650a14
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateTableInput{
		TableName:             aws.String(tableName),
		ProvisionedThroughput: &targetCapacity,
	}

	_, err = cl.dynamodb.UpdateTable(ctx, input)
	return err
}
