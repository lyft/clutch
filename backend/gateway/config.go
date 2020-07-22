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
	"text/template"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	durpb "github.com/golang/protobuf/ptypes/duration"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
)

type Flags struct {
	ConfigPath string
	Template   bool
	Validate   bool
}

// Parse command line arguments.
func ParseFlags() *Flags {
	f := &Flags{}
	flag.StringVar(&f.ConfigPath, "c", "clutch-config.yaml", "path to YAML configuration")
	flag.BoolVar(&f.Template, "template", false, "executes go templates on the configuration file")
	flag.BoolVar(&f.Validate, "validate", false, "validates the configuration file and exits")
	flag.Parse()
	return f
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

	// Interpolate environment variables.
	contents = []byte(os.ExpandEnv(string(contents)))

	return parseYAML(contents, pb)
}

func parseYAML(contents []byte, pb proto.Message) error {
	// Decode YAML.
	var rawConfig map[string]interface{}
	if err := yaml.Unmarshal(contents, &rawConfig); err != nil {
		return err
	}

	// Encode YAML to JSON.
	jsonBuffer := new(bytes.Buffer)
	if err := json.NewEncoder(jsonBuffer).Encode(rawConfig); err != nil {
		return err
	}

	// Unmarshal JSON to proto object.
	if err := jsonpb.Unmarshal(jsonBuffer, pb); err != nil {
		return err
	}

	// All good!
	return nil
}

// Helper "must" function to convert proto durations to Go durations. This should only be called during bootstrap where
// a panic is not a problem. Config validation should also prevent the panic from ever occurring.
func duration(p *durpb.Duration) time.Duration {
	d, err := ptypes.Duration(p)
	if err != nil {
		panic(err)
	}
	return d
}

func newLogger(msg *gatewayv1.Logger) (*zap.Logger, error) {
	var c zap.Config
	if msg.GetPretty() {
		c = zap.NewDevelopmentConfig()
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

	return c.Build()
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

func validateAny(a *any.Any) error {
	if a == nil {
		return nil
	}
	var pb ptypes.DynamicAny
	if err := ptypes.UnmarshalAny(a, &pb); err != nil {
		return err
	}
	if v, ok := pb.Message.(validator); ok {
		return v.Validate()
	}
	return nil
}
