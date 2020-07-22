package envoyadminmock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	envoyadminv1 "github.com/lyft/clutch/backend/api/config/service/envoyadmin/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/envoyadmin"
)

type mockTransport struct{}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp string
	switch req.URL.Path {
	case "/clusters":
		resp = clustersResponse
	case "/config_dump":
		resp = configDumpResponse
	case "/listeners":
		resp = listenersResponse
	case "/runtime":
		resp = runtimeResponse
	case "/server_info":
		resp = serverInfoResponse
	case "/stats":
		resp = statsResponse
	default:
		return nil, fmt.Errorf("path '%s' was not implemented in mock transport", req.URL.Path)
	}

	return &http.Response{
		Status:     "OK",
		StatusCode: 200,
		Request:    req,
		Body:       ioutil.NopCloser(strings.NewReader(resp)),
	}, nil
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	httpClient := &http.Client{Transport: &mockTransport{}}
	a, _ := ptypes.MarshalAny(&envoyadminv1.Config{Secure: false, DefaultRemotePort: 9999})
	return envoyadmin.NewWithHTTPClient(a, nil, nil, httpClient)
}
