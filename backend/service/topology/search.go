package topology

import (
	"context"
	"fmt"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
)

func (c *client) SearchTopology(ctx context.Context, request *topologyv1.SearchTopologyRequest) error {
	return fmt.Errorf("not yet implemented")
}
