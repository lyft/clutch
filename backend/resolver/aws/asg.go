package aws

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	awsv1 "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/resolver"
)

func (r *res) resolveAutoscalingGroupsForInput(ctx context.Context, input proto.Message) (*resolver.Results, error) {
	switch i := input.(type) {
	case *awsv1.AutoscalingGroupName:
		return r.autoscalingGroupResults(ctx, i.Region, []string{i.Name}, 1)
	default:
		return nil, fmt.Errorf("unrecognized input type %T", i)
	}
}

func (r *res) autoscalingGroupResults(ctx context.Context, region string, ids []string, limit uint32) (*resolver.Results, error) {
	ctx, handler := resolver.NewFanoutHandler(ctx)

	for _, id := range ids {
		patternValues, ok, err := meta.ExtractPatternValuesFromString((*ec2v1.AutoscalingGroup)(nil), id)
		if err != nil {
			return nil, err
		}

		asgName := id
		regions := r.determineRegionsForOption(region)
		if ok {
			regions = []string{patternValues["region"]}
			asgName = patternValues["name"]
		}

		for _, region := range regions {
			handler.Add(1)
			go func(region, name string) {
				defer handler.Done()
				groups, err := r.client.DescribeAutoscalingGroups(ctx, region, []string{name})
				select {
				case handler.Channel() <- resolver.NewFanoutResult(groups, err):
					return
				case <-handler.Cancelled():
					return
				}
			}(region, asgName)
		}
	}
	return handler.Results(limit)
}
