package gateway

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
