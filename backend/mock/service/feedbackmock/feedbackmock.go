package feedbackmock

import (
	"context"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/feedback"
)

type svc struct{}

func New() feedback.Service {
	return &svc{}
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}

func (s svc) SubmitFeedback(ctx context.Context, id string, userId string, feedback *feedbackv1.Feedback, metadata *feedbackv1.FeedbackMetadata) error {
	return nil
}
