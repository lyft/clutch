package feedback

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"

	feedbackv1 "github.com/lyft/clutch/backend/api/feedback/v1"
)

func TestCreateOrUpdateSubmission(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	testCases := []struct {
		feedback *feedbackv1.Feedback
		metadata *feedbackv1.FeedbackMetadata
	}{
		{
			feedback: &feedbackv1.Feedback{Rating: "great"},
			metadata: &feedbackv1.FeedbackMetadata{UserSubmitted: true},
		},
		// a sanity check that nil values aren't going to raise errors
		{
			feedback: nil,
			metadata: nil,
		},
	}

	for _, test := range testCases {
		s := &submission{
			id:          "00000000-0000-0000-0000-000000000000",
			userId:      "foo@example.com",
			submittedAt: time.Now(),
			feedback:    test.feedback,
			metadata:    test.metadata,
		}

		feedbackJSON, err := protojson.Marshal(s.feedback)
		assert.NoError(t, err)
		metadataJSON, err := protojson.Marshal(s.metadata)
		assert.NoError(t, err)

		mock.ExpectExec(createOrUpdateSubmissionQuery).
			WithArgs(s.id, s.submittedAt, s.userId, feedbackJSON, metadataJSON).
			WillReturnResult(sqlmock.NewResult(0, 1))

		r := &storage{db: db}
		err = r.createOrUpdateSubmission(context.Background(), s)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}
