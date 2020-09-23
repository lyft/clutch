package topology

import (
	"context"
	"encoding/json"
	"log"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"

	"github.com/lyft/clutch/backend/service"
	"go.uber.org/zap"
)

func (c *client) startCaching() {
	if svc, ok := service.Registry["clutch.service.k8s"]; ok {
		k8sSvc, _ := svc.(CacheableTopology)
		go c.processTopologyObjectChannel(k8sSvc.GetTopologyObjectChannel())
	}
}

func (c *client) processTopologyObjectChannel(objs chan topologyv1.TopologyCacheObject) {
	for {
		obj := <-objs
		switch obj.Action {
		case topologyv1.TopologyCacheObject_CREATE_OR_UPDATE:
			c.SetCache(obj.TopologyObject)
		case topologyv1.TopologyCacheObject_DELETE:
			c.DeleteCache(obj.TopologyObject)
		}
	}
}

func (c *client) DeleteCache(obj *topologyv1.TopologyObject) {
	const deleteQuery = `
		DELETE FROM topology_cache WHERE id = $1
	`

	_, err := c.db.ExecContext(context.Background(), deleteQuery, obj.Id)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

func (c *client) SetCache(obj *topologyv1.TopologyObject) {
	const upsertQuery = `
		INSERT INTO topology_cache (id, resolver_type_url, data, metadata)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			id = EXCLUDED.id,
			resolver_type_url = EXCLUDED.resolver_type_url,
			data = EXCLUDED.data,
			metadata = EXCLUDED.metadata
	`

	metadataJson, err := json.Marshal(obj.Metadata)
	if err != nil {
		c.log.With(zap.Error(err)).Error("unable to marshal metadata")
		return
	}

	dataJson, err := json.Marshal(obj.Pb.Value)
	if err != nil {
		c.log.With(zap.Error(err)).Error("unable to marshal pb data")
		return
	}

	_, err = c.db.ExecContext(
		context.Background(),
		upsertQuery,
		obj.Id,
		obj.Pb.GetTypeUrl(),
		dataJson,
		metadataJson,
	)
	if err != nil {
		c.log.With(zap.Error(err)).Error("unable to upsert cache item")
		return
	}
}
