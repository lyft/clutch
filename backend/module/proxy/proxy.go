package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"

	proxyv1cfg "github.com/lyft/clutch/backend/api/config/module/proxy/v1"
	proxyv1 "github.com/lyft/clutch/backend/api/proxy/v1"
	"github.com/lyft/clutch/backend/module"
)

const Name = "clutch.module.proxy"

func New(cfg *any.Any, log *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &proxyv1cfg.Config{}
	err := cfg.UnmarshalTo(config)
	if err != nil {
		return nil, err
	}

	m := &mod{
		services: config.Services,
		logger:   log,
		scope:    scope,
	}

	return m, nil
}

type mod struct {
	services []*proxyv1cfg.Service
	logger   *zap.Logger
	scope    tally.Scope
}

func (m *mod) Register(r module.Registrar) error {
	proxyv1.RegisterProxyAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(proxyv1.RegisterProxyAPIHandler)
}

func (m *mod) RequestProxy(ctx context.Context, req *proxyv1.RequestProxyRequest) (*proxyv1.RequestProxyResponse, error) {
	if !isAllowedRequest(m.services, req.Service, req.Path, req.HttpMethod) {
		return nil, errors.New("This request is not allowed, check the proxy configuration.")
	}

	// If its allowed lookup the service
	service := &proxyv1cfg.Service{}
	for _, s := range m.services {
		if s.Name == req.Service {
			service = s
		}
	}

	// Set all additional headers if specified
	headers := http.Header{}
	if len(service.Headers) > 0 {
		for k, v := range service.Headers {
			headers.Add(k, v)
		}
	}

	// Parse the URL by joining both the HOST and PATH specifed by the config
	parsedUrl, err := url.Parse(fmt.Sprintf("%s%s", service.Host, req.Path))
	if err != nil {
		return nil, err
	}

	// Constructing the request object
	request := &http.Request{
		Method: req.HttpMethod,
		URL:    parsedUrl,
		Header: headers,
	}

	// Using the default HTTP client make the request
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var bodyData map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&bodyData)
	if err != nil {
		return nil, err
	}

	str, err := structpb.NewStruct(bodyData)
	if err != nil {
		return nil, err
	}

	// Extract headers from response
	resHeaders := map[string]*proxyv1.HeaderValues{}
	for k, v := range response.Header {
		resHeaders[k] = &proxyv1.HeaderValues{
			Values: v,
		}
	}

	return &proxyv1.RequestProxyResponse{
		HttpStatus: int32(response.StatusCode),
		Response:   structpb.NewStructValue(str),
		Headers:    resHeaders,
	}, nil
}

func isAllowedRequest(services []*proxyv1cfg.Service, service, path, method string) bool {
	for _, s := range services {
		if s.Name == service {
			for _, ar := range s.AllowedRequests {
				// if the path has query prams chop them off and eval only the path
				finalPath := path
				if strings.Contains(path, "?") {
					finalPath = strings.Split(path, "?")[0]
				}

				if finalPath == ar.Path && strings.EqualFold(method, ar.Method) {
					return true
				}
			}
			// return early here as were done checking allowed request for this service
			return false
		}
	}
	return false
}
