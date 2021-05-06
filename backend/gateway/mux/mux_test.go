package mux

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
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

func TestCustomHeaderMatcher(t *testing.T) {
	testCases := []struct {
		key          string
		expectedKey  string
		expectedBool bool
	}{
		{
			key:          "X-Foo-Bar",
			expectedKey:  "grpcgateway-X-Foo-Bar",
			expectedBool: true,
		},
		// testing that the headers get uppercased
		{
			key:          "x-foo-bar",
			expectedKey:  "grpcgateway-X-Foo-Bar",
			expectedBool: true,
		},
		// testing the default rule - isPermanentHTTPHeader group
		{
			key:          "Cookie",
			expectedKey:  "grpcgateway-Cookie",
			expectedBool: true,
		},
		// testing the default rule - Grpc-Metadata prefix
		{
			key:          "Grpc-Metadata-Foo",
			expectedKey:  "Foo",
			expectedBool: true,
		},
		// testing the prefix doesn't get applied and doesn't match default rule
		{
			key:          xForwardedFor,
			expectedKey:  "",
			expectedBool: false,
		},
		// testing the prefix doesn't get applied and doesn't match default rule
		{
			key:          xForwardedHost,
			expectedKey:  "",
			expectedBool: false,
		},
		// doesn't match custom or default rules
		{
			key:          "Foo-Bar",
			expectedKey:  "",
			expectedBool: false,
		},
	}

	for _, test := range testCases {
		result, ok := customHeaderMatcher(test.key)
		assert.Equal(t, test.expectedKey, result)
		assert.Equal(t, test.expectedBool, ok)
	}
}

func TestCustomErrorHandler(t *testing.T) {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	marshaler := &runtime.JSONPb{}

	{
		// Business as usual.
		req := &http.Request{}
		rec := httptest.NewRecorder()
		w := mockResponseWriter{ResponseWriter: rec}
		err := status.Error(codes.NotFound, "not found")
		customErrorHandler(ctx, mux, marshaler, w, req, err)
		assert.Equal(t, 404, rec.Code)
	}
	{
		// Auth redirect for browser 401.
		uri := "https://example.com/bar?foo=bar"
		req, _ := http.NewRequest("GET", uri, nil)
		req.RequestURI = uri
		req.Header.Add("Accept", "text/html")
		rec := httptest.NewRecorder()
		w := mockResponseWriter{ResponseWriter: rec}
		err := status.Error(codes.Unauthenticated, "not found")
		customErrorHandler(ctx, mux, marshaler, w, req, err)
		assert.Equal(t, 302, rec.Code)
		assert.Contains(t, rec.Header().Get("Location"), url.QueryEscape(uri))
	}
	{
		// No auth redirect for non-browser 401.
		uri := "https://example.com/bar?foo=bar"
		req, _ := http.NewRequest(http.MethodGet, uri, nil)
		req.RequestURI = uri
		rec := httptest.NewRecorder()
		w := mockResponseWriter{ResponseWriter: rec}
		err := status.Error(codes.Unauthenticated, "not found")
		customErrorHandler(ctx, mux, marshaler, w, req, err)
		assert.Equal(t, 401, rec.Code)
	}
}

func TestCustomResponseForwarder(t *testing.T) {
	ctx := runtime.NewServerMetadataContext(context.Background(), runtime.ServerMetadata{})
	rec := httptest.NewRecorder()
	w := mockResponseWriter{ResponseWriter: rec}
	err := customResponseForwarder(ctx, w, &healthcheckv1.HealthcheckResponse{})
	assert.NoError(t, err)
	assert.Equal(t, 200, rec.Code)
}

func TestCustomResponseForwarderAuthCookies(t *testing.T) {
	ctx := runtime.NewServerMetadataContext(context.Background(), runtime.ServerMetadata{
		HeaderMD: metadata.Pairs(
			"Set-Cookie-Token", "myToken",
			"Location", "https://example.com",
		),
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/bar", nil)
	req.Header.Add("Accept", "text/html") // Is browser.
	w := &mockResponseWriter{ResponseWriter: rec, req: req}
	err := customResponseForwarder(ctx, w, &healthcheckv1.HealthcheckResponse{})
	assert.NoError(t, err)
	assert.Equal(t, 302, rec.Code)
	assert.Equal(t, "token=myToken; Path=/", rec.Header().Get("Set-Cookie"))
	assert.Equal(t, "https://example.com", rec.Header().Get("Location"))
}

func TestCustomResponseForwarderLocationStatusOverrideAndRefreshToken(t *testing.T) {
	ctx := runtime.NewServerMetadataContext(context.Background(), runtime.ServerMetadata{
		HeaderMD: metadata.Pairs(
			"Set-Cookie-Token", "myToken",
			"Set-Cookie-Refresh-Token", "myRefreshToken",
			"Location", "https://example.com",
			"Location-Status", "304",
		),
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/bar", nil)
	req.Header.Add("Accept", "text/html") // Is browser.
	w := &mockResponseWriter{ResponseWriter: rec, req: req}
	err := customResponseForwarder(ctx, w, &healthcheckv1.HealthcheckResponse{})
	assert.NoError(t, err)
	assert.Equal(t, 304, rec.Code)
	assert.Contains(t, rec.Header().Values("Set-Cookie"), "token=myToken; Path=/")
	assert.Contains(t, rec.Header().Values("Set-Cookie"), "refreshToken=myRefreshToken; Path=/v1/authn/login; HttpOnly")
	assert.Equal(t, "https://example.com", rec.Header().Get("Location"))
}

func TestCustomResponseForwarderAuthCookiesNonBrowser(t *testing.T) {
	ctx := runtime.NewServerMetadataContext(context.Background(), runtime.ServerMetadata{
		HeaderMD: metadata.Pairs(
			"Location", "https://example.com",
		),
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/bar", nil)
	req.Header.Add("Accept", "*/*") // Not a browser!
	w := &mockResponseWriter{ResponseWriter: rec, req: req}
	err := customResponseForwarder(ctx, w, &healthcheckv1.HealthcheckResponse{})
	assert.NoError(t, err)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "", rec.Header().Get("Location"))
}
