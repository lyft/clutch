package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func (c *client) SimulateCustomPolicy(ctx context.Context, account, region string, customPolicySimulatorParams *iam.SimulateCustomPolicyInput) (*iam.SimulateCustomPolicyOutput, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	return cl.iam.SimulateCustomPolicy(ctx, customPolicySimulatorParams)
}
