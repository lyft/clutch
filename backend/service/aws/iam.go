package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func (c *client) SimulateCustomPolicy(ctx context.Context, region string, customPolicySimulatorParams *iam.SimulateCustomPolicyInput) (*iam.SimulateCustomPolicyOutput, error) {
	rc, err := c.getRegionalClient(region)
	if err != nil {
		return nil, err
	}

	return rc.iam.SimulateCustomPolicy(ctx, customPolicySimulatorParams)
}
