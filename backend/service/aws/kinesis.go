package aws

import (
	"context"
	"fmt"
	"math"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"

	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
)

func (c *client) DescribeKinesisStream(ctx context.Context, region string, streamName string) (*kinesisv1.Stream, error) {
	cl, ok := c.clients[region]
	if !ok {
		return nil, fmt.Errorf("no client found for region '%s'", region)
	}

	input := &kinesis.DescribeStreamSummaryInput{StreamName: aws.String(streamName)}
	result, err := cl.kinesis.DescribeStreamSummaryWithContext(ctx, input)

	if err != nil {
		return nil, err
	}

	currentShardCount := aws.Int64Value(result.StreamDescriptionSummary.OpenShardCount)
	if currentShardCount < 0 {
		return nil, fmt.Errorf("AWS returned a negative value for the current shard count")
	}

	var ret = &kinesisv1.Stream{
		StreamName:        aws.StringValue(result.StreamDescriptionSummary.StreamName),
		Region:            region,
		CurrentShardCount: uint32(currentShardCount),
	}

	return ret, nil
}

// we limit the resizing of shard sizes to increments of 25% ranging from [50%, 200%] of the current shard count
// https://docs.aws.amazon.com/kinesis/latest/APIReference/API_UpdateShardCount.html
func getRecommendedShardSizes(currentShardCount uint32) map[uint32]bool {
	sizes := make(map[uint32]bool)
	base := 0.5
	for i := 0; i < 7; i++ {
		recommendedValue := uint32(math.Ceil(float64(currentShardCount) * base))
		sizes[recommendedValue] = true
		base += 0.25
	}
	return sizes
}

func isRecommendedChange(ctx context.Context, cl *regionalClient, streamName string, targetShardCount uint32) (bool, error) {
	input := &kinesis.DescribeStreamSummaryInput{StreamName: aws.String(streamName)}
	result, err := cl.kinesis.DescribeStreamSummaryWithContext(ctx, input)

	if err != nil {
		return false, err
	}
	currentShardCount := aws.Int64Value(result.StreamDescriptionSummary.OpenShardCount)
	if currentShardCount < 0 {
		return false, fmt.Errorf("AWS returned a negative value for the current shard count")
	}

	recommendedValues := getRecommendedShardSizes(uint32(currentShardCount))
	return recommendedValues[targetShardCount], nil
}

func (c *client) UpdateKinesisShardCount(ctx context.Context, region string, streamName string, targetShardCount uint32) error {
	cl, ok := c.clients[region]
	if !ok {
		return fmt.Errorf("no client found for region '%s'", region)
	}

	recommended, err := isRecommendedChange(ctx, cl, streamName, targetShardCount)
	if err != nil {
		return err
	}
	if !recommended {
		return fmt.Errorf("new shard count should be a 25%% increment of current shard count ranging from 50-200%%")
	}

	input := &kinesis.UpdateShardCountInput{
		StreamName:       aws.String(streamName),
		TargetShardCount: aws.Int64(int64(targetShardCount)),
		ScalingType:      aws.String(kinesis.ScalingTypeUniformScaling),
	}
	_, err = cl.kinesis.UpdateShardCountWithContext(ctx, input)

	return err
}
