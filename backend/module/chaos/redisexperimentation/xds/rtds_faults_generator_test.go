package xds

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	redisexperimentation "github.com/lyft/clutch/backend/api/chaos/redisexperimentation/v1"
	"github.com/lyft/clutch/backend/module/chaos/experimentation/xds"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

func TestRTDSFaultsGeneration(t *testing.T) {
	tests := []struct {
		experiment              *experimentstore.Experiment
		prefix                  string
		expectedCluster         string
		expectedRuntimeKeyValue *xds.RuntimeKeyValue
	}{
		{
			experiment:              createExperiment(t, "foo1", "bar", Latency, 1),
			prefix:                  "pre",
			expectedCluster:         "foo1",
			expectedRuntimeKeyValue: &xds.RuntimeKeyValue{Key: "pre.bar.delay.fixed_delay_percent", Value: 1},
		},
		{
			experiment:              createExperiment(t, "foo2", "zar", Error, 2),
			prefix:                  "prefix",
			expectedCluster:         "foo2",
			expectedRuntimeKeyValue: &xds.RuntimeKeyValue{Key: "prefix.zar.error.error_percent", Value: 2},
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			g := RTDSFaultsGenerator{
				FaultRuntimePrefix: tt.prefix,
			}
			r, err := g.GenerateResource(tt.experiment)

			assert.NoError(t, err)
			assert.Len(t, r.RuntimeKeyValues, 1)
			assert.Equal(t, tt.expectedRuntimeKeyValue, r.RuntimeKeyValues[0])
		})
	}
}

const (
	Error   = `error`
	Latency = `latency`
)

func createExperiment(t *testing.T, downstreamCluster string, upstreamCluster string, faultType string, faultPercentage uint32) *experimentstore.Experiment {
	fault := &redisexperimentation.FaultConfig{
		FaultTargeting: &redisexperimentation.FaultTargeting{
			DownstreamCluster: &redisexperimentation.SingleCluster{
				Name: downstreamCluster,
			},
			UpstreamCluster: &redisexperimentation.SingleCluster{
				Name: upstreamCluster,
			},
		},
	}

	if faultType == Error {
		fault.Fault = &redisexperimentation.FaultConfig_ErrorFault{
			ErrorFault: &redisexperimentation.ErrorFault{
				Percentage: &redisexperimentation.FaultPercentage{
					Percentage: faultPercentage,
				},
			},
		}
	} else {
		fault.Fault = &redisexperimentation.FaultConfig_LatencyFault{
			LatencyFault: &redisexperimentation.LatencyFault{
				Percentage: &redisexperimentation.FaultPercentage{
					Percentage: faultPercentage,
				},
			},
		}
	}

	anyConfig, err := anypb.New(fault)
	assert.NoError(t, err)
	jsonConfig, err := protojson.Marshal(anyConfig)
	assert.NoError(t, err)
	config, err := experimentstore.NewExperimentConfig("1", string(jsonConfig))
	assert.NoError(t, err)

	return &experimentstore.Experiment{
		Config: config,
	}
}
