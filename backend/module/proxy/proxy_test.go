package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
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

type MockClient struct {
	DoFunc func() (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc()
}

var services []*proxyv1cfg.Service = []*proxyv1cfg.Service{
	{
		Name: "cat",
		Host: "http://cat.cat",
		AllowedRequests: []*proxyv1cfg.AllowRequest{
			{Path: "/meow", Method: "GET"},
			{Path: "/nom", Method: "POST"},
		},
	},
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
		doFunc         func() (*http.Response, error)
		shouldError    bool
		assertStatus   int32
		assertHeaders  map[string]*structpb.ListValue
		assertBodyData bool
	}{
		{
			id: "GET Request with no body return",
			request: &proxyv1.RequestProxyRequest{
				Service:    "cat",
				HttpMethod: "GET",
				Path:       "/meow",
			},
			doFunc: func() (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       http.NoBody,
					Header: http.Header{
						"key1": []string{"value1", "value2"},
					},
				}, nil
			},
			shouldError:  false,
			assertStatus: 200,
			assertHeaders: map[string]*structpb.ListValue{
				"key1": {
					Values: []*structpb.Value{
						structpb.NewStringValue("value1"),
						structpb.NewStringValue("value2"),
					},
				},
			},
			assertBodyData: false,
		},
		{
			id: "POST Request with body data",
			request: &proxyv1.RequestProxyRequest{
				Service:    "cat",
				HttpMethod: "POST",
				Path:       "/nom",
			},
			doFunc: func() (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
					Header: http.Header{
						"key1": []string{"value1", "value2"},
					},
				}, nil
			},
			shouldError:  false,
			assertStatus: 200,
			assertHeaders: map[string]*structpb.ListValue{
				"key1": {
					Values: []*structpb.Value{
						structpb.NewStringValue("value1"),
						structpb.NewStringValue("value2"),
					},
				},
			},
			assertBodyData: true,
		},
	}

	for _, test := range tests {
		m := &mod{
			client: &MockClient{
				DoFunc: test.doFunc,
			},
			services: services,
			logger:   log,
			scope:    scope,
		}

		res, err := m.RequestProxy(context.Background(), test.request)
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.Equal(t, test.assertStatus, res.HttpStatus)
			assert.Equal(t, test.assertHeaders, res.Headers)
			if test.assertBodyData {
				resData, err := test.doFunc()
				assert.NoError(t, err)

				var bodyData map[string]interface{}
				err = json.NewDecoder(resData.Body).Decode(&bodyData)
				assert.NoError(t, err)

				str, err := structpb.NewStruct(bodyData)
				assert.NoError(t, err)

				assert.Equal(t, structpb.NewStructValue(str), res.Response)
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
	}

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
