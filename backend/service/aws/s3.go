package aws

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *client) S3StreamingGet(ctx context.Context, region string, bucket string, key string) (io.ReadCloser, error) {
	rc, err := c.getRegionalClient(region)
	if err != nil {
		return nil, err
	}

	in := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	out, err := rc.s3.GetObject(ctx, in)
	if err != nil {
		return nil, err
	}

	return out.Body, nil
}
