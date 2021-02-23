package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

func TestConvertError(t *testing.T) {
	service := &svc{}

	err := k8serrors.NewUnauthorized("nice try")
	newErr := service.InterceptError(err)

	s, ok := status.FromError(newErr)
	assert.True(t, ok)
	assert.NotNil(t, s)
	assert.Equal(t, codes.Unauthenticated, s.Code())
	assert.Equal(t, "nice try", s.Message())
}
