package slack

import (
	"encoding/json"
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

func TestFormatCustomText(t *testing.T) {
	anyEc2Req, _ := anypb.New(&ec2v1.ResizeAutoscalingGroupRequest{Size: &ec2v1.AutoscalingGroupSize{Min: 2, Max: 4, Desired: 3}})
	anyEc2Resp, _ := anypb.New(&ec2v1.ResizeAutoscalingGroupResponse{})

	anyK8sDescribeReq, _ := anypb.New(&k8sapiv1.DescribePodRequest{Name: "foo"})
	anyK8sDescribeResp, _ := anypb.New(&k8sapiv1.DescribePodResponse{Pod: &k8sapiv1.Pod{PodIp: "000", Labels: map[string]string{}}})

	k8sUpdateReq := &k8sapiv1.UpdatePodRequest{
		ExpectedObjectMetaFields: &k8sapiv1.ExpectedObjectMetaFields{
			Annotations: map[string]*k8sapiv1.NullableString{
				"baz": &k8sapiv1.NullableString{Kind: &k8sapiv1.NullableString_Null{}},
				"foo": &k8sapiv1.NullableString{Kind: &k8sapiv1.NullableString_Value{Value: "new-value"}},
			},
		},
		ObjectMetaFields:       &k8sapiv1.ObjectMetaFields{Labels: map[string]string{"foo": "new-value"}},
		RemoveObjectMetaFields: &k8sapiv1.RemoveObjectMetaFields{Annotations: []string{"foo", "bar"}},
	}

	anyK8sUpdateReq, _ := anypb.New(k8sUpdateReq)
	anyK8UpdateResp, _ := anypb.New(&k8sapiv1.UpdatePodResponse{})

	testCases := []struct {
		text           string
		event          *auditv1.RequestEvent
		expectedErr    bool
		expectedOutput string
	}{
		// metadata from the API request
		{
			text: "`Min size` is {{.Request.size.min}}, `Max size` is {{.Request.size.max}}, `Desired size` is {{.Request.size.desired}}",
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyEc2Req},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyEc2Resp},
			},
			expectedOutput: "`Min size` is 2, `Max size` is 4, `Desired size` is 3",
		},
		// metadata from both the API request and repsonse
		{
			text: "{{.Request.name}} ip address is {{.Response.pod.podIp}}",
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyK8sDescribeReq},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyK8sDescribeResp},
			},
			expectedOutput: "foo ip address is 000",
		},
		// metdata (labels) value is nil
		{
			text: "{{.Request.name}} labels: {{slackList .Response.pod.labels}}",
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyK8sDescribeReq},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyK8sDescribeResp},
			},
			expectedOutput: "foo labels: None",
		},
		// metadata that is a map, uses helper slackList
		{
			text: "*Updated labels*:{{slackList .Request.objectMetaFields.labels}}",
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyK8sUpdateReq},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyK8UpdateResp},
			},
			expectedOutput: "*Updated labels*:\n- foo: new-value",
		},
		// metdata that is a list, uses helper slackList
		{
			text: "*Removed annotations*:{{slackList .Request.removeObjectMetaFields.annotations}}",
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyK8sUpdateReq},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyK8UpdateResp},
			},
			expectedOutput: "*Removed annotations*:\n- foo\n- bar",
		},
		// metdata that is a map, map value is a another map
		// uses the Golang template `range`
		{
			text: "*Expected Preconditions*:{{range $key, $val := .Request.expectedObjectMetaFields.annotations}}\n- {{$key}}: {{range $i, $j := $val}}{{$j}}{{end}}{{end}}",
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyK8sUpdateReq},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyK8UpdateResp},
			},
			expectedOutput: "*Expected Preconditions*:\n- baz: None\n- foo: new-value",
		},
		// invalid field name
		{
			text: "Name is {{.Foo}}",
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyK8sDescribeReq},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyK8sDescribeResp},
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

func TestSlackList(t *testing.T) {
	testCases := []struct {
		input          interface{}
		expectedOutput string
	}{
		{
			input:          "hello",
			expectedOutput: "ERR_INPUT_NOT_SLICE_OR_MAP",
		},
		{
			input:          []string{"foo"},
			expectedOutput: "\n- foo",
		},
		{
			input:          []int{1},
			expectedOutput: "\n- 1",
		},
		{
			input:          map[string]string{"foo": "value"},
			expectedOutput: "\n- foo: value",
		},
		{
			input:          map[string]bool{"foo": true},
			expectedOutput: "\n- foo: true",
		},
		{
			input:          map[string]string{},
			expectedOutput: "None",
		},
		{
			input:          []string{},
			expectedOutput: "None",
		},
		{
			input:          nil,
			expectedOutput: "None",
		},
	}

	for _, test := range testCases {
		result := slackList(test.input)
		assert.Equal(t, test.expectedOutput, result)
	}
}

func TestGetAuditTemplateData(t *testing.T) {
	anyReq, _ := anypb.New(&k8sapiv1.DescribePodRequest{})
	anyResp, _ := anypb.New(&k8sapiv1.DescribePodResponse{})

	testCases := []struct {
		event            *auditv1.RequestEvent
		expectedReqType  string
		expectedRespType string
		expectedEmpty    bool
	}{
		{
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: anyReq},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: anyResp},
			},
			expectedReqType:  "clutch.k8s.v1.DescribePodRequest",
			expectedRespType: "clutch.k8s.v1.DescribePodResponse",
		},
		{
			event: &auditv1.RequestEvent{
				RequestMetadata:  &auditv1.RequestMetadata{Body: (*anypb.Any)(nil)},
				ResponseMetadata: &auditv1.ResponseMetadata{Body: (*anypb.Any)(nil)},
			},
			expectedEmpty: true,
		},
	}

	for _, test := range testCases {
		result, err := getAuditTemplateData(test.event)
		assert.NoError(t, err)
		if test.expectedEmpty {
			assert.Empty(t, result.Request)
			assert.Empty(t, result.Response)
		} else {
			reqB, _ := json.Marshal(result.Request)
			respB, _ := json.Marshal(result.Response)
			assert.Contains(t, string(reqB), test.expectedReqType)
			assert.Contains(t, string(respB), test.expectedRespType)
		}
	}
}
