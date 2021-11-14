package feedbackmock

import (
	"context"

	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
	"github.com/lyft/clutch/backend/service/feedback"
)

type svc struct{}

func (s svc) SubmitFeedback(ctx context.Context, id string, feedback *feedbackv1.Feedback, metadata *feedbackv1.FeedbackMetadata) error {
	return nil
}

func New() feedback.Service {
	return &svc{}
}
