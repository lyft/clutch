package kinesis

import (
	"context"

	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
	"github.com/lyft/clutch/backend/service/aws"
)

func newKinesisAPI(c aws.Client) kinesisv1.KinesisAPIServer {
	return &kinesisAPI{
		client: c,
	}
}

type kinesisAPI struct {
	client aws.Client
}

func (a *kinesisAPI) GetStream(ctx context.Context, request *kinesisv1.GetStreamRequest) (*kinesisv1.GetStreamResponse, error) {
	stream, err := a.client.DescribeKinesisStream(ctx, request.Region, request.StreamName)
	if err != nil {
		return nil, err
	}

	return &kinesisv1.GetStreamResponse{Stream: stream}, nil
}

func (a *kinesisAPI) UpdateShardCount(ctx context.Context, request *kinesisv1.UpdateShardCountRequest) (*kinesisv1.UpdateShardCountResponse, error) {
	err := a.client.UpdateKinesisShardCount(ctx, request.Region, request.StreamName, request.TargetShardCount)
	if err != nil {
		return nil, err
	}

	return &kinesisv1.UpdateShardCountResponse{}, nil
}
