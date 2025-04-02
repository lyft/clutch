package authzmock

import (
	"context"

	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	authzv1 "github.com/lyft/clutch/backend/api/authz/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/authz"
)

type svc struct{}

func (s svc) Check(ctx context.Context, request *authzv1.CheckRequest) (*authzv1.CheckResponse, error) {
	panic("implement me")
}

func New() authz.Client {
	return &svc{}
}

func NewAsService(*anypb.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}
