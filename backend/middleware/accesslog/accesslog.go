package log

// <!-- START clutchdoc -->
// description: Enables access log filtering on responses to filter logs by status code
// <!-- END clutchdoc -->

import (
	"context"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	middlewarev1 "github.com/lyft/clutch/backend/api/config/middleware/v1"
	"github.com/lyft/clutch/backend/gateway/log"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/middleware"
)

const Name = "clutch.middleware.access_log"

func New(config *middlewarev1.AccessLog, logger *zap.Logger, scope tally.Scope) (middleware.Middleware, error) {
	if config == nil || config.StatusCodeFilter == nil {
		// default to logging all status codes
		config = &middlewarev1.AccessLog{
			StatusCodeFilter: &middlewarev1.AccessLog_StatusCodeFilter{
				Operator: middlewarev1.AccessLog_StatusCodeFilter_GE,
				Value:    uint32(codes.OK),
			},
		}
	}
	return &mid{
		logger:     logger,
		operator:   config.StatusCodeFilter.Operator,
		statusCode: codes.Code(config.StatusCodeFilter.Value),
	}, nil
}

type mid struct {
	logger     *zap.Logger
	operator   middlewarev1.AccessLog_StatusCodeFilter_Op
	statusCode codes.Code
}

func (m *mid) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		s := status.Convert(err)
		if s == nil {
			s = status.New(codes.OK, "")
		}
		code := codes.Code(s.Proto().Code)

		if m.validStatusCode(code) {
			respBody, err := meta.APIBody(resp)
			if err != nil {
				return nil, err
			}
			m.logger.Error("GRPC error:",
				zap.Int("status code", int(code)),
				log.ProtoField("response body", respBody))
		}
		return resp, err
	}
}

func (m *mid) validStatusCode(c codes.Code) bool {
	switch m.operator {
	case middlewarev1.AccessLog_StatusCodeFilter_EQ:
		return c == m.statusCode
	case middlewarev1.AccessLog_StatusCodeFilter_GE:
		return c >= m.statusCode
	case middlewarev1.AccessLog_StatusCodeFilter_LE:
		return c <= m.statusCode
	case middlewarev1.AccessLog_StatusCodeFilter_NE:
		return c != m.statusCode
	default:
		return true
	}
}
