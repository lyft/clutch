package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	envoyv1 "github.com/lyft/clutch/backend/api/core/envoy/v1"
	resolverv1 "github.com/lyft/clutch/backend/api/resolver/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/resolver"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/topology"
)

var Name = "clutch.resolver.core"

var typeURLEnvoyCluster = meta.TypeURL((*envoyv1.Cluster)(nil))

var typeSchemas = resolver.TypeURLToSchemaMessagesMap{
	typeURLEnvoyCluster: {},
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (resolver.Resolver, error) {
	var topologyService topology.Service
	if svc, ok := service.Registry[topology.Name]; ok {
		topologyService, ok = svc.(topology.Service)
		if !ok {
			return nil, errors.New("incorrect topology service type")
		}
		logger.Debug("enabling autocomplete api for the lyftcore resolver")
	}

	schemas, err := resolver.InputsToSchemas(typeSchemas)
	if err != nil {
		return nil, err
	}

	r := &res{
		topology: topologyService,
		schemas:  schemas,
	}

	return r, nil
}

type res struct {
	topology topology.Service
	schemas  resolver.TypeURLToSchemasMap
}

func (r *res) Schemas() resolver.TypeURLToSchemasMap { return r.schemas }

func (r *res) Resolve(ctx context.Context, typeURL string, input proto.Message, limit uint32) (*resolver.Results, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (r *res) Search(ctx context.Context, typeURL, query string, limit uint32) (*resolver.Results, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (r *res) Autocomplete(ctx context.Context, typeURL, search string, limit uint64) ([]*resolverv1.AutocompleteResult, error) {
	if r.topology == nil {
		return nil, fmt.Errorf("to use the autocomplete api you must first setup the topology service")
	}

	var resultLimit uint64 = resolver.DefaultAutocompleteLimit
	if limit > 0 {
		resultLimit = limit
	}

	results, err := r.topology.Autocomplete(ctx, typeURL, search, resultLimit)
	if err != nil {
		return nil, err
	}

	autoCompleteValue := make([]*resolverv1.AutocompleteResult, len(results))
	for i, r := range results {
		autoCompleteValue[i] = &resolverv1.AutocompleteResult{
			Id: r.Id,
			// TODO (mcutalo): Add more detailed information to the label
			// the labels value will vary based on resource
			Label: "",
		}
	}

	return autoCompleteValue, nil
}
