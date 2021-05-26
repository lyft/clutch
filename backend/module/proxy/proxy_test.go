package proxy

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"

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
	services := []*proxyv1cfg.Service{
		{
			Name: "cat",
			Host: "http://cat.cat",
			AllowedRequests: []*proxyv1cfg.AllowRequest{
				{Path: "/meow", Method: "GET"},
				{Path: "/nom", Method: "POST"},
			},
		},
	}
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	mockClient := &MockClient{}

	m := &mod{
		client:   mockClient,
		services: services,
		logger:   log,
		scope:    scope,
	}

	mockClient.DoFunc = func() (*http.Response, error) {
		return &http.Response{
			Status: "200",
			Body:   http.NoBody,
		}, nil
	}

	response, err := m.RequestProxy(context.Background(), &proxyv1.RequestProxyRequest{
		Service:    "cat",
		HttpMethod: "GET",
		Path:       "/meow",
	})
	assert.NoError(t, err)
	fmt.Printf("%v", response)
}

func TestIsAllowedRequest(t *testing.T) {
	services := []*proxyv1cfg.Service{
		{
			Name: "cat",
			Host: "http://cat.cat",
			AllowedRequests: []*proxyv1cfg.AllowRequest{
				{Path: "/meow", Method: "GET"},
				{Path: "/nom", Method: "POST"},
			},
		},
	}

	tests := []struct {
		id      string
		service string
		path    string
		method  string
		expect  bool
	}{
		{
			id:      "Allowed request",
			service: "cat",
			path:    "/meow",
			method:  "GET",
			expect:  true,
		},
		{
			id:      "Deined request method does not match",
			service: "cat",
			path:    "/meow",
			method:  "POST",
			expect:  false,
		},
		{
			id:      "Service does not exist",
			service: "foo",
			path:    "/meow",
			method:  "POST",
			expect:  false,
		},
		{
			id:      "Path with query params",
			service: "cat",
			path:    "/nom?food=fancyfeast",
			method:  "POST",
			expect:  true,
		},
	}

	for _, test := range tests {
		isAllowed := isAllowedRequest(services, test.service, test.path, test.method)
		assert.Equal(t, test.expect, isAllowed)
	}
}
