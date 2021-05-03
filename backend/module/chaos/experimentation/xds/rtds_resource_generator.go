package xds

import (
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

var RTDSGeneratorsByTypeUrl = map[string]RTDSResourceGenerator{}

type RuntimeKeyValue struct {
	Key   string
	Value uint32
}

type RTDSResource struct {
	Cluster          string
	RuntimeKeyValues []*RuntimeKeyValue
}

func NewRTDSResource(cluster string, keyValues []*RuntimeKeyValue) (*RTDSResource, error) {
	return &RTDSResource{
		Cluster:          cluster,
		RuntimeKeyValues: keyValues,
	}, nil
}

type RTDSResourceGenerator interface {
	// Generates an RTDS resource for a given experiment. The implementation of this method should
	// return a nil resource and no error if  the receiver is not interested in generating faults for a
	// passed experiment.
	GenerateResource(experiment *experimentstore.Experiment) (*RTDSResource, error)
}
