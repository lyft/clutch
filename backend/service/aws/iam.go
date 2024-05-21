package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"

	iamv1 "github.com/lyft/clutch/backend/api/aws/iam/v1"
)

func (c *client) SimulateCustomPolicy(ctx context.Context, account, region string, customPolicySimulatorParams *iam.SimulateCustomPolicyInput) (*iam.SimulateCustomPolicyOutput, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	return cl.iam.SimulateCustomPolicy(ctx, customPolicySimulatorParams)
}

func (c *client) GetIAMRole(
	ctx context.Context,
	account,
	region,
	roleName string,
) (*iamv1.Role, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	role, err := cl.iam.GetRole(ctx, &iam.GetRoleInput{
		RoleName: &roleName,
	})
	if err != nil {
		return nil, err
	}

	return &iamv1.Role{
		Arn:         *role.Role.Arn,
		Name:        *role.Role.RoleName,
		CreatedDate: role.Role.CreateDate.String(),
		Id:          *role.Role.RoleId,
		Account:     account,
		Region:      region,
	}, nil
}
