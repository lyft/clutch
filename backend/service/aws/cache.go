package aws

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
)

const topologyObjectChanBufferSize = 5000
const topologyLockId = 1

func (c *client) CacheEnabled() bool {
	return true
}

func (c *client) StartTopologyCaching(ctx context.Context, ttl time.Duration) (<-chan *topologyv1.UpdateCacheRequest, error) {
	if !c.topologyLock.TryAcquire(topologyLockId) {
		return nil, errors.New("TopologyCaching is already in progress")
	}

	go func() {
		<-ctx.Done()
		close(c.topologyObjectChan)
		c.topologyLock.Release(topologyLockId)
	}()

	c.processRegionTopologyObjects(ctx)
	return c.topologyObjectChan, nil
}

// We process each resource and region individually, creating tickers for each.
// If a single resource and region takes longer to process it will not block other resources.
// Additionally this give us the flexibility to tune the frequency based on resource.
func (c *client) processRegionTopologyObjects(ctx context.Context) {
	for _, account := range c.accounts {
		for name, client := range account.clients {
			c.log.Info("processing topology objects for account in region", zap.String("account", account.alias), zap.String("region", name))
			go c.startTickerForCacheResource(ctx, time.Duration(time.Minute*5), account.alias, client, c.processAllEC2Instances)
			go c.startTickerForCacheResource(ctx, time.Duration(time.Minute*10), account.alias, client, c.processAllAutoScalingGroups)
			go c.startTickerForCacheResource(ctx, time.Duration(time.Minute*30), account.alias, client, c.processAllKinesisStreams)
			go c.startTickerForCacheResource(ctx, time.Duration(time.Minute*30), account.alias, client, c.processAllDynamoDatabases)
			go c.startTickerForCacheResource(ctx, time.Duration(time.Minute*30), account.alias, client, c.processAllS3Buckets)
			go c.startTickerForCacheResource(ctx, time.Duration(time.Minute*30), account.alias, client, c.processAllS3AccessPoints)
			go c.startTickerForCacheResource(ctx, time.Duration(time.Minute*30), account.alias, client, c.processAllIamRoles)
		}
	}
}

func (c *client) startTickerForCacheResource(ctx context.Context, duration time.Duration, account string, client *regionalClient, resourceFunc func(context.Context, string, *regionalClient)) {
	ticker := time.NewTicker(duration)
	for {
		resourceFunc(ctx, account, client)

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (c *client) processAllAutoScalingGroups(ctx context.Context, account string, client *regionalClient) {
	c.log.Info("starting to process auto scaling groups for region", zap.String("region", client.region))
	// 100 is the maximum amount of records per page allowed for this API
	input := autoscaling.DescribeAutoScalingGroupsInput{
		MaxRecords: aws.Int32(100),
	}

	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(client.autoscaling, &input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			c.log.Error("unable to get next asg page", zap.Error(err))
			break
		}

		for _, asg := range output.AutoScalingGroups {
			protoAsg := newProtoForAutoscalingGroup(account, asg)
			protoAsg.Account = account

			asgAny, err := anypb.New(protoAsg)
			if err != nil {
				c.log.Error("unable to marshal asg proto", zap.Error(err))
				continue
			}

			patternId, err := meta.HydratedPatternForProto(protoAsg)
			if err != nil {
				c.log.Error("unable to get proto id from pattern", zap.Error(err))
				continue
			}

			c.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
				Resource: &topologyv1.Resource{
					Id: patternId,
					Pb: asgAny,
				},
				Action: topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
			}
		}
	}
}

func (c *client) processAllEC2Instances(ctx context.Context, account string, client *regionalClient) {
	c.log.Info("starting to process ec2 instances for region", zap.String("region", client.region))
	// 1000 is the maximum amount of records per page allowed for this API
	input := ec2.DescribeInstancesInput{
		MaxResults: aws.Int32(1000),
	}

	paginator := ec2.NewDescribeInstancesPaginator(client.ec2, &input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			c.log.Error("unable to get next ec2 instance page", zap.Error(err))
			break
		}

		for _, reservation := range output.Reservations {
			for _, instance := range reservation.Instances {
				protoInstance := newProtoForInstance(instance, account)
				protoInstance.Account = account

				instanceAny, err := anypb.New(protoInstance)
				if err != nil {
					c.log.Error("unable to marshal instance proto", zap.Error(err))
					continue
				}

				patternId, err := meta.HydratedPatternForProto(protoInstance)
				if err != nil {
					c.log.Error("unable to get proto id from pattern", zap.Error(err))
					continue
				}

				c.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
					Resource: &topologyv1.Resource{
						Id: patternId,
						Pb: instanceAny,
					},
					Action: topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
				}
			}
		}
	}
}

func (c *client) processAllKinesisStreams(ctx context.Context, account string, client *regionalClient) {
	c.log.Info("starting to process kinesis streams for region", zap.String("region", client.region))
	// 100 is arbatrary, currently this API does not have a per page limit,
	// looking at other aws API limits this value felt safe.
	input := kinesis.ListStreamsInput{
		Limit: aws.Int32(100),
	}

	for {
		output, err := client.kinesis.ListStreams(ctx, &input)
		if err != nil {
			c.log.Error("unable to list kinesis stream", zap.Error(err))
			break
		}

		if aws.ToBool(output.HasMoreStreams) {
			input.ExclusiveStartStreamName = &output.StreamNames[len(output.StreamNames)-1]
		} else {
			break
		}

		for _, stream := range output.StreamNames {
			v1Stream, err := c.DescribeKinesisStream(ctx, account, client.region, stream)
			if err != nil {
				c.log.Error("unable to describe kinesis stream", zap.Error(err))
				continue
			}

			protoStream, err := anypb.New(v1Stream)
			if err != nil {
				c.log.Error("unable to marshal kinesis stream", zap.Error(err))
				continue
			}

			patternId, err := meta.HydratedPatternForProto(v1Stream)
			if err != nil {
				c.log.Error("unable to get proto id from pattern", zap.Error(err))
				continue
			}

			c.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
				Resource: &topologyv1.Resource{
					Id: patternId,
					Pb: protoStream,
				},
				Action: topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
			}
		}
	}
}

func (c *client) processAllDynamoDatabases(ctx context.Context, account string, client *regionalClient) {
	c.log.Info("starting to process dynamodb for region", zap.String("region", client.region))

	input := dynamodb.ListTablesInput{
		Limit: aws.Int32(100),
	}

	paginator := dynamodb.NewListTablesPaginator(client.dynamodb, &input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			c.log.Error("unable to get next dynamodb page", zap.Error(err))
			break
		}

		for _, t := range output.TableNames {
			protoTable, err := c.DescribeTable(ctx, account, client.region, t)
			if err != nil {
				c.log.Error("unable to describe dynamodb table", zap.Error(err))
				continue
			}
			protoTable.Account = account

			tableAny, err := anypb.New(protoTable)
			if err != nil {
				c.log.Error("unable to marshal dynamodb proto", zap.Error(err))
				continue
			}

			patternId, err := meta.HydratedPatternForProto(protoTable)
			if err != nil {
				c.log.Error("unable to get proto id from pattern", zap.Error(err))
				continue
			}

			c.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
				Resource: &topologyv1.Resource{
					Id: patternId,
					Pb: tableAny,
				},
				Action: topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
			}
		}
	}
}

func (c *client) processAllS3Buckets(ctx context.Context, account string, client *regionalClient) {
	c.log.Info("starting to process s3 for region", zap.String("region", client.region))

	input := s3.ListBucketsInput{}

	output, err := client.s3.ListBuckets(ctx, &input)
	if err != nil {
		c.log.Error("unable to list s3 buckets", zap.Error(err))
		return
	}
	for _, bucket := range output.Buckets {
		v1Bucket, err := c.S3DescribeBucket(ctx, account, client.region, *bucket.Name)
		if err != nil {
			c.log.Error("unable to describe s3 bucket", zap.Error(err), zap.String("bucket", *bucket.Name), zap.String("region", client.region))
			continue
		}

		protoBucket, err := anypb.New(v1Bucket)
		if err != nil {
			c.log.Error("unable to marshal s3 bucket", zap.Error(err), zap.String("bucket", *bucket.Name), zap.String("region", client.region))
			continue
		}

		patternId, err := meta.HydratedPatternForProto(v1Bucket)
		if err != nil {
			c.log.Error("unable to get proto id from pattern", zap.Error(err), zap.String("bucket", *bucket.Name), zap.String("region", client.region))
			continue
		}

		c.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
			Resource: &topologyv1.Resource{
				Id: patternId,
				Pb: protoBucket,
			},
			Action: topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
		}
	}
}

func (c *client) processAllS3AccessPoints(ctx context.Context, account string, client *regionalClient) {
	c.log.Info("starting to process s3 access points for region", zap.String("region", client.region))

	callerIdentity, err := client.sts.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		c.log.Error("unable to get caller identity", zap.Error(err))
		return
	}

	accountId := callerIdentity.Account

	input := s3control.ListAccessPointsInput{
		AccountId: accountId,
	}

	output, err := client.s3control.ListAccessPoints(ctx, &input)
	if err != nil {
		c.log.Error("unable to list s3 access points", zap.Error(err))
		return
	}
	for _, accessPoint := range output.AccessPointList {
		v1AccessPoint, err := c.S3GetAccessPoint(ctx, account, client.region, *accessPoint.Name, *accountId)
		if err != nil {
			c.log.Error("unable to describe s3 access point", zap.Error(err))
			continue
		}

		protoAccessPoint, err := anypb.New(v1AccessPoint)
		if err != nil {
			c.log.Error("unable to marshal s3 access point", zap.Error(err))
			continue
		}

		patternId, err := meta.HydratedPatternForProto(v1AccessPoint)
		if err != nil {
			c.log.Error("unable to get proto id from pattern", zap.Error(err))
			continue
		}

		c.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
			Resource: &topologyv1.Resource{
				Id: patternId,
				Pb: protoAccessPoint,
			},
			Action: topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
		}
	}
}

func (c *client) processAllIamRoles(ctx context.Context, account string, client *regionalClient) {
	c.log.Info("starting to process iam roles for region", zap.String("region", client.region))

	input := &iam.ListRolesInput{}
	paginator := iam.NewListRolesPaginator(client.iam, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			c.log.Error("unable to get next iam role page", zap.Error(err))
			break
		}

		for _, role := range output.Roles {
			protoRole, err := c.GetIAMRole(ctx, account, client.region, *role.RoleName)
			if err != nil {
				c.log.Error("unable to get iam role", zap.Error(err))
				continue
			}

			roleAny, err := anypb.New(protoRole)
			if err != nil {
				c.log.Error("unable to marshal iam role proto", zap.Error(err))
				continue
			}

			patternId, err := meta.HydratedPatternForProto(protoRole)
			if err != nil {
				c.log.Error("unable to get proto id from pattern", zap.Error(err))
				continue
			}

			c.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
				Resource: &topologyv1.Resource{
					Id: patternId,
					Pb: roleAny,
				},
				Action: topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
			}
		}
	}
}
