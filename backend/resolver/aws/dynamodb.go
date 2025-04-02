package aws

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	awsv1 "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/resolver"
)

func (r *res) resolveDynamodbTableForInput(ctx context.Context, input proto.Message) (*resolver.Results, error) {
	switch i := input.(type) {
	case *awsv1.DynamodbTableName:
		return r.dynamodbResults(ctx, i.Account, i.Region, i.Name, 1)
	default:
		return nil, status.Errorf(codes.Internal, "unrecognized input type '%T'", i)
	}
}

func (r *res) dynamodbResults(ctx context.Context, account, region string, name string, limit uint32) (*resolver.Results, error) {
	ctx, handler := resolver.NewFanoutHandler(ctx)

	allAccountRegions := r.determineAccountAndRegionsForOption(account, region)
	for account := range allAccountRegions {
		for _, region := range allAccountRegions[account] {
			handler.Add(1)
			go func(account, region string) {
				defer handler.Done()
				table, err := r.client.DescribeTable(ctx, account, region, name)
				select {
				case handler.Channel() <- resolver.NewSingleFanoutResult(table, err):
					return
				case <-handler.Cancelled():
					return
				}
			}(account, region)
		}
	}

	return handler.Results(limit)
}
