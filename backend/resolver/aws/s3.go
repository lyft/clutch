package aws

import (
	"context"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	awsv1resolver "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/resolver"
)

func (r *res) resolveS3BucketForInput(ctx context.Context, input proto.Message) (*resolver.Results, error) {
	switch i := input.(type) {
	case *awsv1resolver.S3BucketName:
		return r.s3BucketResults(ctx, i.Account, i.Region, i.Name, 1)
	default:
		return nil, status.Errorf(codes.Internal, "resolution for type '%T' not implemented", i)
	}
}

func (r *res) resolveS3AccessPointForInput(ctx context.Context, input proto.Message) (*resolver.Results, error) {
	switch i := input.(type) {
	case *awsv1resolver.S3AccessPointName:
		return r.s3AccessPointResults(ctx, i.Account, i.Region, i.Name, 1)
	default:
		return nil, status.Errorf(codes.Internal, "resolution for type '%T' not implemented", i)
	}
}

// Fanout across multiple regions if needed to fetch stream.
func (r *res) s3BucketResults(ctx context.Context, account, region, name string, limit uint32) (*resolver.Results, error) {
	ctx, handler := resolver.NewFanoutHandler(ctx)

	allAccountRegions := r.determineAccountAndRegionsForOption(account, region)

	for account := range allAccountRegions {
		for _, region := range allAccountRegions[account] {
			handler.Add(1)
			go func(account, region string) {
				defer handler.Done()
				bucket, err := r.client.S3DescribeBucket(ctx, account, region, name)
				select {
				case handler.Channel() <- resolver.NewSingleFanoutResult(bucket, err):
					return
				case <-handler.Cancelled():
					return
				}
			}(account, region)
		}
	}

	return handler.Results(limit)
}

// Fanout across multiple regions if needed to fetch stream.
func (r *res) s3AccessPointResults(ctx context.Context, account, region, name string, limit uint32) (*resolver.Results, error) {
	ctx, handler := resolver.NewFanoutHandler(ctx)

	allAccountRegions := r.determineAccountAndRegionsForOption(account, region)
	for account := range allAccountRegions {
		for _, region := range allAccountRegions[account] {
			handler.Add(1)
			go func(account, region string) {
				defer handler.Done()
				callerId, err := r.client.GetCallerIdentity(ctx, account, region)
				if err != nil {
					select {
					case handler.Channel() <- resolver.NewSingleFanoutResult(nil, err):
						return
					case <-handler.Cancelled():
						return
					}
				}
				accessPoint, err := r.client.S3GetAccessPoint(ctx, account, region, name, *callerId.Account)
				select {
				case handler.Channel() <- resolver.NewSingleFanoutResult(accessPoint, err):
					return
				case <-handler.Cancelled():
					return
				}
			}(account, region)
		}
	}

	return handler.Results(limit)
}
