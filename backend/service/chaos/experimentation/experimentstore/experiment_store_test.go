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
	config := &serverexperimentation.TestConfig{
		Target: &serverexperimentation.TestConfig_ClusterPair{
			ClusterPair: &serverexperimentation.ClusterPairTarget{
				DownstreamCluster: "upstreamCluster",
				UpstreamCluster:   "downstreamCluster",
			},
		},
		Fault: &serverexperimentation.TestConfig_Abort{
			Abort: &serverexperimentation.AbortFaultConfig{
				Percent:    100.0,
				HttpStatus: 401,
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
						`{"@type":"type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig","clusterPair":{"downstreamCluster":"upstreamCluster","upstreamCluster":"downstreamCluster"},"abort":{"percent":100,"httpStatus":401}}`,
					},
				},
				{
					sql: `INSERT INTO experiment_run ( id, experiment_config_id, execution_time, scheduled_end_time, creation_time) VALUES ($1, $2, tstzrange($3, $4, '[]'), $4, NOW())`,
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

			es := &experimentStore{db: db}
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

var cancelExperimentsTests = []struct {
	id    string
	runID uint64
	sql   string
	args  []driver.Value
	err   error
}{
	{
		id:    "cancel an experiment run",
		runID: uint64(1),
		sql:   `UPDATE experiment_run SET execution_time = tstzrange(lower(execution_time), NOW(), '[]') WHERE id = $1 AND (upper(execution_time) IS NULL OR NOW() < upper(execution_time))`,
		args:  []driver.Value{1},
	},
}

func TestCancelExperimentRun(t *testing.T) {
	t.Parallel()

	for _, test := range cancelExperimentsTests {
		test := test

		t.Run(test.id, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)

			db, mock, err := sqlmock.New()
			a.NoError(err)

			es := &experimentStore{db: db}
			defer es.Close()

			expected := mock.ExpectExec(regexp.QuoteMeta(test.sql))
			if test.err != nil {
				expected.WillReturnError(test.err)
			} else {
				expected.WithArgs(test.args...).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err = es.CancelExperimentRun(context.Background(), test.runID)
			if test.err != nil {
				a.Equal(test.err, err)
			} else {
				a.NoError(err)
			}
		})
	}
}

var getExperimentsSQLQuery = `SELECT experiment_run.id, details FROM experiment_config, experiment_run WHERE experiment_config.id = experiment_run.experiment_config_id AND ($1 = '' OR $1 = experiment_config.details ->> '@type')`

func TestGetExperimentsUnmarshalsExperimentConfiguration(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.NoError(err)

	counter := &mockCounter{}
	es := &experimentStore{db: db, getExperimentsLoadConfigurationError: counter}
	defer es.Close()

	expected := mock.ExpectQuery(regexp.QuoteMeta(getExperimentsSQLQuery)).WithArgs("foo")
	rows := sqlmock.NewRows(experimentColumns)
	rows.AddRow([]driver.Value{
		1234,
		`{"@type": "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig","clusterPair":{"downstreamCluster":"upstreamCluster","upstreamCluster":"downstreamCluster"},"abort":{"percent":100,"httpStatus":401}}`,
	}...)
	expected.WillReturnRows(rows)

	experiments, err := es.GetExperiments(context.Background(), "foo")
	assert.NoError(err)

	assert.Equal(1, len(experiments))
	experiment := experiments[0]
	assert.Equal(uint64(1234), experiment.GetId())

	config := &serverexperimentation.TestConfig{}
	err2 := ptypes.UnmarshalAny(experiment.GetConfig(), config)
	if err2 != nil {
		t.Errorf("setSnapshot failed %v", err2)
	}
	assert.Nil(config.GetLatency())
	abort := config.GetAbort()
	assert.NotNil(abort)
	assert.Equal(int32(401), abort.GetHttpStatus())
	assert.Equal(float32(100), abort.GetPercent())
}

func TestGetExperimentsIgnoresExperimentWithMalformedConfiguration(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.NoError(err)

	counter := &mockCounter{}
	es := &experimentStore{db: db, getExperimentsLoadConfigurationError: counter}
	defer es.Close()

	expected := mock.ExpectQuery(regexp.QuoteMeta(getExperimentsSQLQuery)).WithArgs("foo")
	rowsData := [][]driver.Value{
		{
			1,
			`{"@type": "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig","clusterPair":{"downstreamCluster":"upstreamCluster","upstreamCluster":"downstreamCluster"},"abort":{"percent":100,"httpStatus":401}}`,
		},
		{
			2,
			`{"@type": "malformed_foo","clusterPair":{"downstreamCluster":"upstreamCluster","upstreamCluster":"downstreamCluster"},"abort":{"percent":100,"httpStatus":401}}`,
		},
		{
			3,
			`{"@type": "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig","clusterPair":{"downstreamCluster":"upstreamCluster","upstreamCluster":"downstreamCluster"},"abort":{"percent":100,"httpStatus":401}}`,
		},
	}

	rows := sqlmock.NewRows(experimentColumns)
	for _, row := range rowsData {
		rows.AddRow(row...)
	}
	expected.WillReturnRows(rows)

	experiments, err := es.GetExperiments(context.Background(), "foo")
	assert.NoError(err)

	experimentIDs := []uint64{}
	for _, experiment := range experiments {
		experimentIDs = append(experimentIDs, experiment.Id)
	}
	assert.Equal([]uint64{1, 3}, experimentIDs)
	assert.Equal(int64(1), counter.count)
}
