package slack

import (
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/service/auditsink/slack/v1"
	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
	"github.com/lyft/clutch/backend/service/auditsink"
)

func TestNew(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	_, err := New(nil, log, scope)
	assert.Error(t, err)

	cfg, _ := ptypes.MarshalAny(&configv1.SlackConfig{})
	svc, err := New(cfg, log, scope)
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	_, ok := svc.(auditsink.Sink)
	assert.True(t, ok)
}

func TestFormat(t *testing.T) {
	t.Parallel()

	username := "foo"
	event := &auditv1.RequestEvent{
		ServiceName: "service",
		MethodName:  "method",
		Type:        apiv1.ActionType_READ,
		Resources: []*auditv1.Resource{{
			TypeUrl: "clutch.aws.v1.Instance",
			Id:      "i-01234567890abcdef0",
		}},
	}
	expected := "`foo` performed `method` via `service` using Clutch on resource(s):\n" +
		"- i-01234567890abcdef0 (`clutch.aws.v1.Instance`)"

	actual := formatText(username, event)
	assert.Equal(t, expected, actual)
}

// TODO: add more test cases (ie map and list)
func TestFormatCustomText(t *testing.T) {
	k8sRequest := &k8sapiv1.DescribePodRequest{Name: "foo"}
	k8sResponse := &k8sapiv1.DescribePodResponse{Pod: &k8sapiv1.Pod{PodIp: "000"}}
	anyK8sReq, _ := anypb.New(k8sRequest)
	anyK8sResp, _ := anypb.New(k8sResponse)

	ec2Request := &ec2v1.ResizeAutoscalingGroupRequest{Size: &ec2v1.AutoscalingGroupSize{Min: 2, Max: 4, Desired: 3}}
	ec2Response := &ec2v1.ResizeAutoscalingGroupResponse{}
	anyEc2Req, _ := anypb.New(ec2Request)
	anyEc2Resp, _ := anypb.New(ec2Response)

	testCases := []struct {
		text           string
		event          *auditv1.RequestEvent
		expectedErr    bool
		expectedOutput string
	}{
		{
			text: "{{.Request.Name}} ip address is {{.Response.Pod.PodIp}}",
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyK8sReq},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyK8sResp},
			},
			expectedOutput: "foo ip address is 000",
		},
		{
			text: "`Min size` is {{.Request.Size.Min}}, `Max size` is {{.Request.Size.Max}}, `Desired size` is {{.Request.Size.Desired}}",
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyEc2Req},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyEc2Resp},
			},
			expectedOutput: "`Min size` is 2, `Max size` is 4, `Desired size` is 3",
		},
		{
			text: "Name is {{.Foo}}",
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyK8sReq},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyK8sResp},
			},
			expectedErr: true,
		},
	}

	for _, test := range testCases {
		result, err := FormatCustomText(test.text, test.event)
		if test.expectedErr {
			assert.Error(t, err)
			assert.Empty(t, result)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedOutput, result)
		}
	}
}

func TestGetAuditMetadata(t *testing.T) {
	request := &k8sapiv1.DescribePodRequest{}
	response := &k8sapiv1.DescribePodResponse{}

	anyReq, _ := anypb.New(request)
	anyResp, _ := anypb.New(response)

	testCases := []struct {
		event       *auditv1.RequestEvent
		expectedErr bool
	}{
		{
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyReq},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyResp},
			},
		},
		{
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyReq},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: (*anypb.Any)(nil)},
			},
			expectedErr: true,
		},
	}

	for _, test := range testCases {
		result, err := getAuditMetadata(test.event)
		if test.expectedErr {
			assert.Error(t, err)
			assert.Nil(t, result)
		} else {
			assert.NoError(t, err)
			assert.IsType(t, request, result.Request)
			assert.IsType(t, response, result.Response)
		}
	}
}
