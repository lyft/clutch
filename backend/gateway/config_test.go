package gateway

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/protobuf/types/known/durationpb"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	"github.com/lyft/clutch/backend/middleware/timeouts"
)

func TestEnsureUnique(t *testing.T) {
	tests := []struct {
		c   *gatewayv1.Config
		err error
	}{
		{
			c: &gatewayv1.Config{
				Services: []*gatewayv1.Service{
					{Name: "foo"},
					{Name: "foo"},
				},
			},
			err: fmt.Errorf("duplicate service found: foo"),
		},
		{
			c: &gatewayv1.Config{
				Modules: []*gatewayv1.Module{
					{Name: "foo"},
					{Name: "foo"},
				},
			},
			err: fmt.Errorf("duplicate module found: foo"),
		},
		{
			c: &gatewayv1.Config{
				Resolvers: []*gatewayv1.Resolver{
					{Name: "foo"},
					{Name: "foo"},
				},
			},
			err: fmt.Errorf("duplicate resolver found: foo"),
		},
		{
			c: &gatewayv1.Config{
				Services:  []*gatewayv1.Service{{Name: "foo"}},
				Modules:   []*gatewayv1.Module{{Name: "foo"}},
				Resolvers: []*gatewayv1.Resolver{{Name: "foo"}},
			},
			err: nil,
		},
	}

	for idx, tt := range tests {
		tc := tt // Pin!
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			err := ensureUnique(tc.c)
			assert.Equal(t, tc.err, err)
		})
	}
}

func tmpFile(filename, content string) *os.File {
	f, err := os.CreateTemp(".", filename)
	if err != nil {
		log.Panic(err)
	}
	_ = os.WriteFile(f.Name(), []byte(content), 0644)
	return f
}

func TestConsolidateConfigs(t *testing.T) {
	baseConfig := `
gateway:
  logger:
    pretty: true
    level: DEBUG
`
	bc := tmpFile("base-config-*.yaml", baseConfig)
	defer os.Remove(bc.Name())

	config := fmt.Sprintf(`
extends:
  - %s
gateway:
  listener:
    tcp:
      address: 0.0.0.0
      port: 8080
      secure: false
`, bc.Name())

	cc := tmpFile("clutch-config-*.yaml", config)
	defer os.Remove(cc.Name())

	var cfg gatewayv1.Config
	var seenCfgs []string
	consolidateConfigs(filepath.Dir(cc.Name()), filepath.Base(cc.Name()), &cfg, &Flags{Template: false}, &seenCfgs)

	assert.Equal(t, true, cfg.GetGateway().GetLogger().GetPretty())
	assert.Equal(t, gatewayv1.Logger_DEBUG, cfg.GetGateway().GetLogger().GetLevel())

	assert.Equal(t, "0.0.0.0", cfg.GetGateway().GetListener().GetTcp().GetAddress())
	assert.Equal(t, uint32(8080), cfg.GetGateway().GetListener().GetTcp().GetPort())
	assert.Equal(t, false, cfg.GetGateway().GetListener().GetTcp().GetSecure())
}

func TestConsolidateConfigsOverrides(t *testing.T) {
	baseConfig := `
gateway:
  logger:
    pretty: true
    level: DEBUG
`
	bc := tmpFile("base-config-*.yaml", baseConfig)
	defer os.Remove(bc.Name())

	logConfig := fmt.Sprintf(`
extends:
  - %s
gateway:
  logger:
    pretty: false
`, bc.Name())
	lc := tmpFile("log-config-*.yaml", logConfig)
	defer os.Remove(lc.Name())

	config := fmt.Sprintf(`
extends:
  - %s
gateway:
  logger:
    level: WARN
`, lc.Name())

	cc := tmpFile("clutch-config-*.yaml", config)
	defer os.Remove(cc.Name())

	var cfg gatewayv1.Config
	var seenCfgs []string
	consolidateConfigs(filepath.Dir(cc.Name()), filepath.Base(cc.Name()), &cfg, &Flags{Template: false}, &seenCfgs)

	assert.Equal(t, false, cfg.GetGateway().GetLogger().GetPretty())
	assert.Equal(t, gatewayv1.Logger_WARN, cfg.GetGateway().GetLogger().GetLevel())
}

func TestConsolidateConfigsIgnoresDuplicateConfigs(t *testing.T) {
	cc := tmpFile("clutch-config-*.yaml", "")
	defer os.Remove(cc.Name())
	bc := tmpFile("base-config-*.yaml", "")
	defer os.Remove(bc.Name())

	baseConfig := fmt.Sprintf(`
extends:
  - %s
  - %s
gateway:
  logger:
    pretty: true
    level: DEBUG
`, cc.Name(), bc.Name())
	_ = os.WriteFile(bc.Name(), []byte(baseConfig), 0644)

	config := fmt.Sprintf(`
extends:
  - %s
gateway:
  logger:
    pretty: true
    level: WARN
`, bc.Name())
	_ = os.WriteFile(cc.Name(), []byte(config), 0644)

	var cfg gatewayv1.Config
	var seenCfgs []string
	consolidateConfigs(filepath.Dir(cc.Name()), filepath.Base(cc.Name()), &cfg, &Flags{Template: false}, &seenCfgs)

	assert.Equal(t, true, cfg.GetGateway().GetLogger().GetPretty())
	assert.Equal(t, gatewayv1.Logger_WARN, cfg.GetGateway().GetLogger().GetLevel())
}

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
		{
			Level:     gatewayv1.Logger_WARN,
			Namespace: "test",
		},
	}

	for idx, tc := range testConfigs {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			l, err := newLogger(tc)
			assert.NotNil(t, l)
			assert.NoError(t, err)
		})
	}
}

func TestNewLoggerNamespace(t *testing.T) {
	core, recorder := observer.New(zapcore.InfoLevel)

	cfg := &gatewayv1.Logger{
		Level:     gatewayv1.Logger_INFO,
		Namespace: "cat",
	}

	l, err := newLoggerWithCore(cfg, core)
	assert.NotNil(t, l)
	assert.NoError(t, err)

	l.Info("meow", zap.String("fields", "zapField"))

	assert.Equal(t, 1, recorder.Len())
	allLogs := recorder.All()
	assert.Equal(t, "meow", allLogs[0].Message)
	assert.ElementsMatch(t, []zap.Field{
		zap.Namespace("cat"),
		zap.String("fields", "zapField"),
	}, allLogs[0].Context)
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
		// Pin!
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
