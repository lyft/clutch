package k8s

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"testing"
)

func TestErrorToStatus(t *testing.T) {
	var err error
	err = k8serrors.NewUnauthorized("nice try")
	err = k8sStatusToGRPCStatusError(err)
	spb, ok :=  status.FromError(err)
	assert.True(t, ok)
	assert.NotNil(t, spb)
	assert.Equal(t, "nice try", spb.Message())
	assert.EqualValues(t, codes.Unauthenticated, spb.Code())
}
