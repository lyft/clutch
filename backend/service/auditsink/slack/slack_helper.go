package slack

import (
	"fmt"

	configv1 "github.com/lyft/clutch/backend/api/config/service/auditsink/slack/v1"
)

// NewOverrideLookup creates a map from the overrides. The generated OverrideLookup is
// used to easily retrieve an override by key (/service/method).
func NewOverrideLookup(overrides []*configv1.CustomMessage) *OverrideLookup {
	if len(overrides) == 0 {
		// no overrides
		return nil
	}

	messages := make(map[string]*configv1.CustomMessage, len(overrides))
	for _, override := range overrides {
		messages[override.FullMethod] = override
	}

	return &OverrideLookup{
		messages,
	}
}

// GetOverrideMessage uses the OverrideLookup and the audit event's service + method name to
// check if an override exists for the key service/method. If found, the custom message is returned.
// Otherwise ok is false.
func (o *OverrideLookup) GetOverrideMessage(service, method string) (string, bool) {
	// /service/method
	// ex. /clutch.k8s.v1.K8sAPI/DescribePod
	pattern := "/%s/%s"

	if o == nil || len(o.messages) == 0 {
		// no overrides
		return "", false
	}

	cm, ok := o.messages[fmt.Sprintf(pattern, service, method)]
	if !ok {
		// no custom message for this /service/method
		return "", false
	}
	return cm.Message, true
}
