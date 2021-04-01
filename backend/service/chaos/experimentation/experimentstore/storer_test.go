package experimentstore

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
)

var experimentConfigTableColumns = []string{
	"id",
	"details",
}

type testQuery struct {
	sql  string
	args []driver.Value
}

type experimentTest struct {
	id        string
	config    *any.Any
	startTime time.Time
	queries   []*testQuery
	err       error
}

type experimentConfigTestData struct {
	stringifiedConfig string
	marshaledConfig   *any.Any
}

func NewExperimentConfigTestData() (*experimentConfigTestData, error) {
	config := &serverexperimentation.HTTPFaultConfig{
		Fault: &serverexperimentation.HTTPFaultConfig_AbortFault{
			AbortFault: &serverexperimentation.AbortFault{
				Percentage:  &serverexperimentation.FaultPercentage{Percentage: 100},
				AbortStatus: &serverexperimentation.FaultAbortStatus{HttpStatusCode: 401},
			},
		},
		FaultTargeting: &serverexperimentation.FaultTargeting{
			Enforcer: &serverexperimentation.FaultTargeting_UpstreamEnforcing{
				UpstreamEnforcing: &serverexperimentation.UpstreamEnforcing{
					UpstreamType: &serverexperimentation.UpstreamEnforcing_UpstreamCluster{
						UpstreamCluster: &serverexperimentation.SingleCluster{
							Name: "upstreamCluster",
						},
					},
					DownstreamType: &serverexperimentation.UpstreamEnforcing_DownstreamCluster{
						DownstreamCluster: &serverexperimentation.SingleCluster{
							Name: "downstreamCluster",
						},
					},
				},
			},
		},
	}

	anyConfig, err := anypb.New(config)
	if err != nil {
		return nil, err
	}

	return &experimentConfigTestData{
		stringifiedConfig: `{"@type":"type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig","faultTargeting":{"upstreamEnforcing":{"upstreamCluster":{"name":"upstreamCluster"},"downstreamCluster":{"name":"downstreamCluster"}}},"abortFault":{"percentage":{"percentage":100},"abortStatus":{"httpStatusCode":401}}}`,
		marshaledConfig:   anyConfig,
	}, nil
}

func createExperimentsTests() ([]experimentTest, error) {
	ctd, err := NewExperimentConfigTestData()
	if err != nil {
		return nil, err
	}

	return []experimentTest{
		{
			id:        "create experiment",
			config:    ctd.marshaledConfig,
			startTime: time.Now(),
			queries: []*testQuery{
				{
					sql: `INSERT INTO experiment_config (id, details) VALUES ($1, $2)`,
					args: []driver.Value{
						sqlmock.AnyArg(),
						ctd.stringifiedConfig,
					},
				},
				{
					sql: `INSERT INTO experiment_run ( id, experiment_config_id, execution_time, creation_time) VALUES ($1, $2, tstzrange($3, $4, '[]'), NOW())`,
					args: []driver.Value{
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					},
				},
			},
			err: nil,
		},
	}, nil
}

func TestCreateExperiments(t *testing.T) {
	t.Parallel()

	tests, err := createExperimentsTests()
	if err != nil {
		t.Errorf("setSnapshot failed %v", err)
	}

	for _, test := range tests {
		test := test

		t.Run(test.id, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)

			db, mock, err := sqlmock.New()
			a.NoError(err)

			es := &storer{db: db}
			defer es.Close()
			mock.ExpectBegin()
			for _, query := range test.queries {
				expected := mock.ExpectExec(regexp.QuoteMeta(query.sql))
				expected.WithArgs(query.args...).WillReturnResult(sqlmock.NewResult(1, 1))
			}
			mock.ExpectCommit()

			s := ExperimentSpecification{RunId: "1", ConfigId: "1", StartTime: test.startTime, EndTime: nil, Config: test.config}
			_, err = es.CreateExperiment(context.Background(), &s)
			a.Equal(test.err, err)
		})
	}
}

func TestCreateOrGetExperiment(t *testing.T) {
	type query struct {
		sql    string
		args   []driver.Value
		result *sqlmock.Rows
	}

	type exec struct {
		sql  string
		args []driver.Value
	}

	a := assert.New(t)
	ctd, err := NewExperimentConfigTestData()
	a.NoError(err)

	tests := []struct {
		runId           string
		expectedQueries []*query
		expectedExecs   []*exec
		expectedOrigin  experimentation.CreateOrGetExperimentResponse_Origin
	}{
		{
			runId: "1",
			expectedQueries: []*query{
				{
					sql:    `SELECT exists (select id from experiment_run where id == $1)`,
					args:   []driver.Value{sqlmock.AnyArg()},
					result: sqlmock.NewRows([]string{"exists"}).AddRow(true),
				},
				{
					sql:    `SELECT experiment_run.id, lower(execution_time), upper(execution_time), cancellation_time, creation_time, experiment_config.id, details FROM experiment_config, experiment_run WHERE experiment_run.id = $1 AND experiment_run.experiment_config_id = experiment_config.id`,
					args:   []driver.Value{"1"},
					result: sqlmock.NewRows([]string{"run_id", "execution_time", "execution_time", "cancellation_time", "creation_time", "config_id", "details"}).AddRow("1", time.Now(), time.Now(), sql.NullTime{}, time.Now(), "2", ctd.stringifiedConfig),
				},
			},
			expectedExecs: []*exec{
				{
					sql: `INSERT INTO experiment_config (id, details) VALUES ($1, $2)`,
					args: []driver.Value{
						sqlmock.AnyArg(),
						ctd.marshaledConfig,
					},
				},
				{
					sql: `INSERT INTO experiment_run ( id, experiment_config_id, execution_time, creation_time) VALUES ($1, $2, tstzrange($3, $4, '[]'), NOW())`,
					args: []driver.Value{
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					},
				},
			},
			expectedOrigin: experimentation.CreateOrGetExperimentResponse_ORIGIN_EXISTING,
		},
		{
			runId: "1",
			expectedQueries: []*query{
				{
					sql:    `SELECT exists (select id from experiment_run where id == $1)`,
					args:   []driver.Value{sqlmock.AnyArg()},
					result: sqlmock.NewRows([]string{"exists"}).AddRow(false),
				},
			},
			expectedExecs: []*exec{
				{
					sql: `INSERT INTO experiment_config (id, details) VALUES ($1, $2)`,
					args: []driver.Value{
						sqlmock.AnyArg(),
						ctd.stringifiedConfig,
					},
				},
				{
					sql: `INSERT INTO experiment_run ( id, experiment_config_id, execution_time, creation_time) VALUES ($1, $2, tstzrange($3, $4, '[]'), NOW())`,
					args: []driver.Value{
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					},
				},
			},
			expectedOrigin: experimentation.CreateOrGetExperimentResponse_ORIGIN_NEW,
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			a.NoError(err)

			es := &storer{db: db}
			defer es.Close()
			for _, query := range tt.expectedQueries {
				expected := mock.ExpectQuery(regexp.QuoteMeta(query.sql))
				expected.WithArgs(query.args...)
				expected.WillReturnRows(query.result)
			}

			if len(tt.expectedExecs) > 0 {
				mock.ExpectBegin()
				for _, exec := range tt.expectedExecs {
					expected := mock.ExpectExec(regexp.QuoteMeta(exec.sql))
					expected.WithArgs(exec.args...).WillReturnResult(sqlmock.NewResult(1, 1))
				}
				mock.ExpectCommit()
			}

			s := ExperimentSpecification{RunId: tt.runId, StartTime: time.Now(), Config: ctd.marshaledConfig}
			result, err := es.CreateOrGetExperiment(context.Background(), &s)
			a.NoError(err)
			a.Equal(tt.expectedOrigin, result.Origin)
		})
	}
}

func TestCancelExperimentRun(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.NoError(err)

	es := &storer{db: db}
	defer es.Close()

	expected := mock.ExpectExec(regexp.QuoteMeta(`UPDATE experiment_run SET cancellation_time = NOW(), termination_reason = $2 WHERE id = $1 AND cancellation_time IS NULL`))
	expected.WithArgs([]driver.Value{"1", ""}...).WillReturnResult(sqlmock.NewResult(1, 1))

	err = es.CancelExperimentRun(context.Background(), "1", "")
	assert.NoError(err)
}

var getExperimentsSQLQuery = `SELECT experiment_run.id, details FROM experiment_config, experiment_run WHERE experiment_config.id = experiment_run.experiment_config_id AND ($1 = '' OR $1 = experiment_config.details ->> '@type')`

func TestGetExperimentsUnmarshalsExperimentConfiguration(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.NoError(err)

	es := &storer{db: db}
	defer es.Close()

	expected := mock.ExpectQuery(regexp.QuoteMeta(getExperimentsSQLQuery)).WithArgs("foo", "STATUS_UNSPECIFIED")
	rows := sqlmock.NewRows(experimentConfigTableColumns)
	rows.AddRow([]driver.Value{
		1234,
		`{"@type": "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig","faultTargeting":{"upstreamEnforcing":{"upstreamCluster":{"name":"upstreamCluster"},"downstreamCluster":{"name":"downstreamCluster"}}},"abortFault":{"percentage":{"percentage":100},"abortStatus":{"httpStatusCode":401}}}`,
	}...)
	expected.WillReturnRows(rows)

	experiments, err := es.GetExperiments(context.Background(), "foo", experimentation.GetExperimentsRequest_STATUS_UNSPECIFIED)
	assert.NoError(err)

	assert.Equal(1, len(experiments))
	experiment := experiments[0]
	assert.Equal("1234", experiment.GetRunId())

	config := &serverexperimentation.HTTPFaultConfig{}
	err = ptypes.UnmarshalAny(experiment.GetConfig(), config)
	assert.NoError(err)
	assert.Nil(config.GetLatencyFault())
	abort := config.GetAbortFault()
	assert.NotNil(abort)
	assert.Equal(uint32(401), abort.GetAbortStatus().GetHttpStatusCode())
	assert.Equal(uint32(100), abort.GetPercentage().GetPercentage())
}

func TestGetExperimentsFailsIfItReadsExperimentWithMalformedConfiguration(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.NoError(err)

	es := &storer{db: db}
	defer es.Close()

	expected := mock.ExpectQuery(regexp.QuoteMeta(getExperimentsSQLQuery)).WithArgs("foo", "STATUS_UNSPECIFIED")
	rowsData := [][]driver.Value{
		{
			1,
			`{"@type": "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig","faultTargeting":{"upstreamEnforcing":{"upstreamCluster":{"name":"upstreamCluster"},"downstreamCluster":{"name":"downstreamCluster"}}},"abortFault":{"percentage":{"percentage":100},"abortStatus":{"httpStatusCode":401}}}`,
		},
		{
			2,
			`{"@type": "malformed_foo","faultTargeting":{"upstreamEnforcing":{"upstreamCluster":{"name":"upstreamCluster"},"downstreamCluster":{"name":"downstreamCluster"}}},"abortFault":{"percentage":{"percentage":100},"abortStatus":{"httpStatusCode":401}}}`,
		},
	}

	rows := sqlmock.NewRows(experimentConfigTableColumns)
	for _, row := range rowsData {
		rows.AddRow(row...)
	}
	expected.WillReturnRows(rows)

	experiments, err := es.GetExperiments(context.Background(), "foo", experimentation.GetExperimentsRequest_STATUS_UNSPECIFIED)
	assert.Nil(experiments)
	assert.Error(err)
}
