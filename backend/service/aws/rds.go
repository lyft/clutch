package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rdsv1 "github.com/lyft/clutch/backend/api/aws/rds/v1"
)

// RDSDescribeCluster fetches a single cluster description and returns it as proto.
func (c *client) RDSDescribeCluster(ctx context.Context, account, region, clusterName string) (*rdsv1.Cluster, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	in := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(clusterName),
	}

	out, err := cl.rds.DescribeDBClusters(ctx, in)
	if err != nil {
		return nil, err
	}

	if len(out.DBClusters) == 0 {
		return nil, status.Error(codes.NotFound, "cluster not found")
	}

	sdkCluster := out.DBClusters[0]

	proto := &rdsv1.Cluster{
		ClusterIdentifier: aws.ToString(sdkCluster.DBClusterIdentifier),
		Account:           account,
		Region:            region,
		Engine:            aws.ToString(sdkCluster.Engine),
		Status:            aws.ToString(sdkCluster.Status),
	}

	return proto, nil
}
