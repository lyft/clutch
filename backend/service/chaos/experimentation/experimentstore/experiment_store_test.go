package experimentstore

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
)

type ExperimentTest struct {
	id          string
	experiments []*experimentation.Experiment
	sql         string
	args        []driver.Value
	err         error
}


func createExperimentsTests() ([]ExperimentTest, error) {
	specification := &serverexperimentation.TestSpecification{
		Config: &serverexperimentation.TestSpecification_Abort{
			Abort: &serverexperimentation.AbortFault{
				Target: &serverexperimentation.AbortFault_ClusterPair{
					ClusterPair: &serverexperimentation.ClusterPairTarget{
						DownstreamCluster: "upstreamCluster",
						UpstreamCluster:   "downstreamCluster",
					},
				},
				Percent:    100.0,
				HttpStatus: 401,
			},
		},
	}

	anyConfig, err := ptypes.MarshalAny(specification)
	if err != nil {
		return []ExperimentTest{}, err
	}

	return []ExperimentTest {
		ExperimentTest{
			id: "create experiment",
			experiments: []*experimentation.Experiment{
				{
					TestConfig: anyConfig,
				},
			},
			sql: "INSERT INTO experiments (id,details) VALUES ($1,$2)",
			args: []driver.Value{
				1,
				`{"@type":"type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestSpecification","abort":{"clusterPair":{"downstreamCluster":"upstreamCluster","upstreamCluster":"downstreamCluster"},"percent":100,"httpStatus":401}}`,
			},
		},
		ExperimentTest {
			id:  "create empty experiments",
			err: errors.New("insert statements must have at least one set of values or select clause"),
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

			expected := mock.ExpectExec(regexp.QuoteMeta(test.sql))
			if test.err != nil {
				expected.WillReturnError(test.err)
			} else {
				expected.WithArgs(sqlmock.AnyArg(), test.args[1]).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err = es.CreateExperiments(context.Background(), test.experiments)
			if test.err != nil {
				a.Equal(test.err, err)
			} else {
				a.NoError(err)
			}
		})
	}
}

var deleteExperimentsTests = []struct {
	id   string
	ids  []uint64
	sql  string
	args []driver.Value
	err  error
}{
	{
		id:  "delete all experiments",
		sql: `DELETE FROM experiments`,
	},
	{
		id:   "delete specific experiment",
		ids:  []uint64{1, 2, 3},
		sql:  `DELETE FROM experiments WHERE id IN ($1,$2,$3)`,
		args: []driver.Value{1, 2, 3},
	},
}

func TestDeleteExperiments(t *testing.T) {
	t.Parallel()

	for _, test := range deleteExperimentsTests {
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

			err = es.DeleteExperiments(context.Background(), test.ids)
			if test.err != nil {
				a.Equal(test.err, err)
			} else {
				a.NoError(err)
			}
		})
	}
}

var getExperimentsTests = []struct {
	id   string
	sql  string
	args []driver.Value
	rows [][]driver.Value
	err  error
}{
	{
		id:   "get all experiments",
		sql:  `SELECT id, details FROM experiments`,
		args: []driver.Value{"upstreamCluster", "downstreamCluster", 1},
		rows: [][]driver.Value{
			{
				1234,
				`{"@type": "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestSpecification", "abort":{"clusterPair":{"downstreamCluster":"upstreamCluster","upstreamCluster":"downstreamCluster"},"percent":100,"httpStatus":401}}`,
			},
		},
	},
}

func TestGetExperiments(t *testing.T) {
	t.Parallel()

	for _, test := range getExperimentsTests {
		test := test

		t.Run(test.id, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)

			db, mock, err := sqlmock.New()
			a.NoError(err)

			es := &experimentStore{db: db}
			defer es.Close()

			expected := mock.ExpectQuery(regexp.QuoteMeta(test.sql))
			if test.err != nil {
				expected.WillReturnError(test.err)
			} else {
				rows := sqlmock.NewRows(experimentColumns)
				for _, row := range test.rows {
					rows.AddRow(row...)
				}
				expected.WillReturnRows(rows)
			}

			experiments, err := es.GetExperiments(context.Background())
			if test.err != nil {
				a.Equal(test.err, err)
			} else {
				a.NoError(err)
			}

			a.Equal(1, len(experiments))
			experiment := experiments[0]
			a.NotEqual(0, experiment.GetId())

			specification := &serverexperimentation.TestSpecification{}
			err2 := ptypes.UnmarshalAny(experiment.GetTestConfig(), specification)
			if err2 != nil {
				t.Errorf("setSnapshot failed %v", err2)
			}
			a.Nil(specification.GetLatency())
			abort := specification.GetAbort()
			a.NotNil(abort)
			a.Equal(int32(401), abort.GetHttpStatus())
			a.Equal(float32(100), abort.GetPercent())
		})
	}
}
