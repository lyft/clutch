package metricsmock

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

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
	randNums := make([]float64, 0)
	for i := 0; i < 5; i++ {
		ran, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			return nil, fmt.Errorf("failed to generate random number: %s", err)
		}
		randNums = append(randNums, float64(ran.Int64()))
	}

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
								Value:     randNums[0],
								Timestamp: 123456789,
							},
							{
								Value:     randNums[1],
								Timestamp: 123456890,
							},
							{
								Value:     randNums[2],
								Timestamp: 123457000,
							},
							{
								Value:     randNums[3],
								Timestamp: 123457100,
							},
							{
								Value:     randNums[4],
								Timestamp: 123457200,
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
