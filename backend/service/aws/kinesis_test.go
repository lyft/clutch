package aws

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/stretchr/testify/assert"

	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
)

var testAwsStream = &kinesis.StreamDescriptionSummary{
	EnhancedMonitoring:      []*kinesis.EnhancedMetrics{},
	OpenShardCount:          aws.Int64(100),
	RetentionPeriodHours:    aws.Int64(24),
	StreamARN:               aws.String("test-arn"),
	StreamCreationTimestamp: aws.Time(time.Unix(1449952498, 0)),
	StreamName:              aws.String("test-stream"),
	StreamStatus:            aws.String("ACTIVE"),
}

var testAwsStreamWithBadData = &kinesis.StreamDescriptionSummary{
	EnhancedMonitoring:      []*kinesis.EnhancedMetrics{},
	OpenShardCount:          aws.Int64(-100),
	RetentionPeriodHours:    aws.Int64(24),
	StreamARN:               aws.String("test-arn"),
	StreamCreationTimestamp: aws.Time(time.Unix(1449952498, 0)),
	StreamName:              aws.String("test-stream"),
	StreamStatus:            aws.String("ACTIVE"),
}

var testStreamOutPut = &kinesisv1.Stream{
	StreamName:        "test-stream",
	Region:            "us-east-1",
	CurrentShardCount: 100,
}

func TestDescribeStreamWithGoodData(t *testing.T) {
	m := &mockKinesis{
		stream: testAwsStream,
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", kinesis: m}},
	}

	result, err := c.DescribeKinesisStream(context.Background(), "us-east-1", "test-stream")
	assert.NoError(t, err)
	assert.Equal(t, testStreamOutPut, result)

	m.streamErr = errors.New("error")
	_, err1 := c.DescribeKinesisStream(context.Background(), "us-east-1", "test-stream")
	assert.EqualError(t, err1, "error")
}

func TestDescribeStreamWithBadData(t *testing.T) {
	m := &mockKinesis{
		stream: testAwsStreamWithBadData,
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", kinesis: m}},
	}

	_, err := c.DescribeKinesisStream(context.Background(), "us-east-1", "test-stream")
	assert.EqualError(t, err, "AWS returned a negative value for the current shard count")
}

func TestUpdateShardCount(t *testing.T) {
	m := &mockKinesis{
		stream: testAwsStream,
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", kinesis: m}},
	}

	err := c.UpdateKinesisShardCount(context.Background(), "us-east-1", "test-stream", 50)
	assert.NoError(t, err)

	m.updateErr = errors.New("error")
	err1 := c.UpdateKinesisShardCount(context.Background(), "us-east-1", "test-stream", 50)
	assert.EqualError(t, err1, "error")

	err2 := c.UpdateKinesisShardCount(context.Background(), "us-east-1", "test-stream", 10)
	assert.EqualError(t, err2, "new shard count should be a 25% increment of current shard count ranging from 50-200%")
}

func TestGetRecommendedShardSizes(t *testing.T) {
	m1 := make(map[uint32]bool)
	m1[50] = true
	m1[75] = true
	m1[100] = true
	m1[125] = true
	m1[150] = true
	m1[175] = true
	m1[200] = true
	assert.Equal(t, m1, getRecommendedShardSizes(100))
	m2 := make(map[uint32]bool)
	m2[5] = true
	m2[8] = true
	m2[10] = true
	m2[13] = true
	m2[15] = true
	m2[18] = true
	m2[20] = true
	assert.Equal(t, m2, getRecommendedShardSizes(10))
	m3 := make(map[uint32]bool)
	m3[0] = true
	assert.Equal(t, m3, getRecommendedShardSizes(0))
	m4 := make(map[uint32]bool)
	m4[1] = true
	m4[2] = true
	assert.Equal(t, m4, getRecommendedShardSizes(1))
	m5 := make(map[uint32]bool)
	m5[1] = true
	m5[2] = true
	m5[3] = true
	m5[4] = true
	assert.Equal(t, m5, getRecommendedShardSizes(2))
}

func TestIsRecommendedChangeWithGoodData(t *testing.T) {
	m := &mockKinesis{
		stream: testAwsStream,
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", kinesis: m}},
	}
	bool1, err1 := isRecommendedChange(context.Background(), c.clients["us-east-1"], "test-stream", uint32(100))
	assert.True(t, bool1)
	assert.NoError(t, err1)

	bool2, err2 := isRecommendedChange(context.Background(), c.clients["us-east-1"], "test-stream", uint32(10))
	assert.False(t, bool2)
	assert.NoError(t, err2)
}

func TestIsRecommendedChangeWithBadData(t *testing.T) {
	m := &mockKinesis{
		stream: testAwsStreamWithBadData,
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", kinesis: m}},
	}
	_, err1 := isRecommendedChange(context.Background(), c.clients["us-east-1"], "test-stream", uint32(100))
	assert.EqualError(t, err1, "AWS returned a negative value for the current shard count")
}

type mockKinesis struct {
	kinesisiface.KinesisAPI

	streamErr error
	stream    *kinesis.StreamDescriptionSummary

	updateErr error
	update    *kinesis.UpdateShardCountOutput
}

func (m *mockKinesis) DescribeStreamSummaryWithContext(ctx context.Context, input *kinesis.DescribeStreamSummaryInput, opts ...request.Option) (*kinesis.DescribeStreamSummaryOutput, error) {
	if m.streamErr != nil {
		return nil, m.streamErr
	}

	ret := &kinesis.DescribeStreamSummaryOutput{
		StreamDescriptionSummary: m.stream,
	}

	return ret, nil
}

func (m *mockKinesis) UpdateShardCountWithContext(ctx context.Context, input *kinesis.UpdateShardCountInput, opts ...request.Option) (*kinesis.UpdateShardCountOutput, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	ret := m.update

	return ret, nil
}
