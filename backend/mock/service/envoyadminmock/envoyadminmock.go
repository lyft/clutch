package envoyadminmock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

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

func NewAsService(*anypb.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	httpClient := &http.Client{Transport: &mockTransport{}}
	a, _ := anypb.New(&envoyadminv1.Config{Secure: false, DefaultRemotePort: 9999})
	return envoyadmin.NewWithHTTPClient(a, nil, nil, httpClient)
}
