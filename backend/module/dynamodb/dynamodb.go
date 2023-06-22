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
	table, err := a.client.DescribeTable(ctx, req.Account, req.Region, req.TableName)
	if err != nil {
		return nil, err
	}

	return &dynamodbv1.DescribeTableResponse{Table: table}, nil
}

func (a *dynamodbAPI) DescribeContinuousBackups(ctx context.Context, req *dynamodbv1.DescribeContinuousBackupsRequest) (*dynamodbv1.DescribeContinuousBackupsResponse, error) {
	backups, err := a.client.DescribeContinuousBackups(ctx, req.Account, req.Region, req.TableName)
	if err != nil {
		return nil, err
	}

	return &dynamodbv1.DescribeContinuousBackupsResponse{ContinuousBackups: backups}, nil
}

func (a *dynamodbAPI) UpdateCapacity(ctx context.Context, req *dynamodbv1.UpdateCapacityRequest) (*dynamodbv1.UpdateCapacityResponse, error) {
	result, err := a.client.UpdateCapacity(ctx, req.Account, req.Region, req.TableName, req.TableThroughput, req.GsiUpdates, req.IgnoreMaximums)
	if err != nil {
		return nil, err
	}

	return &dynamodbv1.UpdateCapacityResponse{Table: result}, nil
}
