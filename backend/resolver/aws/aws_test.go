package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	resolverv1 "github.com/lyft/clutch/backend/api/resolver/v1"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/resolver"

	"github.com/lyft/clutch/backend/mock/service/awsmock"
	"github.com/lyft/clutch/backend/mock/service/topologymock"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/topology"
)

func TestNewAwsResolver(t *testing.T) {
	service.Registry["clutch.service.aws"] = awsmock.New()
	service.Registry["clutch.service.topology"] = topologymock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	aws, err := New(nil, log, scope)
	assert.NoError(t, err)
	assert.NotNil(t, aws)

	// Test successful construction without the topology service
	delete(service.Registry, "clutch.service.topology")
	awsNoTopology, err2 := New(nil, log, scope)
	assert.NoError(t, err2)
	assert.NotNil(t, awsNoTopology)
}

func TestAutoCompleteErrorHandling(t *testing.T) {
	service.Registry["clutch.service.aws"] = awsmock.New()
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	aws, err := New(nil, log, scope)
	assert.NoError(t, err)
	assert.NotNil(t, aws)

	// Test error handling for topology service not found
	_, err2 := aws.Autocomplete(context.Background(), "type_url", "search", 0)
	assert.Error(t, err2)

	// Test error handling for a topology search failure
	awsResolver := res{
		topology: &mockTopologySearch{
			autoCompleteError: fmt.Errorf("error"),
		},
	}
	_, err3 := awsResolver.Autocomplete(context.Background(), "type_url", "search", 0)
	assert.Error(t, err3)
}

func TestAutoCompleteResults(t *testing.T) {
	awsResolver := res{
		topology: &mockTopologySearch{
			autoCompleteResults: []*topologyv1.Resource{
				{
					Id: "meow",
				},
				{
					Id: "cat",
				},
				{
					Id: "yawn",
				},
			},
		},
	}

	expect := []*resolverv1.AutocompleteResult{
		{
			Id: "meow",
		},
		{
			Id: "cat",
		},
		{
			Id: "yawn",
		},
	}

	results, err := awsResolver.Autocomplete(context.Background(), "type_url", "search", 0)
	assert.NoError(t, err)
	assert.Equal(t, expect, results)
}

func TestDetermineAccountAndRegionsForOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		id           string
		accountInput string
		regionInput  string
		expect       map[string][]string
	}{
		{
			id:           "all accounts and regions",
			accountInput: resolver.OptionAll,
			regionInput:  resolver.OptionAll,
			expect: map[string][]string{
				"default":    {"us-mock-1"},
				"staging":    {"us-mock-1", "us-mock-2"},
				"production": {"us-mock-1", "us-mock-2", "us-mock-3"},
			},
		},
		{
			id:           "Only pull from a specific region",
			accountInput: resolver.OptionAll,
			regionInput:  "us-mock-2",
			expect: map[string][]string{
				"staging":    {"us-mock-2"},
				"production": {"us-mock-2"},
			},
		},
		{
			id:           "Only pull from a specific account",
			accountInput: "staging",
			regionInput:  resolver.OptionAll,
			expect: map[string][]string{
				"staging": {"us-mock-1", "us-mock-2"},
			},
		},
	}

	res := &res{
		client: awsmock.New(),
	}

	for _, test := range tests {
		actual := res.determineAccountAndRegionsForOption(test.accountInput, test.regionInput)
		for account, region := range test.expect {
			val, ok := actual[account]
			assert.True(t, ok)
			assert.ElementsMatch(t, region, val)
		}
	}
}

type mockTopologySearch struct {
	topology.Service

	autoCompleteError   error
	autoCompleteResults []*topologyv1.Resource
}

func (m *mockTopologySearch) Autocomplete(ctx context.Context, typeURL, search string, limit uint64) ([]*topologyv1.Resource, error) {
	if m.autoCompleteError != nil {
		return nil, m.autoCompleteError
	}

	return m.autoCompleteResults, nil
}
