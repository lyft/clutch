package validate

// <!-- START clutchdoc -->
// description: Enforces input validation annotations from the proto definition on incoming requests.
// <!-- END clutchdoc -->

import (
	"github.com/golang/protobuf/ptypes/any"
	validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/lyft/clutch/backend/middleware"
)

const Name = "clutch.middleware.validate"

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (middleware.Middleware, error) {
	return &mid{}, nil
}

type mid struct{}

func (m *mid) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return validator.UnaryServerInterceptor()
}
