package kinesis

import (
	"errors"

	"github.com/uber-go/tally"
	"go.uber.org/zap"

	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"

	"github.com/golang/protobuf/ptypes/any"

	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/aws"
)

const (
	Name = "clutch.module.kinesis"
)

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	awsClient, ok := service.Registry["clutch.service.aws"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	c, ok := awsClient.(aws.Client)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	mod := &mod{
		kinesis: newKinesisAPI(c),
	}

	return mod, nil
}

type mod struct {
	kinesis kinesisv1.KinesisAPIServer
}

func (m *mod) Register(r module.Registrar) error {
	kinesisv1.RegisterKinesisAPIServer(r.GRPCServer(), m.kinesis)
	return r.RegisterJSONGateway(kinesisv1.RegisterKinesisAPIHandler)
}
