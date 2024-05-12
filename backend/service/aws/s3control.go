package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
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
