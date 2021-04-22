package xds

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lyft/clutch/backend/module/chaos/experimentation/xds"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

func TestServerFaultsRTDSResourceGenerating(t *testing.T) {
	tests := []struct {
		experiment               *experimentstore.Experiment
		ingressRuntimeKeyPrefix  string
		egressRuntimeKeyPrefix   string
		expectedCluster          string
		expectedRuntimeKeyValues []*xds.RuntimeKeyValue
	}{
		{
			// Abort - Service B -> Service A (Internal)
			experiment:              createExperiment(t, "serviceA", "serviceB", 10, 404, INTERNAL, ABORT),
			ingressRuntimeKeyPrefix: "ingressfoo",
			expectedCluster:         "serviceA",
			expectedRuntimeKeyValues: []*xds.RuntimeKeyValue{
				{Key: "ingressfoo.serviceB.abort.abort_percent", Value: 10},
				{Key: "ingressfoo.serviceB.abort.http_status", Value: 404},
			},
		},
		{
			// Latency - Service D -> Service A (Internal)
			experiment:              createExperiment(t, "serviceA", "serviceD", 30, 100, INTERNAL, LATENCY),
			ingressRuntimeKeyPrefix: "ingressfoo",
			expectedCluster:         "serviceA",
			expectedRuntimeKeyValues: []*xds.RuntimeKeyValue{
				{Key: "ingressfoo.serviceD.delay.fixed_delay_percent", Value: 30},
				{Key: "ingressfoo.serviceD.delay.fixed_duration_ms", Value: 100},
			},
		},
		{
			// Abort - Service A -> Service X (External)
			experiment:             createExperiment(t, "serviceX", "serviceA", 65, 400, EXTERNAL, ABORT),
			egressRuntimeKeyPrefix: "egressfoo",
			expectedCluster:        "serviceA",
			expectedRuntimeKeyValues: []*xds.RuntimeKeyValue{
				{Key: "egressfoo.serviceX.abort.abort_percent", Value: 65},
				{Key: "egressfoo.serviceX.abort.http_status", Value: 400},
			},
		},
		{
			// Latency - Service F -> Service Y (External)
			experiment:             createExperiment(t, "serviceY", "serviceF", 40, 200, EXTERNAL, LATENCY),
			egressRuntimeKeyPrefix: "egressfoo",
			expectedCluster:        "serviceF",
			expectedRuntimeKeyValues: []*xds.RuntimeKeyValue{
				{Key: "egressfoo.serviceY.delay.fixed_delay_percent", Value: 40},
				{Key: "egressfoo.serviceY.delay.fixed_duration_ms", Value: 200},
			},
		},
	}

	containsKeyValue := func(keyValues []*xds.RuntimeKeyValue, keyValue *xds.RuntimeKeyValue) bool {
		for _, kv := range keyValues {
			if kv.Key == keyValue.Key && kv.Value == keyValue.Value {
				return true
			}
		}
		return false
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			g := RTDSServerFaultsResourceGenerator{
				IngressFaultRuntimePrefix: tt.ingressRuntimeKeyPrefix,
				EgressFaultRuntimePrefix:  tt.egressRuntimeKeyPrefix,
			}
			r, err := g.GenerateResource(tt.experiment)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCluster, r.Cluster)
			assert.Len(t, tt.expectedRuntimeKeyValues, len(r.RuntimeKeyValues))
			for _, kv := range tt.expectedRuntimeKeyValues {
				assert.True(t, containsKeyValue(r.RuntimeKeyValues, kv))
			}
		})
	}
}
