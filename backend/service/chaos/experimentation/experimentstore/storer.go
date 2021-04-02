package experimentstore

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const Name = "clutch.service.chaos.experimentation.store"

// Storer stores experiment data
type Storer interface {
	CreateExperiment(context.Context, *ExperimentSpecification) (*experimentation.Experiment, error)
	CreateOrGetExperiment(context.Context, *ExperimentSpecification) (*CreateOrGetExperimentResult, error)
	CancelExperimentRun(ctx context.Context, id string, reason string) error
	GetExperiments(ctx context.Context, configType string, status experimentation.GetExperimentsRequest_Status) ([]*experimentation.Experiment, error)
	GetExperimentRunDetails(ctx context.Context, id string) (*experimentation.ExperimentRunDetails, error)
	GetListView(ctx context.Context) ([]*experimentation.ListViewItem, error)
	RegisterTransformation(transformation Transformation) error
	Close()
}

type storer struct {
	db          *sql.DB
	logger      *zap.SugaredLogger
	transformer *Transformer
}

var _ Storer = (*storer)(nil)

// New returns a new NewExperimentStore instance.
func New(_ *any.Any, logger *zap.Logger, _ tally.Scope) (service.Service, error) {
	p, ok := service.Registry[pgservice.Name]
	if !ok {
		return nil, errors.New("could not find database service")
	}

	client, ok := p.(pgservice.Client)
	if !ok {
		return nil, errors.New("experiment store wrong type")
	}

	sugaredLogger := logger.Sugar()
	transformer := NewTransformer(sugaredLogger)
	return &storer{
		client.DB(),
		sugaredLogger,
		&transformer,
	}, nil
}

func (s *storer) CreateExperiment(ctx context.Context, es *ExperimentSpecification) (*experimentation.Experiment, error) {
	// This API call will eventually be broken into 2 separate calls:
	// 1) creating the config
	// 2) starting a new experiment with the config

	// All experiments are created in a single transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	configJson, err := marshalConfig(es.Config)
	if err != nil {
		return nil, err
	}

	configSql := `INSERT INTO experiment_config (id, details) VALUES ($1, $2)`
	_, err = s.db.ExecContext(ctx, configSql, es.ConfigId, configJson)
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

	_, err = s.db.ExecContext(ctx, runSql, es.RunId, es.ConfigId, es.StartTime, es.EndTime)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// TODO(bgallagher) temporarily returning the experiment run ID. Eventually, the CreateExperiment function
	// will be split into CreateExperimentConfig and CreateExperimentRun in which case they will each return
	// their respective IDs
	return es.toExperiment()
}

func (s *storer) CreateOrGetExperiment(ctx context.Context, es *ExperimentSpecification) (*CreateOrGetExperimentResult, error) {
	var exists bool
	query := `SELECT exists (select id from experiment_run where id == $1)`
	err := s.db.QueryRow(query, es.RunId).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	if exists {
		runConfigPair, err := s.getExperiment(ctx, es.RunId)
		if err != nil {
			return nil, err
		}

		experiment, err := runConfigPair.toProto()
		if err != nil {
			return nil, err
		}

		return &CreateOrGetExperimentResult{Experiment: experiment, Origin: experimentation.CreateOrGetExperimentResponse_ORIGIN_EXISTING}, nil
	} else {
		experiment, err := s.CreateExperiment(ctx, es)
		if err != nil {
			return nil, err
		}

		return &CreateOrGetExperimentResult{Experiment: experiment, Origin: experimentation.CreateOrGetExperimentResponse_ORIGIN_NEW}, nil
	}
}

func (s *storer) CancelExperimentRun(ctx context.Context, id string, reason string) error {
	if len(reason) > 32 {
		reason = reason[0:32]
	}
	sql :=
		`UPDATE experiment_run
         SET cancellation_time = NOW(),
		 termination_reason = $2
         WHERE id = $1 AND cancellation_time IS NULL AND (upper(execution_time) IS NULL OR NOW() < upper(execution_time))`

	_, err := s.db.ExecContext(ctx, sql, id, reason)
	return err
}

// GetExperiments experiments with a given type of the configuration. Returns all experiments if provided configuration type
// parameter is an empty string.
func (s *storer) GetExperiments(ctx context.Context, configType string, status experimentation.GetExperimentsRequest_Status) ([]*experimentation.Experiment, error) {
	query := `
		SELECT
			experiment_run.id,
			details
		FROM experiment_config, experiment_run
		WHERE
			experiment_config.id = experiment_run.experiment_config_id` +
		// Return only experiments of a given `configType` or all of them if configType is equal to an empty string.
		` AND ($1 = '' OR $1 = experiment_config.details ->> '@type')` +
		// Return only running experiments if `status` is equal to `Running`, return all experiments otherwise.
		` AND ($2 = 'STATUS_UNSPECIFIED' OR (experiment_run.cancellation_time is NULL AND NOW() > lower(experiment_run.execution_time) AND (upper(experiment_run.execution_time) IS NULL OR NOW() < upper(experiment_run.execution_time))))`

	rows, err := s.db.QueryContext(ctx, query, configType, status.String())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var experiments []*experimentation.Experiment
	for rows.Next() {
		var experiment experimentation.Experiment
		var details string

		err = rows.Scan(&experiment.RunId, &details)
		if err != nil {
			return nil, err
		}

		anyConfig := &any.Any{}
		err = protojson.Unmarshal([]byte(details), anyConfig)
		if err != nil {
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

func (s *storer) GetListView(ctx context.Context) ([]*experimentation.ListViewItem, error) {
	query := `
		SELECT
			experiment_run.id,
			lower(execution_time),
			upper(execution_time),
			cancellation_time,
			creation_time,
			experiment_config.id,
			details
		FROM experiment_config, experiment_run
		WHERE
			experiment_config.id = experiment_run.experiment_config_id`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listViewItems []*experimentation.ListViewItem
	for rows.Next() {
		var details string
		run := ExperimentRun{}
		config := ExperimentConfig{Config: &any.Any{}}
		err = rows.Scan(&run.Id, &run.StartTime, &run.EndTime, &run.CancellationTime, &run.CreationTime, &config.id, &details)
		if err != nil {
			return nil, err
		}
		
		if err = protojson.Unmarshal([]byte(details), config.Config); err != nil {
			return nil, err
		}

		item, err := NewRunListView(&run, &config, s.transformer, time.Now())
		if err != nil {
			return nil, err
		}

		listViewItems = append(listViewItems, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return listViewItems, nil
}

func (s *storer) GetExperimentRunDetails(ctx context.Context, runId string) (*experimentation.ExperimentRunDetails, error) {
	e, err := s.getExperiment(ctx, runId)
	if err != nil {
		return nil, err
	}

	return NewRunDetails(e.Run, e.Config, s.transformer, time.Now())
}

func (s *storer) getExperiment(ctx context.Context, runId string) (*Experiment, error) {
	sqlQuery := `
		SELECT
			experiment_run.id,
			lower(execution_time),
			upper(execution_time),
			cancellation_time,
			creation_time,
			experiment_config.id,
			details FROM experiment_config, experiment_run
        WHERE experiment_run.id = $1 AND experiment_run.experiment_config_id = experiment_config.id`

	row := s.db.QueryRowContext(ctx, sqlQuery, runId)

	var details string
	run := ExperimentRun{}
	config := ExperimentConfig{Config: &any.Any{}}
	err := row.Scan(&run.Id, &run.StartTime, &run.EndTime, &run.CancellationTime, &run.CreationTime, &config.id, &details)
	if err != nil {
		return nil, err
	}

	err = protojson.Unmarshal([]byte(details), config.Config)
	if err != nil {
		return nil, err
	}

	return &Experiment{Run: &run, Config: &config}, nil
}

// Close closes all resources held.
func (s *storer) Close() {
	s.db.Close()
}

func (s *storer) RegisterTransformation(transformation Transformation) error {
	err := s.transformer.Register(transformation)
	if err != nil {
		s.logger.Fatal("Could not register transformation %v", transformation)
	}

	return err
}

func marshalConfig(config *any.Any) (string, error) {
	b, err := protojson.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
