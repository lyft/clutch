package aws

// <!-- START clutchdoc -->
// description: Endpoints for interacting with resources in the Amazon Web Services (AWS) cloud.
// <!-- END clutchdoc -->

import (
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/aws"
)

const (
	Name = "clutch.module.aws"
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
		ec2: newEC2API(c),
	}

	return mod, nil
}

type mod struct {
	ec2 ec2v1.EC2APIServer
}

func (m *mod) Register(r module.Registrar) error {
	ec2v1.RegisterEC2APIServer(r.GRPCServer(), m.ec2)
	return r.RegisterJSONGateway(ec2v1.RegisterEC2APIHandler)
}
