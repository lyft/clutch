package project

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	projectv1core "github.com/lyft/clutch/backend/api/core/project/v1"
	projectv1 "github.com/lyft/clutch/backend/api/project/v1"
	"github.com/lyft/clutch/backend/mock/service/dbmock"
	"github.com/lyft/clutch/backend/service"
)

func TestNew(t *testing.T) {
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	// Should throw error looking for postgres
	_, err := New(&anypb.Any{}, log, scope)
	assert.Error(t, err)

	// Should not error with postgres in the registry
	service.Registry["clutch.service.db.postgres"] = dbmock.NewMockDB()
	_, err = New(&anypb.Any{}, log, scope)
	assert.NoError(t, err)
}

func TestGetProjects(t *testing.T) {
	m := dbmock.NewMockDB()
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	c := &client{
		db:    m.DB(),
		log:   log,
		scope: scope,
	}

	projectAny, err := anypb.New(&projectv1core.Project{
		Name: "foo",
	})
	assert.NoError(t, err)

	projectJson, err := protojson.Marshal(projectAny)
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"data"})
	rows.AddRow(projectJson)

	m.Mock.ExpectQuery("SELECT .*? FROM topology_cache WHERE.*").
		WithArgs(
			"foo", "type.googleapis.com/clutch.core.project.v1.Project",
		).WillReturnRows(rows)

	_, err = c.GetProjects(context.Background(), &projectv1.GetProjectsRequest{
		Names: []string{"foo"},
	})
	assert.NoError(t, err)

	m.MustMeetExpectations()
}
