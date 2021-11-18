package aws

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *client) S3GetBucketPolicy(ctx context.Context, account, region, bucket, accountID string) (*s3.GetBucketPolicyOutput, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	in := &s3.GetBucketPolicyInput{
		ExpectedBucketOwner: aws.String(accountID),
		Bucket:              aws.String(bucket),
	}

	return cl.s3.GetBucketPolicy(ctx, in)
}

func (c *client) S3StreamingGet(ctx context.Context, account, region, bucket, key string) (io.ReadCloser, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	in := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	out, err := cl.s3.GetObject(ctx, in)
	if err != nil {
		return nil, err
	}

	return out.Body, nil
}
