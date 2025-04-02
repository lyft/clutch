package aws

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	awsv1resolver "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/resolver"
)

func (r *res) resolveKinesisStreamForInput(ctx context.Context, input proto.Message) (*resolver.Results, error) {
	switch i := input.(type) {
	case *awsv1resolver.KinesisStreamName:
		return r.kinesisResults(ctx, i.Account, i.Region, i.Name, 1)
	default:
		return nil, status.Errorf(codes.Internal, "resolution for type '%T' not implemented", i)
	}
}

// Fanout across multiple regions if needed to fetch stream.
func (r *res) kinesisResults(ctx context.Context, account, region, id string, limit uint32) (*resolver.Results, error) {
	ctx, handler := resolver.NewFanoutHandler(ctx)

	allAccountRegions := r.determineAccountAndRegionsForOption(account, region)
	for account := range allAccountRegions {
		for _, region := range allAccountRegions[account] {
			handler.Add(1)
			go func(account, region string) {
				defer handler.Done()
				stream, err := r.client.DescribeKinesisStream(ctx, account, region, id)
				select {
				case handler.Channel() <- resolver.NewSingleFanoutResult(stream, err):
					return
				case <-handler.Cancelled():
					return
				}
			}(account, region)
		}
	}

	return handler.Results(limit)
}
