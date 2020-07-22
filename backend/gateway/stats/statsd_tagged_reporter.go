package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/uber-go/tally"
	tallystatsd "github.com/uber-go/tally/statsd"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
)

func NewStatsdReporter(cfg *gatewayv1.Stats_StatsdReporter) (tally.StatsReporter, error) {
	client, err := statsd.NewClientWithConfig(&statsd.ClientConfig{
		Address: cfg.Address,
	})
	if err != nil {
		return nil, err
	}

	baseReporter := tallystatsd.NewReporter(client, tallystatsd.Options{})

	switch cfg.TagMode.(type) {
	case nil:
		return baseReporter, nil
	case *gatewayv1.Stats_StatsdReporter_PointTags_:
		return &pointTagReporter{
			StatsReporter: baseReporter,
			separator:     cfg.GetPointTags().Separator,
		}, nil
	default:
		return nil, fmt.Errorf("unspported tag mode for statsd reporter")
	}
}

type pointTagReporter struct {
	tally.StatsReporter

	separator string
}

// https://github.com/influxdata/telegraf/blob/master/plugins/inputs/statsd/README.md#influx-statsd
func (r *pointTagReporter) taggedName(name string, tags map[string]string) string {
	var b strings.Builder
	b.WriteString(name)
	for k, v := range tags {
		b.WriteString(r.separator)
		b.WriteString(replaceChars(k))
		b.WriteByte('=')
		b.WriteString(replaceChars(v))
	}
	return b.String()
}

func (r *pointTagReporter) ReportCounter(name string, tags map[string]string, value int64) {
	r.StatsReporter.ReportCounter(r.taggedName(name, tags), nil, value)
}

func (r *pointTagReporter) ReportGauge(name string, tags map[string]string, value float64) {
	r.StatsReporter.ReportGauge(r.taggedName(name, tags), nil, value)
}

func (r *pointTagReporter) ReportTimer(name string, tags map[string]string, interval time.Duration) {
	r.StatsReporter.ReportTimer(r.taggedName(name, tags), nil, interval)
}

func (r *pointTagReporter) ReportHistogramValueSamples(name string, tags map[string]string, buckets tally.Buckets, bucketLowerBound, bucketUpperBound float64, samples int64) {
	r.StatsReporter.ReportHistogramValueSamples(r.taggedName(name, tags), nil, buckets, bucketLowerBound, bucketUpperBound, samples)
}

func (r *pointTagReporter) ReportHistogramDurationSamples(name string, tags map[string]string, buckets tally.Buckets, bucketLowerBound, bucketUpperBound time.Duration, samples int64) {
	r.StatsReporter.ReportHistogramDurationSamples(r.taggedName(name, tags), nil, buckets, bucketLowerBound, bucketUpperBound, samples)
}

// Replace problematic characters in tags.
func replaceChars(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.', ':', '|', '-', '=':
			b.WriteByte('_')
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
}
