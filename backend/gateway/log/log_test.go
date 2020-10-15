package log

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
)

func Test(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(w), zap.DebugLevel),
	)

	r := &healthcheckv1.HealthcheckRequest{}
	a, _ := ptypes.MarshalAny(r)

	logger.Info("test", ProtoField("key", a))
	assert.NoError(t, logger.Sync())
	assert.NoError(t, w.Flush())

	o := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(b.Bytes(), &o))
	assert.Contains(t, b.String(), `{"@type":"type.googleapis.com/clutch.healthcheck.v1.HealthcheckRequest"}`)
}
