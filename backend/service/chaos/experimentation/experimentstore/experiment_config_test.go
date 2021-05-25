package experimentstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	serverexperimentationv1 "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
)

func TestExperimentConfigInitializationSuccess(t *testing.T) {
	config := &serverexperimentationv1.HTTPFaultConfig{
		Fault: &serverexperimentationv1.HTTPFaultConfig_AbortFault{
			AbortFault: &serverexperimentationv1.AbortFault{
				Percentage:  &serverexperimentationv1.FaultPercentage{Percentage: 100},
				AbortStatus: &serverexperimentationv1.FaultAbortStatus{HttpStatusCode: 401},
			},
		},
		FaultTargeting: &serverexperimentationv1.FaultTargeting{
			Enforcer: &serverexperimentationv1.FaultTargeting_UpstreamEnforcing{
				UpstreamEnforcing: &serverexperimentationv1.UpstreamEnforcing{
					UpstreamType: &serverexperimentationv1.UpstreamEnforcing_UpstreamCluster{
						UpstreamCluster: &serverexperimentationv1.SingleCluster{
							Name: "upstreamCluster",
						},
					},
					DownstreamType: &serverexperimentationv1.UpstreamEnforcing_DownstreamCluster{
						DownstreamCluster: &serverexperimentationv1.SingleCluster{
							Name: "downstreamCluster",
						},
					},
				},
			},
		},
	}

	anyConfig, err := anypb.New(config)
	assert.NoError(t, err)
	jsonConfig, err := protojson.Marshal(anyConfig)
	assert.NoError(t, err)
	message, err := anypb.UnmarshalNew(anyConfig, proto.UnmarshalOptions{})
	assert.NoError(t, err)

	ec, err := NewExperimentConfig("123", string(jsonConfig))
	assert.NoError(t, err)
	_, ok := ec.Message.(*serverexperimentationv1.HTTPFaultConfig)
	assert.True(t, ok)
	assert.Equal(t, "123", ec.Id)
	assert.Equal(t, anyConfig, ec.Config)
	assert.Equal(t, message, ec.Message)
}

func TestExperimentConfigInitializationFailure(t *testing.T) {
	config, err := NewExperimentConfig("123", "foo")
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestExperimentConfigProperties(t *testing.T) {
	config := &serverexperimentationv1.HTTPFaultConfig{
		Fault: &serverexperimentationv1.HTTPFaultConfig_AbortFault{
			AbortFault: &serverexperimentationv1.AbortFault{
				Percentage:  &serverexperimentationv1.FaultPercentage{Percentage: 100},
				AbortStatus: &serverexperimentationv1.FaultAbortStatus{HttpStatusCode: 401},
			},
		},
		FaultTargeting: &serverexperimentationv1.FaultTargeting{
			Enforcer: &serverexperimentationv1.FaultTargeting_UpstreamEnforcing{
				UpstreamEnforcing: &serverexperimentationv1.UpstreamEnforcing{
					UpstreamType: &serverexperimentationv1.UpstreamEnforcing_UpstreamCluster{
						UpstreamCluster: &serverexperimentationv1.SingleCluster{
							Name: "upstreamCluster",
						},
					},
					DownstreamType: &serverexperimentationv1.UpstreamEnforcing_DownstreamCluster{
						DownstreamCluster: &serverexperimentationv1.SingleCluster{
							Name: "downstreamCluster",
						},
					},
				},
			},
		},
	}

	anyConfig, err := anypb.New(config)
	assert.NoError(t, err)
	jsonConfig, err := protojson.Marshal(anyConfig)
	assert.NoError(t, err)

	ec, err := NewExperimentConfig("1", string(jsonConfig))
	assert.NoError(t, err)

	properties, err := ec.CreateProperties()
	assert.NoError(t, err)
	assert.Len(t, properties, 1)
	assert.Equal(t, "config_identifier", properties[0].Id)
	assert.Equal(t, "1", properties[0].GetStringValue())
}
