package temporal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	temporalv1 "github.com/lyft/clutch/backend/api/config/service/temporal/v1"
)

func TestNew(t *testing.T) {
	cfg := &temporalv1.Config{
		Host: "dns:///example.com",
		Port: 9233,
	}
	anycfg, _ := anypb.New(cfg)

	c, err := New(anycfg, zap.NewNop(), tally.NoopScope)
	assert.NoError(t, err)
	impl := c.(*clientImpl)
	assert.NotNil(t, impl.logger)
	assert.NotNil(t, impl.metricsHandler)
	assert.Nil(t, impl.copts.TLS)
	assert.True(t, impl.copts.DisableHealthCheck)
	assert.Equal(t, impl.hostPort, "dns:///example.com:9233")
}

func TestNewClientWithConnectionOptions(t *testing.T) {
	cfg := &temporalv1.Config{
		Host:              "dns:///example.com",
		Port:              9233,
		ConnectionOptions: &temporalv1.ConnectionOptions{UseSystemCaBundle: true},
	}
	c, err := newClient(cfg, zap.NewNop(), tally.NoopScope)
	assert.NoError(t, err)

	impl := c.(*clientImpl)
	assert.NotNil(t, impl.copts.TLS.RootCAs)
	assert.True(t, impl.copts.DisableHealthCheck)
}

func TestGetNamespaceClient(t *testing.T) {
	cfg := &temporalv1.Config{Host: "example.com", Port: 9233}
	c, _ := newClient(cfg, zap.NewNop(), tally.NoopScope)
	client, err := c.GetNamespaceClient("foo-namespace")
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
