package auditmock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/audit"
)

type svc struct {
	sync.RWMutex

	events []*auditv1.Event
}

func (s *svc) WriteRequestEvent(_ context.Context, req *auditv1.RequestEvent) (int64, error) {
	s.Lock()
	defer s.Unlock()

	i := int64(len(s.events))
	event := &auditv1.Event{
		Id:         i,
		OccurredAt: timestamppb.Now(),
		EventType: &auditv1.Event_Event{
			Event: req,
		},
	}
	s.events = append(s.events, event)
	return i, nil
}

func (s *svc) UpdateRequestEvent(_ context.Context, id int64, event *auditv1.RequestEvent) error {
	s.Lock()
	defer s.Unlock()

	proto.Merge(s.events[id].GetEvent(), event)
	return nil
}

func (s *svc) ReadEvents(_ context.Context, start time.Time, end *time.Time, options *audit.ReadOptions) ([]*auditv1.Event, error) {
	s.RLock()
	defer s.RUnlock()

	var stop time.Time
	if end == nil {
		stop = time.Now()
	} else {
		stop = *end
	}

	var events []*auditv1.Event
	for _, event := range s.events {
		eventTime := event.OccurredAt.AsTime()

		if start.After(eventTime) {
			continue
		}
		if stop.Before(eventTime) {
			break
		}
		events = append(events, event)
	}

	if options != nil {
		if options.Offset > int64(len(events)) {
			return []*auditv1.Event{}, nil
		}
		startIdx := int64(0)
		endIdx := int64(len(events))
		if options.Offset != 0 && options.Offset > 0 {
			startIdx = options.Offset
		}
		if options.Limit != 0 {
			endIdx = options.Offset + options.Limit
			if endIdx > int64(len(events)) {
				endIdx = int64(len(events))
			}
		}
		return events[startIdx:endIdx], nil
	}
	return events, nil
}

func (s *svc) ReadEvent(_ context.Context, id int64) (*auditv1.Event, error) {
	s.RLock()
	defer s.RUnlock()

	if id >= int64(len(s.events)) || id < 0 {
		return nil, fmt.Errorf("event with id %d not found", id)
	}

	return s.events[id], nil
}

func (s *svc) UnsentEvents(_ context.Context) ([]*auditv1.Event, error) {
	return s.events, nil
}

func (s *svc) Filter(_ *auditv1.Event) bool {
	return true
}

func New() audit.Auditor {
	return &svc{}
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}
