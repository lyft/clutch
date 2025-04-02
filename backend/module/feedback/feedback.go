package feedback

// <!-- START clutchdoc -->
// description: Exposes endpoints to return survey questions for feedback components and to submit feedback submissions.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"

	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	feedbackv1cfg "github.com/lyft/clutch/backend/api/config/module/feedback/v1"
	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/feedback"
)

const (
	Name = "clutch.module.feedback"
)

func New(cfg *anypb.Any, log *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &feedbackv1cfg.Config{}
	err := cfg.UnmarshalTo(config)
	if err != nil {
		return nil, err
	}

	feedbackClient, ok := service.Registry["clutch.service.feedback"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	c, ok := feedbackClient.(feedback.Service)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	m := &mod{
		surveyMap: newSurveyLookup(config.Origins),
		client:    c,
		logger:    log,
		scope:     scope,
	}

	return m, nil
}

// SurveyLookup stores a map of the origin to survey mapping that are
// provided in the feedback module config
type surveyLookup struct {
	// key is the origin name
	surveys map[string]*feedbackv1cfg.SurveyOrigin
}

type mod struct {
	surveyMap surveyLookup
	client    feedback.Service
	logger    *zap.Logger
	scope     tally.Scope
}

func (m *mod) Register(r module.Registrar) error {
	feedbackv1.RegisterFeedbackAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(feedbackv1.RegisterFeedbackAPIHandler)
}

func (m *mod) GetSurveys(tx context.Context, req *feedbackv1.GetSurveysRequest) (*feedbackv1.GetSurveysResponse, error) {
	results := make(map[string]*feedbackv1.Survey)

	// this scenario shouldn't happen as the feedback config fields cannot be empty
	if len(m.surveyMap.surveys) == 0 {
		m.logger.Error("survey questions were not found in the config",
			zap.Any("config", m.surveyMap.surveys),
		)
		return nil, status.Errorf(codes.NotFound, "survey questions were not found")
	}

	for _, origin := range req.Origins {
		v, ok := m.surveyMap.getConfigSurveys(origin)
		if !ok {
			continue
		}
		results[origin.String()] = &feedbackv1.Survey{
			Prompt:         v.Prompt,
			FreeformPrompt: v.FreeformPrompt,
			RatingLabels:   v.RatingLabels,
		}
	}

	// this scenario shouldn't happen as the request and config have to provide defined enum types
	// TODO: if multiple orgins are requested at once and only some surveys are found, return partial error
	if len(results) == 0 {
		msg := "surveys not found for the requested origin(s)"
		m.logger.Error(
			msg,
			zap.Any("request origins", req.Origins),
			zap.Any("config origins", m.surveyMap.surveys),
		)
		return nil, status.Errorf(codes.InvalidArgument, "%s", msg)
	}

	return &feedbackv1.GetSurveysResponse{
		OriginSurvey: results,
	}, nil
}

func newSurveyLookup(origins []*feedbackv1cfg.SurveyOrigin) surveyLookup {
	if len(origins) == 0 {
		return surveyLookup{}
	}

	surveys := make(map[string]*feedbackv1cfg.SurveyOrigin, len(origins))
	for _, origin := range origins {
		surveys[origin.Origin.String()] = origin
	}

	return surveyLookup{
		surveys,
	}
}

// GetConfigSurveys uses the SurveyLookup and the origin passed in the API request to
// check if a survey exists for the origin. If found, the survey is returned.
// Otherwise ok is false.
func (sl surveyLookup) getConfigSurveys(origin feedbackv1.Origin) (*feedbackv1cfg.Survey, bool) {
	if len(sl.surveys) == 0 {
		return nil, false
	}
	v, ok := sl.surveys[origin.String()]
	if !ok {
		return nil, false
	}
	return v.Survey, true
}

func (m *mod) SubmitFeedback(ctx context.Context, req *feedbackv1.SubmitFeedbackRequest) (*feedbackv1.SubmitFeedbackResponse, error) {
	if err := m.client.SubmitFeedback(ctx, req.Id, req.UserId, req.Feedback, req.Metadata); err != nil {
		m.logger.Error("failed to submit feedback", zap.Error(err))
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	return &feedbackv1.SubmitFeedbackResponse{}, nil
}
