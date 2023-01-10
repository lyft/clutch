package audit

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	"github.com/lyft/clutch/backend/mock/service/auditmock"
)

func TestGetEvents(t *testing.T) {
	testCases := []struct {
		id                 string
		eventCount         int
		req                *auditv1.GetEventsRequest
		expectedEventCount int
		expectedNextToken  string
		expectedErr        error
	}{
		{
			id:          "without time window",
			req:         &auditv1.GetEventsRequest{},
			expectedErr: errors.New("no time window requested"),
		},
		{
			id:         "with time window and events",
			eventCount: 11,
			req: &auditv1.GetEventsRequest{
				Window: &auditv1.GetEventsRequest_Range{
					Range: &auditv1.TimeRange{
						StartTime: timestamppb.New(time.Now().Add(-1 * time.Hour)),
						EndTime:   timestamppb.New(time.Now().Add(6 * time.Hour)),
					},
				},
			},
			expectedEventCount: 11,
			expectedNextToken:  "",
		},
		{
			id:         "with time window in future",
			eventCount: 11,
			req: &auditv1.GetEventsRequest{
				Window: &auditv1.GetEventsRequest_Range{
					Range: &auditv1.TimeRange{
						StartTime: timestamppb.New(time.Now().Add(1 * time.Hour)),
						EndTime:   timestamppb.New(time.Now().Add(2 * time.Hour)),
					},
				},
			},
			expectedEventCount: 0,
		},
		{
			id:         "with time since an hour ago",
			eventCount: 11,
			req: &auditv1.GetEventsRequest{
				Window: &auditv1.GetEventsRequest_Since{
					Since: durationpb.New(1 * time.Hour),
				},
			},
			expectedEventCount: 11,
			expectedNextToken:  "",
		},
		{
			id:         "with time since a microsecond ago",
			eventCount: 11,
			req: &auditv1.GetEventsRequest{
				Window: &auditv1.GetEventsRequest_Since{
					Since: durationpb.New(1 * time.Microsecond),
				},
			},
			expectedEventCount: 0,
		},
		// invalid page token should return error on conversion
		{
			id:         "with invalid page token",
			eventCount: 11,
			req: &auditv1.GetEventsRequest{
				Window: &auditv1.GetEventsRequest_Since{
					Since: durationpb.New(1 * time.Hour),
				},
				PageToken: "invalid",
			},
			expectedErr: errors.New("invalid page token: invalid"),
		},
		{
			id:         "with negative page token",
			eventCount: 11,
			req: &auditv1.GetEventsRequest{
				Window: &auditv1.GetEventsRequest_Since{
					Since: durationpb.New(1 * time.Hour),
				},
				PageToken: "-1",
			},
			expectedErr: errors.New("invalid page token: -1"),
		},
		{
			id:         "with page token",
			eventCount: 11,
			req: &auditv1.GetEventsRequest{
				Window: &auditv1.GetEventsRequest_Since{
					Since: durationpb.New(1 * time.Hour),
				},
				PageToken: "0",
			},
			expectedEventCount: 11,
			expectedNextToken:  "",
		},
		{
			id:         "with page token and limit",
			eventCount: 11,
			req: &auditv1.GetEventsRequest{
				Window: &auditv1.GetEventsRequest_Since{
					Since: durationpb.New(1 * time.Hour),
				},
				PageToken: "0",
				Limit:     5,
			},
			expectedEventCount: 5,
			expectedNextToken:  "1",
		},
		{
			id:         "with page token and no next page",
			eventCount: 10,
			req: &auditv1.GetEventsRequest{
				Window: &auditv1.GetEventsRequest_Since{
					Since: durationpb.New(1 * time.Hour),
				},
				PageToken: "0",
			},
			expectedEventCount: 10,
			expectedNextToken:  "",
		},
		{
			id:         "with too large of page token",
			eventCount: 1,
			req: &auditv1.GetEventsRequest{
				Window: &auditv1.GetEventsRequest_Since{
					Since: durationpb.New(1 * time.Hour),
				},
				PageToken: "5",
			},
			expectedEventCount: 0,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(fmt.Sprintf(test.id), func(t *testing.T) {
			t.Parallel()
			m := &mod{
				client: auditmock.New(),
			}
			for i := 0; i < test.eventCount; i++ {
				_, err := m.client.WriteRequestEvent(context.Background(), &auditv1.RequestEvent{})
				if err != nil {
					assert.Fail(t, err.Error())
				}
			}
			time.Sleep(1 * time.Second)
			resp, err := m.GetEvents(context.Background(), test.req)
			if test.expectedErr != nil {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Equal(t, err, test.expectedErr)
			} else {
				assert.Equal(t, test.expectedEventCount, len(resp.Events))
				assert.Equal(t, test.expectedNextToken, resp.NextPageToken)
			}
		})
	}
}

func TestGetEvent(t *testing.T) {
	testCases := []struct {
		req         *auditv1.GetEventRequest
		expectedId  int64
		expectedErr error
	}{
		{
			req:        &auditv1.GetEventRequest{EventId: 0},
			expectedId: 0,
		},
		{
			req:         &auditv1.GetEventRequest{EventId: 1},
			expectedErr: errors.New("event with id 1 not found"),
		},
		{
			req:         &auditv1.GetEventRequest{EventId: -1},
			expectedErr: errors.New("event with id -1 not found"),
		},
	}

	for idx, test := range testCases {
		test := test
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			m := &mod{
				client: auditmock.New(),
			}
			_, err := m.client.WriteRequestEvent(context.Background(), &auditv1.RequestEvent{})
			if err != nil {
				assert.Fail(t, err.Error())
			}
			resp, err := m.GetEvent(context.Background(), test.req)
			if test.expectedErr != nil {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Equal(t, err, test.expectedErr)
			} else {
				assert.Equal(t, test.expectedId, resp.Event.Id)
			}
		})
	}
}
