package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	proxyv1cfg "github.com/lyft/clutch/backend/api/config/module/proxy/v1"
	proxyv1 "github.com/lyft/clutch/backend/api/proxy/v1"
	"github.com/lyft/clutch/backend/module/moduletest"
)

func generateServicesConfig(host string) []*proxyv1cfg.Service {
	return []*proxyv1cfg.Service{
		{
			Name: "cat",
			Host: host,
			AllowedRequests: []*proxyv1cfg.AllowRequest{
				{Path: "/meow", Method: "GET"},
				{Path: "/nom", Method: "POST"},
			},
		},
		{
			Name: "meow",
			Host: host,
			AllowedRequests: []*proxyv1cfg.AllowRequest{
				{Path: "/meow", Method: "GET"},
				{Path: "/nom", Method: "POST"},
			},
		},
	}
}

func structpbFromBody(body []byte) *structpb.Value {
	var bodyData interface{}
	err := json.NewDecoder(bytes.NewReader(body)).Decode(&bodyData)
	if err != nil {
		panic(err)
	}

	str, err := structpb.NewValue(bodyData)
	if err != nil {
		panic(err)
	}

	return str
}

func TestModule(t *testing.T) {
	config, _ := anypb.New(&proxyv1cfg.Config{})
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	m, err := New(config, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("clutch.proxy.v1.ProxyAPI"))
	assert.True(t, r.JSONRegistered())
}

func TestRequestProxy(t *testing.T) {
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	tests := []struct {
		id             string
		request        *proxyv1.RequestProxyRequest
		handler        func(http.ResponseWriter, *http.Request)
		shouldError    bool
		assertStatus   int32
		assertHeaders  map[string]*structpb.ListValue
		assertBodyData *structpb.Value
	}{
		{
			id: "GET Request with no body return with headers",
			request: &proxyv1.RequestProxyRequest{
				Service:    "cat",
				HttpMethod: "GET",
				Path:       "/meow",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Key1", "value1")
				w.Header().Set("Key1", "value2")
				w.Header().Set("Content-Length", "0")
				w.Header().Set("Date", "0")
				w.WriteHeader(200)
			},
			shouldError:  false,
			assertStatus: 200,
			assertHeaders: map[string]*structpb.ListValue{
				"Key1": {
					Values: []*structpb.Value{
						structpb.NewStringValue("value1"),
						structpb.NewStringValue("value2"),
					},
				},
				"Content-Length": {
					Values: []*structpb.Value{
						structpb.NewStringValue("0"),
					},
				},
				"Date": {
					Values: []*structpb.Value{
						structpb.NewStringValue("0"),
					},
				},
			},
		},
		{
			id: "POST Request with body data",
			request: &proxyv1.RequestProxyRequest{
				Service:    "cat",
				HttpMethod: "POST",
				Path:       "/nom",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte("{}"))
			},
			shouldError:    false,
			assertStatus:   200,
			assertBodyData: structpbFromBody([]byte("{}")),
		},
	}

	for _, test := range tests {
		srv := httptest.NewServer(http.HandlerFunc(test.handler))
		defer srv.Close()

		m := &mod{
			client:   srv.Client(),
			services: generateServicesConfig(srv.URL),
			logger:   log,
			scope:    scope,
		}

		res, err := m.RequestProxy(context.Background(), test.request)
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.Equal(t, test.assertStatus, res.HttpStatus)

			if test.assertHeaders != nil {
				assert.Equal(t, len(test.assertHeaders), len(res.Headers))
			}

			if test.assertBodyData != nil {
				assert.Equal(t, test.assertBodyData, res.Response)
			}
		}
	}
}

func TestIsAllowedRequest(t *testing.T) {
	tests := []struct {
		id          string
		service     string
		path        string
		method      string
		expect      bool
		shouldError bool
	}{
		{
			id:          "Allowed request",
			service:     "cat",
			path:        "/meow",
			method:      "GET",
			expect:      true,
			shouldError: false,
		},
		{
			id:          "Deined request method does not match",
			service:     "cat",
			path:        "/meow",
			method:      "POST",
			expect:      false,
			shouldError: false,
		},
		{
			id:          "Service does not exist",
			service:     "foo",
			path:        "/meow",
			method:      "POST",
			expect:      false,
			shouldError: false,
		},
		{
			id:          "Path with query params",
			service:     "cat",
			path:        "/nom?food=fancyfeast",
			method:      "POST",
			expect:      true,
			shouldError: false,
		},
		{
			id:          "Path with no /",
			service:     "cat",
			path:        "nom?food=fancyfeast",
			method:      "POST",
			expect:      false,
			shouldError: false,
		},
	}

	services := generateServicesConfig("http://test.test")

	for _, test := range tests {
		isAllowed, err := isAllowedRequest(services, test.service, test.path, test.method)
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expect, isAllowed)
		}
	}
}
