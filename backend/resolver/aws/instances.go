package aws

import (
	"context"
	"fmt"
	"regexp"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ec2v1api "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	awsv1 "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/resolver"
)

var instanceIDPattern = regexp.MustCompile("[a-fA-F0-9]{17}")

func normalizeInstanceID(input string) (string, error) {
	instanceID := instanceIDPattern.FindString(input)
	if instanceID == "" {
		return "", status.Error(codes.InvalidArgument, "did not understand input")
	}

	return fmt.Sprintf("i-%s", instanceID), nil
}

// Take any inputs that can get instance IDs and normalize them for the client.
func (r *res) resolveInstancesForInput(ctx context.Context, input proto.Message) (*resolver.Results, error) {
	switch i := input.(type) {
	case *awsv1.InstanceID:
		return r.instanceResults(ctx, i.Region, []string{i.Id}, 1)
	default:
		return nil, fmt.Errorf("unrecognized input type %T", i)
	}
}

// Fanout across multiple regions if needed to fetch instances.
func (r *res) instanceResults(ctx context.Context, region string, ids []string, limit uint32) (*resolver.Results, error) {
	ctx, handler := resolver.NewFanoutHandler(ctx)

	for _, id := range ids {
		patternValues, ok, err := meta.ExtractPatternValuesFromString((*ec2v1api.Instance)(nil), id)
		if err != nil {
			return nil, err
		}

		if ok {
			handler.Add(1)
			go func(region, instanceId string) {
				defer handler.Done()
				instances, err := r.client.DescribeInstances(ctx, region, []string{instanceId})
				select {
				case handler.Channel() <- resolver.NewFanoutResult(instances, err):
					return
				case <-handler.Cancelled():
					return
				}
			}(patternValues["region"], patternValues["instance_id"])
		} else {
			normalizedId, err := normalizeInstanceID(id)
			if err != nil {
				return nil, err
			}
			regions := r.determineRegionsForOption(region)
			for _, region := range regions {
				handler.Add(1)
				go func(region string) {
					defer handler.Done()
					instances, err := r.client.DescribeInstances(ctx, region, []string{normalizedId})
					select {
					case handler.Channel() <- resolver.NewFanoutResult(instances, err):
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
