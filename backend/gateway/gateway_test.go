package gateway

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
)

func TestLoadEnv(t *testing.T) {
	type File struct {
		name  string
		value string
	}
	testCases := []struct {
		files         []File
		envVar        string
		expectedValue string
	}{
		{
			files: []File{
				{
					name:  ".env.dev",
					value: "FOOBAR1=true",
				},
				{
					name:  ".env",
					value: "FOOBAR1=false",
				},
			},
			envVar:        "FOOBAR1",
			expectedValue: "true",
		},
		{
			files: []File{
				{
					name:  ".env.dev",
					value: "",
				},
				{
					name:  ".env",
					value: "FOOBAR2=true",
				},
			},
			envVar:        "FOOBAR2",
			expectedValue: "true",
		},
	}
	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			tmpVar := os.Getenv(tc.envVar)
			fileNames := []string{}
			for _, f := range tc.files {
				envFile, err := os.CreateTemp(".", f.name)
				if err != nil {
					log.Fatal(err)
				}
				defer os.Remove(envFile.Name())
				_ = os.WriteFile(envFile.Name(), []byte(f.value), 0644)
				fileNames = append(fileNames, envFile.Name())
			}

			flags := &Flags{EnvFiles: fileNames}

			loadEnv(flags)
			v := os.Getenv(tc.envVar)
			assert.Equal(t, tc.expectedValue, v)
			os.Setenv(tc.envVar, tmpVar)
		})
	}
}

func TestGetStatsReporterConfiguration(t *testing.T) {
	testCases := []struct {
		cfg                *gatewayv1.Config
		prefix             string
		nullReporter       bool
		prometheusreporter bool
	}{
		{
			cfg:          &gatewayv1.Config{Gateway: &gatewayv1.GatewayOptions{Stats: &gatewayv1.Stats{}}},
			prefix:       "",
			nullReporter: true,
		},
		{
			cfg: &gatewayv1.Config{Gateway: &gatewayv1.GatewayOptions{
				Stats: &gatewayv1.Stats{Prefix: "gateway", Reporter: &gatewayv1.Stats_LogReporter_{}},
			}},
			prefix: "gateway",
		},
		{
			cfg: &gatewayv1.Config{Gateway: &gatewayv1.GatewayOptions{
				Stats: &gatewayv1.Stats{
					Reporter: &gatewayv1.Stats_StatsdReporter_{StatsdReporter: &gatewayv1.Stats_StatsdReporter{Address: "0.0.0.0:000"}},
				},
			}},
			prefix: "clutch",
		},
		{
			cfg: &gatewayv1.Config{Gateway: &gatewayv1.GatewayOptions{
				Stats: &gatewayv1.Stats{
					Reporter: &gatewayv1.Stats_PrometheusReporter_{PrometheusReporter: &gatewayv1.Stats_PrometheusReporter{HandlerPath: "foo.com/path"}},
				},
			}},
			prefix:             "clutch",
			prometheusreporter: true,
		},
	}

	for _, test := range testCases {
		logger := zaptest.NewLogger(t)
		scopeOptions, metricsHandler := getStatsReporterConfiguration(test.cfg, logger)
		assert.Equal(t, test.prefix, scopeOptions.Prefix)
		if test.prometheusreporter {
			assert.Empty(t, scopeOptions.Reporter)
			assert.NotNil(t, metricsHandler)
			assert.NotEmpty(t, scopeOptions.CachedReporter)
			assert.NotEmpty(t, scopeOptions.SanitizeOptions)
		} else {
			assert.Nil(t, metricsHandler)
			assert.Empty(t, scopeOptions.CachedReporter)
			assert.Empty(t, scopeOptions.SanitizeOptions)
			if test.nullReporter {
				assert.Empty(t, scopeOptions.Reporter)
			} else {
				assert.NotEmpty(t, scopeOptions.Reporter)
			}
		}
	}
}
