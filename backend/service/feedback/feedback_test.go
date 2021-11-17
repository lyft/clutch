package feedback

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
)

func TestProcessSubmission(t *testing.T) {
	id := "00000000-0000-0000-0000-000000000000"
	userId := "foo@example.com"
	validFeedbackTestCase := &feedbackv1.Feedback{UrlPath: "/k8s/deletePod", Rating: "great"}
	vailidMetadataTestCase := &feedbackv1.FeedbackMetadata{
		Origin: feedbackv1.Origin_WIZARD,
		Survey: &feedbackv1.Survey{
			Prompt:        "Rate your experience",
			RatingOptions: &feedbackv1.RatingOptions{Three: "great"},
		},
		UserSubmitted: true,
	}

	testErrorCases := []struct {
		id       string
		userId   string
		feedback *feedbackv1.Feedback
		metadata *feedbackv1.FeedbackMetadata
	}{
		// feedback is nil
		{
			id:       id,
			userId:   userId,
			feedback: nil,
			metadata: vailidMetadataTestCase,
		},
		// id is not the required length
		{
			id:       "00000000-0000-0000-0000",
			userId:   userId,
			feedback: validFeedbackTestCase,
			metadata: vailidMetadataTestCase,
		},
		// userId is empty
		{
			id:       id,
			userId:   "",
			feedback: &feedbackv1.Feedback{UrlPath: "/k8s/deletePod", Rating: "great"},
			metadata: vailidMetadataTestCase,
		},
		// metadata is nil
		{
			id:       id,
			userId:   userId,
			feedback: validFeedbackTestCase,
			metadata: nil,
		},
		// metadata survey prompt is empty
		{
			id:       id,
			userId:   userId,
			feedback: validFeedbackTestCase,
			metadata: &feedbackv1.FeedbackMetadata{
				Origin: feedbackv1.Origin_WIZARD,
				Survey: &feedbackv1.Survey{
					Prompt:        "",
					RatingOptions: &feedbackv1.RatingOptions{Three: "great"},
				}, UserSubmitted: true},
		},
		// metadata rating options is nil
		{
			id:       id,
			userId:   userId,
			feedback: validFeedbackTestCase,
			metadata: &feedbackv1.FeedbackMetadata{
				Origin: feedbackv1.Origin_WIZARD,
				Survey: &feedbackv1.Survey{
					Prompt:        "Rate your experience",
					RatingOptions: nil,
				},
				UserSubmitted: true,
			},
		},
	}

	logger := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	storage := &storage{}
	s := &svc{storage, logger, scope}

	// happy path
	submission, err := s.processSubmission(id, userId, validFeedbackTestCase, vailidMetadataTestCase)
	assert.NoError(t, err)
	assert.NotNil(t, submission)

	// error scenarios
	for _, test := range testErrorCases {
		submission, err := s.processSubmission(test.id, test.userId, test.feedback, test.metadata)
		assert.Error(t, err)
		assert.Nil(t, submission)
	}
}
