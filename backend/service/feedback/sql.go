package feedback

import (
	"context"

	"google.golang.org/protobuf/encoding/protojson"
)

const createOrUpdateSubmissionQuery = `
INSERT INTO feedback (client_id, submitted_at, details, metadata) VALUES ($1, $2, $3, $4)
ON CONFLICT (client_id) DO UPDATE SET
		client_id = EXCLUDED.client_id,
		submitted_at = EXCLUDED.submitted_at,
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
	_, err = s.db.ExecContext(ctx, createOrUpdateSubmissionQuery, submission.id, submission.submittedAt, feedbackJSON, metadataJSON)
	return err
}
