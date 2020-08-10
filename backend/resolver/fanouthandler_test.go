package resolver

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
)

type basicValue struct {
	value proto.Message
	err   error
}

var values = []basicValue{
	{value: nil, err: errors.New("whoa")},
	{value: &healthcheckv1.HealthcheckResponse{}, err: nil},
}

func TestFanoutHandler(t *testing.T) {
	_, handler := NewFanoutHandler(context.Background())

	for _, value := range values {
		handler.Add(1)
		go func(v basicValue) {
			defer handler.Done()
			select {
			case handler.Channel() <- NewSingleFanoutResult(v.value, v.err):
				return
			case <-handler.Cancelled():
				return
			}
		}(value)
	}

	results, err := handler.Results(0)
	assert.NoError(t, err)
	assert.Len(t, results.Messages, 1)
	assert.Len(t, results.PartialFailures, 1)
	assert.Contains(t, results.Messages, values[1].value)
}
