package logger

// <!-- START clutchdoc -->
// description: Logs audit events to the Gateway's logger.
// <!-- END clutchdoc -->

import (
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
	"github.com/lyft/clutch/backend/gateway/log"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/auditsink"
)

const Name = "clutch.service.audit.sink.logger"

func New(cfg *anypb.Any, logger *zap.Logger, _ tally.Scope) (service.Service, error) {
	var filter *configv1.Filter

	config := &configv1.SinkConfig{}
	if err := cfg.UnmarshalTo(config); err == nil {
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
	if auditsink.Filter(s.filter, event) {
		s.logger.Info("new audit event", log.ProtoField("event", event))
	}
	return nil
}
