package metricsmock

import (
	"context"
	"math/rand"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

	metricsv1 "github.com/lyft/clutch/backend/api/metrics/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/metrics"
)

type svc struct{}

func New() metrics.Service {
	return &svc{}
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}

func (s *svc) GetMetrics(ctx context.Context, req *metricsv1.GetMetricsRequest) (*metricsv1.GetMetricsResponse, error) {
	return &metricsv1.GetMetricsResponse{
		QueryResults: map[string]*metricsv1.MetricsResult{
			"req": {
				Metrics: []*metricsv1.Metrics{
					{
						DataPoints: []*metricsv1.MetricDataPoint{
							{
								Value:     10,
								Timestamp: 123456789,
							},
							{
								Value:     30,
								Timestamp: 123456890,
							},
							{
								Value:     20,
								Timestamp: 123457000,
							},
						},
						Label: "foo",
						Tags: map[string]string{
							"meow": "woof",
						},
					},
				},
			},
			"random": {
				Metrics: []*metricsv1.Metrics{
					{
						DataPoints: []*metricsv1.MetricDataPoint{
							{
								Value:     float64(rand.Intn(100)),
								Timestamp: 123456789,
							},
							{
								Value:     float64(rand.Intn(100)),
								Timestamp: 123456890,
							},
							{
								Value:     float64(rand.Intn(100)),
								Timestamp: 123457000,
							},
						},
						Label: "random",
						Tags: map[string]string{
							"cat sounds": "dog sounds",
						},
					},
				},
			},
		},
	}, nil
}
