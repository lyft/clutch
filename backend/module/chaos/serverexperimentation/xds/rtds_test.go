package xds

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

func TestCreateRuntimeKeys(t *testing.T) {
	testDataList := mockGenerateFaultData(t)
	ingressPrefix := "ingress"
	egressPrefix := "egress"

	rtdsConfig := RTDSConfig{
		layerName:     "testLayerName",
		ingressPrefix: ingressPrefix,
		egressPrefix:  egressPrefix,
	}
	runtimePrefixes := &experimentstore.RuntimePrefixes{
		IngressPrefix: rtdsConfig.ingressPrefix,
		EgressPrefix:  rtdsConfig.egressPrefix,
	}

	for _, testExperiment := range testDataList {
		var expectedPercentageKey string
		var expectedPercentageValue uint32
		var expectedFaultKey string
		var expectedFaultValue uint32

		config := &serverexperimentation.HTTPFaultConfig{}
		err := testExperiment.Config.Config.UnmarshalTo(config)
		if err != nil {
			t.Errorf("unmarshalAny failed %v", err)
		}

		upstream, downstream, err := getClusterPair(config)
		assert.NoError(t, err)

		switch config.GetFault().(type) {
		case *serverexperimentation.HTTPFaultConfig_AbortFault:
			abort := config.GetAbortFault()
			expectedFaultValue = abort.GetAbortStatus().GetHttpStatusCode()
			expectedPercentageValue = abort.GetPercentage().GetPercentage()

			switch config.GetFaultTargeting().GetEnforcer().(type) {
			case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
				expectedPercentageKey = fmt.Sprintf(HTTPPercentageForEgress, egressPrefix, upstream)
				expectedFaultKey = fmt.Sprintf(HTTPStatusForEgress, egressPrefix, upstream)
			case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
				if downstream == "" {
					expectedPercentageKey = fmt.Sprintf(HTTPPercentageWithoutDownstream, ingressPrefix)
					expectedFaultKey = fmt.Sprintf(HTTPStatusWithoutDownstream, ingressPrefix)
				} else {
					expectedPercentageKey = fmt.Sprintf(HTTPPercentageWithDownstream, ingressPrefix, downstream)
					expectedFaultKey = fmt.Sprintf(HTTPStatusWithDownstream, ingressPrefix, downstream)
				}
			}
		case *serverexperimentation.HTTPFaultConfig_LatencyFault:
			latency := config.GetLatencyFault()
			expectedFaultValue = latency.GetLatencyDuration().GetFixedDurationMs()
			expectedPercentageValue = latency.GetPercentage().GetPercentage()

			switch config.GetFaultTargeting().GetEnforcer().(type) {
			case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
				expectedPercentageKey = fmt.Sprintf(LatencyPercentageForEgress, egressPrefix, upstream)
				expectedFaultKey = fmt.Sprintf(LatencyDurationForEgress, egressPrefix, upstream)
			case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
				if downstream == "" {
					expectedPercentageKey = fmt.Sprintf(LatencyPercentageWithoutDownstream, ingressPrefix)
					expectedFaultKey = fmt.Sprintf(LatencyDurationWithoutDownstream, ingressPrefix)
				} else {
					expectedPercentageKey = fmt.Sprintf(LatencyPercentageWithDownstream, ingressPrefix, downstream)
					expectedFaultKey = fmt.Sprintf(LatencyDurationWithDownstream, ingressPrefix, downstream)
				}
			}
		}

		percentageKey, percentageValue, faultKey, faultValue, err := createRuntimeKeys(upstream, downstream, config, runtimePrefixes)
		assert.NoError(t, err)

		assert.Equal(t, expectedPercentageKey, percentageKey)
		assert.Equal(t, expectedPercentageValue, percentageValue)
		assert.Equal(t, expectedFaultKey, faultKey)
		assert.Equal(t, expectedFaultValue, faultValue)
	}
}
