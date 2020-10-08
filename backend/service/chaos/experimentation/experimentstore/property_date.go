package experimentstore

import (
	"database/sql"
	"github.com/golang/protobuf/ptypes"
	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TimeToPropertyDateValue(t sql.NullTime) (*experimentation.Property_DateValue, error) {
	if t.Valid {
		timestamp, err := ptypes.TimestampProto(t.Time)
		if err == nil {
			return &experimentation.Property_DateValue{DateValue: timestamp}, nil
		}
	}

	return nil, nil
}
