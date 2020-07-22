package aws

import (
	"context"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	"github.com/lyft/clutch/backend/service/aws"
)

func newEC2API(c aws.Client) ec2v1.EC2APIServer {
	return &ec2API{
		client: c,
	}
}

type ec2API struct {
	client aws.Client
}

func (a *ec2API) ResizeAutoscalingGroup(ctx context.Context, request *ec2v1.ResizeAutoscalingGroupRequest) (*ec2v1.ResizeAutoscalingGroupResponse, error) {
	// TODO: additional validation
	err := a.client.ResizeAutoscalingGroup(ctx, request.Region, request.Name, request.Size)
	if err != nil {
		return nil, err
	}
	return &ec2v1.ResizeAutoscalingGroupResponse{}, nil
}

func (a *ec2API) TerminateInstance(ctx context.Context, req *ec2v1.TerminateInstanceRequest) (*ec2v1.TerminateInstanceResponse, error) {
	err := a.client.TerminateInstances(ctx, req.Region, []string{req.InstanceId})
	if err != nil {
		return nil, err
	}

	return &ec2v1.TerminateInstanceResponse{}, nil
}

func (a *ec2API) GetInstance(ctx context.Context, req *ec2v1.GetInstanceRequest) (*ec2v1.GetInstanceResponse, error) {
	instances, err := a.client.DescribeInstances(ctx, req.Region, []string{req.InstanceId})
	if err != nil {
		return nil, err
	}

	return &ec2v1.GetInstanceResponse{Instance: instances[0]}, nil
}
