package gateway

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v3"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	"github.com/lyft/clutch/backend/middleware/timeouts"
)

type Flags struct {
	ConfigPath string
	Template   bool
	Validate   bool
}

// Link register the struct vars globally for parsing by the flag library.
func (f *Flags) Link() {
	flag.StringVar(&f.ConfigPath, "c", "clutch-config.yaml", "path to YAML configuration")
	flag.BoolVar(&f.Template, "template", false, "executes go templates on the configuration file")
	flag.BoolVar(&f.Validate, "validate", false, "validates the configuration file and exits")
}

// Parse command line arguments.
func ParseFlags() *Flags {
	f := &Flags{}
	f.Link()
	flag.Parse()
	return f
}

func MustReadOrValidateConfig(f *Flags) *gatewayv1.Config {
	// Use a temporary logger to parse the configuration and output.
	tmpLogger := newTmpLogger().With(zap.String("file", f.ConfigPath))

	var cfg gatewayv1.Config
	if err := parseFile(f.ConfigPath, &cfg, f.Template); err != nil {
		tmpLogger.Fatal("parsing configuration failed", zap.Error(err))
	}
	if err := cfg.Validate(); err != nil {
		tmpLogger.Fatal("validating configuration failed", zap.Error(err))
	}

	if f.Validate {
		tmpLogger.Info("configuration validation was successful")
		os.Exit(0)
	}

	return &cfg
}

func executeTemplate(contents []byte) ([]byte, error) {
	tmpl := template.New("config").Funcs(map[string]interface{}{
		"getenv": os.Getenv,
		"getboolenv": func(key string) bool {
			b, _ := strconv.ParseBool(os.Getenv(key))
			return b
		},
	})

	tmpl, err := tmpl.Parse(string(contents))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, nil); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func parseFile(path string, pb proto.Message, template bool) error {
	// Get absolute path representation for better error message in case file not found.
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Read file.
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Execute templates if enabled.
	if template {
		contents, err = executeTemplate(contents)
		if err != nil {
			return err
		}
	}

	// Interpolate environment variables
	expandedData := os.ExpandEnv(string(contents))

	contents = []byte(processTemplateToken(expandedData))

	return parseYAML(contents, pb)
}

// Replace the Clutch-specific templating token with the Go $ template syntax
func processTemplateToken(data string) string {
	return strings.ReplaceAll(data, "%%", "$")
}

func parseYAML(contents []byte, pb proto.Message) error {
	// Decode YAML.
	var rawConfig map[string]interface{}
	if err := yaml.Unmarshal(contents, &rawConfig); err != nil {
		return err
	}

	// Encode YAML to JSON.
	rawJSON, err := json.Marshal(rawConfig)
	if err != nil {
		return err
	}

	// Unmarshal JSON to proto object.
	if err := protojson.Unmarshal(rawJSON, pb); err != nil {
		return err
	}

	// All good!
	return nil
}

func newLogger(msg *gatewayv1.Logger) (*zap.Logger, error) {
	var c zap.Config
	var opts []zap.Option
	if msg.GetPretty() {
		c = zap.NewDevelopmentConfig()
		opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	} else {
		c = zap.NewProductionConfig()
	}

	level := zap.NewAtomicLevel()

	levelName := "INFO"
	if msg.Level != gatewayv1.Logger_UNSPECIFIED {
		levelName = msg.Level.String()
	}

	if err := level.UnmarshalText([]byte(levelName)); err != nil {
		return nil, fmt.Errorf("could not parse log level %s", msg.Level.String())
	}
	c.Level = level

	return c.Build(opts...)
}

func newTmpLogger() *zap.Logger {
	c := zap.NewProductionConfig()
	c.DisableStacktrace = true
	l, err := c.Build()
	if err != nil {
		panic(err)
	}
	return l
}

type validator interface {
	Validate() error
}

// Returns maximum timeout, where 0 is considered maximum (i.e. no timeout).
func computeMaximumTimeout(cfg *gatewayv1.Timeouts) time.Duration {
	if cfg == nil {
		return timeouts.DefaultTimeout
	}

	ret := cfg.Default.AsDuration()
	for _, e := range cfg.Overrides {
		override := e.Timeout.AsDuration()
		if ret == 0 || override == 0 {
			return 0
		}

		if override > ret {
			ret = override
		}
	}

	return ret
}

func validateAny(a *anypb.Any) error {
	if a == nil {
		return nil
	}

	m, err := a.UnmarshalNew()
	if err != nil {
		return err
	}

	if v, ok := m.(validator); ok {
		return v.Validate()
	}
	return nil
}
