package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	proxyv1cfg "github.com/lyft/clutch/backend/api/config/module/proxy/v1"
	proxyv1 "github.com/lyft/clutch/backend/api/proxy/v1"
	"github.com/lyft/clutch/backend/module"
)

const (
	Name          = "clutch.module.proxy"
	HostHeaderKey = "Host"
)

func New(cfg *anypb.Any, log *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &proxyv1cfg.Config{}
	err := cfg.UnmarshalTo(config)
	if err != nil {
		return nil, err
	}

	err = validateConfigPaths(config)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("unable to parse the configured URL for service [%s]", req.Service)
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
		return nil, fmt.Errorf("unable to parse the configured URL for service [%s]", service.Name)
	}

	// Constructing the request object
	request := &http.Request{
		Method: req.HttpMethod,
		URL:    parsedUrl,
		Header: headers,
	}

	addExcludedHeaders(request)

	if req.Request != nil {
		requestJSON, err := protojson.Marshal(req.Request)
		if err != nil {
			return nil, err
		}
		buff := bytes.NewBuffer(requestJSON)
		request.Body = io.NopCloser(buff)
	}

	response, err := m.client.Do(request)
	if err != nil {
		m.scope.Tagged(map[string]string{
			"service": service.Name,
			"path":    req.Path,
		}).Counter("request.error").Inc(1)
		m.logger.Error("proxy request error", zap.Error(err))
		return nil, err
	}

	m.scope.Tagged(map[string]string{
		"service":     service.Name,
		"path":        req.Path,
		"status_code": fmt.Sprintf("%d", response.StatusCode),
	}).Counter("request").Inc(1)

	// Extract headers from response
	// TODO: It might make sense to provide a list of allowed headers, as there can be a lot.
	resHeaders := make(map[string]*structpb.ListValue, len(response.Header))
	for key, headers := range response.Header {
		headerValues := make([]*structpb.Value, len(headers))
		for i, h := range headers {
			headerValues[i] = structpb.NewStringValue(h)
		}

		resHeaders[key] = &structpb.ListValue{
			Values: headerValues,
		}
	}

	proxyResponse := &proxyv1.RequestProxyResponse{
		HttpStatus: int32(response.StatusCode), //nolint
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

func (m *mod) RequestProxyGet(ctx context.Context, req *proxyv1.RequestProxyGetRequest) (*proxyv1.RequestProxyGetResponse, error) {
	// Validate that it's a GET.
	if req.HttpMethod != http.MethodGet {
		return nil, status.Errorf(codes.InvalidArgument, "non-GET request passed to GET specific endpoint")
	}

	// Proxy the call to the original proxy method and return the response.
	resp, err := m.RequestProxy(ctx, getRequestToRequest(req))
	if err != nil {
		return nil, err
	}
	return responseToGetResponse(resp), nil
}

func responseToGetResponse(resp *proxyv1.RequestProxyResponse) *proxyv1.RequestProxyGetResponse {
	return &proxyv1.RequestProxyGetResponse{
		HttpStatus: resp.HttpStatus,
		Headers:    resp.Headers,
		Response:   resp.Response,
	}
}

func getRequestToRequest(req *proxyv1.RequestProxyGetRequest) *proxyv1.RequestProxyRequest {
	return &proxyv1.RequestProxyRequest{
		Service:    req.Service,
		HttpMethod: req.HttpMethod,
		Path:       req.Path,
		Request:    req.Request,
	}
}

func isAllowedRequest(services []*proxyv1cfg.Service, service, path, method string) (bool, error) {
	for _, s := range services {
		if s.Name == service {
			for _, ar := range s.AllowedRequests {
				switch t := ar.PathType.(type) {
				case *proxyv1cfg.AllowRequest_Path:
					parsedUrl, err := url.Parse(fmt.Sprintf("%s%s", s.Host, path))
					if err != nil {
						return false, err
					}
					if parsedUrl.Path == t.Path && strings.EqualFold(method, ar.Method) {
						return true, nil
					}
				case *proxyv1cfg.AllowRequest_PathRegex:
					r, err := regexp.Compile(t.PathRegex)
					if err != nil {
						return false, err
					}
					if r.MatchString(path) {
						return true, nil
					}
				default:
					return false, fmt.Errorf("path type not supported: %T", t)
				}
			}
			// return early here as were done checking allowed request for this service
			return false, nil
		}
	}
	return false, nil
}

/*
For headers that get ignored in the header map, this helper adds their values back to the designated
fields on the Request struct.
Context:
	https://github.com/golang/go/issues/29865
	https://github.com/golang/go/blob/8c94aa40e6f5e61e8a570e9d20b7d0d4ad8c382d/src/net/http/request.go#L88
*/
// TODO: add the other headers that get excluded from the request
func addExcludedHeaders(request *http.Request) {
	// Get() is case insensitive
	if hostHeader := request.Header.Get(HostHeaderKey); hostHeader != "" {
		request.Host = hostHeader
	}
}

func validateConfigPaths(config *proxyv1cfg.Config) error {
	for _, service := range config.Services {
		for _, ar := range service.AllowedRequests {
			switch t := ar.PathType.(type) {
			case *proxyv1cfg.AllowRequest_Path:
				// For exact path type, validate that string constructs a parsable URL
				_, err := url.Parse(fmt.Sprintf("%s%s", service.Host, t.Path))
				if err != nil {
					return fmt.Errorf("unable to parse the configured URL for service [%s]", service.Name)
				}
			case *proxyv1cfg.AllowRequest_PathRegex:
				// For path regex type, validate that expression can be parsed
				_, err := regexp.Compile(t.PathRegex)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("path type not supported: %T", t)
			}
		}
	}
	return nil
}
