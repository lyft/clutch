package log

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
)

func TestProtoField(t *testing.T) {
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

func TestNamedErrorField(t *testing.T) {
	// a Status with no detailed appended
	s1 := status.New(codes.PermissionDenied, "Permission denied")
	err1 := s1.Err()

	// a Status with details appended
	s2 := status.New(codes.NotFound, "Resource not found")
	s2, _ = s2.WithDetails(
		&errdetails.ResourceInfo{ResourceType: "ConfigMap", Description: "configMap-test-1 not found"},
		&errdetails.ResourceInfo{ResourceType: "ConfigMap", Description: "configMap-test-2 not found"},
	)
	err2 := s2.Err()

	tests := []struct {
		err             error
		expectedMsg     string
		expectedCode    int
		expectedDetails []string
	}{
		{
			err:         errors.New("yikes"),
			expectedMsg: "yikes",
		},
		{
			err:          err1,
			expectedMsg:  "Permission denied",
			expectedCode: 7,
		},
		{
			err:             err2,
			expectedCode:    5,
			expectedMsg:     "Resource not found",
			expectedDetails: []string{"configMap-test-1 not found", "configMap-test-2 not found"},
		},
	}

	for _, test := range tests {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)

		logger := zap.New(
			zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(w), zap.DebugLevel),
		)
		logger.Info("test", NamedErrorField("key", test.err))
		assert.NoError(t, logger.Sync())
		assert.NoError(t, w.Flush())

		o := make(map[string]interface{})
		assert.NoError(t, json.Unmarshal(b.Bytes(), &o))
		assert.Contains(t, b.String(), test.expectedMsg)

		if test.expectedCode != 0 {
			assert.Contains(t, b.String(), strconv.Itoa(test.expectedCode))
		}
		if test.expectedDetails != nil {
			for _, detail := range test.expectedDetails {
				assert.Contains(t, b.String(), detail)
			}
		}
	}
}

func TestErrorField(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(w), zap.DebugLevel),
	)

	logger.Info("test", ErrorField(errors.New("yikes")))
	assert.NoError(t, logger.Sync())
	assert.NoError(t, w.Flush())

	o := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(b.Bytes(), &o))
	assert.Contains(t, b.String(), "error")
	assert.Contains(t, b.String(), "yikes")
}
