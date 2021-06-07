package project

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/any"
	projectv1core "github.com/lyft/clutch/backend/api/core/project/v1"
	projectv1 "github.com/lyft/clutch/backend/api/project/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	Name = "clutch.service.project"
)

type Service interface {
	GetProjects(ctx context.Context, req *projectv1.GetProjectsRequest) (*projectv1.GetProjectsResponse, error)
}

type client struct {
	db    *sql.DB
	log   *zap.Logger
	scope tally.Scope
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	p, ok := service.Registry[pgservice.Name]
	if !ok {
		return nil, errors.New("Please config the datastore [clutch.service.db.postgres] to use the topology service")
	}

	dbClient, ok := p.(pgservice.Client)
	if !ok {
		return nil, errors.New("Unable to get the datastore client")
	}

	c := &client{
		db:    dbClient.DB(),
		log:   logger,
		scope: scope,
	}

	return c, nil
}

func (c *client) GetProjects(ctx context.Context, req *projectv1.GetProjectsRequest) (*projectv1.GetProjectsResponse, error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("data").
		From("topology_cache").
		Where(
			sq.And{
				sq.Eq{"id": req.Names},
				sq.Eq{"resolver_type_url": meta.TypeURL((*projectv1core.Project)(nil))},
			},
		)

	rows, err := query.RunWith(c.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []*projectv1core.Project{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			c.log.Error("Error scanning row", zap.Error(err))
			return nil, err
		}

		var dataAny any.Any
		if err := protojson.Unmarshal(data, &dataAny); err != nil {
			c.log.Error("Error unmarshaling data field", zap.Error(err))
			return nil, err
		}

		var project projectv1core.Project
		if err := anypb.UnmarshalTo(&dataAny, &project, proto.UnmarshalOptions{}); err != nil {
			c.log.Error("Error unmarshaling project", zap.Error(err))
			return nil, err
		}

		results = append(results, &project)
	}

	return &projectv1.GetProjectsResponse{
		Projects: results,
	}, nil
}
