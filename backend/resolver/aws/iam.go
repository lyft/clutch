package aws

import (
	"context"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	awsv1resolver "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/resolver"
)

func (r *res) resolveIAMRoleForInput(ctx context.Context, input proto.Message) (*resolver.Results, error) {
	switch i := input.(type) {
	case *awsv1resolver.IAMRoleName:
		return r.iamRoleResults(ctx, i.Account, i.Region, i.Name, 1)
	default:
		return nil, status.Errorf(codes.Internal, "resolution for type '%T' not implemented", i)
	}
}

// Fanout across multiple regions if needed to fetch stream.
func (r *res) iamRoleResults(ctx context.Context, account, region, name string, limit uint32) (*resolver.Results, error) {
	ctx, handler := resolver.NewFanoutHandler(ctx)

	allAccountRegions := r.determineAccountAndRegionsForOption(account, region)
	for account := range allAccountRegions {
		for _, region := range allAccountRegions[account] {
			handler.Add(1)
			go func(account, region string) {
				defer handler.Done()
				role, err := r.client.GetIAMRole(ctx, account, region, name)
				select {
				case handler.Channel() <- resolver.NewSingleFanoutResult(role, err):
					return
				case <-handler.Cancelled():
					return
				}
			}(account, region)
		}
	}

	return handler.Results(limit)
}
