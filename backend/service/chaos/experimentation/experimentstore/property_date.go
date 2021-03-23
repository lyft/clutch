package experimentstore

import (
	"database/sql"

	"google.golang.org/protobuf/types/known/timestamppb"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TimeToPropertyDateValue(t sql.NullTime) (*experimentation.Property_DateValue, error) {
	if t.Valid {
		timestamp := timestamppb.New(t.Time)
		if err := timestamp.CheckValid(); err != nil {
			return nil, err
		}

		return &experimentation.Property_DateValue{DateValue: timestamp}, nil
	}

	return nil, nil
}
