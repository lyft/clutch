package gateway

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
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
