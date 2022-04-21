package aws

// <!-- START clutchdoc -->
// description: Multi-region client for Amazon Web Services.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	astypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/iancoleman/strcase"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
	awsv1 "github.com/lyft/clutch/backend/api/config/service/aws/v1"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/service"
)

const (
	Name = "clutch.service.aws"
)

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	ac := &awsv1.Config{}
	err := cfg.UnmarshalTo(ac)
	if err != nil {
		return nil, err
	}

	// aws_config_profile_name is not currently implemented
	// if this is set will error out to let the user know what they are trying to do will not work
	if ac.AwsConfigProfileName != "" {
		return nil, errors.New("AWS config field [aws_config_profile_name] is not implemented")
	}

	accountAlias := ac.PrimaryAccountAliasDisplayName
	if ac.PrimaryAccountAliasDisplayName == "" {
		accountAlias = "default"
	}

	c := &client{
		accounts:            make(map[string]*accountClients),
		topologyObjectChan:  make(chan *topologyv1.UpdateCacheRequest, topologyObjectChanBufferSize),
		topologyLock:        semaphore.NewWeighted(1),
		currentAccountAlias: accountAlias,
		log:                 logger,
		scope:               scope,
	}

	clientRetries := 0
	if ac.ClientConfig != nil && ac.ClientConfig.Retries >= 0 {
		clientRetries = int(ac.ClientConfig.Retries)
	}

	ds := getScalingLimits(ac)
	awsHTTPClient := &http.Client{}
	awsClientCommonOptions := []func(*config.LoadOptions) error{
		config.WithHTTPClient(awsHTTPClient),
		config.WithRetryer(func() aws.Retryer {
			customRetryer := retry.NewStandard(func(so *retry.StandardOptions) {
				so.MaxAttempts = clientRetries
			})
			return customRetryer
		}),
	}

	for _, region := range ac.Regions {
		regionCfg, err := config.LoadDefaultConfig(context.TODO(),
			append(awsClientCommonOptions, config.WithRegion(region))...,
		)
		if err != nil {
			return nil, err
		}

		c.createRegionalClients(c.currentAccountAlias, region, ac.Regions, ds, regionCfg)
	}

	if err := c.configureAdditionalAccountClient(ac.AdditionalAccounts, ds, awsHTTPClient, awsClientCommonOptions); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) configureAdditionalAccountClient(accounts []*awsv1.AWSAccount, ds *awsv1.ScalingLimits, awsHTTPClient *http.Client, awsClientOptions []func(*config.LoadOptions) error) error {
	for _, account := range accounts {
		accountRoleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", account.AccountNumber, account.IamRole)
		// For doing STS calls it does not matter which region client we are using, as they are not bounded by region
		// we choose just the first region client
		stsClient := c.accounts[c.currentAccountAlias].clients[c.accounts[c.currentAccountAlias].regions[0]].sts
		assumeRoleProvider := stscreds.NewAssumeRoleProvider(stsClient, accountRoleARN)
		credsCache := aws.NewCredentialsCache(assumeRoleProvider)

		for _, region := range account.Regions {
			regionCfg, err := config.LoadDefaultConfig(context.TODO(),
				append(awsClientOptions, config.WithRegion(region))...,
			)
			if err != nil {
				return err
			}

			regionCfg.Credentials = credsCache

			c.createRegionalClients(account.Alias, region, account.Regions, ds, regionCfg)
		}
	}

	return nil
}

func (c *client) createRegionalClients(accountAlias, region string, regions []string, ds *awsv1.ScalingLimits, regionCfg aws.Config) {
	if _, ok := c.accounts[accountAlias]; !ok {
		c.accounts[accountAlias] = &accountClients{
			alias:   accountAlias,
			regions: regions,
			clients: map[string]*regionalClient{},
		}
	}

	c.accounts[accountAlias].clients[region] = &regionalClient{
		region: region,
		dynamodbCfg: &awsv1.DynamodbConfig{
			ScalingLimits: &awsv1.ScalingLimits{
				MaxReadCapacityUnits:  ds.MaxReadCapacityUnits,
				MaxWriteCapacityUnits: ds.MaxWriteCapacityUnits,
				MaxScaleFactor:        ds.MaxScaleFactor,
				EnableOverride:        ds.EnableOverride,
			},
		},

		s3:          s3.NewFromConfig(regionCfg),
		kinesis:     kinesis.NewFromConfig(regionCfg),
		ec2:         ec2.NewFromConfig(regionCfg),
		autoscaling: autoscaling.NewFromConfig(regionCfg),
		dynamodb:    dynamodb.NewFromConfig(regionCfg),
		sts:         sts.NewFromConfig(regionCfg),
		iam:         iam.NewFromConfig(regionCfg),
	}
}

type Client interface {
	DescribeInstances(ctx context.Context, account, region string, ids []string) ([]*ec2v1.Instance, error)
	TerminateInstances(ctx context.Context, account, region string, ids []string) error
	RebootInstances(ctx context.Context, account, region string, ids []string) error

	DescribeAutoscalingGroups(ctx context.Context, account, region string, names []string) ([]*ec2v1.AutoscalingGroup, error)
	ResizeAutoscalingGroup(ctx context.Context, account, region string, name string, size *ec2v1.AutoscalingGroupSize) error

	DescribeKinesisStream(ctx context.Context, account, region, streamName string) (*kinesisv1.Stream, error)
	UpdateKinesisShardCount(ctx context.Context, account, region, streamName string, targetShardCount int32) error

	S3GetBucketPolicy(ctx context.Context, account, region, bucket, accountID string) (*s3.GetBucketPolicyOutput, error)
	S3StreamingGet(ctx context.Context, account, region, bucket, key string) (io.ReadCloser, error)

	DescribeTable(ctx context.Context, account, region, tableName string) (*dynamodbv1.Table, error)
	UpdateCapacity(ctx context.Context, account, region, tableName string, targetTableCapacity *dynamodbv1.Throughput, indexUpdates []*dynamodbv1.IndexUpdateAction, ignoreMaximums bool) (*dynamodbv1.Table, error)

	GetCallerIdentity(ctx context.Context, account, region string) (*sts.GetCallerIdentityOutput, error)

	SimulateCustomPolicy(ctx context.Context, account, region string, customPolicySimulatorParams *iam.SimulateCustomPolicyInput) (*iam.SimulateCustomPolicyOutput, error)

	Accounts() []string
	AccountsAndRegions() map[string][]string
	GetAccountsInRegion(region string) []string
	GetPrimaryAccountAlias() string
	Regions() []string

	GetDirectClient(account string, region string) (DirectClient, error)
}

// DirectClient gives access to the underlying AWS clients from the Golang SDK.
// This allows arbitrary feature development on top of AWS from other services and modules without having to
// contribute to the upstream interface. Using these clients will make mocking extremely difficult since it returns the
// AWS SDK's struct types and not an interface that can be substituted for. It is recommended following initial
// development of a feature that you add the calls to a service interface so they can be tested more easily.
type DirectClient interface {
	Autoscaling() *autoscaling.Client
	DynamoDB() *dynamodb.Client
	EC2() *ec2.Client
	IAM() *iam.Client
	Kinesis() *kinesis.Client
	S3() *s3.Client
	STS() *sts.Client
}

type client struct {
	accounts            map[string]*accountClients
	topologyObjectChan  chan *topologyv1.UpdateCacheRequest
	topologyLock        *semaphore.Weighted
	currentAccountAlias string
	log                 *zap.Logger
	scope               tally.Scope
}

type regionalClient struct {
	region string

	dynamodbCfg *awsv1.DynamodbConfig

	autoscaling autoscalingClient
	dynamodb    dynamodbClient
	ec2         ec2Client
	iam         iamClient
	kinesis     kinesisClient
	s3          s3Client
	sts         stsClient
}

func (r *regionalClient) Autoscaling() *autoscaling.Client {
	return r.autoscaling.(*autoscaling.Client)
}

func (r *regionalClient) DynamoDB() *dynamodb.Client {
	return r.dynamodb.(*dynamodb.Client)
}

func (r *regionalClient) EC2() *ec2.Client {
	return r.ec2.(*ec2.Client)
}

func (r *regionalClient) IAM() *iam.Client {
	return r.iam.(*iam.Client)
}

func (r *regionalClient) Kinesis() *kinesis.Client {
	return r.kinesis.(*kinesis.Client)
}

func (r *regionalClient) S3() *s3.Client {
	return r.s3.(*s3.Client)
}

func (r *regionalClient) STS() *sts.Client {
	return r.sts.(*sts.Client)
}

type accountClients struct {
	alias   string
	regions []string

	clients map[string]*regionalClient
}

func (c *client) GetDirectClient(account string, region string) (DirectClient, error) {
	return c.getAccountRegionClient(account, region)
}

// Implement the interface provided by errorintercept, so errors are caught at middleware and converted to gRPC status.
func (c *client) InterceptError(e error) error {
	return ConvertError(e)
}

func (c *client) getAccountRegionClient(account, region string) (*regionalClient, error) {
	accountClients, ok := c.accounts[account]
	if !ok || accountClients == nil {
		return nil, status.Errorf(codes.NotFound, "account %s not found", account)
	}
	cl, ok := accountClients.clients[region]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "no client found for account '%s' in region '%s'", account, region)
	}
	return cl, nil
}

func (c *client) ResizeAutoscalingGroup(ctx context.Context, account, region, name string, size *ec2v1.AutoscalingGroupSize) error {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return err
	}

	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		DesiredCapacity:      aws.Int32(int32(size.Desired)),
		MaxSize:              aws.Int32(int32(size.Max)),
		MinSize:              aws.Int32(int32(size.Min)),
	}

	_, err = cl.autoscaling.UpdateAutoScalingGroup(ctx, input)
	return err
}

func (c *client) DescribeAutoscalingGroups(ctx context.Context, account, region string, names []string) ([]*ec2v1.AutoscalingGroup, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: names,
	}
	result, err := cl.autoscaling.DescribeAutoScalingGroups(ctx, input)
	if err != nil {
		return nil, err
	}

	ret := make([]*ec2v1.AutoscalingGroup, len(result.AutoScalingGroups))
	for idx, group := range result.AutoScalingGroups {
		ret[idx] = newProtoForAutoscalingGroup(account, group)
	}

	return ret, nil
}

// Shave off the trailing zone identifier to get the region
func zoneToRegion(zone string) string {
	if zone == "" {
		return "UNKNOWN"
	}
	return zone[:len(zone)-1]
}

func protoForTerminationPolicy(policy string) ec2v1.AutoscalingGroup_TerminationPolicy {
	policy = strcase.ToScreamingSnake(policy)
	val, ok := ec2v1.AutoscalingGroup_TerminationPolicy_value[policy]
	if !ok {
		return ec2v1.AutoscalingGroup_UNKNOWN
	}
	return ec2v1.AutoscalingGroup_TerminationPolicy(val)
}

func protoForAutoscalingGroupInstanceLifecycleState(state string) ec2v1.AutoscalingGroup_Instance_LifecycleState {
	state = strcase.ToScreamingSnake(strings.ReplaceAll(state, ":", ""))
	val, ok := ec2v1.AutoscalingGroup_Instance_LifecycleState_value[state]
	if !ok {
		return ec2v1.AutoscalingGroup_Instance_UNKNOWN
	}
	return ec2v1.AutoscalingGroup_Instance_LifecycleState(val)
}

func newProtoForAutoscalingGroupInstance(instance astypes.Instance) *ec2v1.AutoscalingGroup_Instance {
	return &ec2v1.AutoscalingGroup_Instance{
		Id:                      aws.ToString(instance.InstanceId),
		Zone:                    aws.ToString(instance.AvailabilityZone),
		LaunchConfigurationName: aws.ToString(instance.LaunchConfigurationName),
		Healthy:                 aws.ToString(instance.HealthStatus) == "HEALTHY",
		LifecycleState:          protoForAutoscalingGroupInstanceLifecycleState(string(instance.LifecycleState)),
	}
}

func newProtoForAutoscalingGroup(account string, group astypes.AutoScalingGroup) *ec2v1.AutoscalingGroup {
	pb := &ec2v1.AutoscalingGroup{
		Name:    aws.ToString(group.AutoScalingGroupName),
		Account: account,
		Zones:   group.AvailabilityZones,
		Size: &ec2v1.AutoscalingGroupSize{
			Min:     uint32(aws.ToInt32(group.MinSize)),
			Max:     uint32(aws.ToInt32(group.MaxSize)),
			Desired: uint32(aws.ToInt32(group.DesiredCapacity)),
		},
	}

	if len(pb.Zones) > 0 {
		pb.Region = zoneToRegion(pb.Zones[0])
	}

	pb.TerminationPolicies = make([]ec2v1.AutoscalingGroup_TerminationPolicy, len(group.TerminationPolicies))
	for idx, p := range group.TerminationPolicies {
		pb.TerminationPolicies[idx] = protoForTerminationPolicy(p)
	}

	pb.Instances = make([]*ec2v1.AutoscalingGroup_Instance, len(group.Instances))
	for idx, i := range group.Instances {
		pb.Instances[idx] = newProtoForAutoscalingGroupInstance(i)
	}

	return pb
}

func (c *client) Regions() []string {
	uniqueRegions := map[string]bool{}

	for _, account := range c.accounts {
		for region := range account.clients {
			uniqueRegions[region] = true
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

func (c *client) Accounts() []string {
	accounts := []string{}
	for account := range c.accounts {
		accounts = append(accounts, account)
	}
	return accounts
}

func (c *client) AccountsAndRegions() map[string][]string {
	ar := make(map[string][]string)
	for name, account := range c.accounts {
		ar[name] = account.regions
	}
	return ar
}

// Get all accounts that exist in a specific region
func (c *client) GetAccountsInRegion(region string) []string {
	accounts := []string{}
	for _, a := range c.accounts {
		for _, r := range a.regions {
			if r == region {
				accounts = append(accounts, a.alias)
			}
		}
	}

	return accounts
}

func (c *client) GetPrimaryAccountAlias() string {
	return c.currentAccountAlias
}

func (c *client) DescribeInstances(ctx context.Context, account, region string, ids []string) ([]*ec2v1.Instance, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeInstancesInput{InstanceIds: ids}
	result, err := cl.ec2.DescribeInstances(ctx, input)

	if err != nil {
		return nil, err
	}

	var ret []*ec2v1.Instance
	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			ret = append(ret, newProtoForInstance(i, account))
		}
	}

	return ret, nil
}

func (c *client) TerminateInstances(ctx context.Context, account, region string, ids []string) error {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return err
	}

	input := &ec2.TerminateInstancesInput{InstanceIds: ids}
	_, err = cl.ec2.TerminateInstances(ctx, input)

	return err
}

func (c *client) RebootInstances(ctx context.Context, account, region string, ids []string) error {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return err
	}

	input := &ec2.RebootInstancesInput{InstanceIds: ids}
	_, err = cl.ec2.RebootInstances(ctx, input)

	return err
}

func protoForInstanceState(state string) ec2v1.Instance_State {
	// Transform kebab case 'shutting-down' to upper snake case 'SHUTTING_DOWN'.
	state = strings.ReplaceAll(strings.ToUpper(state), "-", "_")
	// Look up value in generated enum map.
	val, ok := ec2v1.Instance_State_value[state]
	if !ok {
		return ec2v1.Instance_UNKNOWN
	}

	return ec2v1.Instance_State(val)
}

func newProtoForInstance(i ec2types.Instance, account string) *ec2v1.Instance {
	ret := &ec2v1.Instance{
		InstanceId:       aws.ToString(i.InstanceId),
		Account:          account,
		State:            protoForInstanceState(string(i.State.Name)),
		InstanceType:     string(i.InstanceType),
		PublicIpAddress:  aws.ToString(i.PublicIpAddress),
		PrivateIpAddress: aws.ToString(i.PrivateIpAddress),
		AvailabilityZone: aws.ToString(i.Placement.AvailabilityZone),
	}

	ret.Region = zoneToRegion(ret.AvailabilityZone)

	// Transform tag list to map.
	ret.Tags = make(map[string]string, len(i.Tags))
	for _, tag := range i.Tags {
		ret.Tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return ret
}
