package auditsink

import (
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
)

// Required functions to register successfully with the configured
// Auditor in order to process audit events.
type Sink interface {
	// Write an event out to whatever this sinks into.
	Write(event *auditv1.Event) error
}

// Returns true if the filter matched the event, false if not.
// Because of how it interprets the denylist flag, auditors or sinks
// should check if auditsink.Filter(...) to see if the event should be passed
// down.
func Filter(filter *configv1.Filter, event *auditv1.Event) bool {
	if filter == nil {
		return true
	}

	req := event.GetEvent()
	if req == nil {
		return false
	}

	// If a denylist, return false when it matches. Else, true.
	rval := !filter.Denylist
	for _, filter := range filter.Rules {
		if ok := RunRequestFilter(filter, req); ok {
			return rval
		}
	}

	// The filter didn't apply, so allow it if it was a denylist and block
	// it if it was an allowlist.
	return filter.Denylist
}

func RunRequestFilter(filter *configv1.EventFilter, event *auditv1.RequestEvent) bool {
	switch filter.Value.(type) {
	case *configv1.EventFilter_Text:
		return textComparison(filter.GetText(), filter.GetField(), event)
	default:
		return false
	}
}

func textComparison(text string, field configv1.EventFilter_FilterType, event *auditv1.RequestEvent) bool {
	switch field {
	case configv1.EventFilter_SERVICE:
		return event.ServiceName == text
	case configv1.EventFilter_METHOD:
		return event.MethodName == text
	case configv1.EventFilter_TYPE:
		return event.Type.String() == text
	case configv1.EventFilter_UNSPECIFIED:
		fallthrough
	default:
		return true
	}
}
