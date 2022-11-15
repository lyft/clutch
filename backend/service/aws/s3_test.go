package aws

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
)

func TestS3StreamGet(t *testing.T) {
	s3Client := &mockS3{
		getObjectOutput: &s3.GetObjectOutput{
			Body: io.NopCloser(strings.NewReader("choice")),
		},
	}
	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", s3: s3Client},
				},
			},
		},
	}

	output, err := c.S3StreamingGet(context.Background(), "default", "us-east-1", "clutch", "key")
	assert.NoError(t, err)
	assert.Equal(t, output, io.NopCloser(strings.NewReader("choice")))
}

func TestS3StreamGetErrorHandling(t *testing.T) {
	s3Client := &mockS3{
		getObjectErr: fmt.Errorf("error"),
	}
	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", s3: s3Client},
				},
			},
		},
	}

	output1, err1 := c.S3StreamingGet(context.Background(), "default", "us-east-1", "clutch", "key")
	assert.Nil(t, output1)
	assert.Error(t, err1)

	// Test unknown region
	output2, err2 := c.S3StreamingGet(context.Background(), "default", "choice-region-1", "clutch", "key")
	assert.Nil(t, output2)
	assert.Error(t, err2)
}

func TestS3GetBucketPolicy(t *testing.T) {
	s3Client := &mockS3{
		getObjectPolicyOutput: &s3.GetBucketPolicyOutput{
			Policy:         aws.String("{}"),
			ResultMetadata: middleware.Metadata{},
		},
	}
	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", s3: s3Client},
				},
			},
		},
	}

	output, err := c.S3GetBucketPolicy(context.Background(), "default", "us-east-1", "clutch", "000000000000")
	assert.NoError(t, err)
	assert.Equal(t, output, &s3.GetBucketPolicyOutput{
		Policy:         aws.String("{}"),
		ResultMetadata: middleware.Metadata{},
	})
}

func TestS3GetBucketPolicyErrorHandling(t *testing.T) {
	s3Client := &mockS3{
		getObjectPolicyErr: fmt.Errorf("error"),
	}
	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", s3: s3Client},
				},
			},
		},
	}

	output1, err1 := c.S3GetBucketPolicy(context.Background(), "default", "us-east-1", "clutch", "000000000000")
	assert.Nil(t, output1)
	assert.Error(t, err1)

	// Test unknown region
	output2, err2 := c.S3GetBucketPolicy(context.Background(), "default", "choice-region-1", "clutch", "000000000000")
	assert.Nil(t, output2)
	assert.Error(t, err2)
}

type mockS3 struct {
	s3Client

	getObjectErr    error
	getObjectOutput *s3.GetObjectOutput

	getObjectPolicyErr    error
	getObjectPolicyOutput *s3.GetBucketPolicyOutput
}

func (m *mockS3) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if m.getObjectErr != nil {
		return nil, m.getObjectErr
	}

	return m.getObjectOutput, nil
}

func (m *mockS3) GetBucketPolicy(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
	if m.getObjectPolicyErr != nil {
		return nil, m.getObjectPolicyErr
	}

	return m.getObjectPolicyOutput, nil
}
