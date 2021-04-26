package xds

import (
	"fmt"
	"testing"

	gcpRoute "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	gcpFilterCommon "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/common/fault/v3"
	gcpFilterFault "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	gcpType "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

func TestECDSFaultsGeneration(t *testing.T) {
	tests := []struct {
		experiment                       *experimentstore.Experiment
		expectedCluster                  string
		expectedResourceName             string
		expectedAbort                    *gcpFilterFault.FaultAbort
		expectedDelay                    *gcpFilterCommon.FaultDelay
		expectedHeadersDownstreamCluster string
	}{
		{
			// Abort - Service B -> Service A (Internal)
			experiment:           createExperiment(t, "serviceA", "serviceB", 10, 404, faultUpstreamServiceTypeInternal, faultTypeAbort),
			expectedCluster:      "serviceA",
			expectedResourceName: "envoy.extension_config",
			expectedAbort: &gcpFilterFault.FaultAbort{
				ErrorType: &gcpFilterFault.FaultAbort_HttpStatus{
					HttpStatus: 404,
				},
				Percentage: &gcpType.FractionalPercent{
					Numerator:   10,
					Denominator: gcpType.FractionalPercent_HUNDRED,
				},
			},
			expectedDelay:                    nil,
			expectedHeadersDownstreamCluster: "serviceB",
		},
		{
			// Latency - Service D -> Service A (Internal)
			experiment:           createExperiment(t, "serviceA", "serviceD", 30, 100, faultUpstreamServiceTypeInternal, faultTypeLatency),
			expectedCluster:      "serviceA",
			expectedResourceName: "envoy.extension_config",
			expectedAbort:        nil,
			expectedDelay: &gcpFilterCommon.FaultDelay{
				FaultDelaySecifier: &gcpFilterCommon.FaultDelay_FixedDelay{
					FixedDelay: &duration.Duration{
						Nanos: 100000000, // 0.001 second
					},
				},
				Percentage: &gcpType.FractionalPercent{
					Numerator:   30,
					Denominator: gcpType.FractionalPercent_HUNDRED,
				},
			},
			expectedHeadersDownstreamCluster: "serviceD",
		},
		{
			// Abort - Service A -> Service X (External)
			experiment:           createExperiment(t, "serviceX", "serviceA", 65, 400, faultUpstreamServiceTypeExternal, faultTypeAbort),
			expectedCluster:      "serviceA",
			expectedResourceName: "envoy.egress.extension_config.serviceX",
			expectedAbort: &gcpFilterFault.FaultAbort{
				ErrorType: &gcpFilterFault.FaultAbort_HttpStatus{
					HttpStatus: 400,
				},
				Percentage: &gcpType.FractionalPercent{
					Numerator:   65,
					Denominator: gcpType.FractionalPercent_HUNDRED,
				},
			},
			expectedDelay: nil,
		},
		{
			// Latency - Service F -> Service Y (External)
			experiment:           createExperiment(t, "serviceY", "serviceF", 40, 200, faultUpstreamServiceTypeExternal, faultTypeLatency),
			expectedCluster:      "serviceF",
			expectedResourceName: "envoy.egress.extension_config.serviceY",
			expectedAbort:        nil,
			expectedDelay: &gcpFilterCommon.FaultDelay{
				FaultDelaySecifier: &gcpFilterCommon.FaultDelay_FixedDelay{
					FixedDelay: &duration.Duration{
						Nanos: 200000000, // 0.001 second
					},
				},
				Percentage: &gcpType.FractionalPercent{
					Numerator:   40,
					Denominator: gcpType.FractionalPercent_HUNDRED,
				},
			},
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			g := &ECDSFaultsGenerator{}
			r, err := g.GenerateResource(tt.experiment)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCluster, r.Cluster)
			assert.Equal(t, tt.expectedResourceName, r.ExtensionConfig.Name)
			faultFilter := &gcpFilterFault.HTTPFault{}
			err = r.ExtensionConfig.TypedConfig.UnmarshalTo(faultFilter)
			assert.NoError(t, err)
			assert.NotNil(t, faultFilter.AbortPercentRuntime)
			assert.NotNil(t, faultFilter.AbortHttpStatusRuntime)
			assert.NotNil(t, faultFilter.DelayPercentRuntime)
			assert.NotNil(t, faultFilter.DelayDurationRuntime)
			if tt.expectedAbort != nil {
				assert.Equal(t, tt.expectedAbort, faultFilter.Abort)
			}
			if tt.expectedDelay != nil {
				assert.Equal(t, tt.expectedDelay, faultFilter.Delay)
			}

			if tt.expectedHeadersDownstreamCluster != "" {
				expectedHeaders := []*gcpRoute.HeaderMatcher{
					{
						Name: envoyDownstreamServiceClusterHeader,
						HeaderMatchSpecifier: &gcpRoute.HeaderMatcher_ExactMatch{
							ExactMatch: tt.expectedHeadersDownstreamCluster,
						},
					},
				}
				assert.Equal(t, expectedHeaders, faultFilter.Headers)
			}
		})
	}
}

func TestECDSDefaultFaultsGeneration(t *testing.T) {
	tests := []struct {
		cluster         string
		resourceName    string
		expectedCluster string
	}{
		{
			cluster:         "foo",
			resourceName:    "envoy.extension_config",
			expectedCluster: "foo",
		},
		{
			cluster:         "foo",
			resourceName:    "envoy.egress.extension_config.bar",
			expectedCluster: "foo",
		},
		{
			cluster:         "foo",
			resourceName:    "not_existing_resource",
			expectedCluster: "",
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			g, err := NewECDSFaultsGenerator()
			assert.NoError(t, err)
			r, err := g.GenerateDefaultResource(tt.cluster, tt.resourceName)
			assert.NoError(t, err)
			if tt.expectedCluster != "" {
				assert.Equal(t, tt.expectedCluster, r.Cluster)
				faultFilter := &gcpFilterFault.HTTPFault{}
				err = r.ExtensionConfig.TypedConfig.UnmarshalTo(faultFilter)
				assert.NoError(t, err)
				assert.NotNil(t, faultFilter.AbortPercentRuntime)
				assert.NotNil(t, faultFilter.AbortHttpStatusRuntime)
				assert.NotNil(t, faultFilter.DelayPercentRuntime)
				assert.NotNil(t, faultFilter.DelayDurationRuntime)
			} else {
				assert.Nil(t, r)
			}
		})
	}
}

const (
	faultUpstreamServiceTypeInternal = "internal"
	faultUpstreamServiceTypeExternal = "external"
	faultTypeAbort                   = "abort"
	faultTypeLatency                 = "latency"
)

func createExperiment(t *testing.T, upstreamCluster string, downstreamCluster string, faultPercent uint32, faultValue uint32, faultInjectorEnforcing string, faultType string) *experimentstore.Experiment {
	httpConfig := &serverexperimentation.HTTPFaultConfig{}

	upstreamSingleCluster := &serverexperimentation.SingleCluster{
		Name: upstreamCluster,
	}

	downstreamSingleCluster := &serverexperimentation.SingleCluster{
		Name: downstreamCluster,
	}

	switch faultType {
	case faultTypeAbort:
		httpConfig.Fault = &serverexperimentation.HTTPFaultConfig_AbortFault{
			AbortFault: &serverexperimentation.AbortFault{
				Percentage:  &serverexperimentation.FaultPercentage{Percentage: faultPercent},
				AbortStatus: &serverexperimentation.FaultAbortStatus{HttpStatusCode: faultValue},
			},
		}
	case faultTypeLatency:
		httpConfig.Fault = &serverexperimentation.HTTPFaultConfig_LatencyFault{
			LatencyFault: &serverexperimentation.LatencyFault{
				Percentage:      &serverexperimentation.FaultPercentage{Percentage: faultPercent},
				LatencyDuration: &serverexperimentation.FaultLatencyDuration{FixedDurationMs: faultValue},
			},
		}
	}

	switch faultInjectorEnforcing {
	case faultUpstreamServiceTypeInternal:
		httpConfig.FaultTargeting = &serverexperimentation.FaultTargeting{
			Enforcer: &serverexperimentation.FaultTargeting_UpstreamEnforcing{
				UpstreamEnforcing: &serverexperimentation.UpstreamEnforcing{
					UpstreamType: &serverexperimentation.UpstreamEnforcing_UpstreamCluster{
						UpstreamCluster: upstreamSingleCluster,
					},
					DownstreamType: &serverexperimentation.UpstreamEnforcing_DownstreamCluster{
						DownstreamCluster: downstreamSingleCluster,
					},
				},
			},
		}
	case faultUpstreamServiceTypeExternal:
		httpConfig.FaultTargeting = &serverexperimentation.FaultTargeting{
			Enforcer: &serverexperimentation.FaultTargeting_DownstreamEnforcing{
				DownstreamEnforcing: &serverexperimentation.DownstreamEnforcing{
					UpstreamType: &serverexperimentation.DownstreamEnforcing_UpstreamCluster{
						UpstreamCluster: upstreamSingleCluster,
					},
					DownstreamType: &serverexperimentation.DownstreamEnforcing_DownstreamCluster{
						DownstreamCluster: downstreamSingleCluster,
					},
				},
			},
		}
	}

	anyConfig, err := anypb.New(httpConfig)
	assert.NoError(t, err)
	jsonConfig, err := protojson.Marshal(anyConfig)
	assert.NoError(t, err)
	config, err := experimentstore.NewExperimentConfig("1", string(jsonConfig))
	assert.NoError(t, err)

	return &experimentstore.Experiment{
		Config: config,
	}
}
