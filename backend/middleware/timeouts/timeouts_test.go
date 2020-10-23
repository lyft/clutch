package timeouts

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
)

func TestNew(t *testing.T) {
	tests := []struct {
		config *gatewayv1.Timeouts
	}{
		{config: nil},
		{config: &gatewayv1.Timeouts{Default: durationpb.New(time.Second)}},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			assert.NoError(t, tt.config.Validate())

			m, err := New(tt.config, nil, nil)
			assert.NoError(t, err)
			assert.NotNil(t, m)
		})
	}
}
