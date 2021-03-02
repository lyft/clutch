package aws

import (
	"context"
	"fmt"
	"regexp"

	"github.com/golang/protobuf/proto"

	awsv1 "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/resolver"
)

var regionASGInputSyntax = regexp.MustCompile(`(.*)[/](.*)`)

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
		// if we have a region in the search query then try it
		if regionASGInputSyntax.MatchString(id) {
			handler.Add(1)
			result := regionASGInputSyntax.FindAllStringSubmatch(id, -1)

			go func(region, asg string) {
				defer handler.Done()
				groups, err := r.client.DescribeAutoscalingGroups(ctx, region, []string{asg})
				select {
				case handler.Channel() <- resolver.NewFanoutResult(groups, err):
					return
				case <-handler.Cancelled():
					return
				}
			}(result[0][1], result[0][2])
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

	// regions := r.determineRegionsForOption(region)
	// for _, region := range regions {
	// 	handler.Add(1)
	// 	go func(region string) {
	// 		defer handler.Done()
	// 		groups, err := r.client.DescribeAutoscalingGroups(ctx, region, ids)
	// 		select {
	// 		case handler.Channel() <- resolver.NewFanoutResult(groups, err):
	// 			return
	// 		case <-handler.Cancelled():
	// 			return
	// 		}
	// 	}(region)
	// }

	return handler.Results(limit)
}
