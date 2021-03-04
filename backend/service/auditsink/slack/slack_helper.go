package slack

import (
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/service/auditsink/slack/v1"
	"github.com/lyft/clutch/backend/middleware"
)

// GetOverrideMessage returns the custom slack message from the slack config if the FullMethod matches
// the service and method in the audit event. Otherwise ok is false and an empty string is returned.
func GetOverrideMessage(overrides []*configv1.CustomMessage, event *auditv1.RequestEvent) (string, bool) {
	if overrides == nil {
		// no overrides
		return "", false
	}

	for _, customSlack := range overrides {
		service, method, ok := middleware.SplitFullMethod(customSlack.FullMethod)
		if !ok {
			return "", false
		}
		if service == event.ServiceName && method == event.MethodName {
			return customSlack.Message, true
		}
	}

	return "", false
}
