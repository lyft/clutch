package slack

import (
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/service/auditsink/slack/v1"
	"github.com/lyft/clutch/backend/middleware"
)

// GetSlackOverrideText returns the custom slack message from the slack override config if the FullMethod matches
// the serive and method in the audit event. Otherwise ok is false and an empty string is returned.
func GetSlackOverrideText(override *configv1.Override, event *auditv1.RequestEvent) (string, bool) {
	if override == nil {
		// no slack overrides
		return "", false
	}

	for _, customSlack := range override.CustomSlackMessages {
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
