package types

import "github.com/golang/protobuf/ptypes/any"

// TopologyCacheAction signifies to the Topology service what to do with an incomeing TopologyObject
//
// The topology service gets TopologyObjects off of the `GetTopologyObjectChannel` which is processed
// and stored in the topology_cache table.
//
type TopologyCacheAction int

const (
	// UPSERT creates or updates items in the topology cache table
	UPSERT TopologyCacheAction = iota
	DELETE
)

// CacheableTopology is implemented by a service that wishes to enable the topology API feature set
//
// By implmenting this interface the topology service will automatically setup all services which
// implement this interface. Automatically ingesting TopologyObjects via the `GetTopologyObjectChannel()` function.
// This enables users to make use of the Topology APIs with these new TopologyObjects.
//
type CacheableTopology interface {
	CacheEnabled() bool
	GetTopologyObjectChannel() chan TopologyObject
}

// TopologyObject is the representation of a single object, it could be a K8s pod or EC2 instance.
type TopologyObject struct {
	// Id is the unique identifer of the object.
	Id string
	// Pb is the clutch proto object.
	Pb *any.Any
	// Metadata is set by the service which produces TopologyObjects, for example k8s would extract
	// relevent metadata that gives the TopologyAPI the ability to query against it.
	Metadata map[string]string
	// Action denotes what the topology service should do with this object.
	Action TopologyCacheAction
}
