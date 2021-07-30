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

func (a *dynamodbAPI) UpdateTableCapacity(ctx context.Context, req *dynamodbv1.UpdateTableCapacityRequest) (*dynamodbv1.UpdateTableCapacityResponse, error) {
	// TO DO

	return &dynamodbv1.UpdateTableCapacityResponse{}, nil
}

func (a *dynamodbAPI) UpdateGSICapacity(ctx context.Context, req *dynamodbv1.UpdateGSICapacityRequest) (*dynamodbv1.UpdateGSICapacityResponse, error) {
	// TO DO

	return &dynamodbv1.UpdateGSICapacityResponse{}, nil
}
