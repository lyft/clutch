package authn

import (
	"testing"

	"github.com/stretchr/testify/assert"

	apimock "github.com/lyft/clutch/backend/mock/api"
)

func TestNotImpl(t *testing.T) {
	a := apimock.AnyFromYAML(`
"@type": types.google.com/clutch.config.service.authn.v1.Config
session_secret: my_session_secret
`)
	svc, err := New(a, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, svc)
}

// TODO: until services have plumbed context in the factory, this would require the test to call out to the internet to succeeed, at least for the OIDC provider
// func TestReturnsCorrectProvider(t *testing.T)
