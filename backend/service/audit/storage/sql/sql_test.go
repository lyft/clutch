package sql

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
	"github.com/lyft/clutch/backend/mock/service/dbmock"
)

func TestConvertAPIBody(t *testing.T) {
	// set up for TestConvertAPIBody
	a1 := (*anypb.Any)(nil)

	p1 := &ec2v1.Instance{InstanceId: "i-123456789abcdef0"}
	a2, _ := anypb.New(p1)

	p2 := &k8sapiv1.ResizeHPAResponse{}
	a3, _ := anypb.New(p2)

	tests := []struct {
		input *anypb.Any
	}{
		// case: untyped nil
		{input: nil},
		// case: input is a typed nil
		{input: a1},
		// case: input is typed with non-nil value
		{input: a2},
		// case: input is typed with non-nil value
		{input: a3},
	}

	for _, test := range tests {
		b, err := convertAPIBody(test.input)
		assert.NotNil(t, b)
		assert.NoError(t, err)
	}
}

func TestAPIBodyProto(t *testing.T) {
	var nilJSON json.RawMessage

	tests := []struct {
		input     json.RawMessage
		expectNil bool
	}{
		{input: nil, expectNil: true},
		{input: nilJSON, expectNil: true},
		{input: []byte(`{}`), expectNil: false},
		{input: []byte(`{"@type":"type.googleapis.com/clutch.k8s.v1.Pod"}`), expectNil: false},
	}

	for _, test := range tests {
		a, err := apiBodyProto(test.input)
		if test.expectNil {
			assert.Nil(t, a)
		} else {
			assert.NotNil(t, a)
		}
		assert.NoError(t, err)
	}
}

func TestGetAdvisoryConn(t *testing.T) {
	dbm := dbmock.NewMockDB()
	c := &client{
		db: dbm.DB(),
	}

	conn, err := c.getAdvisoryConn()
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	assert.NotNil(t, c.advisoryLockConn)
}

func TestReadEvent(t *testing.T) {
	dbm := dbmock.NewMockDB()
	dbm.Register()

	c := &client{
		db: dbm.DB(),
	}

	dbm.Mock.ExpectExec("SELECT id, occurred_at, details FROM audit_events WHERE id = $1").WillReturnResult(
		sqlmock.NewRows(string{"{id: 1}"}),
	)

	anyReq, _ := anypb.New(&ec2v1.ResizeAutoscalingGroupRequest{Size: &ec2v1.AutoscalingGroupSize{Min: 2, Max: 4, Desired: 3}})
	writeEvent := &auditv1.RequestEvent{
		RequestMetadata: &auditv1.RequestMetadata{Body: anyReq},
	}
	_, err := c.WriteRequestEvent(context.Background(), writeEvent)
	assert.NoError(t, err)
	_, err = c.WriteRequestEvent(context.Background(), writeEvent)
	assert.NoError(t, err)
	_, err = c.WriteRequestEvent(context.Background(), writeEvent)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	event, err := c.ReadEvent(context.Background(), 1)
	assert.NoError(t, err)

	assert.Equal(t, int64(1), event.Id)
}
