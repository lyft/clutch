package gateway

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	"github.com/lyft/clutch/backend/middleware/timeouts"
)

func TestExecuteTemplate(t *testing.T) {
	config := `
foo: bar
options:
  - yes
{{- if (getboolenv "ENABLE_CHOICE") }}
  - no
  - maybe
{{ end }}
{{- if eq (getenv "SPEED") "walk" }}
use: shoes
{{ end }}
`
	out, err := executeTemplate([]byte(config))
	assert.NoError(t, err)
	assert.NotContains(t, string(out), "- no")
	assert.NotContains(t, string(out), "- maybe")

	os.Setenv("ENABLE_CHOICE", "true")
	out, err = executeTemplate([]byte(config))
	assert.NoError(t, err)
	assert.Contains(t, string(out), "  - yes\n  - no\n  - maybe")

	os.Setenv("SPEED", "run")
	out, err = executeTemplate([]byte(config))
	assert.NoError(t, err)
	assert.NotContains(t, string(out), "use: shoes")

	os.Setenv("SPEED", "walk")
	out, err = executeTemplate([]byte(config))
	assert.NoError(t, err)
	assert.Contains(t, string(out), "use: shoes")
}

func TestNewLogger(t *testing.T) {
	testConfigs := []*gatewayv1.Logger{
		{
			Level:  gatewayv1.Logger_INFO,
			Format: &gatewayv1.Logger_Pretty{Pretty: true},
		},
		{
			Level: gatewayv1.Logger_WARN,
		},
	}

	for idx, tc := range testConfigs {
		tc := tc
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			l, err := newLogger(tc)
			assert.NotNil(t, l)
			assert.NoError(t, err)
		})
	}
}

func TestComputeMaximumTimeout(t *testing.T) {
	tests := []struct {
		c        *gatewayv1.Timeouts
		expected time.Duration
	}{
		{
			c:        nil,
			expected: timeouts.DefaultTimeout,
		},
		{
			c:        &gatewayv1.Timeouts{Default: durationpb.New(0)},
			expected: 0,
		},
		{
			c:        &gatewayv1.Timeouts{Default: durationpb.New(time.Second)},
			expected: time.Second,
		},
		{
			c: &gatewayv1.Timeouts{
				Default: durationpb.New(time.Second),
				Overrides: []*gatewayv1.Timeouts_Entry{
					{Timeout: durationpb.New(time.Second * 10)},
				},
			},
			expected: 10 * time.Second,
		},
		{
			c: &gatewayv1.Timeouts{
				Default: durationpb.New(time.Second),
				Overrides: []*gatewayv1.Timeouts_Entry{
					{Timeout: durationpb.New(0)},
				},
			},
			expected: 0,
		},
		{
			c: &gatewayv1.Timeouts{
				Default: durationpb.New(time.Second),
				Overrides: []*gatewayv1.Timeouts_Entry{
					{Timeout: durationpb.New(time.Millisecond)},
				},
			},
			expected: time.Second,
		},
	}

	for idx, tt := range tests {
		tt := tt // Pin!
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			result := computeMaximumTimeout(tt.c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBulkReplaceTemplateTokens(t *testing.T) {
	config := `
	foo: bar
	message: [[range $$v, $$k := .Bar]][[$$k]]: [[$$v]][[end]]
	`
	expected := `
	foo: bar
	message: {{range @#@v, @#@k := .Bar}}{{@#@k}}: {{@#@v}}{{end}}
	`
	contents := bulkReplaceTemplateTokens(config)
	assert.Equal(t, expected, contents)
}

func TestReplaceVarTemplateToken(t *testing.T) {
	config := `
	foo: bar
	message: {{range @#@v, @#@k := .Bar}}{{@#@k}}: {{@#@v}}{{end}}
	`

	expected := `
	foo: bar
	message: {{range $v, $k := .Bar}}{{$k}}: {{$v}}{{end}}
	`

	contents := replaceVarTemplateToken(config)
	assert.Equal(t, expected, contents)
}
