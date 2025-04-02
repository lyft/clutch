package aws

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	awsv1 "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/resolver"
)

func (r *res) resolveAutoscalingGroupsForInput(ctx context.Context, input proto.Message) (*resolver.Results, error) {
	switch i := input.(type) {
	case *awsv1.AutoscalingGroupName:
		return r.autoscalingGroupResults(ctx, i.Account, i.Region, []string{i.Name}, 1)
	default:
		return nil, status.Errorf(codes.Internal, "unrecognized input type '%T'", i)
	}
}

func (r *res) autoscalingGroupResults(ctx context.Context, account, region string, ids []string, limit uint32) (*resolver.Results, error) {
	ctx, handler := resolver.NewFanoutHandler(ctx)

	allAccountRegions := r.determineAccountAndRegionsForOption(account, region)
	for account := range allAccountRegions {
		for _, region := range allAccountRegions[account] {
			handler.Add(1)
			go func(account, region string) {
				defer handler.Done()
				groups, err := r.client.DescribeAutoscalingGroups(ctx, account, region, ids)
				select {
				case handler.Channel() <- resolver.NewFanoutResult(groups, err):
					return
				case <-handler.Cancelled():
					return
				}
			}(account, region)
		}
	}

	return handler.Results(limit)
}
