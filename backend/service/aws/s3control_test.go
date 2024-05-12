package aws

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
)

func TestS3ControlGetAccessPointPolicy(t *testing.T) {
	s3ControlClient := &mockS3Control{
		getAccessPointPolicyOutput: &s3control.GetAccessPointPolicyOutput{
			Policy: aws.String("policy"),
		},
	}

	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", s3control: s3ControlClient},
				},
			},
		},
	}

	output, err := c.S3GetAccessPointPolicy(context.Background(), "default", "us-east-1", "access-point", "accountID")
	assert.NoError(t, err)
	assert.Equal(t, output.Policy, aws.String("policy"))
}

func TestS3ControlGetAccessPointPolicyErrorHandling(t *testing.T) {
	s3ControlClient := &mockS3Control{
		getAccessPointPolicyErr: fmt.Errorf("error"),
	}

	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", s3control: s3ControlClient},
				},
			},
		},
	}

	output1, err1 := c.S3GetAccessPointPolicy(context.Background(), "default", "us-east-1", "access-point", "accountID")
	assert.Nil(t, output1)
	assert.Error(t, err1)

	// Test unknown region
	output2, err2 := c.S3GetAccessPointPolicy(context.Background(), "default", "choice-region-1", "access-point", "accountID")
	assert.Nil(t, output2)
	assert.Error(t, err2)
}

type mockS3Control struct {
	getAccessPointPolicyOutput *s3control.GetAccessPointPolicyOutput
	getAccessPointPolicyErr    error
}

func (m *mockS3Control) GetAccessPointPolicy(ctx context.Context, params *s3control.GetAccessPointPolicyInput, optFns ...func(*s3control.Options)) (*s3control.GetAccessPointPolicyOutput, error) {
	if m.getAccessPointPolicyErr != nil {
		return nil, m.getAccessPointPolicyErr
	}
	return m.getAccessPointPolicyOutput, nil
}
