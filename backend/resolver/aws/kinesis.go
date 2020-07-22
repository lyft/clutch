package aws

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"

	awsv1resolver "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/resolver"
)

func (r *res) resolveKinesisStreamForInput(ctx context.Context, input proto.Message) (*resolver.Results, error) {
	switch i := input.(type) {
	case *awsv1resolver.KinesisStreamName:
		return r.kinesisResults(ctx, i.Region, i.Name, 1)
	default:
		return nil, fmt.Errorf("unrecognized input type %T", i)
	}
}

// Fanout across multiple regions if needed to fetch stream.
func (r *res) kinesisResults(ctx context.Context, region, id string, limit uint32) (*resolver.Results, error) {
	ctx, handler := resolver.NewFanoutHandler(ctx)

	regions := r.determineRegionsForOption(region)
	for _, region := range regions {
		handler.Add(1)
		go func(region string) {
			defer handler.Done()
			stream, err := r.client.DescribeKinesisStream(ctx, region, id)
			select {
			case handler.Channel() <- resolver.NewSingleFanoutResult(stream, err):
				return
			case <-handler.Cancelled():
				return
			}
		}(region)
	}

	return handler.Results(limit)
}
