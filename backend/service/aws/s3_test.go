package aws

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestS3StreamGet(t *testing.T) {
	s3Client := &mockS3{
		getObjectOutput: &s3.GetObjectOutput{
			Body: ioutil.NopCloser(strings.NewReader("choice")),
		},
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", s3: s3Client}},
	}

	output, err := c.S3StreamingGet(context.Background(), "us-east-1", "clutch", "key")
	assert.NoError(t, err)
	assert.Equal(t, output, ioutil.NopCloser(strings.NewReader("choice")))
}

func TestS3StreamGetErrorHandling(t *testing.T) {
	s3Client := &mockS3{
		getObjectErr: fmt.Errorf("error"),
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", s3: s3Client}},
	}

	output1, err1 := c.S3StreamingGet(context.Background(), "us-east-1", "clutch", "key")
	assert.Nil(t, output1)
	assert.Error(t, err1)

	// Test unknown region
	output2, err2 := c.S3StreamingGet(context.Background(), "choice-region-1", "clutch", "key")
	assert.Nil(t, output2)
	assert.Error(t, err2)
}

type mockS3 struct {
	s3Client

	getObjectErr    error
	getObjectOutput *s3.GetObjectOutput
}

func (m *mockS3) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if m.getObjectErr != nil {
		return nil, m.getObjectErr
	}

	return m.getObjectOutput, nil
}
