package shortlinkmock

import (
	"context"

	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	shortlinkv1 "github.com/lyft/clutch/backend/api/shortlink/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/shortlink"
)

type svc struct{}

func New() shortlink.Service {
	return &svc{}
}

func NewAsService(*anypb.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}

func (s *svc) Create(ctx context.Context, path string, state []*shortlinkv1.ShareableState) (string, error) {
	return "mockhash", nil
}

func (s *svc) Get(ctx context.Context, hash string) (string, []*shortlinkv1.ShareableState, error) {
	return "/mockpath", []*shortlinkv1.ShareableState{
		{
			Key: "mock",
			State: &structpb.Value{
				Kind: &structpb.Value_StringValue{StringValue: "mock string"},
			},
		},
	}, nil
}
