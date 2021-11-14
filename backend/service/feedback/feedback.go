package feedback

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const (
	Name = "clutch.service.feedback"
)

func New(_ *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	p, ok := service.Registry[pgservice.Name]
	if !ok {
		return nil, fmt.Errorf("could not find the %v database service", pgservice.Name)
	}

	dbClient, ok := p.(pgservice.Client)
	if !ok {
		return nil, errors.New("database does not implement the required interface")
	}

	return &svc{storage: &storage{db: dbClient.DB()}, logger: logger, scope: scope}, nil
}

type svc struct {
	storage *storage
	logger  *zap.Logger
	scope   tally.Scope
}

type storage struct {
	db *sql.DB
}

type submission struct {
	id          string
	submittedAt time.Time
	feedback    *feedbackv1.Feedback
	metadata    *feedbackv1.FeedbackMetadata
}

type Service interface {
	SubmitFeedback(ctx context.Context, id string, feedback *feedbackv1.Feedback, metadata *feedbackv1.FeedbackMetadata) error
}

// TODO: check field values
func (s *svc) SubmitFeedback(ctx context.Context, id string, feedback *feedbackv1.Feedback, metadata *feedbackv1.FeedbackMetadata) error {
	feedbackSubmission := &submission{
		id:          id,
		submittedAt: time.Now(),
		feedback:    feedback,
		metadata:    metadata,
	}

	return s.storage.createOrUpdateSubmission(ctx, feedbackSubmission)
}
