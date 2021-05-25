package resolver

import (
	"context"
	"fmt"
	"reflect"

	protodeprecated "github.com/golang/protobuf/proto"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	resolverv1 "github.com/lyft/clutch/backend/api/resolver/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
)

const (
	OptionAll = "__ALL__"
	// TODO: Layout the ground work for resolver configurations
	// allowing a user to set the default autocomplete limit
	DefaultAutocompleteLimit = 50
)

// Input map from resolvers with the type URL of the resolved object mapped to the schema's message.
// e.g. {"v1.Instance": [v1.Name, v1.IPAddress]}
type TypeURLToSchemaMessagesMap map[string][]proto.Message

// Output map used by the resolver after the input TypeURLToSchemaMessagesMap was used to read metadata from the
// protobuf objects.
// e.g. {"v1.Instance": [{"DisplayName": "Name", "Searchable": true}, {"DisplayName": "IP Address", "Searchable": true}]}
type TypeURLToSchemasMap map[string][]*resolverv1.Schema

type Factory map[string]func(*anypb.Any, *zap.Logger, tally.Scope) (Resolver, error)

var Registry = map[string]Resolver{}

type Results struct {
	Messages        []proto.Message
	PartialFailures []*status.Status
}

type Resolver interface {
	Schemas() TypeURLToSchemasMap

	Search(ctx context.Context, typeURL, query string, limit uint32) (*Results, error)
	// ValidateSearch(typeURL string, query string) error for async validation from frontend

	Resolve(ctx context.Context, typeURL string, input protodeprecated.Message, limit uint32) (*Results, error)
	// ValidateResolveInput(typeURL string, input proto.Message) for async validation from frontend

	Autocomplete(ctx context.Context, typeURL, search string, limit uint64) ([]*resolverv1.AutocompleteResult, error)
}

const TypePrefix = "type.googleapis.com/"

// Deprecated: use meta.TypeURL instead, will require moving to new proto APIs.
func TypeURL(m protodeprecated.Message) string {
	return TypePrefix + string(protodeprecated.MessageReflect(m).Descriptor().FullName())
}

func MarshalProtoSlice(pbs interface{}) ([]*anypb.Any, error) {
	if pbs == nil {
		return nil, nil
	}

	switch t := reflect.TypeOf(pbs).Kind(); t {
	case reflect.Slice:
		// OK.
	default:
		return nil, fmt.Errorf("tried to marshal slice but received %s", t)
	}

	s := reflect.ValueOf(pbs)
	ret := make([]*anypb.Any, s.Len())
	for i := 0; i < s.Len(); i++ {
		item := s.Index(i)

		v, ok := item.Interface().(proto.Message)
		if !ok {
			return nil, fmt.Errorf("could not use %s as proto.Message", item.Kind())
		}
		a, err := anypb.New(v)
		if err != nil {
			return nil, err
		}
		ret[i] = a
	}

	return ret, nil
}

func HydrateDynamicOptions(schemas TypeURLToSchemasMap, options map[string][]*resolverv1.Option) {
	for _, schemasForType := range schemas {
		for _, schema := range schemasForType {
			for _, field := range schema.Fields {
				// Check each option field's annotation for include_dynamic_options and hydrate if there's a match.
				if f, ok := field.Metadata.Type.(*resolverv1.FieldMetadata_OptionField); ok {
					for _, include := range f.OptionField.IncludeDynamicOptions {
						if opts, ok := options[include]; ok {
							f.OptionField.Options = append(f.OptionField.Options, opts...)
						}
					}
				}
			}
		}
	}
}

// Pass in annotated resolver input objects and return schemas for them.
func InputsToSchemas(typeSchemas TypeURLToSchemaMessagesMap) (TypeURLToSchemasMap, error) {
	schemas := make(TypeURLToSchemasMap, len(typeSchemas))

	for typeURL, inputObjects := range typeSchemas {
		schemas[typeURL] = make([]*resolverv1.Schema, len(inputObjects))
		for i, inputObject := range inputObjects {
			desc := inputObject.ProtoReflect().Descriptor()
			ext := proto.GetExtension(desc.Options(), resolverv1.E_Schema)
			md := ext.(*resolverv1.SchemaMetadata)

			fds := desc.Fields()

			schema := &resolverv1.Schema{
				TypeUrl:  meta.TypeURL(inputObject),
				Fields:   make([]*resolverv1.Field, fds.Len()),
				Metadata: md,
			}

			// Fill fields from per-field annotations.
			for j := 0; j < fds.Len(); j++ {
				fd := fds.Get(j)
				fext := proto.GetExtension(fd.Options(), resolverv1.E_SchemaField)
				fieldMeta := fext.(*resolverv1.FieldMetadata)
				// Clone the fieldMeta since it's mutable (i.e. dynamic options).
				fieldMeta = proto.Clone(fieldMeta).(*resolverv1.FieldMetadata)

				name := string(fd.Name())
				if fd.HasJSONName() {
					// TODO(maybe): this should probably always respond with Name instead of JsonName for gRPC clients.
					// Would need to check context and add a flag.
					name = fd.JSONName()
				}

				// Use default display name of field name if none was provided.
				if fieldMeta.DisplayName == "" {
					fieldMeta.DisplayName = name
				}

				schema.Fields[j] = &resolverv1.Field{
					Name:     name,
					Metadata: fieldMeta,
				}
			}
			schemas[typeURL][i] = schema
		}
	}
	return schemas, nil
}
