package {{ .ModuleName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap/zaptest"

	"github.com/lyft/clutch/backend/module/moduletest"
	"github.com/lyft/clutch/backend/service"
{{ if .ServiceName }}
	{{ .ServiceName }}svc "{{ .BackendGoPackage }}/service/{{ .ServiceName }}"
{{ end }}
)

func TestModule(t *testing.T) {
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
{{ if .ServiceName }}
	// If service is not unit testable as-is, a mock service will be needed
	svc, err := {{ .ServiceName }}svc.New(nil, log, scope)
	assert.NoError(t, err)
	service.Registry[{{ .ServiceName }}svc.Name] = svc
{{ end }}
	m, err := New(nil, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("{{ .RepoOwner }}.{{ .ModuleName }}.{{ .ApiVersion }}.{{ .ApiName }}API"))
	assert.True(t, r.JSONRegistered())
}
