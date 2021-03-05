package aws

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"

	kinesisv1api "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
	awsv1resolver "github.com/lyft/clutch/backend/api/resolver/aws/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
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

	patternValues, err := meta.ExtractPatternValuesFromString((*kinesisv1api.Stream)(nil), id)
	if err != nil && len(patternValues["region"]) > 0 && len(patternValues["stream_name"]) > 0 {
		handler.Add(1)
		go func(region, id string) {
			defer handler.Done()
			stream, err := r.client.DescribeKinesisStream(ctx, region, id)
			select {
			case handler.Channel() <- resolver.NewSingleFanoutResult(stream, err):
				return
			case <-handler.Cancelled():
				return
			}
		}(patternValues["region"], patternValues["stream_name"])
	} else {
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
	}

	return handler.Results(limit)
}
