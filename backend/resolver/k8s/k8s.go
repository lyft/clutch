package k8s

// <!-- START clutchdoc -->
// description: Locates resources in Kubernetes.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sv1api "github.com/lyft/clutch/backend/api/k8s/v1"
	k8sv1resolver "github.com/lyft/clutch/backend/api/resolver/k8s/v1"
	resolverv1 "github.com/lyft/clutch/backend/api/resolver/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/resolver"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/k8s"
	"github.com/lyft/clutch/backend/service/topology"
)

const Name = "clutch.resolver.k8s"

var typeURLPod = meta.TypeURL((*k8sv1api.Pod)(nil))
var typeURLHPA = meta.TypeURL((*k8sv1api.HPA)(nil))
var typeURLNode = meta.TypeURL((*k8sv1api.Node)(nil))

var typeSchemas = resolver.TypeURLToSchemaMessagesMap{
	typeURLPod: {
		(*k8sv1resolver.PodID)(nil),
	},
	typeURLHPA: {
		(*k8sv1resolver.HPAName)(nil),
	},
	typeURLNode: {
		(*k8sv1resolver.Node)(nil),
	},
}

// Loosely https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names
var idPattern = regexp.MustCompile(`[a-fA-F0-9-\.]{1,253}`)

func makeClientsetOptions(clientsets []string) []*resolverv1.Option {
	ret := make([]*resolverv1.Option, len(clientsets))
	for i, name := range clientsets {
		ret[i] = &resolverv1.Option{
			Value: &resolverv1.Option_StringValue{
				StringValue: name,
			},
		}
	}
	return ret
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (resolver.Resolver, error) {
	k8sRegistered, ok := service.Registry["clutch.service.k8s"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	svc, ok := k8sRegistered.(k8s.Service)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	var topologyService topology.Service
	if svc, ok := service.Registry[topology.Name]; ok {
		topologyService, ok = svc.(topology.Service)
		if !ok {
			return nil, errors.New("incorrect topology service type")
		}
		logger.Debug("enabling autocomplete api for the k8s resolver")
	}

	schemas, err := resolver.InputsToSchemas(typeSchemas)
	if err != nil {
		return nil, err
	}

	clientsets, err := svc.Clientsets(context.Background())
	if err != nil {
		return nil, err
	}

	resolver.HydrateDynamicOptions(schemas, map[string][]*resolverv1.Option{
		"clientset": makeClientsetOptions(clientsets),
	})

	r := &res{
		svc:      svc,
		topology: topologyService,
		schemas:  schemas,
	}
	return r, nil
}

type res struct {
	svc      k8s.Service
	topology topology.Service
	schemas  resolver.TypeURLToSchemasMap
}

func (r *res) Schemas() resolver.TypeURLToSchemasMap { return r.schemas }

func (r *res) locateByPodID(ctx context.Context, in *k8sv1resolver.PodID) ([]*k8sv1api.Pod, error) {
	// Only possible to get one at a time by PodID.
	pod, err := r.svc.DescribePod(ctx, in.Clientset, "", in.Namespace, in.Name)
	if err != nil {
		return nil, err
	}
	return []*k8sv1api.Pod{pod}, nil
}

func (r *res) resolveForPod(ctx context.Context, input proto.Message) ([]*k8sv1api.Pod, error) {
	switch i := input.(type) {
	case *k8sv1resolver.PodID:
		return r.locateByPodID(ctx, i)
	default:
		// TODO: IP address via List
		return nil, status.Errorf(codes.Internal, "unrecognized input type %T", i)
	}
}

func (r *res) locateByHPAName(ctx context.Context, in *k8sv1resolver.HPAName) ([]*k8sv1api.HPA, error) {
	// Only possible to get one at a time by name.
	hpa, err := r.svc.DescribeHPA(ctx, in.Clientset, "", in.Namespace, in.Name)
	if err != nil {
		return nil, err
	}
	return []*k8sv1api.HPA{hpa}, nil
}

func (r *res) resolveForHPA(ctx context.Context, input proto.Message) ([]*k8sv1api.HPA, error) {
	switch i := input.(type) {
	case *k8sv1resolver.HPAName:
		return r.locateByHPAName(ctx, i)
	default:
		return nil, status.Errorf(codes.Internal, "unrecognized input type '%T'", i)
	}
}

func (r *res) locateByNodeName(ctx context.Context, in *k8sv1resolver.Node) ([]*k8sv1api.Node, error) {
	// Only possible to get one at a time by name.
	node, err := r.svc.DescribeNode(ctx, in.Clientset, in.Cluster, in.Name)
	if err != nil {
		return nil, err
	}
	return []*k8sv1api.Node{node}, nil
}

func (r *res) resolveForNode(ctx context.Context, input proto.Message) ([]*k8sv1api.Node, error) {
	switch i := input.(type) {
	case *k8sv1resolver.Node:
		return r.locateByNodeName(ctx, i)
	default:
		return nil, status.Errorf(codes.Internal, "unrecognized input type '%T'", i)
	}
}

func (r *res) Resolve(ctx context.Context, typeURL string, input proto.Message, limit uint32) (*resolver.Results, error) {
	switch typeURL {
	case typeURLPod:
		result, err := r.resolveForPod(ctx, input)
		if err != nil {
			return nil, err
		}
		return &resolver.Results{Messages: resolver.MessageSlice(result)}, nil
	case typeURLHPA:
		result, err := r.resolveForHPA(ctx, input)
		if err != nil {
			return nil, err
		}
		return &resolver.Results{Messages: resolver.MessageSlice(result)}, nil
	case typeURLNode:
		result, err := r.resolveForNode(ctx, input)
		if err != nil {
			return nil, err
		}
		return &resolver.Results{Messages: resolver.MessageSlice(result)}, nil
	default:
		return nil, status.Errorf(codes.Internal, "don't know how to resolve type %s", typeURL)
	}
}

func (r *res) Search(ctx context.Context, typeURL, query string, limit uint32) (*resolver.Results, error) {
	clientsets, err := r.svc.Clientsets(ctx)
	if err != nil {
		return nil, err
	}

	ctx, handler := resolver.NewFanoutHandler(ctx)
	switch typeURL {
	case typeURLPod:
		if idPattern.MatchString(query) {
			patternValues, ok, err := meta.ExtractPatternValuesFromString((*k8sv1api.Pod)(nil), query)
			if err != nil {
				return nil, err
			}

			namespace := metav1.NamespaceAll
			podQuery := query
			cluster := ""

			if ok {
				namespace = patternValues["namespace"]
				podQuery = patternValues["name"]
				cluster = patternValues["cluster"]
			}

			for _, clientset := range clientsets {
				handler.Add(1)
				go func(clientset, cluster, namespace, name string) {
					defer handler.Done()
					pod, err := r.svc.DescribePod(ctx, clientset, cluster, namespace, name)
					select {
					case handler.Channel() <- resolver.NewFanoutResult([]*k8sv1api.Pod{pod}, err):
						return
					case <-handler.Cancelled():
						return
					}
				}(clientset, cluster, namespace, podQuery)
			}
		} else {
			return nil, status.Error(codes.InvalidArgument, "did not understand input")
		}
	case typeURLHPA:
		if idPattern.MatchString(query) {
			patternValues, ok, err := meta.ExtractPatternValuesFromString((*k8sv1api.HPA)(nil), query)
			if err != nil {
				return nil, err
			}

			namespace := metav1.NamespaceAll
			hpaQuery := query
			cluster := ""

			if ok {
				namespace = patternValues["namespace"]
				hpaQuery = patternValues["name"]
				cluster = patternValues["cluster"]
			}

			for _, clientset := range clientsets {
				handler.Add(1)
				go func(clientset, cluster, namespace, query string) {
					defer handler.Done()
					hpa, err := r.svc.DescribeHPA(ctx, clientset, cluster, namespace, query)
					select {
					case handler.Channel() <- resolver.NewFanoutResult([]*k8sv1api.HPA{hpa}, err):
						return
					case <-handler.Cancelled():
						return
					}
				}(clientset, cluster, namespace, hpaQuery)
			}
		} else {
			return nil, status.Error(codes.InvalidArgument, "did not understand input")
		}
	case typeURLNode:
		if idPattern.MatchString(query) {
			patternValues, ok, err := meta.ExtractPatternValuesFromString((*k8sv1api.Node)(nil), query)
			if err != nil {
				return nil, err
			}

			nodeQuery := query
			cluster := ""

			if ok {
				nodeQuery = patternValues["name"]
				cluster = patternValues["cluster"]
			}

			for _, clientset := range clientsets {
				handler.Add(1)
				go func(clientset, cluster, query string) {
					defer handler.Done()
					node, err := r.svc.DescribeNode(ctx, clientset, cluster, query)
					select {
					case handler.Channel() <- resolver.NewFanoutResult([]*k8sv1api.Node{node}, err):
						return
					case <-handler.Cancelled():
						return
					}
				}(clientset, cluster, nodeQuery)
			}
		} else {
			return nil, status.Error(codes.InvalidArgument, "did not understand input")
		}
	default:
		return nil, status.Error(codes.Internal, fmt.Sprintf("cannot search for type '%s'", typeURL))
	}

	return handler.Results(limit)
}

func (r *res) Autocomplete(ctx context.Context, typeURL, search string, limit uint64, caseSensitive bool) ([]*resolverv1.AutocompleteResult, error) {
	if r.topology == nil {
		return nil, status.Error(codes.FailedPrecondition, "topology service must be enabled to use the K8s autocomplete API")
	}

	var resultLimit uint64 = resolver.DefaultAutocompleteLimit
	if limit > 0 {
		resultLimit = limit
	}

	results, err := r.topology.Autocomplete(ctx, typeURL, search, resultLimit, caseSensitive)
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
