package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/stretchr/testify/assert"
)

type mockIAM struct {
	getSimulationResultsErr error
	getSimulationResults    *iam.SimulateCustomPolicyOutput
	getGetIAMRoleResultsErr error
	getGetIAMRoleResults    *iam.GetRoleOutput
}

func TestIAMSimulateCustomPolicy(t *testing.T) {
	iamClient := &mockIAM{
		getSimulationResults: &iam.SimulateCustomPolicyOutput{
			EvaluationResults: []types.EvaluationResult{
				{
					EvalActionName:       aws.String("s3:GetBucketPolicy"),
					EvalDecision:         types.PolicyEvaluationDecisionTypeImplicitDeny,
					EvalResourceName:     aws.String("arn:aws:s3:::a/*"),
					MissingContextValues: []string{},
					MatchedStatements:    []types.Statement{},
				},
			},
		},
	}

	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", iam: iamClient},
				},
			},
		},
	}

	simulationInput := &iam.SimulateCustomPolicyInput{
		ActionNames:     []string{"s3:GetObject"},
		PolicyInputList: []string{`{"Version":"2012-10-17","Statement":[{"Sid":"Stmt1635280953828","Action":["s3:AbortMultipartUpload"],"Effect":"Allow","Resource":"arn:aws:s3:::a/*"}]}`},
		CallerArn:       aws.String("arn:aws:iam::000000000000:user/fake-user"),
		ResourceArns:    []string{"arn:aws:s3:::a/*"},
		ResourcePolicy:  aws.String(`{}`),
	}

	output, err := c.SimulateCustomPolicy(context.Background(), "default", "us-east-1", simulationInput)
	assert.NoError(t, err)
	assert.Equal(t, output, &iam.SimulateCustomPolicyOutput{
		EvaluationResults: []types.EvaluationResult{
			{
				EvalActionName:       aws.String("s3:GetBucketPolicy"),
				EvalDecision:         types.PolicyEvaluationDecisionTypeImplicitDeny,
				EvalResourceName:     aws.String("arn:aws:s3:::a/*"),
				MissingContextValues: []string{},
				MatchedStatements:    []types.Statement{},
			},
		},
	})
}

func TestIAMSimulateCustomPolicyErrorHandling(t *testing.T) {
	iamClient := &mockIAM{
		getSimulationResultsErr: fmt.Errorf("error"),
	}
	c := &client{
		currentAccountAlias: "default",
		accounts: map[string]*accountClients{
			"default": {
				clients: map[string]*regionalClient{
					"us-east-1": {region: "us-east-1", iam: iamClient},
				},
			},
		},
	}

	simulationInput := &iam.SimulateCustomPolicyInput{
		ActionNames:     []string{"s3:GetObject"},
		PolicyInputList: []string{`{"Version":"2012-10-17","Statement":[{"Sid":"Stmt1635280953828","Action":["s3:AbortMultipartUpload"],"Effect":"Allow","Resource":"arn:aws:s3:::a/*"}]}`},
		CallerArn:       aws.String("arn:aws:iam::000000000000:role/fake-user"),
		ResourceArns:    []string{"arn:aws:s3:::a/*"},
		ResourcePolicy:  aws.String(`{}`),
	}

	output1, err1 := c.SimulateCustomPolicy(context.Background(), "default", "us-east-1", simulationInput)
	assert.Nil(t, output1)
	assert.Error(t, err1)

	// Test unknown region
	output2, err2 := c.SimulateCustomPolicy(context.Background(), "default", "invalid-region-1", simulationInput)
	assert.Nil(t, output2)
	assert.Error(t, err2)
}

func (m *mockIAM) SimulateCustomPolicy(ctx context.Context, params *iam.SimulateCustomPolicyInput, optFns ...func(*iam.Options)) (*iam.SimulateCustomPolicyOutput, error) {
	if m.getSimulationResultsErr != nil {
		return nil, m.getSimulationResultsErr
	}
	return m.getSimulationResults, nil
}

func (m *mockIAM) GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error) {
	if m.getGetIAMRoleResultsErr != nil {
		return nil, m.getGetIAMRoleResultsErr
	}
	return m.getGetIAMRoleResults, nil
}
