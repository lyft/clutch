package feedback

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
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

	return &svc{db: dbClient.DB(), logger: logger, scope: scope}, nil
}

type svc struct {
	db     *sql.DB
	logger *zap.Logger
	scope  tally.Scope
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

const createOrUpdateSubmissionQuery = `
INSERT INTO feedback(client_id, submitted_at, feedback, metadata) VALUES ($1, $2, $3, $4)
ON CONFLICT (client_id) DO UPDATE SET
		client_id = EXCLUDED.client_id,
		submitted_at = EXCLUDED.submitted_at,
		feedback = EXCLUDED.feedback,
		metadata = EXCLUDED.metadata,
`

func (s *svc) createOrUpdateSubmission(ctx context.Context, submission *submission) error {
	feedbackJSON, err := protojson.Marshal(submission.feedback)
	if err != nil {
		return err
	}
	metadataJSON, err := protojson.Marshal(submission.metadata)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, createOrUpdateSubmissionQuery, submission.id, submission.submittedAt, feedbackJSON, metadataJSON)
	return err
}

func (s *svc) SubmitFeedback(ctx context.Context, id string, feedback *feedbackv1.Feedback, metadata *feedbackv1.FeedbackMetadata) error {
	feedbackSubmission := &submission{
		id:          id,
		submittedAt: time.Now(),
		feedback:    feedback,
		metadata:    metadata,
	}
	// TODO: create a sql file?
	return s.createOrUpdateSubmission(ctx, feedbackSubmission)
}
