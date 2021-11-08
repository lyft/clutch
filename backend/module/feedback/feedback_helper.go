package feedback

import (
	feedbackv1cfg "github.com/lyft/clutch/backend/api/config/module/feedback/v1"
	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
)

// SurveyLookup stores a map of the origin to survey mapping that are
// provided in the feedback module config
type SurveyLookup struct {
	// key is the origin name
	surveys map[string]*feedbackv1cfg.SurveyOrigin
}

// TODO: add test
func NewSurveyLookup(origins []*feedbackv1cfg.SurveyOrigin) SurveyLookup {
	if len(origins) == 0 {
		return SurveyLookup{}
	}

	surveys := make(map[string]*feedbackv1cfg.SurveyOrigin, len(origins))
	for _, origin := range origins {
		surveys[origin.Origin.String()] = origin
	}

	return SurveyLookup{
		surveys,
	}
}

// TODO: add test
// GetConfigSurveys uses the SurveyLookup and the origin passed in the API request to
// check if a survey exists for the origin. If found, the survey is returned.
// Otherwise ok is false.
func (sl SurveyLookup) GetConfigSurveys(origin feedbackv1.Origin) (*feedbackv1cfg.Survey, bool) {
	if len(sl.surveys) == 0 {
		return nil, false
	}
	v, ok := sl.surveys[origin.String()]
	if !ok {
		return nil, false
	}
	return v.Survey, true
}
