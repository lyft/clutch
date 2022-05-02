package metrics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

	metricsv1cfg "github.com/lyft/clutch/backend/api/config/service/metrics/v1"
	metricsv1 "github.com/lyft/clutch/backend/api/metrics/v1"
	"github.com/lyft/clutch/backend/service"
)

const (
	Name           = "clutch.service.metrics"
	QueryRangePath = "/api/v1/query_range"
	QueryScheme    = "http"
	StepDefaultMs  = 60000
)

type Service interface {
	GetMetrics(context.Context, *metricsv1.GetMetricsRequest) (*metricsv1.GetMetricsResponse, error)
}
type PrometheusResponse struct {
	Status string
	Data   PrometheusResponseData
}

type PrometheusResponseData struct {
	ResultType string
	Result     []PrometheusResult
}

type PrometheusResult struct {
	Metric map[string]string
	Values []PrometheusValue
}

type PrometheusValue struct {
	timestamp int64
	value     string
}

type client struct {
	prometheusAPI         *http.Client
	prometheusAPIEndpoint string
	log                   *zap.Logger
	scope                 tally.Scope
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	return NewWithHTTPClient(cfg, logger, scope, &http.Client{})
}

func NewWithHTTPClient(cfg *any.Any, logger *zap.Logger, scope tally.Scope, prometheusAPI *http.Client) (service.Service, error) {
	config := &metricsv1cfg.Config{}
	if err := cfg.UnmarshalTo(config); err != nil {
		return nil, err
	}

	if config.PrometheusApiEndpoint == "" {
		return nil, fmt.Errorf("prometheus api endpoint is required")
	}

	return &client{
		prometheusAPI:         prometheusAPI,
		prometheusAPIEndpoint: config.PrometheusApiEndpoint,
		log:                   logger,
		scope:                 scope,
	}, nil
}

func constructRequest(query, host string, start, end int64, step int64) (*http.Request, error) {
	baseURL := url.URL{Scheme: QueryScheme, Host: host, Path: QueryRangePath}
	params := url.Values{}

	params.Add("query", query)
	// We are given timestamps in milliseconds, but Prometheus expects seconds.
	params.Add("start", strconv.FormatInt(start/1000, 10))
	params.Add("end", strconv.FormatInt(end/1000, 10))

	// If step is 0, update it to be 1 minute
	if step == 0 {
		step = StepDefaultMs
	}
	params.Add("step", strconv.FormatInt(step/1000, 10))

	baseURL.RawQuery = params.Encode()
	url := baseURL.String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// TODO(smonero): Do we need these?
	//	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

func (c *client) formatPrometheusResponseData(response PrometheusResponseData) ([]*metricsv1.Metrics, error) {
	metricsList := make([]*metricsv1.Metrics, 0)

	for _, r := range response.Result {
		var timeSeries metricsv1.Metrics

		for k, v := range r.Metric {
			if k == "__name__" {
				timeSeries.Label = v
			} else {
				timeSeries.Tags[k] = v
			}
		}

		for _, s := range r.Values {
			var datapoint metricsv1.MetricDataPoint

			datapoint.Timestamp = s.timestamp

			v, err := strconv.ParseFloat(s.value, 32)
			if err != nil {
				return nil, err
			}
			if math.IsInf(v, 0) || math.IsNaN(v) {
				// While +/-Inf and NaN are valid floating point values, they are
				// not useful to our alerting system.
				continue
			}
			datapoint.Value = float64(v)

			timeSeries.DataPoints = append(timeSeries.DataPoints, &datapoint)
		}

		metricsList = append(metricsList, &timeSeries)
	}

	return metricsList, nil
}

// The expected data from the Prometheus api is
// [ [ <unix_time>, "<sample_value>" ], ... ] where sample_value is numeric. Go doesn't support
// slices of mixed type, so we do the un-marshaling ourselves.
func (v *PrometheusValue) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	data = bytes.Trim(data, "[]")
	parts := bytes.Split(data, []byte(","))

	// Check for two parts.
	if err := json.Unmarshal(parts[0], &v.timestamp); err != nil {
		return err
	}
	if err := json.Unmarshal(parts[1], &v.value); err != nil {
		return err
	}

	return nil
}

/**
 * GetMetrics returns the metrics for the given prometheus queries.
 * For each query, it uses `query_range` as defined:
 * https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries
 * It goes through the following steps:
 * 1. Construct the request
 * 2. Send the request
 * 3. Parse the response
 */
func (c *client) GetMetrics(ctx context.Context, req *metricsv1.GetMetricsRequest) (*metricsv1.GetMetricsResponse, error) {
	queryResults := make(map[string]*metricsv1.MetricsResult)
	for _, q := range req.MetricQueries {
		req, err := constructRequest(q.Expression, c.prometheusAPIEndpoint, q.StartTimeMs, q.EndTimeMs, q.StepMs)
		if err != nil {
			return nil, err
		}

		resp, err := c.prometheusAPI.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf(resp.Status)
		}

		var queryResponse PrometheusResponse
		err = json.Unmarshal(bodyBytes, &queryResponse)
		if err != nil {
			return nil, err
		}

		//TODO: check status?

		qResult, err := c.formatPrometheusResponseData(queryResponse.Data)
		if err != nil {
			return nil, err
		}
		queryResults[q.Expression] = &metricsv1.MetricsResult{
			Metrics: qResult,
		}
	}

	return &metricsv1.GetMetricsResponse{
		QueryResults: queryResults,
	}, nil
}
