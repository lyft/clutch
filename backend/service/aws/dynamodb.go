package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
)

func (c *client) DescribeTable(ctx context.Context, region string, tableName string) (*dynamodbv1.Table, error) {
	cl, err := c.getRegionalClient(region)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DescribeTableInput{TableName: aws.String(tableName)}
	result, err := cl.dynamodb.DescribeTable(ctx, input)
	if err != nil {
		return nil, err
	}

	currentCapacity := &dynamodbv1.ProvisionedThroughput{
		WriteCapacityUnits: aws.ToInt64(result.Table.ProvisionedThroughput.WriteCapacityUnits),
		ReadCapacityUnits:  aws.ToInt64(result.Table.ProvisionedThroughput.ReadCapacityUnits),
	}

	globalSecondaryIndexes := getGlobalSecondaryIndexes(result.Table.GlobalSecondaryIndexes)

	var ret = &dynamodbv1.Table{
		Name:                   aws.ToString(result.Table.TableName),
		Region:                 region,
		GlobalSecondaryIndexes: globalSecondaryIndexes,
		ProvisionedThroughput:  currentCapacity,
	}
	return ret, nil

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
