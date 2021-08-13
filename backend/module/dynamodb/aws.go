package dynamodb

import (
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/aws"
)

const (
	Name = "clutch.module.dynamodb"
)

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	awsClient, ok := service.Registry["clutch.service.aws"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	c, ok := awsClient.(aws.Client)
	if !ok {
		return nil, errors.New("dynamodb: service was not the correct type")
	}

	mod := &mod{
		dynamodb: newDDBAPI(c),
	}

	return mod, nil
}

type mod struct {
	dynamodb dynamodbv1.DDBAPIServer
}

func (m *mod) Register(r module.Registrar) error {
	dynamodbv1.RegisterDDBAPIServer(r.GRPCServer(), m.dynamodb)
	return r.RegisterJSONGateway(dynamodbv1.RegisterDDBAPIHandler)
}
