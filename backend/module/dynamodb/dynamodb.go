package dynamodb

import (
	"context"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
	"github.com/lyft/clutch/backend/service/aws"
)

func newDDBAPI(c aws.Client) dynamodbv1.DDBAPIServer {
	return &dynamodbAPI{
		client: c,
	}
}

type dynamodbAPI struct {
	client aws.Client
}

func (a *dynamodbAPI) DescribeTable(ctx context.Context, req *dynamodbv1.DescribeTableRequest) (*dynamodbv1.DescribeTableResponse, error) {
	table, err := a.client.DescribeTable(ctx, req.Region, req.TableName)
	if err != nil {
		return nil, err
	}

	return &dynamodbv1.DescribeTableResponse{Table: table}, nil
}

func (a *dynamodbAPI) IncreaseCapacity(ctx context.Context, req *dynamodbv1.IncreaseCapacityRequest) (*dynamodbv1.IncreaseCapacityResponse, error) {
	result, err := a.client.IncreaseCapacity(ctx, req.Region, req.TableName, req.TableThroughput, req.GsiUpdates, req.IgnoreMaximums)
	if err != nil {
		return nil, err
	}

	return &dynamodbv1.IncreaseCapacityResponse{Table: result}, nil
}
