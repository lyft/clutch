package feedback

import (
	"testing"

	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap/zaptest"
)

func TestProcessSubmission(t *testing.T) {
	id := "00000000-0000-0000-0000-000000000000"
	userId := "foo@example.com"
	validFeedbackTestCase := &feedbackv1.Feedback{FeedbackType: "/k8s/deletePod", RatingLabel: "great", RatingScale: &feedbackv1.RatingScale{
		Type: &feedbackv1.RatingScale_Emoji{Emoji: feedbackv1.EmojiRating_HAPPY},
	}}
	vailidMetadataTestCase := &feedbackv1.FeedbackMetadata{
		Origin: feedbackv1.Origin_WIZARD,
		Survey: &feedbackv1.Survey{
			Prompt:       "Rate your experience",
			RatingLabels: []*feedbackv1.RatingLabel{{Type: &feedbackv1.RatingLabel_Emoji{Emoji: feedbackv1.EmojiRating_HAPPY}, Label: "great"}},
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
			feedback: &feedbackv1.Feedback{FeedbackType: "/k8s/deletePod", RatingLabel: "great"},
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
					Prompt:       "",
					RatingLabels: []*feedbackv1.RatingLabel{{Type: &feedbackv1.RatingLabel_Emoji{Emoji: feedbackv1.EmojiRating_HAPPY}, Label: "great"}},
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
					Prompt:       "Rate your experience",
					RatingLabels: nil,
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

func TestCalculateRatingScore(t *testing.T) {
	testCases := []struct {
		scale         *feedbackv1.RatingScale
		expectedScore int64
		expectedErr   bool
	}{
		{
			scale:         &feedbackv1.RatingScale{Type: &feedbackv1.RatingScale_Emoji{Emoji: feedbackv1.EmojiRating_SAD}},
			expectedScore: 30,
		},
		{
			scale:         &feedbackv1.RatingScale{Type: &feedbackv1.RatingScale_Emoji{Emoji: feedbackv1.EmojiRating_NEUTRAL}},
			expectedScore: 70,
		},
		{
			scale:         &feedbackv1.RatingScale{Type: &feedbackv1.RatingScale_Emoji{Emoji: feedbackv1.EmojiRating_HAPPY}},
			expectedScore: 100,
		},
		{
			scale:         &feedbackv1.RatingScale{Type: &feedbackv1.RatingScale_Emoji{Emoji: feedbackv1.EmojiRating_EMOJI_UNSPECIFIED}},
			expectedScore: -1,
			expectedErr:   true,
		},
	}

	for _, test := range testCases {
		result, err := calculateRatingScore(test.scale)
		if test.expectedErr {
			assert.Error(t, err)
			assert.Equal(t, test.expectedScore, result)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedScore, result)
		}
	}
}
