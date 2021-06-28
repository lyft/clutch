package meta

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
)

var (
	// TODO: Lock after startup.
	methodDescriptors map[string]*desc.MethodDescriptor

	fieldNameRegexp = regexp.MustCompile(`{(\w+)}`)

	actionTypeDescriptor          = apiv1.E_Action.TypeDescriptor()
	auditDisabledTypeDescriptor   = apiv1.E_DisableAudit.TypeDescriptor()
	identifierTypeDescriptor      = apiv1.E_Id.TypeDescriptor()
	redactedMessageTypeDescriptor = apiv1.E_Redacted.TypeDescriptor()
	referenceTypeDescriptor       = apiv1.E_Reference.TypeDescriptor()
)

const typePrefix = "type.googleapis.com/"

func TypeURL(pb proto.Message) string {
	return typePrefix + string(pb.ProtoReflect().Descriptor().FullName())
}

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

func IsAuditDisabled(method string) bool {
	md, ok := methodDescriptors[method]
	if !ok {
		return false
	}
	opts := md.GetMethodOptions().ProtoReflect()
	return opts.Has(auditDisabledTypeDescriptor) && opts.Get(auditDisabledTypeDescriptor).Bool()
}

// If fields have the option log set to false,
func ClearLogDisabledFields(m proto.Message) proto.Message {
	if m == nil {
		return m
	}

	pb := m.ProtoReflect()
	pb.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		opts := fd.Options().(*descriptorpb.FieldOptions)
		if proto.HasExtension(opts, apiv1.E_Log) && !proto.GetExtension(opts, apiv1.E_Log).(bool) {
			pb.Clear(fd)
			return true // Continue.
		}

		// Handle nested types.
		switch t := v.Interface().(type) {
		case protoreflect.Message:
			ClearLogDisabledFields(t.Interface())
		case protoreflect.Map:
			t.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
				if _, ok := v.Interface().(protoreflect.Message); ok {
					ClearLogDisabledFields(v.Message().Interface())
				}
				return true
			})
		case protoreflect.List: // i.e. `repeated`.
			for i := 0; i < t.Len(); i++ {
				if _, ok := t.Get(i).Interface().(protoreflect.Message); ok {
					ClearLogDisabledFields(t.Get(i).Message().Interface())
				}
			}
		}
		return true
	})

	return m
}

func IsRedacted(pb proto.Message) bool {
	m := pb.ProtoReflect()
	opts := m.Descriptor().Options().ProtoReflect()
	return opts.Has(redactedMessageTypeDescriptor) && opts.Get(redactedMessageTypeDescriptor).Bool()
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

// HydratedPatternForProto takes a proto and returns its pattern populated with values
func HydratedPatternForProto(pb proto.Message) (string, error) {
	m := pb.ProtoReflect()
	opts := m.Descriptor().Options().ProtoReflect()

	populatedPattern := []string{}

	if opts.Has(identifierTypeDescriptor) {
		v := opts.Get(identifierTypeDescriptor)
		id := v.Message().Interface().(*apiv1.Identifier)

		for _, pattern := range id.Patterns {
			rs := resolvePattern(pb, pattern)
			populatedPattern = append(populatedPattern, rs.Id)
		}

		// At the time of writing there is only support for a single pattern
		// this list should only have one item to return
		return populatedPattern[0], nil
	}

	return "", fmt.Errorf("the supplied proto does not have a pattern: [%T]", pb)
}

// ExtractPatternValuesFromString takes a string value and maps the patterns from a proto pattern
// this is utilized by the resolver search api
//
// For example given the following proto pattern
// option (clutch.api.v1.id).patterns = {
//  pattern : "{cluster}/{namespace}/{name}"
// };
//
// And the value of "mycluster/mynamespace/nameofresource"
// we transform the pattern into a regex and map the values to the pattern names
//
// The output for this example is:
// map[string]string{
//  cluster: mycluster
//  namespace: mynamespace
//  name: nameofresource
// }
func ExtractPatternValuesFromString(pb proto.Message, value string) (map[string]string, bool, error) {
	m := pb.ProtoReflect()
	opts := m.Descriptor().Options().ProtoReflect()

	// Field and Value result map
	result := map[string]string{}

	if opts.Has(identifierTypeDescriptor) {
		v := opts.Get(identifierTypeDescriptor)
		id := v.Message().Interface().(*apiv1.Identifier)

		for _, pattern := range id.Patterns {
			// The variable names on the pattern
			patternFields := extractProtoPatternFieldNames(pattern)

			// Convert the pattern into a regex
			convertedRegex := fmt.Sprintf("^%s$", fieldNameRegexp.ReplaceAllString(pattern.Pattern, "(.*)"))
			patternRegex, err := regexp.Compile(convertedRegex)
			if err != nil {
				return nil, false, err
			}

			// Extract the regex groups, index 0 is always the input string
			subStringGroups := patternRegex.FindAllStringSubmatch(value, -1)
			if subStringGroups != nil {
				for i, name := range patternFields {
					// Plus one here because the first value is the input string
					result[name] = subStringGroups[0][i+1]
				}
			}
		}
	}

	// If we dont have any results then we can just return false
	if len(result) == 0 {
		return result, false, nil
	}

	// Check that all of the fields have values
	for _, value := range result {
		if len(value) == 0 {
			return result, false, nil
		}
	}

	return result, true, nil
}

func extractProtoPatternFieldNames(pattern *apiv1.Pattern) []string {
	variableNames := fieldNameRegexp.FindAllStringSubmatch(pattern.Pattern, -1)
	results := make([]string, 0, len(variableNames))
	for _, name := range variableNames {
		results = append(results, name[1])
	}
	return results
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

// APIBody returns a API request/response interface as an anypb.Any message, with any redaction or clearing based on
// message or field annotations.
func APIBody(body interface{}) (*anypb.Any, error) {
	m, ok := body.(proto.Message)
	if !ok {
		// body is not the type/value we want to process
		return nil, nil
	}

	if IsRedacted(m) {
		return anypb.New(&apiv1.Redacted{RedactedTypeUrl: TypeURL(m)})
	}

	// Deep copy before field redaction so we do not unintentionally remove fields
	// from the original object that were passed by reference
	m = proto.Clone(m)
	return anypb.New(ClearLogDisabledFields(m))
}

/* ToValue converts custom types to a structpb.Value. This helper was added
since structpb.NewValue has a limited set of types that it supports.
More details here: https://github.com/golang/protobuf/issues/1302#issuecomment-805453221
*/
func ToValue(data interface{}) (*structpb.Value, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	v := &structpb.Value{}
	err = v.UnmarshalJSON(b)
	if err != nil {
		return nil, err
	}

	return v, nil
}
