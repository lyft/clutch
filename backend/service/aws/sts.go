package aws

import (
"context"
"github.com/aws/aws-sdk-go-v2/service/sts"
)

func (c *client) GetCallerIdentity(ctx context.Context, region string) (*sts.GetCallerIdentityOutput, error) {
	rc, err := c.getRegionalClient(region)
	if err != nil {
		return nil, err
	}

	in := &sts.GetCallerIdentityInput{}


	out, err := rc.sts.GetCallerIdentity(ctx, in)
	if err != nil {
		return nil, err
	}

	return out, nil
}
