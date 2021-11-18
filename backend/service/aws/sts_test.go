package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/stretchr/testify/assert"
)

func TestSTSGetCallerIdentity(t *testing.T) {
	stsClient := &mockSTS{
		getIdentityOutput: &sts.GetCallerIdentityOutput{
			Account: aws.String("000000000000"),
			Arn:     aws.String("arn:aws:sts::000000000000:fake-role"),
			UserId:  aws.String("some_special_id"),
		},
	}
	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", sts: stsClient},
				},
			},
		},
	}

	output, err := c.GetCallerIdentity(context.Background(), "default", "us-east-1")
	assert.NoError(t, err)
	assert.Equal(t, output, &sts.GetCallerIdentityOutput{
		Account: aws.String("000000000000"),
		Arn:     aws.String("arn:aws:sts::000000000000:fake-role"),
		UserId:  aws.String("some_special_id"),
	})
}

func TestSTSGetCallerIdentityErrorHandling(t *testing.T) {
	stsClient := &mockSTS{
		getIdentityErr: fmt.Errorf("error"),
	}
	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", sts: stsClient},
				},
			},
		},
	}

	output1, err1 := c.GetCallerIdentity(context.Background(), "default", "us-east-1")
	assert.Nil(t, output1)
	assert.Error(t, err1)

	// Test unknown region
	output2, err2 := c.GetCallerIdentity(context.Background(), "default", "choice-region-1")
	assert.Nil(t, output2)
	assert.Error(t, err2)
}

type mockSTS struct {
	stsClient

	getIdentityErr    error
	getIdentityOutput *sts.GetCallerIdentityOutput
}

func (m *mockSTS) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	if m.getIdentityErr != nil {
		return nil, m.getIdentityErr
	}

	return m.getIdentityOutput, nil
}
