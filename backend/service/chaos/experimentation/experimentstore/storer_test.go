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
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentationv1 "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
)

type experimentConfigTestData struct {
	stringifiedConfig string
	marshaledConfig   *any.Any
	message           proto.Message
}

func NewExperimentConfigTestData() (*experimentConfigTestData, error) {
	config := &serverexperimentationv1.HTTPFaultConfig{
		Fault: &serverexperimentationv1.HTTPFaultConfig_AbortFault{
			AbortFault: &serverexperimentationv1.AbortFault{
				Percentage:  &serverexperimentationv1.FaultPercentage{Percentage: 100},
				AbortStatus: &serverexperimentationv1.FaultAbortStatus{HttpStatusCode: 401},
			},
		},
		FaultTargeting: &serverexperimentationv1.FaultTargeting{
			Enforcer: &serverexperimentationv1.FaultTargeting_UpstreamEnforcing{
				UpstreamEnforcing: &serverexperimentationv1.UpstreamEnforcing{
					UpstreamType: &serverexperimentationv1.UpstreamEnforcing_UpstreamCluster{
						UpstreamCluster: &serverexperimentationv1.SingleCluster{
							Name: "upstreamCluster",
						},
					},
					DownstreamType: &serverexperimentationv1.UpstreamEnforcing_DownstreamCluster{
						DownstreamCluster: &serverexperimentationv1.SingleCluster{
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

	message, err := anypb.UnmarshalNew(anyConfig, proto.UnmarshalOptions{})
	if err != nil {
		return nil, err
	}

	jsonConfig, err := protojson.Marshal(anyConfig)
	if err != nil {
		return nil, err
	}

	return &experimentConfigTestData{
		stringifiedConfig: string(jsonConfig),
		marshaledConfig:   anyConfig,
		message:           message,
	}, nil
}

func TestCreateExperiment(t *testing.T) {
	ctd, err := NewExperimentConfigTestData()
	assert.NoError(t, err)

	type exec struct {
		sql  string
		args []driver.Value
	}

	type query struct {
		sql    string
		args   []driver.Value
		result *sqlmock.Rows
	}

	tests := []struct {
		id              string
		config          *any.Any
		startTime       time.Time
		expectedExecs   []*exec
		expectedQueries []*query
		expectedError   error
	}{
		{
			id:        "create experiment",
			config:    ctd.marshaledConfig,
			startTime: time.Now(),
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
			expectedQueries: []*query{
				{
					sql:    `SELECT experiment_run.id, lower(execution_time), upper(execution_time), cancellation_time, creation_time, termination_reason, experiment_config.id, details FROM experiment_config, experiment_run WHERE experiment_run.id = $1 AND experiment_run.experiment_config_id = experiment_config.id`,
					args:   []driver.Value{"1"},
					result: sqlmock.NewRows([]string{"run_id", "execution_time", "execution_time", "cancellation_time", "creation_time", "termination_reason", "config_id", "details"}).AddRow("1", time.Now(), time.Now(), sql.NullTime{}, time.Now(), "", "2", ctd.stringifiedConfig),
				},
			},
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			es := &storer{db: db}
			defer es.Close()
			mock.ExpectBegin()
			for _, query := range tt.expectedExecs {
				expected := mock.ExpectExec(regexp.QuoteMeta(query.sql))
				expected.WithArgs(query.args...).WillReturnResult(sqlmock.NewResult(1, 1))
			}
			mock.ExpectCommit()

			for _, query := range tt.expectedQueries {
				expected := mock.ExpectQuery(regexp.QuoteMeta(query.sql))
				expected.WithArgs(query.args...)
				expected.WillReturnRows(query.result)
			}

			s := ExperimentSpecification{RunId: "1", ConfigId: "1", StartTime: tt.startTime, EndTime: nil, Config: tt.config}
			_, err = es.CreateExperiment(context.Background(), &s)
			assert.Equal(t, tt.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateOrGetExperiment(t *testing.T) {
	ctd, err := NewExperimentConfigTestData()
	assert.NoError(t, err)

	type query struct {
		sql    string
		args   []driver.Value
		result *sqlmock.Rows
	}

	type exec struct {
		sql  string
		args []driver.Value
	}

	tests := []struct {
		runId              string
		beforeExecsQueries []*query
		expectedExecs      []*exec
		afterExecsQueries  []*query
		expectedOrigin     experimentationv1.CreateOrGetExperimentResponse_Origin
	}{
		{
			runId: "1",
			beforeExecsQueries: []*query{
				{
					sql:    `SELECT exists (select id from experiment_run where id == $1)`,
					args:   []driver.Value{sqlmock.AnyArg()},
					result: sqlmock.NewRows([]string{"exists"}).AddRow(true),
				},
				{
					sql:    `SELECT experiment_run.id, lower(execution_time), upper(execution_time), cancellation_time, creation_time, termination_reason, experiment_config.id, details FROM experiment_config, experiment_run WHERE experiment_run.id = $1 AND experiment_run.experiment_config_id = experiment_config.id`,
					args:   []driver.Value{"1"},
					result: sqlmock.NewRows([]string{"run_id", "execution_time", "execution_time", "cancellation_time", "creation_time", "termination_reason", "config_id", "details"}).AddRow("1", time.Now(), time.Now(), sql.NullTime{}, time.Now(), "", "2", ctd.stringifiedConfig),
				},
			},
			expectedOrigin: experimentationv1.CreateOrGetExperimentResponse_ORIGIN_EXISTING,
		},
		{
			runId: "1",
			beforeExecsQueries: []*query{
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
			afterExecsQueries: []*query{
				{
					sql:    `SELECT experiment_run.id, lower(execution_time), upper(execution_time), cancellation_time, creation_time, termination_reason, experiment_config.id, details FROM experiment_config, experiment_run WHERE experiment_run.id = $1 AND experiment_run.experiment_config_id = experiment_config.id`,
					args:   []driver.Value{"1"},
					result: sqlmock.NewRows([]string{"run_id", "execution_time", "execution_time", "cancellation_time", "creation_time", "termination_reason", "config_id", "details"}).AddRow("1", time.Now(), time.Now(), sql.NullTime{}, time.Now(), "", "2", ctd.stringifiedConfig),
				},
			},
			expectedOrigin: experimentationv1.CreateOrGetExperimentResponse_ORIGIN_NEW,
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			es := &storer{db: db}
			defer es.Close()
			for _, query := range tt.beforeExecsQueries {
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

			for _, query := range tt.afterExecsQueries {
				expected := mock.ExpectQuery(regexp.QuoteMeta(query.sql))
				expected.WithArgs(query.args...)
				expected.WillReturnRows(query.result)
			}

			s := ExperimentSpecification{RunId: tt.runId, StartTime: time.Now(), Config: ctd.marshaledConfig}
			result, err := es.CreateOrGetExperiment(context.Background(), &s)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOrigin, result.Origin)
			assert.NoError(t, mock.ExpectationsWereMet())
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
	assert.NoError(mock.ExpectationsWereMet())
}

func TestGetExperiments(t *testing.T) {
	ctd, err := NewExperimentConfigTestData()
	assert.NoError(t, err)

	now := time.Now()
	future := time.Now().Add(2 * time.Minute)

	run := &ExperimentRun{
		Id:                "1",
		StartTime:         now,
		EndTime:           &future,
		CancellationTime:  &future,
		CreationTime:      now,
		TerminationReason: "foo",
	}
	config := &ExperimentConfig{
		Id:      "2",
		Config:  ctd.marshaledConfig,
		Message: ctd.message,
	}

	tests := []struct {
		queryResult         [][]driver.Value
		expectedExperiments []*Experiment
	}{
		{
			queryResult: [][]driver.Value{
				{
					run.Id, run.StartTime, run.EndTime, run.CancellationTime, run.CreationTime, run.TerminationReason, config.Id, ctd.stringifiedConfig,
				},
			},
			expectedExperiments: []*Experiment{
				{
					Run: run, Config: config,
				},
			},
		},
		{
			queryResult: [][]driver.Value{
				{
					run.Id, run.StartTime, run.EndTime, run.CancellationTime, run.CreationTime, run.TerminationReason, config.Id, ctd.stringifiedConfig,
				},
				{
					run.Id, run.StartTime, run.EndTime, run.CancellationTime, run.CreationTime, run.TerminationReason, config.Id, "{malformed dictionary}",
				},
			},
			expectedExperiments: []*Experiment{
				{
					Run: run, Config: config,
				},
			},
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			expectedQuery := `SELECT experiment_run.id, lower(execution_time), upper(execution_time), cancellation_time, creation_time, termination_reason, experiment_config.id, details FROM experiment_config, experiment_run WHERE experiment_config.id = experiment_run.experiment_config_id AND ($1 = '' OR $1 = experiment_config.details ->> '@type')`
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			l, err := zap.NewDevelopment()
			assert.NoError(t, err)

			testScope := tally.NewTestScope("", map[string]string{})

			es := &storer{
				db:                              db,
				logger:                          l.Sugar(),
				configDeserializationErrorCount: testScope.Counter("config_initialization"),
			}

			defer es.Close()

			rows := sqlmock.NewRows([]string{
				"run_id",
				"execution_time",
				"execution_time",
				"cancellation_time",
				"creation_time",
				"termination_reason",
				"config_id",
				"details",
			})

			for _, row := range tt.queryResult {
				rows.AddRow(row...)
			}

			expected := mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs("foo", "STATUS_UNSPECIFIED")
			expected.WillReturnRows(rows)
			experiments, err := es.GetExperiments(context.Background(), "foo", experimentationv1.GetExperimentsRequest_STATUS_UNSPECIFIED)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedExperiments, experiments)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
