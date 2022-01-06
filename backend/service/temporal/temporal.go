package temporal

// <!-- START clutchdoc -->
// description: Workflow client for temporal.io.
// <!-- END clutchdoc -->

import (
	"fmt"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	temporalv1 "github.com/lyft/clutch/backend/api/config/service/temporal/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.temporal"

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &temporalv1.Config{}
	if err := cfg.UnmarshalTo(config); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("implement me")
}
