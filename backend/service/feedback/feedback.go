package feedback

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	userId      string
	feedback    *feedbackv1.Feedback
	metadata    *feedbackv1.FeedbackMetadata
}

type Service interface {
	SubmitFeedback(ctx context.Context, id string, userId string, feedback *feedbackv1.Feedback, metadata *feedbackv1.FeedbackMetadata) error
}

func (s *svc) processSubmission(id string, userId string, feedback *feedbackv1.Feedback, metadata *feedbackv1.FeedbackMetadata) (*submission, error) {
	if len(id) < 36 {
		return nil, status.Error(codes.InvalidArgument, "client id was not a uuid length")
	}

	if userId == "" {
		return nil, status.Error(codes.InvalidArgument, "user id was empty")
	}

	if feedback == nil || metadata == nil {
		return nil, status.Errorf(codes.InvalidArgument, "feedback: %v or metadata: %v provided was nil", feedback, metadata)
	}

	if metadata.Survey == nil {
		return nil, status.Error(codes.InvalidArgument, "metadata survey provided was nil")
	}

	// the main question that was asked in the feedback component
	if metadata.Survey.Prompt == "" {
		return nil, status.Error(codes.InvalidArgument, "metadata survey prompt was empty")
	}

	if metadata.Survey.RatingOptions == nil {
		return nil, status.Error(codes.InvalidArgument, "metadata rating options was nil")
	}

	return &submission{
		id:          id,
		submittedAt: time.Now(),
		userId:      userId,
		feedback:    feedback,
		metadata:    metadata,
	}, nil
}

func (s *svc) SubmitFeedback(ctx context.Context, id string, userId string, feedback *feedbackv1.Feedback, metadata *feedbackv1.FeedbackMetadata) error {
	feedbackSubmission, err := s.processSubmission(id, userId, feedback, metadata)
	if err != nil {
		return err
	}
	return s.storage.createOrUpdateSubmission(ctx, feedbackSubmission)
}
