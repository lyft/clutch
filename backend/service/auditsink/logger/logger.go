package logger

// <!-- START clutchdoc -->
// description: Logs audit events to the Gateway's logger.
// <!-- END clutchdoc -->

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/auditsink"
)

const Name = "clutch.service.audit.sink.logger"

func New(cfg *any.Any, logger *zap.Logger, _ tally.Scope) (service.Service, error) {
	var filter *configv1.Filter

	config := &configv1.SinkConfig{}
	if err := ptypes.UnmarshalAny(cfg, config); err == nil {
		filter = config.Filter
	}

	s := &svc{logger: logger, filter: filter}
	return s, nil
}

type svc struct {
	logger *zap.Logger
	filter *configv1.Filter
}

func (s *svc) Write(event *auditv1.Event) error {
	if !auditsink.Filter(s.filter, event) {
		s.logger.Info("new audit event", zap.Any("event", event))
	}
	return nil
}
