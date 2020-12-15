package accesslog

// <!-- START clutchdoc -->
// description: Logs gRPC requests and responses, optionally filtered by status code.
// <!-- END clutchdoc -->

import (
	"context"
	"fmt"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accesslogv1 "github.com/lyft/clutch/backend/api/config/middleware/accesslog/v1"
	"github.com/lyft/clutch/backend/gateway/log"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/middleware"
)

const Name = "clutch.middleware.accesslog"

func New(config *accesslogv1.Config, logger *zap.Logger, scope tally.Scope) (middleware.Middleware, error) {
	var statusCodes []codes.Code
	// if no filter is provided default to logging all status codes
	if config != nil {
		for _, filter := range config.StatusCodeFilters {
			switch t := filter.GetFilterType().(type) {
			case *accesslogv1.Config_StatusCodeFilter_Equals:
				statusCode := filter.GetEquals()
				statusCodes = append(statusCodes, codes.Code(statusCode))
			default:
				return nil, fmt.Errorf("status code filter `%T` not supported", t)
			}
		}
	}
	return &mid{
		logger:      logger,
		statusCodes: statusCodes,
	}, nil
}

type mid struct {
	logger *zap.Logger
	// TODO(perf): improve lookup efficiency using a lookup table
	statusCodes []codes.Code
}

func (m *mid) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		s := status.Convert(err)
		if s == nil {
			s = status.New(codes.OK, "")
		}
		code := s.Code()
		// common logger context fields
		fields := []zap.Field{
			zap.Int("statusCode", int(code)),
			zap.String("status", code.String()),
		}

		if m.validStatusCode(code) {
			// if err is returned from handler, log error details only
			// as response body will be nil
			if err != nil {
				fields = append(fields, zap.String("message", s.Message()))
				m.logger.Error("gRPC", fields...)
			} else {
				respBody, err := meta.APIBody(resp)
				if err != nil {
					return nil, err
				}
				fields = append(fields, log.ProtoField("responseBody", respBody))
				m.logger.Info("GRPC:", fields...)
			}
		}
		return resp, err
	}
}

func (m *mid) validStatusCode(c codes.Code) bool {
	// If no filter is provided all status codes are valid
	if len(m.statusCodes) == 0 {
		return true
	}
	for _, code := range m.statusCodes {
		if c == code {
			return true
		}
	}
	return false
}
