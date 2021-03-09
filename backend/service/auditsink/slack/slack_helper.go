package slack

import (
	"fmt"

	configv1 "github.com/lyft/clutch/backend/api/config/service/auditsink/slack/v1"
)

// OverrideLookup stores a map of the FullMethod to CustomMessage(s) that are
// provided in the slack config
type OverrideLookup struct {
	messages map[string]*configv1.CustomMessage
}

// NewOverrideLookup creates a map from the overrides. The generated OverrideLookup is
// used to easily retrieve an override by key (/service/method).
func NewOverrideLookup(overrides []*configv1.CustomMessage) OverrideLookup {
	if len(overrides) == 0 {
		// no overrides
		return OverrideLookup{}
	}

	messages := make(map[string]*configv1.CustomMessage, len(overrides))
	for _, override := range overrides {
		messages[override.FullMethod] = override
	}

	return OverrideLookup{
		messages,
	}
}

// GetOverrideMessage uses the OverrideLookup and the audit event's service + method name to
// check if an override exists for the key service/method. If found, the custom message is returned.
// Otherwise ok is false.
func (o OverrideLookup) GetOverrideMessage(service, method string) (string, bool) {
	if len(o.messages) == 0 {
		// no overrides
		return "", false
	}

	// /service/method
	// ex. /clutch.k8s.v1.K8sAPI/DescribePod
	pattern := "/%s/%s"
	cm, ok := o.messages[fmt.Sprintf(pattern, service, method)]
	if !ok {
		// no custom message for this /service/method
		return "", false
	}
	return cm.Message, true
}
