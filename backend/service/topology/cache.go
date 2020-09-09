package topology

import (
	"context"
	"log"

	"github.com/lyft/clutch/backend/service"
	k8sservice "github.com/lyft/clutch/backend/service/k8s"
)

func (c *client) startCaching() {
	if _, ok := service.Registry["clutch.service.k8s"]; ok {
		c.populateCacheFromKubernetes()
	}
}

// pretend we are the leader
func (c *client) populateCacheFromKubernetes() {
	k8sClient, _ := service.Registry["clutch.service.k8s"]
	k8sSvc, _ := k8sClient.(k8sservice.Service)
	log.Print("topology is enabled and starting k8s cache.")

	stop := make(chan struct{})

	for name, cs := range k8sSvc.GetClientSets() {
		log.Printf("starting informer for cluster: %s", name)
		c.startInformers(cs, stop)
	}
}

func (c *client) DeleteCache(id string) {
	const deleteQuery = `
		DELETE FROM topology_cache WHERE id = $1
	`
	_, err := c.db.ExecContext(context.Background(), deleteQuery, id)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

func (c *client) SetCache(id string, resolver_type_url string, data []byte) {
	const upsertQuery = `
		INSERT INTO topology_cache (id, data, resolver_type_url)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET
			id = EXCLUDED.id,
			data = EXCLUDED.data,
			resolver_type_url = EXCLUDED.resolver_type_url
	`

	_, err := c.db.ExecContext(context.Background(), upsertQuery, id, data, resolver_type_url)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}
