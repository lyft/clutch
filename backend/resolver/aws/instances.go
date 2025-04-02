package aws

import (
	"context"
	"fmt"
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	awsv1 "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/resolver"
)

var instanceIDPattern = regexp.MustCompile("[a-fA-F0-9]{17}")

func normalizeInstanceID(input string) (string, error) {
	instanceID := instanceIDPattern.FindString(input)
	if instanceID == "" {
		return "", status.Error(codes.InvalidArgument, "the expected format for instance IDs is i-0123456789abcdef0")
	}

	return fmt.Sprintf("i-%s", instanceID), nil
}

// Take any inputs that can get instance IDs and normalize them for the client.
func (r *res) resolveInstancesForInput(ctx context.Context, input proto.Message) (*resolver.Results, error) {
	switch i := input.(type) {
	case *awsv1.InstanceID:
		return r.instanceResults(ctx, i.Account, i.Region, []string{i.Id}, 1)
	default:
		return nil, status.Errorf(codes.Internal, "unrecognized input type '%T'", i)
	}
}

// Fanout across multiple regions if needed to fetch instances.
func (r *res) instanceResults(ctx context.Context, account, region string, ids []string, limit uint32) (*resolver.Results, error) {
	ctx, handler := resolver.NewFanoutHandler(ctx)

	allAccountRegions := r.determineAccountAndRegionsForOption(account, region)
	for account := range allAccountRegions {
		for _, region := range allAccountRegions[account] {
			handler.Add(1)
			go func(account, region string) {
				defer handler.Done()
				instances, err := r.client.DescribeInstances(ctx, account, region, ids)
				select {
				case handler.Channel() <- resolver.NewFanoutResult(instances, err):
					return
				case <-handler.Cancelled():
					return
				}
			}(account, region)
		}
	}

	return handler.Results(limit)
}
