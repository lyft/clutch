package experimentstore

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TimeToPropertyDateValue(t *time.Time) (*experimentation.Property_DateValue, error) {
	if t == nil {
		return nil, nil
	}

    timestamp := timestamppb.New(t.Time)
    if err := timestamp.CheckValid(); err != nil {
        return nil, err
    }

	return &experimentation.Property_DateValue{DateValue: timestamp}, nil
}
