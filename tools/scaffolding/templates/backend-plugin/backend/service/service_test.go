package {{ .ServiceName }}

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	service, err := New(nil, nil, nil)
	assert.NoError(t, err)
	assert.Implements(t, (*Client)(nil), service)
}

var greetTests = []struct {
	name    string
	message string
}{
	{
		name:    "John",
		message: "Hello, John",
	},
	{
		name:    "Jane",
		message: "Hello, Jane",
	},
}

func TestGreet(t *testing.T) {
	for _, tt := range greetTests {
		t.Run(tt.name, func(t *testing.T) {
			service := &svc{}
			res, err := service.Greet(context.Background(), tt.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.message, res.Message)
		})
	}
}
