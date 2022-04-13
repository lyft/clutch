package awsmock

import (
	"context"
	"fmt"
	"io"
	"math/rand"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go/middleware"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
	"github.com/lyft/clutch/backend/service"
	clutchawsclient "github.com/lyft/clutch/backend/service/aws"
)

type svc struct{}

func (s *svc) GetDirectClient(account string, region string) (clutchawsclient.DirectClient, error) {
	panic("implement me")
}

func (s *svc) DescribeKinesisStream(ctx context.Context, account, region, streamName string) (*kinesisv1.Stream, error) {
	ret := &kinesisv1.Stream{
		StreamName:        streamName,
		Region:            region,
		Account:           "default",
		CurrentShardCount: int32(32),
	}
	return ret, nil
}

func (s *svc) UpdateKinesisShardCount(ctx context.Context, account, region, streamName string, targetShardCount int32) error {
	return nil
}

func (s *svc) ResizeAutoscalingGroup(ctx context.Context, account, region, name string, size *ec2v1.AutoscalingGroupSize) error {
	return nil
}

func (s *svc) DescribeAutoscalingGroups(ctx context.Context, account, region string, names []string) ([]*ec2v1.AutoscalingGroup, error) {
	ret := make([]*ec2v1.AutoscalingGroup, len(names))
	for idx, name := range names {
		ret[idx] = &ec2v1.AutoscalingGroup{
			Name:    name,
			Region:  region,
			Account: "default",
			Zones:   []string{region + "a", region + "b"},
			Size: &ec2v1.AutoscalingGroupSize{
				Min:     25,
				Max:     200,
				Desired: 100,
			},
			TerminationPolicies: []ec2v1.AutoscalingGroup_TerminationPolicy{
				ec2v1.AutoscalingGroup_OLDEST_INSTANCE,
			},
			Instances: []*ec2v1.AutoscalingGroup_Instance{
				{
					Id:                      "i-12345892010",
					Zone:                    region + "b",
					LaunchConfigurationName: "myLaunchConfig",
					Healthy:                 true,
					LifecycleState:          ec2v1.AutoscalingGroup_Instance_IN_SERVICE,
				},
				{
					Id:                      "i-22345892010",
					Zone:                    region + "a",
					LaunchConfigurationName: "myLaunchConfig",
					Healthy:                 false,
					LifecycleState:          ec2v1.AutoscalingGroup_Instance_PENDING,
				},
			},
		}
	}

	return ret, nil
}

func (s *svc) DescribeInstances(ctx context.Context, account, region string, ids []string) ([]*ec2v1.Instance, error) {
	var ret []*ec2v1.Instance
	for _, id := range ids {
		ret = append(ret, &ec2v1.Instance{
			InstanceId:       id,
			Region:           region,
			Account:          "default",
			State:            ec2v1.Instance_State(rand.Intn(len(ec2v1.Instance_State_value))),
			InstanceType:     "c3.2xlarge",
			PublicIpAddress:  "8.1.1.1",
			PrivateIpAddress: "10.1.1.1",
			AvailabilityZone: fmt.Sprintf("%sz", region),
			Tags:             map[string]string{"Name": fmt.Sprintf("my-instance-%s", id)},
		})
	}
	return ret, nil
}

func (s *svc) TerminateInstances(ctx context.Context, account, region string, ids []string) error {
	return nil
}

func (s *svc) RebootInstances(ctx context.Context, account, region string, ids []string) error {
	return nil
}

func (s *svc) S3StreamingGet(ctx context.Context, account, region, bucket, key string) (io.ReadCloser, error) {
	panic("implement me")
}

func (s *svc) S3GetBucketPolicy(ctx context.Context, account, region, bucket, accountID string) (*s3.GetBucketPolicyOutput, error) {
	return &s3.GetBucketPolicyOutput{
		Policy:         aws.String("{}"),
		ResultMetadata: middleware.Metadata{},
	}, nil
}

func (s *svc) DescribeTable(ctx context.Context, account, region, tableName string) (*dynamodbv1.Table, error) {
	currentThroughput := &dynamodbv1.Throughput{
		ReadCapacityUnits:  100,
		WriteCapacityUnits: 200,
	}
	gsis := []*dynamodbv1.GlobalSecondaryIndex{
		{
			Name: "test-gsi",
			ProvisionedThroughput: &dynamodbv1.Throughput{
				ReadCapacityUnits:  10,
				WriteCapacityUnits: 20,
			},
			Status: dynamodbv1.GlobalSecondaryIndex_Status(5),
		},
	}

	ret := &dynamodbv1.Table{
		Name:                   tableName,
		Region:                 region,
		ProvisionedThroughput:  currentThroughput,
		Status:                 dynamodbv1.Table_Status(5),
		GlobalSecondaryIndexes: gsis,
	}
	return ret, nil
}

func (s *svc) UpdateCapacity(ctx context.Context, account, region, tableName string, targetTableCapacity *dynamodbv1.Throughput, indexUpdates []*dynamodbv1.IndexUpdateAction, ignoreMaximums bool) (*dynamodbv1.Table, error) {
	currentThroughput := &dynamodbv1.Throughput{
		ReadCapacityUnits:  100,
		WriteCapacityUnits: 200,
	}
	gsis := []*dynamodbv1.GlobalSecondaryIndex{
		{
			Name: "test-gsi",
			ProvisionedThroughput: &dynamodbv1.Throughput{
				ReadCapacityUnits:  10,
				WriteCapacityUnits: 20,
			},
			Status: dynamodbv1.GlobalSecondaryIndex_Status(3),
		},
	}

	ret := &dynamodbv1.Table{
		Name:                   tableName,
		Region:                 region,
		ProvisionedThroughput:  currentThroughput,
		Status:                 dynamodbv1.Table_Status(3),
		GlobalSecondaryIndexes: gsis,
	}
	return ret, nil
}

var accountsAndRegions = map[string][]string{
	"default":    {"us-mock-1"},
	"staging":    {"us-mock-1", "us-mock-2"},
	"production": {"us-mock-1", "us-mock-2", "us-mock-3"},
}

func (s *svc) Regions() []string {
	uniqueRegions := map[string]bool{}

	for _, regions := range accountsAndRegions {
		for _, r := range regions {
			uniqueRegions[r] = true
		}
	}

	regions := make([]string, len(uniqueRegions))
	i := 0
	for region := range uniqueRegions {
		regions[i] = region
		i++
	}

	return regions
}

func (s *svc) Accounts() []string {
	accounts := []string{}
	for account := range accountsAndRegions {
		accounts = append(accounts, account)
	}
	return accounts
}

func (s *svc) AccountsAndRegions() map[string][]string {
	return accountsAndRegions
}

func (s *svc) GetAccountsInRegion(region string) []string {
	accounts := []string{}
	for a, regions := range accountsAndRegions {
		for _, r := range regions {
			if r == region {
				accounts = append(accounts, a)
			}
		}
	}

	return accounts
}

func (s *svc) GetPrimaryAccountAlias() string {
	return "default"
}

func (s *svc) GetCallerIdentity(ctx context.Context, account, region string) (*sts.GetCallerIdentityOutput, error) {
	return &sts.GetCallerIdentityOutput{
		Account: aws.String("000000000000"),
		Arn:     aws.String("arn:aws:sts::000000000000:fake-role"),
		UserId:  aws.String("some_special_id"),
	}, nil
}

func (s *svc) SimulateCustomPolicy(ctx context.Context, account, region string, customPolicySimulatorParams *iam.SimulateCustomPolicyInput) (*iam.SimulateCustomPolicyOutput, error) {
	return &iam.SimulateCustomPolicyOutput{
		EvaluationResults: []types.EvaluationResult{
			{EvalActionName: aws.String("s3:GetBucketPolicy"),
				EvalDecision:         types.PolicyEvaluationDecisionTypeImplicitDeny,
				EvalResourceName:     aws.String("arn:aws:s3:::a/*"),
				MissingContextValues: []string{},
				MatchedStatements:    []types.Statement{},
			},
		},
	}, nil
}

func New() clutchawsclient.Client {
	return &svc{}
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}
