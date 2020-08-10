package experimentstore

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/id"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const Name = "clutch.service.chaos.experimentation.store"

var experimentColumns = []string{
	"id",
	"details",
}

const experimentsTableName = "experiments"

// ExperimentStore stores experiment data
type ExperimentStore interface {
	CreateExperiments(context.Context, []*experimentation.Experiment) error
	DeleteExperiments(context.Context, []uint64) error
	GetExperiments(context.Context) ([]*experimentation.Experiment, error)
	Close()
}

type experimentStore struct {
	db *sql.DB
}

// New returns a new NewExperimentStore instance.
func New(cfg *any.Any, _ *zap.Logger, _ tally.Scope) (service.Service, error) {
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

func (fs *experimentStore) CreateExperiments(ctx context.Context, experiments []*experimentation.Experiment) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Insert(experimentsTableName).Columns(experimentColumns...)

	for _, experiment := range experiments {
		marshaler := jsonpb.Marshaler{}
		buf := &bytes.Buffer{}
		err := marshaler.Marshal(buf, experiment.GetTestSpecification())
		if err != nil {
			return err
		}
		s := buf.String()
		builder = builder.Values(
			id.NewID(),
			s,
		)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = fs.db.ExecContext(ctx, sql, args...)
	return err
}

// DeleteExperiments deletes the specified experiments from the store.
func (fs *experimentStore) DeleteExperiments(ctx context.Context, ids []uint64) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Delete("").From(experimentsTableName)

	if len(ids) > 0 {
		builder = builder.Where(sq.Eq{"id": ids})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = fs.db.ExecContext(ctx, sql, args...)
	return err
}

// GetExperiments gets all experiments
func (fs *experimentStore) GetExperiments(ctx context.Context) ([]*experimentation.Experiment, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(experimentColumns...).From(experimentsTableName)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := fs.db.QueryContext(ctx, sql, args...)
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

		spec := &experimentation.TestSpecification{}
		if nil != jsonpb.Unmarshal(strings.NewReader(details), spec) {
			return nil, err
		}
		experiment.TestSpecification = spec

		experiments = append(experiments, &experiment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return experiments, nil
}

// Close closes all resources held.
func (fs *experimentStore) Close() {
	fs.db.Close()
}
