package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	proxyv1cfg "github.com/lyft/clutch/backend/api/config/module/proxy/v1"
	proxyv1 "github.com/lyft/clutch/backend/api/proxy/v1"
	"github.com/lyft/clutch/backend/module"
)

const (
	Name = "clutch.module.proxy"
)

func New(cfg *any.Any, log *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &proxyv1cfg.Config{}
	err := cfg.UnmarshalTo(config)
	if err != nil {
		return nil, err
	}

	// Validate that each services constructs a parsable URL
	for _, service := range config.Services {
		for _, ar := range service.AllowedRequests {
			_, err := url.Parse(fmt.Sprintf("%s%s", service.Host, ar.Path))
			if err != nil {
				return nil, fmt.Errorf("Unable to parse the configured URL for service [%s]", service.Name)
			}
		}
	}

	m := &mod{
		client:   &http.Client{},
		services: config.Services,
		logger:   log,
		scope:    scope,
	}

	return m, nil
}

type mod struct {
	client   *http.Client
	services []*proxyv1cfg.Service
	logger   *zap.Logger
	scope    tally.Scope
}

func (m *mod) Register(r module.Registrar) error {
	proxyv1.RegisterProxyAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(proxyv1.RegisterProxyAPIHandler)
}

func (m *mod) RequestProxy(ctx context.Context, req *proxyv1.RequestProxyRequest) (*proxyv1.RequestProxyResponse, error) {
	isAllowed, err := isAllowedRequest(m.services, req.Service, req.Path, req.HttpMethod)
	if err != nil {
		m.logger.Error("Unable to parse the configured URL", zap.Error(err))
		return nil, fmt.Errorf("Unable to parse the configured URL for service [%s]", req.Service)
	}

	if !isAllowed {
		return nil, status.Error(codes.InvalidArgument, "This request is not allowed, check the proxy configuration.")
	}

	// If its allowed lookup the service
	var service *proxyv1cfg.Service
	for _, s := range m.services {
		if s.Name == req.Service {
			service = s
		}
	}

	// Set all additional headers if specified
	headers := http.Header{}
	for k, v := range service.Headers {
		headers.Add(k, v)
	}

	// Parse the URL by joining both the HOST and PATH specifed by the config
	parsedUrl, err := url.Parse(fmt.Sprintf("%s%s", service.Host, req.Path))
	if err != nil {
		m.logger.Error("Unable to parse the configured URL", zap.Error(err))
		return nil, fmt.Errorf("Unable to parse the configured URL for service [%s]", service.Name)
	}

	// Constructing the request object
	request := &http.Request{
		Method: req.HttpMethod,
		URL:    parsedUrl,
		Header: headers,
	}

	response, err := m.client.Do(request)
	if err != nil {
		m.logger.Error("proxy request error", zap.Error(err))
		return nil, err
	}

	// Extract headers from response
	// TODO: It might make sense to provide a list of allowed headers, as there can be a lot.
	resHeaders := make(map[string]*structpb.ListValue, len(response.Header))
	for key, headers := range response.Header {
		headerValues := make([]*structpb.Value, len(headers))
		for _, h := range headers {
			headerValues = append(headerValues, structpb.NewStringValue(h))
		}

		resHeaders[key] = &structpb.ListValue{
			Values: headerValues,
		}
	}

	proxyResponse := &proxyv1.RequestProxyResponse{
		HttpStatus: int32(response.StatusCode),
		Headers:    resHeaders,
	}

	var bodyData interface{}
	err = json.NewDecoder(response.Body).Decode(&bodyData)
	switch {
	// There is no body data so do nothing
	case err == io.EOF:
	case err != nil:
		m.logger.Error("Unable to decode response body", zap.Error(err))
		return nil, err
	default:
		bodyStruct, err := structpb.NewValue(bodyData)
		if err != nil {
			m.logger.Error("Unable to create structpb from body data", zap.Error(err))
			return nil, err
		}
		proxyResponse.Response = bodyStruct
	}

	return proxyResponse, nil
}

func isAllowedRequest(services []*proxyv1cfg.Service, service, path, method string) (bool, error) {
	for _, s := range services {
		if s.Name == service {
			for _, ar := range s.AllowedRequests {
				parsedUrl, err := url.Parse(fmt.Sprintf("%s%s", s.Host, path))
				if err != nil {
					return false, err
				}

				if parsedUrl.Path == ar.Path && strings.EqualFold(method, ar.Method) {
					return true, nil
				}
			}
			// return early here as were done checking allowed request for this service
			return false, nil
		}
	}
	return false, nil
}
