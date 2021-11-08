package feedback

import (
	"context"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	feedbackv1cfg "github.com/lyft/clutch/backend/api/config/module/feedback/v1"
	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
	"github.com/lyft/clutch/backend/module"
)

const (
	Name = "clutch.module.feedback"
)

func New(cfg *any.Any, log *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &feedbackv1cfg.Config{}
	err := cfg.UnmarshalTo(config)
	if err != nil {
		return nil, err
	}

	m := &mod{
		surveyMap: NewSurveyLookup(config.Origins),
		logger:    log,
		scope:     scope,
	}

	return m, nil
}

type mod struct {
	surveyMap SurveyLookup
	logger    *zap.Logger
	scope     tally.Scope
}

func (m *mod) Register(r module.Registrar) error {
	feedbackv1.RegisterFeedbackAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(feedbackv1.RegisterFeedbackAPIHandler)
}

// TODO: tests
func (m *mod) GetSurveys(tx context.Context, req *feedbackv1.GetSurveysRequest) (*feedbackv1.GetSurveysResponse, error) {
	results := make(map[string]*feedbackv1.Survey)

	// this scenario shouldn't happen as the feedback config cannot be empty
	if len(m.surveyMap.surveys) == 0 {
		return nil, status.Errorf(codes.NotFound, "survey questions were not found")
	}

	for _, origin := range req.Origins {
		v, ok := m.surveyMap.GetConfigSurveys(origin)
		if ok {
			results[origin.String()] = &feedbackv1.Survey{
				Prompt:         v.Prompt,
				FreeformPrompt: v.FreeformPrompt,
				RatingOptions: &feedbackv1.RatingOptions{
					One:   v.RatingOptions.One,
					Two:   v.RatingOptions.Two,
					Three: v.RatingOptions.Three,
				},
			}
		}
	}

	// this scenario shouldn't happen as the request and config have to provide defined enum types
	// TODO: if multiple orgins are requested at once and only some surveys are found, return partial error
	if len(results) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "origins must be one of HEADER or WIZARD")
	}

	return &feedbackv1.GetSurveysResponse{
		OriginSurvey: results,
	}, nil
}
