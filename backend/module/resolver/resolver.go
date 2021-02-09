package resolver

// <!-- START clutchdoc -->
// description: Exposes registered resolvers and their schemas for use in structured or free-form search.
// <!-- END clutchdoc -->

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	resolverv1 "github.com/lyft/clutch/backend/api/resolver/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/resolver"
)

const Name = "clutch.module.resolver"

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	m := &mod{
		api: newAPI(),
	}
	return m, nil
}

type mod struct {
	api resolverv1.ResolverAPIServer
}

func (m *mod) Register(r module.Registrar) error {
	resolverv1.RegisterResolverAPIServer(r.GRPCServer(), m.api)
	return r.RegisterJSONGateway(resolverv1.RegisterResolverAPIHandler)
}

func newAPI() resolverv1.ResolverAPIServer {
	return &resolverAPI{}
}

type resolverAPI struct{}

func (r *resolverAPI) Resolve(ctx context.Context, req *resolverv1.ResolveRequest) (*resolverv1.ResolveResponse, error) {
	resp := newResponse()

	var searchedSchemas []string
	for _, res := range resolver.Registry {
		// TODO: fan-out, fan-in for speeeed.
		// TODO: dedupe results, as technically multiple resolvers could
		//  resolve the same input schema (not yet though), and return the same object.
		inputSchemas, ok := res.Schemas()[req.Want]
		if !ok {
			continue
		}

		for _, schema := range inputSchemas {
			if schema.TypeUrl == req.Have.TypeUrl {
				searchedSchemas = append(searchedSchemas, schema.Metadata.DisplayName)
				a := &ptypes.DynamicAny{}
				if err := ptypes.UnmarshalAny(req.Have, a); err != nil {
					return nil, err
				}

				results, err := res.Resolve(ctx, req.Want, a.Message, req.Limit)
				if err != nil {
					return nil, err
				}

				err = resp.marshalResults(results)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	resp.truncate(req.Limit)
	if err := resp.isError(req.Want, searchedSchemas); err != nil {
		return nil, err
	}

	return &resolverv1.ResolveResponse{
		Results:         resp.Results,
		PartialFailures: resp.PartialFailures,
	}, nil
}

func (r *resolverAPI) Search(ctx context.Context, req *resolverv1.SearchRequest) (*resolverv1.SearchResponse, error) {
	resp := newResponse()

	var searchedSchemas []string
	for _, res := range resolver.Registry {
		resSchemas := res.Schemas()
		if schemas, ok := resSchemas[req.Want]; ok {
			for _, ss := range schemas {
				if ss.Metadata.Searchable || (ss.Metadata.Search != nil && ss.Metadata.Search.Enabled) {
					searchedSchemas = append(searchedSchemas, ss.Metadata.DisplayName)
				}
			}

			results, err := res.Search(ctx, req.Want, req.Query, req.Limit)
			if err != nil {
				return nil, err
			}

			err = resp.marshalResults(results)
			if err != nil {
				return nil, err
			}
		}
	}

	resp.truncate(req.Limit)
	if err := resp.isError(req.Want, searchedSchemas); err != nil {
		return nil, err
	}

	return &resolverv1.SearchResponse{
		Results:         resp.Results,
		PartialFailures: resp.PartialFailures,
	}, nil
}

func (r *resolverAPI) GetObjectSchemas(ctx context.Context, req *resolverv1.GetObjectSchemasRequest) (*resolverv1.GetObjectSchemasResponse, error) {
	var schemas []*resolverv1.Schema

	// Find schemas that match the requested type.
	for _, res := range resolver.Registry {
		resSchemas := res.Schemas()
		if typeSchemas, ok := resSchemas[req.TypeUrl]; ok {
			for _, typeSchema := range typeSchemas {
				// Make a clone of each schema in the event we need to modify below.
				schemas = append(schemas, proto.Clone(typeSchema).(*resolverv1.Schema))
			}
		}
	}

	// Handle option presentation for matching schemas.
	for _, schema := range schemas {
		for _, field := range schema.Fields {
			fm, ok := field.Metadata.Type.(*resolverv1.FieldMetadata_OptionField)
			if !ok {
				continue
			}
			for _, opt := range fm.OptionField.Options {
				if opt.DisplayName == "" {
					switch t := opt.Value.(type) {
					case *resolverv1.Option_StringValue:
						opt.DisplayName = t.StringValue
					default:
						opt.DisplayName = fmt.Sprintf("%s (auto display name TODO)", opt.Value)
					}
				}
			}
			// include_all_option.
			if fm.OptionField.IncludeAllOption && len(fm.OptionField.Options) != 1 {
				all := &resolverv1.Option{
					DisplayName: "All",
					Value:       &resolverv1.Option_StringValue{StringValue: resolver.OptionAll},
				}
				fm.OptionField.Options = append([]*resolverv1.Option{all}, fm.OptionField.Options...)
			}

			// add error if needed
			updateSchemaError(schema)
		}
	}

	return &resolverv1.GetObjectSchemasResponse{
		TypeUrl: req.TypeUrl,
		Schemas: schemas,
	}, nil
}

func (r *resolverAPI) AutoComplete(ctx context.Context, req *resolverv1.AutocompleteRequest) (*resolverv1.AutocompleteResponse, error) {
	var err error
	results := []string{}

	// Iterate through all of the available resolvers & schemas to find the one requested
	// If that schema exists then we call the associated autocomplete function for that resolver
	for _, res := range resolver.Registry {
		resSchema := res.Schemas()
		if _, ok := resSchema[req.Want]; ok {
			results, err = res.AutoComplete(ctx, req.Want, req.Search, req.ResultLimit)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	return &resolverv1.AutocompleteResponse{
		Results: results,
	}, nil
}

// Add error information to the schema if it's broken in some way.
func updateSchemaError(input *resolverv1.Schema) {
	for _, field := range input.Fields {
		switch t := field.Metadata.Type.(type) {
		case *resolverv1.FieldMetadata_OptionField:
			if field.Metadata.Required && len(t.OptionField.Options) == 0 {
				s := status.New(codes.OutOfRange, fmt.Sprintf("missing required options for field '%s'", field.Name))
				input.Error = s.Proto()
				return
			}
		}
	}
}
