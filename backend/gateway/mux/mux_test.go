package mux

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
)

func TestCopyHTTPResponse(t *testing.T) {
	status := http.StatusBadGateway
	headers := http.Header{"Foo": []string{"bar", "baz"}}
	body := "bang"

	resp := &http.Response{
		Status:     http.StatusText(status),
		StatusCode: status,
		Header:     headers,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
	}

	rec := httptest.NewRecorder()
	copyHTTPResponse(resp, rec)
	result := rec.Result()
	assert.Equal(t, status, result.StatusCode)
	assert.Equal(t, headers, rec.Header())
	assert.Equal(t, body, rec.Body.String())
}

func TestAssetProviderS3Handler(t *testing.T) {
	handler := &assetHandler{
		assetCfg: &gatewayv1.Assets{
			Provider: &gatewayv1.Assets_S3{
				S3: &gatewayv1.Assets_S3Provider{
					Region: "us-east-1",
					Bucket: "clutch",
					Key:    "static",
				},
			},
		},
	}

	// Test that the aws service must be configured to use the S3 handler
	_, err := handler.assetProviderHandler(context.TODO(), "clutch.sh/static/main.js")
	assert.Error(t, err)
}

func TestGetAssetProivderService(t *testing.T) {
	assetCfg := &gatewayv1.Assets{
		Provider: &gatewayv1.Assets_S3{
			S3: &gatewayv1.Assets_S3Provider{
				Region: "us-east-1",
				Bucket: "clutch",
				Key:    "static",
			},
		},
	}

	// Test that the aws service must be configured to use the S3 handler
	_, err := getAssetProviderService(assetCfg)
	assert.Error(t, err)
}

func TestCustomMatcher(t *testing.T) {
	testCases := []struct {
		key          string
		expectedKey  string
		expectedBool bool
	}{
		{key: "X-Foo-Bar", expectedKey: "grpcgateway-X-Foo-Bar", expectedBool: true},
		// testing the default rule - a part of the X group
		{key: "Cookie", expectedKey: "grpcgateway-Cookie", expectedBool: true},
		// testing the default rule - has the Grpc-Metadata prefix
		{key: "Grpc-Metadata-Foo", expectedKey: "Foo", expectedBool: true},
		// doesn't match custom or default rules
		{key: "Foo-Bar", expectedKey: "", expectedBool: false},
	}

	for _, test := range testCases {
		result, ok := customMatcher(test.key)
		assert.Equal(t, test.expectedKey, result)
		assert.Equal(t, test.expectedBool, ok)
	}
}
