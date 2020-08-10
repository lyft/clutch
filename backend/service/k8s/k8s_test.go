package k8s

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	restclient "k8s.io/client-go/rest"

	k8sv1 "github.com/lyft/clutch/backend/api/config/service/k8s/v1"
)

var testConfig = `
apiVersion: v1
clusters:
- cluster:
    server: test-server
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-user@test-cluster
current-context: test-user@test-cluster
kind: Config
preferences: {}
users:
- name: test-user
`

func TestNew(t *testing.T) {
	tempfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tempfile.Name())
	_ = ioutil.WriteFile(tempfile.Name(), []byte(testConfig), 0500)

	paths := []string{tempfile.Name()}

	cfg, _ := ptypes.MarshalAny(&k8sv1.Config{
		Kubeconfigs: paths,
	})
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	s, err := New(cfg, log, scope)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Check public interface compliance.
	_, ok := s.(Service)
	assert.True(t, ok)

	// Check private interface compliance.
	c, ok := s.(*svc)
	assert.True(t, ok)
	assert.NotNil(t, c.log)
	assert.NotNil(t, c.scope)
	assert.Len(t, c.manager.Clientsets(), 1)
}

func TestNewWithWrongConfig(t *testing.T) {
	_, err := New(&any.Any{TypeUrl: "foobar"}, nil, nil)
	assert.EqualError(t, err, `mismatched message type: got "foobar" want "clutch.config.service.k8s.v1.Config"`)
}

func TestApplyRestClientConfig(t *testing.T) {
	t.Parallel()
	var testCases = []struct {
		id                 string
		restConfig         *restclient.Config
		expectedRestConfig restclient.Config
		restClientConfig   *k8sv1.RestClientConfig
	}{
		{
			id:                 "no config override",
			restConfig:         &restclient.Config{},
			expectedRestConfig: restclient.Config{},
			restClientConfig:   &k8sv1.RestClientConfig{},
		},
		{
			id:         "all config override",
			restConfig: &restclient.Config{},
			expectedRestConfig: restclient.Config{
				Timeout: 10000000000, // This is in nanoseconds
				QPS:     100,
				Burst:   1000,
			},
			restClientConfig: &k8sv1.RestClientConfig{
				Timeout: &duration.Duration{
					Seconds: 10,
				},
				Qps:   100,
				Burst: 1000,
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			err := ApplyRestClientConfig(tt.restConfig, tt.restClientConfig)
			assert.NoError(t, err)
			assert.Equal(t, tt.restConfig.Timeout, tt.expectedRestConfig.Timeout)
			assert.Equal(t, tt.restConfig.QPS, tt.expectedRestConfig.QPS)
			assert.Equal(t, tt.restConfig.Burst, tt.expectedRestConfig.Burst)
		})
	}
}
