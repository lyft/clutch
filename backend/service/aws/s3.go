package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"io"
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

func (c *client) S3GetBucketPolicy(ctx context.Context, region string, bucket string) (*string, error) {
	rc, err := c.getRegionalClient(region)
	if err != nil {
		return nil, err
	}

	currentAccountInfo, err := rc.sts.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})

	if err != nil {
		return nil, err
	}


	in := &s3.GetBucketPolicyInput{
		ExpectedBucketOwner: currentAccountInfo.Account,
		Bucket:              aws.String(bucket),
	}

	bucketPolicy, err := rc.s3.GetBucketPolicy(ctx, in)


	if err != nil {
		return nil, err
	}

	return bucketPolicy.Policy, nil
}
