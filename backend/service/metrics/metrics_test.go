package metrics

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortlink(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		host     string
		start    int64
		end      int64
		step     int64
		expected *http.Request
	}{
		{
			name:  "test default",
			query: "avg(rate(container_cpu_usage_seconds_total{container_name!=\"POD\"}[1m]))",
			host:  "localhost",
			start: 1546272000000,
			end:   1546358400000,
			step:  60,
			expected: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "localhost",
					Path:   "/api/v1/query_range",
				},
			},
		},
		{
			name:  "test with step 0",
			query: "avg(rate(container_cpu_usage_seconds_total{container_name!=\"POD\"}[1m]))",
			host:  "foo",
			start: 1601097600000,
			end:   1601184000000,
			step:  0,
			expected: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "foo",
					Path:   "/api/v1/query_range",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := constructRequest(test.query, test.host, test.start, test.end, test.step)
			assert.Nil(t, err)
			assert.Equal(t, test.expected.URL.Host, req.URL.Host)
			assert.Equal(t, test.expected.URL.Path, req.URL.Path)
			assert.Equal(t, test.expected.Method, req.Method)
			assert.Equal(t, test.expected.URL.Scheme, req.URL.Scheme)
			if test.step == 0 {
				assert.Equal(t, "60", req.FormValue("step"))
			}
		})
	}
}
