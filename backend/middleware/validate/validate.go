package validate

// <!-- START clutchdoc -->
// description: Enforces input validation annotations from the proto definition on incoming requests.
// <!-- END clutchdoc -->

import (
	validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/lyft/clutch/backend/middleware"
)

const Name = "clutch.middleware.validate"

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (middleware.Middleware, error) {
	return &mid{}, nil
}

type mid struct{}

func (m *mid) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return validator.UnaryServerInterceptor()
}
