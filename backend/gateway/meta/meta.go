package meta

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
)

var (
	// TODO: Lock after startup.
	methodDescriptors map[string]*desc.MethodDescriptor

	fieldNameRegexp = regexp.MustCompile(`{(\w+)}`)

	actionTypeDescriptor     = apiv1.E_Action.TypeDescriptor()
	identifierTypeDescriptor = apiv1.E_Id.TypeDescriptor()
	redactedTypeDescriptor   = apiv1.E_Redacted.TypeDescriptor()
	referenceTypeDescriptor  = apiv1.E_Reference.TypeDescriptor()
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
	opts := md.GetMethodOptions().ProtoReflect()

	if !opts.Has(actionTypeDescriptor) {
		return apiv1.ActionType_UNSPECIFIED
	}
	return opts.Get(actionTypeDescriptor).Message().Interface().(*apiv1.Action).Type
}

func IsRedacted(pb proto.Message) bool {
	m := pb.ProtoReflect()
	opts := m.Descriptor().Options().ProtoReflect()
	return opts.Has(redactedTypeDescriptor) && opts.Get(redactedTypeDescriptor).Bool()
}

func ResourceNames(pb proto.Message) []*auditv1.Resource {
	m := pb.ProtoReflect()
	opts := m.Descriptor().Options().ProtoReflect()

	if opts.Has(identifierTypeDescriptor) {
		v := opts.Get(identifierTypeDescriptor)
		id := v.Message().Interface().(*apiv1.Identifier)

		names := make([]*auditv1.Resource, 0, len(id.Patterns))
		for _, pattern := range id.Patterns {
			if newName := resolvePattern(pb, pattern); newName != nil {
				names = append(names, newName)
			}
		}
		return names
	}

	if opts.Has(referenceTypeDescriptor) {
		v := opts.Get(referenceTypeDescriptor)
		ref := v.Message().Interface().(*apiv1.Reference)

		// Best effort sizing to avoid reallocations.
		names := make([]*auditv1.Resource, 0, len(ref.Fields))
		for _, field := range ref.Fields {
			for _, resolved := range resolveField(pb, field) {
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

func resolveField(pb proto.Message, name string) []*auditv1.Resource {
	m := pb.ProtoReflect()
	fd := m.Descriptor().Fields().ByName(protoreflect.Name(name))
	if fd == nil {
		return nil
	}

	v := m.Get(fd)

	if fd.IsList() {
		return resolveSlice(v.List())
	}

	return ResourceNames(v.Message().Interface())
}

func resolveSlice(list protoreflect.List) []*auditv1.Resource {
	var resources []*auditv1.Resource
	for i := 0; i < list.Len(); i++ {
		v := list.Get(i)
		resources = append(resources, ResourceNames(v.Message().Interface())...)
	}
	return resources
}

func resolvePattern(pb proto.Message, pattern *apiv1.Pattern) *auditv1.Resource {
	m := pb.ProtoReflect()
	fields := m.Descriptor().Fields()

	resourceName := pattern.Pattern

	substitutions := fieldNameRegexp.FindAllStringSubmatch(pattern.Pattern, -1)
	for _, name := range substitutions {
		fd := fields.ByName(protoreflect.Name(name[1]))
		if fd == nil {
			continue
		}
		v := m.Get(fd)
		resourceName = strings.Replace(resourceName, name[0], v.String(), 1)
	}
	return &auditv1.Resource{TypeUrl: pattern.TypeUrl, Id: resourceName}
}

// APIBody returns a API request/response interface as an anypb.Any message.
func APIBody(body interface{}) (*anypb.Any, error) {
	m, ok := body.(proto.Message)
	if !ok {
		// body is not the type/value we want to process
		return nil, nil
	}

	return anypb.New(m)
}
