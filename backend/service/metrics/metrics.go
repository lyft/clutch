package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

	"net/http"
	"net/url"

	metricsv1cfg "github.com/lyft/clutch/backend/api/config/service/metrics/v1"
	metricsv1 "github.com/lyft/clutch/backend/api/metrics/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.metrics"

type Service interface {
	GetMetrics(context.Context, *metricsv1.GetMetricsRequest) (*metricsv1.GetMetricsResponse, error)
}

const ()

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

	return &client{
		prometheusAPI:         prometheusAPI,
		prometheusAPIEndpoint: config.PrometheusApiEndpoint,
		log:                   logger,
		scope:                 scope,
	}, nil
}

func (c *client) constructRequest(query string, start, end int64, step int64) (*http.Request, error) {
	baseURL := url.URL{Scheme: "http", Host: c.prometheusAPIEndpoint, Path: "/api/v1/query_range"}
	params := url.Values{}

	params.Add("query", query)
	params.Add("start", strconv.FormatInt(start/1000, 10))
	params.Add("end", strconv.FormatInt(end/1000, 10))
	params.Add("step", strconv.FormatInt(step/1000, 10))

	baseURL.RawQuery = params.Encode()
	url := baseURL.String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// TODO(smonero): Do we need these?
	//	req.Header.Set("M3-Host", string(m3.m3Host))
	//	req.Header.Set("User-Agent", userAgent)

	return req, nil
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
	Values []PrometheusDataPoint `json:"values"`
}
type PrometheusDataPoint struct {
	timestamp int64
	value     string
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

	queryResults := make(map[string]*metricsv1.Metrics)
	for _, q := range req.MetricQueries {

		req, err := c.constructRequest(q.Expression, q.StartTimeMs, q.EndTimeMs, q.StepMs)
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

		err = json.Unmarshal(bodyBytes, &queryResponse)

		resStr := nil

		queryResults[q.Expression] = resStr
	}

	return &metricsv1.GetMetricsResponse{
		QueryResults: queryResults,
	}, nil
}
