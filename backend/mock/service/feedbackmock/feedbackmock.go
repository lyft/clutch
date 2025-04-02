package feedbackmock

import (
	"context"

	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/feedback"
)

type svc struct{}

func New() feedback.Service {
	return &svc{}
}

func NewAsService(*anypb.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}

// TODO: add error handling scenarios for mock testing
func (s svc) SubmitFeedback(ctx context.Context, id string, userId string, feedback *feedbackv1.Feedback, metadata *feedbackv1.FeedbackMetadata) error {
	return nil
}
