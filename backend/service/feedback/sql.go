package feedback

import (
	"context"

	"google.golang.org/protobuf/encoding/protojson"
)

const createOrUpdateSubmissionQuery = `
INSERT INTO feedback (client_id, submitted_at, user_id, score, details, metadata) VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (client_id) DO UPDATE SET
		client_id = EXCLUDED.client_id,
		submitted_at = EXCLUDED.submitted_at,
		user_id = EXCLUDED.user_id,
		score = EXCLUDED.score,
		details = EXCLUDED.details,
		metadata = EXCLUDED.metadata
`

func (s *storage) createOrUpdateSubmission(ctx context.Context, submission *submission) error {
	feedbackJSON, err := protojson.Marshal(submission.feedback)
	if err != nil {
		return err
	}
	metadataJSON, err := protojson.Marshal(submission.metadata)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, createOrUpdateSubmissionQuery, submission.id, submission.submittedAt, submission.userId, submission.score, feedbackJSON, metadataJSON)
	return err
}
