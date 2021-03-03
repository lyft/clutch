package aws

import (
	"context"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	awsv1 "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/resolver"
)

func (r *res) resolveAutoscalingGroupsForInput(ctx context.Context, input proto.Message) (*resolver.Results, error) {
	switch i := input.(type) {
	case *awsv1.AutoscalingGroupName:
		return r.autoscalingGroupResults(ctx, i.Region, []string{i.Name}, 1)
	default:
		return nil, status.Errorf(codes.Internal, "unrecognized input type '%T'", i)
	}
}

func (r *res) autoscalingGroupResults(ctx context.Context, region string, ids []string, limit uint32) (*resolver.Results, error) {
	ctx, handler := resolver.NewFanoutHandler(ctx)

	regions := r.determineRegionsForOption(region)
	for _, region := range regions {
		handler.Add(1)
		go func(region string) {
			defer handler.Done()
			groups, err := r.client.DescribeAutoscalingGroups(ctx, region, ids)
			select {
			case handler.Channel() <- resolver.NewFanoutResult(groups, err):
				return
			case <-handler.Cancelled():
				return
			}
		}(region)
	}

	return handler.Results(limit)
}
