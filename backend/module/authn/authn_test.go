package authn

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	authnv1 "github.com/lyft/clutch/backend/api/authn/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
)

func TestAllRequestResponseAreRedacted(t *testing.T) {
	obj := []proto.Message{
		&authnv1.CallbackRequest{},
		&authnv1.CallbackResponse{},
		&authnv1.LoginRequest{},
		&authnv1.LoginResponse{},
	}

	for _, o := range obj {
		assert.True(t, meta.IsRedacted(o))
	}
}
