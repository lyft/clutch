package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func (c *client) GetCallerIdentity(ctx context.Context, account, region string) (*sts.GetCallerIdentityOutput, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	in := &sts.GetCallerIdentityInput{}
	return cl.sts.GetCallerIdentity(ctx, in)
}
