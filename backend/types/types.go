package types

import "github.com/golang/protobuf/ptypes/any"

type Action int

const (
	CREATE Action = iota
	UPDATE
	DELETE
)

type CacheableTopology interface {
	CacheEnabled() bool
	GetTopologyObjectChannel() chan TopologyObject
}

type TopologyObject struct {
	Id              string
	ResolverTypeURL string
	Pb              *any.Any
	Metadata        map[string]string
	Action          Action
}
