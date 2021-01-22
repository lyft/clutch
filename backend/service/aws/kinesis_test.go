package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/stretchr/testify/assert"

	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
)

var testStreamOutPut = &kinesisv1.Stream{
	StreamName:        "test-stream",
	Region:            "us-east-1",
	CurrentShardCount: 100,
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
	err1 := isRecommendedChange(&kinesis.DescribeStreamSummaryOutput{}, int32(100))
	assert.NoError(t, err1)

	err2 := isRecommendedChange(&kinesis.DescribeStreamSummaryOutput{}, int32(10))
	assert.NoError(t, err2)
}

func TestIsRecommendedChangeWithBadData(t *testing.T) {
	err1 := isRecommendedChange(&kinesis.DescribeStreamSummaryOutput{}, int32(100))
	assert.EqualError(t, err1, "AWS returned a negative value for the current shard count")
}
