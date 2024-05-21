package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3control"

	s3v1 "github.com/lyft/clutch/backend/api/aws/s3/v1"
)

func (c *client) S3GetAccessPointPolicy(ctx context.Context, account, region, accessPointName, accountID string) (*s3control.GetAccessPointPolicyOutput, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	in := &s3control.GetAccessPointPolicyInput{
		AccountId: aws.String(accountID),
		Name:      aws.String(accessPointName),
	}

	return cl.s3control.GetAccessPointPolicy(ctx, in)
}

func (c *client) S3GetAccessPoint(ctx context.Context, account, region, accessPointName, accountId string) (*s3v1.AccessPoint, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	in := &s3control.GetAccessPointInput{
		Name:      aws.String(accessPointName),
		AccountId: aws.String(accountId),
	}

	out, err := cl.s3control.GetAccessPoint(ctx, in)
	if err != nil {
		return nil, err
	}

	return &s3v1.AccessPoint{
		Name:            *out.Name,
		AccessPointArn:  *out.AccessPointArn,
		Bucket:          *out.Bucket,
		Alias:           *out.Alias,
		BucketAccountId: *out.BucketAccountId,
		CreationDate:    out.CreationDate.String(),
		Account:         account,
		Region:          region,
	}, nil
}
