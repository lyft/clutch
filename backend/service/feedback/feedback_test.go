package feedback

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
)

func TestProcessSubmission(t *testing.T) {
	id := "uniquefeedbackid"
	validFeedbackTestCase := &feedbackv1.Feedback{UserId: "foo@example.com", UrlPath: "/k8s/deletePod", Rating: "great"}
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
		feedback *feedbackv1.Feedback
		metadata *feedbackv1.FeedbackMetadata
	}{
		// feedback is nil
		{
			id:       id,
			feedback: nil,
			metadata: vailidMetadataTestCase,
		},
		// id is empty
		{
			id:       "",
			feedback: validFeedbackTestCase,
			metadata: vailidMetadataTestCase,
		},
		// userId is empty
		{
			id:       id,
			feedback: &feedbackv1.Feedback{UrlPath: "/k8s/deletePod", Rating: "great"},
			metadata: vailidMetadataTestCase,
		},
		// urlPath is empty
		{
			id:       id,
			feedback: &feedbackv1.Feedback{UserId: "foo@example.com", Rating: "great"},
			metadata: vailidMetadataTestCase,
		},
		// rating is empty
		{
			id:       id,
			feedback: &feedbackv1.Feedback{UserId: "foo@example.com", UrlPath: "/k8s/deletePod"},
			metadata: vailidMetadataTestCase,
		},
		// metadata is nil
		{
			id:       id,
			feedback: validFeedbackTestCase,
			metadata: nil,
		},
		// survey is nil
		{
			id:       id,
			feedback: validFeedbackTestCase,
			metadata: &feedbackv1.FeedbackMetadata{
				Origin:        feedbackv1.Origin_WIZARD,
				Survey:        nil,
				UserSubmitted: true,
			},
		},
		// origin is unspecified
		{
			id:       id,
			feedback: validFeedbackTestCase,
			metadata: &feedbackv1.FeedbackMetadata{
				Origin: feedbackv1.Origin_ORIGIN_UNSPECIFIED,
				Survey: &feedbackv1.Survey{
					Prompt:        "Rate your experience",
					RatingOptions: &feedbackv1.RatingOptions{Three: "great"},
				},
				UserSubmitted: true,
			},
		},
		// metadata survey prompt is empty
		{
			id:       id,
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
	submission, err := s.processSubmission(id, validFeedbackTestCase, vailidMetadataTestCase)
	assert.NoError(t, err)
	assert.NotNil(t, submission)

	// error scenarios
	for _, test := range testErrorCases {
		submission, err := s.processSubmission(test.id, test.feedback, test.metadata)
		assert.Error(t, err)
		assert.Nil(t, submission)
	}
}
