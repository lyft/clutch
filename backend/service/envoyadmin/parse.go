package envoyadmin

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

	envoy_admin_v3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"

	envoytriagev1 "github.com/lyft/clutch/backend/api/envoytriage/v1"
)

func unmarshal(v []byte, pb proto.Message) error {
	um := &jsonpb.Unmarshaler{AllowUnknownFields: true}
	return um.Unmarshal(bytes.NewReader(v), pb)
}

func nodeMetadataFromResponse(resp []byte) (*envoytriagev1.NodeMetadata, error) {
	pb := &envoy_admin_v3.ServerInfo{}
	if err := unmarshal(resp, pb); err != nil {
		return nil, err
	}

	return &envoytriagev1.NodeMetadata{
		ServiceNode:    pb.CommandLineOptions.ServiceNode,
		ServiceCluster: pb.CommandLineOptions.ServiceCluster,
		ServiceZone:    pb.CommandLineOptions.ServiceZone,
		Version:        pb.Version,
	}, nil
}

func addrAsString(address *envoy_config_core_v3.Address) string {
	switch in := address.Address.(type) {
	case *envoy_config_core_v3.Address_SocketAddress:
		sa := in.SocketAddress
		return fmt.Sprintf("%s://%s:%d", strings.ToLower(sa.GetProtocol().String()), sa.Address, sa.GetPortValue())
	case *envoy_config_core_v3.Address_Pipe:
		return fmt.Sprintf("unix://%s", in.Pipe.Path)
	default:
		return "unknown address format"
	}
}

type RuntimeEntry struct {
	FinalValue  string   `json:"final_value"`
	LayerValues []string `json:"layer_values"`
}

type Runtime struct {
	Entries map[string]*RuntimeEntry `json:"entries"`
}

var scalarStatPattern = regexp.MustCompile(`^([\w.]+): (\d+)$`)

func statsFromResponse(resp []byte) (*envoytriagev1.Stats, error) {
	scanner := bufio.NewScanner(bytes.NewReader(resp))

	var stats []*envoytriagev1.Stats_Stat
	for scanner.Scan() {
		matches := scalarStatPattern.FindStringSubmatch(scanner.Text())
		if len(matches) == 3 {
			v, _ := strconv.ParseUint(matches[2], 10, 64)
			stats = append(stats, &envoytriagev1.Stats_Stat{Key: matches[1], Value: v})
		}
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Key < stats[j].Key
	})

	return &envoytriagev1.Stats{
		Stats: stats,
	}, nil
}

func runtimeFromResponse(resp []byte) (*envoytriagev1.Runtime, error) {
	r := &Runtime{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	entries := make([]*envoytriagev1.Runtime_Entry, 0, len(r.Entries))
	for key, entry := range r.Entries {
		entries = append(entries, &envoytriagev1.Runtime_Entry{
			Key:  key,
			Type: &envoytriagev1.Runtime_Entry_Value{Value: entry.FinalValue},
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return &envoytriagev1.Runtime{Entries: entries}, nil
}

func listenersFromResponse(resp []byte) (*envoytriagev1.Listeners, error) {
	pb := &envoy_admin_v3.Listeners{}
	if err := unmarshal(resp, pb); err != nil {
		return nil, err
	}

	ret := &envoytriagev1.Listeners{
		ListenerStatuses: make([]*envoytriagev1.ListenerStatus, len(pb.ListenerStatuses)),
	}

	for idx, ls := range pb.ListenerStatuses {
		ret.ListenerStatuses[idx] = &envoytriagev1.ListenerStatus{
			Name:         ls.Name,
			LocalAddress: addrAsString(ls.LocalAddress),
		}
	}

	return ret, nil
}

func serverInfoFromResponse(resp []byte) (*envoytriagev1.ServerInfo, error) {
	ret := &envoytriagev1.ServerInfo{
		Value: &structpb.Value{},
	}
	if err := unmarshal(resp, ret.Value); err != nil {
		return nil, err
	}
	return ret, nil
}

func configDumpFromResponse(resp []byte) (*envoytriagev1.ConfigDump, error) {
	ret := &envoytriagev1.ConfigDump{
		Value: &structpb.Value{},
	}
	if err := unmarshal(resp, ret.Value); err != nil {
		return nil, err
	}
	return ret, nil
}

func healthy(status *envoy_admin_v3.HostHealthStatus) bool {
	unhealthy := status.FailedActiveDegradedCheck ||
		status.FailedActiveHealthCheck ||
		status.FailedOutlierCheck ||
		status.PendingActiveHc ||
		status.PendingDynamicRemoval ||
		(status.EdsHealthStatus != envoy_config_core_v3.HealthStatus_HEALTHY)
	return !unhealthy
}

func clustersFromResponse(resp []byte) (*envoytriagev1.Clusters, error) {
	pb := &envoy_admin_v3.Clusters{}
	if err := unmarshal(resp, pb); err != nil {
		return nil, err
	}

	ret := &envoytriagev1.Clusters{
		ClusterStatuses: make([]*envoytriagev1.ClusterStatus, len(pb.ClusterStatuses)),
	}

	for i, cluster := range pb.ClusterStatuses {
		hostStatuses := make([]*envoytriagev1.HostStatus, len(cluster.HostStatuses))
		for j, hs := range cluster.HostStatuses {
			hostStatuses[j] = &envoytriagev1.HostStatus{
				Address: addrAsString(hs.Address),
				Healthy: healthy(hs.HealthStatus),
			}
		}

		ret.ClusterStatuses[i] = &envoytriagev1.ClusterStatus{
			Name:         cluster.Name,
			HostStatuses: hostStatuses,
		}
	}

	return ret, nil
}
