package mux

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	"github.com/stretchr/testify/assert"
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
					Bucket: "clutch",
					Key:    "static",
				},
			},
		},
	}

	// Test that the aws service must be configured to use the S3 handler
	_, err := handler.assetProviderHandler("clutch.sh/static/main.js")
	assert.Error(t, err)
}
