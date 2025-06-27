package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stretchr/testify/assert"

	rdsv1 "github.com/lyft/clutch/backend/api/aws/rds/v1"
)

func TestRDSDescribeCluster(t *testing.T) {
	rdsClient := &mockRDS{
		describeClusterOutput: &rds.DescribeDBClustersOutput{
			DBClusters: []types.DBCluster{
				{
					DBClusterIdentifier: aws.String("test-cluster"),
					Engine:              aws.String("aurora-mysql"),
					Status:              aws.String("available"),
				},
			},
		},
	}

	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", rds: rdsClient},
				},
			},
		},
	}

	output, err := c.RDSDescribeCluster(context.Background(), "default", "us-east-1", "test-cluster")
	assert.NoError(t, err)
	assert.Equal(t, &rdsv1.Cluster{
		ClusterIdentifier: "test-cluster",
		Account:           "default",
		Region:            "us-east-1",
		Engine:            "aurora-mysql",
		Status:            "available",
	}, output)
}

func TestRDSDescribeClusterErrorHandling(t *testing.T) {
	rdsClient := &mockRDS{
		describeClusterErr: fmt.Errorf("error"),
	}
	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", rds: rdsClient},
				},
			},
		},
	}

	output1, err1 := c.RDSDescribeCluster(context.Background(), "default", "us-east-1", "test-cluster")
	assert.Nil(t, output1)
	assert.Error(t, err1)

	// Test unknown region
	output2, err2 := c.RDSDescribeCluster(context.Background(), "default", "invalid-region-1", "test-cluster")
	assert.Nil(t, output2)
	assert.Error(t, err2)
}

type mockRDS struct {
	describeClusterErr    error
	describeClusterOutput *rds.DescribeDBClustersOutput
}

func (m *mockRDS) DescribeDBClusters(ctx context.Context, params *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error) {
	if m.describeClusterErr != nil {
		return nil, m.describeClusterErr
	}
	return m.describeClusterOutput, nil
}
