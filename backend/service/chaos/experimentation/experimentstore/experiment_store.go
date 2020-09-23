package experimentstore

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/id"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const Name = "clutch.service.chaos.experimentation.store"

// ExperimentStore stores experiment data
type ExperimentStore interface {
	CreateExperiment(context.Context, *any.Any, *time.Time, *time.Time) (*experimentation.Experiment, error)
	CancelExperimentRun(context.Context, uint64) error
	GetExperiments(ctx context.Context, configType string, status experimentation.GetExperimentsStatus) ([]*experimentation.Experiment, error)
	GetExperimentRunDetails(ctx context.Context, id uint64) (*experimentation.ExperimentRunDetails, error)
	Close()
}

type experimentStore struct {
	db *sql.DB
}

// New returns a new NewExperimentStore instance.
func New(_ *any.Any, _ *zap.Logger, _ tally.Scope) (service.Service, error) {
	p, ok := service.Registry[pgservice.Name]
	if !ok {
		return nil, errors.New("could not find database service")
	}

	client, ok := p.(pgservice.Client)
	if !ok {
		return nil, errors.New("experiment store wrong type")
	}

	return &experimentStore{
		client.DB(),
	}, nil
}

func (fs *experimentStore) CreateExperiment(ctx context.Context, config *any.Any, startTime *time.Time, endTime *time.Time) (*experimentation.Experiment, error) {
	// This API call will eventually be broken into 2 separate calls:
	// 1) creating the config
	// 2) starting a new experiment with the config

	// All experiments are created in a single transaction
	tx, err := fs.db.Begin()
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, errors.New("empty config")
	}

	// Step 1) create the config
	configID := id.NewID()

	configJson, err := marshalConfig(config)
	if err != nil {
		return nil, err
	}

	configSql := `INSERT INTO experiment_config (id, details) VALUES ($1, $2)`
	_, err = fs.db.ExecContext(ctx, configSql, configID, configJson)
	if err != nil {
		return nil, err
	}

	// Step 2) start a new experiment with the config
	runSql := `
			INSERT INTO experiment_run (
				id,
				experiment_config_id,
				execution_time,
				creation_time)
			VALUES ($1, $2, tstzrange($3, $4, '[]'), NOW())`

	runId := id.NewID()
	_, err = fs.db.ExecContext(ctx, runSql, runId, configID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	st, err := toProto(startTime)
	if err != nil {
		return nil, err
	}

	et, err := toProto(endTime)
	if err != nil {
		return nil, err
	}

	return &experimentation.Experiment{
		// TODO(bgallagher) temporarily returning the experiment run ID. Eventually, the CreateExperiments function
		// will be split into CreateExperimentConfig and CreateExperimentRun in which case they will each return
		// their respective IDs
		Id:        uint64(runId),
		Config:    config,
		StartTime: st,
		EndTime:   et,
	}, nil
}

func toProto(t *time.Time) (*timestamp.Timestamp, error) {
	if t == nil {
		return nil, nil
	}

	timestampProto, err := ptypes.TimestampProto(*t)
	if err != nil {
		return nil, err
	}

	return timestampProto, nil
}

func (fs *experimentStore) CancelExperimentRun(ctx context.Context, id uint64) error {
	sql :=
		`UPDATE experiment_run 
         SET cancellation_time = NOW()
         WHERE id = $1 AND cancellation_time IS NULL AND (upper(execution_time) IS NULL OR NOW() < upper(execution_time))`

	_, err := fs.db.ExecContext(ctx, sql, id)
	return err
}

// GetExperiments experiments with a given type of the configuration. Returns all experiments if provided configuration type
// parameter is an emtpy string.
func (fs *experimentStore) GetExperiments(ctx context.Context, configType string, status experimentation.GetExperimentsStatus) ([]*experimentation.Experiment, error) {
	query := `
        SELECT 
			experiment_run.id, 
			details
		FROM experiment_config, experiment_run
		WHERE
			experiment_config.id = experiment_run.experiment_config_id
			AND ($1 = '' OR $1 = experiment_config.details ->> '@type')
			AND ($2 != 'RUNNING' OR (experiment_run.cancellation_time is NULL AND NOW() > lower(experiment_run.execution_time) AND (upper(experiment_run.execution_time) IS NULL OR NOW() < upper(experiment_run.execution_time))))`

	rows, err := fs.db.QueryContext(ctx, query, configType, status.String())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var experiments []*experimentation.Experiment
	for rows.Next() {
		var experiment experimentation.Experiment
		var details string

		err = rows.Scan(&experiment.Id, &details)
		if err != nil {
			return nil, err
		}

		anyConfig := &any.Any{}
		if nil != jsonpb.Unmarshal(strings.NewReader(details), anyConfig) {
			return nil, err
		}

		experiment.Config = anyConfig
		experiments = append(experiments, &experiment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return experiments, nil
}

func (fs *experimentStore) GetExperimentRunDetails(ctx context.Context, id uint64) (*experimentation.ExperimentRunDetails, error) {
	sqlQuery := `
        SELECT experiment_run.id, lower(execution_time), upper(execution_time), cancellation_time, creation_time, details FROM experiment_config, experiment_run
        WHERE experiment_run.id = $1 AND experiment_run.experiment_config_id = experiment_config.id`

	row := fs.db.QueryRowContext(ctx, sqlQuery, id)

	var fetchedID uint64
	var startTime, endTime, cancellationTime sql.NullTime
	var creationTime time.Time
	var details string

	err := row.Scan(&fetchedID, &startTime, &endTime, &cancellationTime, &creationTime, &details)
	if err != nil {
		return nil, err
	}

	return NewRunDetails(fetchedID, creationTime, startTime, endTime, cancellationTime, time.Now(), details)
}

// Close closes all resources held.
func (fs *experimentStore) Close() {
	fs.db.Close()
}

func marshalConfig(config *any.Any) (string, error) {
	marshaler := jsonpb.Marshaler{}
	buf := &bytes.Buffer{}
	err := marshaler.Marshal(buf, config)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
