package meta

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/iancoleman/strcase"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
)

var (
	// TODO: Lock after startup.
	methodDescriptors map[string]*desc.MethodDescriptor

	fieldNameRegexp = regexp.MustCompile(`{(\w+)}`)
)

func GenerateGRPCMetadata(server *grpc.Server) error {
	serviceDescriptors, err := grpcreflect.LoadServiceDescriptors(server)
	if err != nil {
		return err
	}

	mds := make(map[string]*desc.MethodDescriptor)
	for _, sd := range serviceDescriptors {
		for _, md := range sd.GetMethods() {
			methodName := fmt.Sprintf("/%s/%s", sd.GetFullyQualifiedName(), md.GetName())
			mds[methodName] = md
		}
	}

	methodDescriptors = mds
	return nil
}

func GetAction(method string) apiv1.ActionType {
	md, ok := methodDescriptors[method]
	if !ok {
		return apiv1.ActionType_UNSPECIFIED
	}

	opts := md.GetMethodOptions()
	ext, err := proto.GetExtension(opts, apiv1.E_Action)
	if err != nil {
		return apiv1.ActionType_UNSPECIFIED
	}

	action := ext.(*apiv1.Action)
	return action.Type
}

func ResourceNames(message descriptor.Message) []*auditv1.Resource {
	_, descriptorMeta := descriptor.ForMessage(message)
	if proto.HasExtension(descriptorMeta.Options, apiv1.E_Id) {
		idExt, err := proto.GetExtension(descriptorMeta.Options, apiv1.E_Id)
		if err != nil {
			return nil
		}

		id := idExt.(*apiv1.Identifier)

		names := make([]*auditv1.Resource, 0, len(id.Patterns))
		for _, pattern := range id.Patterns {
			if newName := resolvePattern(message, pattern); newName != nil {
				names = append(names, resolvePattern(message, pattern))
			}
		}

		return names
	}

	if proto.HasExtension(descriptorMeta.Options, apiv1.E_Reference) {
		refExt, err := proto.GetExtension(descriptorMeta.Options, apiv1.E_Reference)
		if err != nil {
			return nil
		}

		ref := refExt.(*apiv1.Reference)

		// Best effort to avoid reallocations.
		names := make([]*auditv1.Resource, 0, len(ref.Fields))
		for _, field := range ref.Fields {
			for _, resolved := range resolveField(message, field) {
				if resolved == nil {
					continue
				}
				names = append(names, resolved)
			}
		}

		return names
	}

	return nil
}

func resolveField(message descriptor.Message, field string) []*auditv1.Resource {
	// Loop through fields by name looking for field
	rvalue := reflect.ValueOf(message)
	if rvalue.Kind() == reflect.Ptr {
		rvalue = rvalue.Elem()
	}

	// Somehow we haven't ended up with a proto message, kick out.
	if rvalue.Kind() != reflect.Struct {
		return nil
	}

	title := strcase.ToCamel(field)
	value := rvalue.FieldByName(title)

	// If we're looking at a slice (i.e. a repeated field), then resolve the name for each element.
	if value.Kind() == reflect.Slice {
		return resolveSlice(value)
	}

	// Otherwise, resolve for the plain type field here.
	message, ok := value.Interface().(descriptor.Message)
	if !ok {
		return nil
	}
	return ResourceNames(message)
}

func resolveSlice(value reflect.Value) []*auditv1.Resource {
	var resources []*auditv1.Resource
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		message, ok := item.Interface().(descriptor.Message)
		if !ok {
			continue
		}

		if resolved := ResourceNames(message); resolved != nil {
			resources = append(resources, resolved...)
		}
	}
	return resources
}

func resolvePattern(message descriptor.Message, pattern *apiv1.Pattern) *auditv1.Resource {
	rvalue := reflect.ValueOf(message)
	if rvalue.Kind() == reflect.Ptr {
		rvalue = rvalue.Elem()
	}

	// Somehow we haven't ended up with a proto message, kick out.
	if rvalue.Kind() != reflect.Struct {
		return nil
	}

	resourceName := pattern.Pattern

	substitutions := fieldNameRegexp.FindAllStringSubmatch(resourceName, -1)
	for _, name := range substitutions {
		// TODO: precompute name to field number? Would speed this up.
		title := strcase.ToCamel(name[1])
		// Get field by title
		resolved := rvalue.FieldByName(title).String()
		resourceName = strings.Replace(resourceName, name[0], resolved, 1)
	}

	return &auditv1.Resource{TypeUrl: pattern.TypeUrl, Id: resourceName}
}

// APIMetadata returns the API request/response interface as an anypb.Any message.
func APIMetadata(metadata interface{}) *any.Any {
	if metadata == nil {
		return nil
	}

	switch t := reflect.TypeOf(metadata).Kind(); t {
	// API request and response are type Ptr
	case reflect.Ptr:
		// OK.
	default:
		return nil
	}

	result := reflect.ValueOf(metadata)

	protomsg, ok := result.Interface().(proto.Message)
	if !ok {
		return nil
	}

	a, err := ptypes.MarshalAny(protomsg)
	if err != nil {
		return nil
	}

	return a
}
