package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
)

const maxDynamodbReadCapacityUnits = 35000
const maxDynamodbWriteCapacityUnits = 50000
const maxDynamodbCapacityScaleFactor = 2

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
		return status.Error(codes.FailedPrecondition, "Target write capacity exceeds scale limit. If you really intend to more than double the capacity, please reach out to #core-datastores oncall to request change.")
	}
	return nil
}

func (c *client) UpdateTableCapacity(ctx context.Context, region string, tableName string, targetTableRcu int64, targetTableWcu int64) error {
	cl, err := c.getRegionalClient(region)
	if err != nil {
		return err
	}

	currentTable, err := getTable(ctx, cl, tableName)
	if err != nil {
		return err
	}

	targetCapacity := types.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(targetTableRcu),
		WriteCapacityUnits: aws.Int64(targetTableWcu),
	}

	err = isValidIncrease(currentTable.Table.ProvisionedThroughput, targetCapacity)
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
