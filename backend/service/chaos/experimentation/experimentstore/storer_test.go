package experimentstore

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
)

var experimentColumns = []string{
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

func createExperimentsTests() ([]experimentTest, error) {
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

	anyConfig, err := ptypes.MarshalAny(config)
	if err != nil {
		return []experimentTest{}, err
	}

	return []experimentTest{
		{
			id:        "create experiment",
			config:    anyConfig,
			startTime: time.Now(),
			queries: []*testQuery{
				{
					sql: `INSERT INTO experiment_config (id, details) VALUES ($1, $2)`,
					args: []driver.Value{
						sqlmock.AnyArg(),
						`{"@type":"type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig","fault":{"abortFault":{"percentage":{"percentage"":100}}, "abortStatus":{"httpStatusCode":401}},"faultTargeting": {"enforcer":{"upstreamEnforcing":{"upstreamType":{"upstreamCluster":{"name":"uCluster"}}, "downstreamType":{"downstreamCluster":{"name":"dCluster"}}}}}}`,
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
		{
			id:        "create empty experiments",
			startTime: time.Now(),
			err:       errors.New("empty config"),
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

			_, err = es.CreateExperiment(context.Background(), test.config, &test.startTime, nil)
			a.Equal(test.err, err)
		})
	}
}

func TestCancelExperimentRun(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.NoError(err)

	es := &storer{db: db}
	defer es.Close()

	expected := mock.ExpectExec(regexp.QuoteMeta(`UPDATE experiment_run SET cancellation_time = NOW() WHERE id = $1 AND cancellation_time IS NULL`))
	expected.WithArgs([]driver.Value{1}...).WillReturnResult(sqlmock.NewResult(1, 1))

	err = es.CancelExperimentRun(context.Background(), uint64(1))
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
	rows := sqlmock.NewRows(experimentColumns)
	rows.AddRow([]driver.Value{
		1234,
		`{"@type":"type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig","fault":{"abortFault":{"percentage":{"percentage":100}}, "abortStatus":{"httpStatusCode":401}},"faultTargeting": {"enforcer":{"upstreamEnforcing":{"upstreamType":{"upstreamCluster":{"name":"uCluster"}}, "downstreamType":{"downstreamCluster":{"name":"dCluster"}}}}}}`,
	}...)
	expected.WillReturnRows(rows)

	experiments, err := es.GetExperiments(context.Background(), "foo", experimentation.GetExperimentsRequest_STATUS_UNSPECIFIED)
	assert.NoError(err)

	assert.Equal(1, len(experiments))
	experiment := experiments[0]
	assert.Equal(uint64(1234), experiment.GetId())

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
			`{"@type":"type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig","fault":{"abortFault":{"percentage":{"percentage":100}}, "abortStatus":{"httpStatusCode":401}},"faultTargeting": {"enforcer":{"upstreamEnforcing":{"upstreamType":{"upstreamCluster":{"name":"uCluster"}}, "downstreamType":{"downstreamCluster":{"name":"dCluster"}}}}}}`,
		},
		{
			2,
			`{"@type": "malformed_foo","fault":{"abortFault":{"percentage":{"percentage":100}}, "abortStatus":{"httpStatusCode":401}},"faultTargeting": {"enforcer":{"upstreamEnforcing":{"upstreamType":{"upstreamCluster":{"name":"uCluster"}}, "downstreamType":{"downstreamCluster":{"name":"dCluster"}}}}}}`,
		},
	}

	rows := sqlmock.NewRows(experimentColumns)
	for _, row := range rowsData {
		rows.AddRow(row...)
	}
	expected.WillReturnRows(rows)

	experiments, err := es.GetExperiments(context.Background(), "foo", experimentation.GetExperimentsRequest_STATUS_UNSPECIFIED)
	assert.Nil(experiments)
	assert.Error(err)
}
