package experimentstore

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TimeToPropertyDateValue(t *time.Time) (*experimentation.Property_DateValue, error) {
	if t == nil {
		return nil, nil
	}

	timestamp, err := ptypes.TimestampProto(*t)
	if err != nil {
		return nil, err
	}

	return &experimentation.Property_DateValue{DateValue: timestamp}, nil
}
