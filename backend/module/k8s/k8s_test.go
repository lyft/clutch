package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
	"github.com/lyft/clutch/backend/mock/service/k8smock"
	"github.com/lyft/clutch/backend/module/moduletest"
	"github.com/lyft/clutch/backend/service"
)

func TestModule(t *testing.T) {
	service.Registry["clutch.service.k8s"] = k8smock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	m, err := New(nil, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("clutch.k8s.v1.K8sAPI"))
	assert.True(t, r.JSONRegistered())
}

func TestK8SAPIDescribePod(t *testing.T) {
	c := k8smock.New()
	api := newK8sAPI(c)
	resp, err := api.DescribePod(context.Background(), &k8sapiv1.DescribePodRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestK8SAPIListPods(t *testing.T) {
	c := k8smock.New()
	api := newK8sAPI(c)
	resp, err := api.ListPods(context.Background(), &k8sapiv1.ListPodsRequest{Options: &k8sapiv1.ListOptions{}})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestK8SAPIUpdatePod(t *testing.T) {
	c := k8smock.New()
	api := newK8sAPI(c)
	resp, err := api.UpdatePod(context.Background(), &k8sapiv1.UpdatePodRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestK8SAPIResizeHPA(t *testing.T) {
	c := k8smock.New()
	api := newK8sAPI(c)
	resp, err := api.ResizeHPA(context.Background(), &k8sapiv1.ResizeHPARequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestK8SAPIDescribeConfigMap(t *testing.T) {
	c := k8smock.New()
	api := newK8sAPI(c)
	resp, err := api.DescribeConfigMap(context.Background(), &k8sapiv1.DescribeConfigMapRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestK8SAPIDeleteConfigMap(t *testing.T) {
	c := k8smock.New()
	api := newK8sAPI(c)
	resp, err := api.DeleteConfigMap(context.Background(), &k8sapiv1.DeleteConfigMapRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestK8SAPIListConfigMaps(t *testing.T) {
	c := k8smock.New()
	api := newK8sAPI(c)
	resp, err := api.ListConfigMaps(context.Background(), &k8sapiv1.ListConfigMapsRequest{Options: &k8sapiv1.ListOptions{}})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
