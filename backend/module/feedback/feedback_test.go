package feedback

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"

	feedbackv1cfg "github.com/lyft/clutch/backend/api/config/module/feedback/v1"
	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
	"github.com/lyft/clutch/backend/module/moduletest"
)

func TestModule(t *testing.T) {
	config, _ := anypb.New(&feedbackv1cfg.Config{})
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	m, err := New(config, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("clutch.feedback.v1.FeedbackAPI"))
	assert.True(t, r.JSONRegistered())
}

func TestNewSurveyLookup(t *testing.T) {
	testCases := []struct {
		surveyOrigin  []*feedbackv1cfg.SurveyOrigin
		expectedEmpty bool
	}{
		{surveyOrigin: []*feedbackv1cfg.SurveyOrigin{}, expectedEmpty: true},
		{surveyOrigin: []*feedbackv1cfg.SurveyOrigin{{
			Origin: feedbackv1.Origin_WIZARD,
			Survey: &feedbackv1cfg.Survey{Prompt: "bar"},
		}},
		},
	}

	for _, test := range testCases {
		result := newSurveyLookup(test.surveyOrigin)
		if test.expectedEmpty {
			assert.Empty(t, result)
		} else {
			assert.Equal(t, 1, len(result.surveys))
			v, ok := result.surveys["WIZARD"]
			assert.True(t, ok)
			assert.Equal(t, "bar", v.Survey.Prompt)
		}
	}
}

func TestGetConfigSurveys(t *testing.T) {
	// match
	testCases := []struct {
		surveyMap            SurveyLookup
		expectedOk           bool
		expectedSurveyPrompt string
	}{
		{
			surveyMap:  SurveyLookup{surveys: map[string]*feedbackv1cfg.SurveyOrigin{}},
			expectedOk: false,
		},
		{
			surveyMap: SurveyLookup{surveys: map[string]*feedbackv1cfg.SurveyOrigin{
				"FOO": &feedbackv1cfg.SurveyOrigin{},
			}},
			expectedOk: false,
		},
		{
			surveyMap: SurveyLookup{surveys: map[string]*feedbackv1cfg.SurveyOrigin{
				"WIZARD": &feedbackv1cfg.SurveyOrigin{Survey: &feedbackv1cfg.Survey{Prompt: "bar"}},
			}},
			expectedOk: true,
		},
	}

	for _, test := range testCases {
		v, ok := test.surveyMap.getConfigSurveys(feedbackv1.Origin_WIZARD)
		if !test.expectedOk {
			assert.False(t, ok)
			assert.Nil(t, v)
		} else {
			assert.True(t, ok)
			assert.Equal(t, "bar", v.Prompt)
		}
	}
}
