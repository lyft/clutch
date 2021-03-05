package aws

import (
	"context"
	"math"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
)

func (c *client) DescribeKinesisStream(ctx context.Context, region string, streamName string) (*kinesisv1.Stream, error) {
	cl, err := c.getRegionalClient(region)
	if err != nil {
		return nil, err
	}

	input := &kinesis.DescribeStreamSummaryInput{StreamName: aws.String(streamName)}
	result, err := cl.kinesis.DescribeStreamSummary(ctx, input)
	if err != nil {
		return nil, err
	}

	currentShardCount := aws.ToInt32(result.StreamDescriptionSummary.OpenShardCount)
	if currentShardCount < 0 {
		return nil, status.Error(codes.Internal, "AWS returned a negative value for the current shard count")
	}

	var ret = &kinesisv1.Stream{
		StreamName:        aws.ToString(result.StreamDescriptionSummary.StreamName),
		Region:            region,
		CurrentShardCount: aws.ToInt32(result.StreamDescriptionSummary.OpenShardCount),
	}
	return ret, nil
}

// we limit the resizing of shard sizes to increments of 25% ranging from [50%, 200%] of the current shard count
// https://docs.aws.amazon.com/kinesis/latest/APIReference/API_UpdateShardCount.html
func getRecommendedShardSizes(currentShardCount int32) map[int32]bool {
	sizes := make(map[int32]bool)
	base := 0.5
	for i := 0; i < 7; i++ {
		recommendedValue := int32(math.Ceil(float64(currentShardCount) * base))
		sizes[recommendedValue] = true
		base += 0.25
	}
	return sizes
}

func isRecommendedChange(stream *kinesis.DescribeStreamSummaryOutput, targetShardCount int32) error {
	currentShardCount := aws.ToInt32(stream.StreamDescriptionSummary.OpenShardCount)
	if currentShardCount < 0 {
		return status.Error(codes.Internal, "AWS returned a negative value for the current shard count")
	}

	recommendedValues := getRecommendedShardSizes(currentShardCount)
	if !recommendedValues[targetShardCount] {
		return status.Error(codes.FailedPrecondition, "new shard count should be a 25% increment of current shard count ranging from 50-200%")
	}

	return nil
}

func (c *client) UpdateKinesisShardCount(ctx context.Context, region string, streamName string, targetShardCount int32) error {
	cl, err := c.getRegionalClient(region)
	if err != nil {
		return err
	}

	describeInput := &kinesis.DescribeStreamSummaryInput{StreamName: aws.String(streamName)}
	result, err := cl.kinesis.DescribeStreamSummary(ctx, describeInput)
	if err != nil {
		return err
	}

	err = isRecommendedChange(result, targetShardCount)
	if err != nil {
		return err
	}

	input := &kinesis.UpdateShardCountInput{
		StreamName:       aws.String(streamName),
		TargetShardCount: aws.Int32(targetShardCount),
		ScalingType:      types.ScalingTypeUniformScaling,
	}

	_, err = cl.kinesis.UpdateShardCount(ctx, input)
	return err
}
