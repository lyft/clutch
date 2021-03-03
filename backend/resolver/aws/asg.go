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
		// if the pattern matches then proceeded
		mappedValues, err := meta.PatternValueMapping(&ec2v1.AutoscalingGroup{}, id)
		if err == nil && len(mappedValues["region"]) > 0 && len(mappedValues["name"]) > 0 {
			handler.Add(1)
			go func(region string, name []string) {
				defer handler.Done()
				groups, err := r.client.DescribeAutoscalingGroups(ctx, region, name)
				select {
				case handler.Channel() <- resolver.NewFanoutResult(groups, err):
					return
				case <-handler.Cancelled():
					return
				}
			}(mappedValues["region"], []string{mappedValues["name"]})
		} else {
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
		}
	}
	return handler.Results(limit)
}
