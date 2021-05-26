package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"

	proxyv1cfg "github.com/lyft/clutch/backend/api/config/module/proxy/v1"
	"github.com/lyft/clutch/backend/module/moduletest"
)

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
