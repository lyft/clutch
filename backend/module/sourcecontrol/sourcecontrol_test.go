package sourcecontrol

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	sourcecontrolv1 "github.com/lyft/clutch/backend/api/sourcecontrol/v1"
	"github.com/lyft/clutch/backend/mock/service/githubmock"
	"github.com/lyft/clutch/backend/module/moduletest"
	"github.com/lyft/clutch/backend/service"
)

func TestModule(t *testing.T) {
	service.Registry["clutch.service.github"] = githubmock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	m, err := New(nil, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("clutch.sourcecontrol.v1.SourceControlAPI"))
	assert.True(t, r.JSONRegistered())

	a, ok := m.(sourcecontrolv1.SourceControlAPIServer)
	assert.Equal(t, true, ok)

	resp, err := a.CreateRepository(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

var containsTests = []struct {
	source         []string
	expectedResult bool
}{
	{
		source:         []string{"foo", "bar"},
		expectedResult: true,
	},
	{
		source:         []string{"bar"},
		expectedResult: false,
	},
}

func TestContains(t *testing.T) {
	for idx, tt := range containsTests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			result := contains(tt.source, "foo")

			assert.Equal(t, result, tt.expectedResult)
		})
	}
}
