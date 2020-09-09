// package cache

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"log"
// 	"time"

// 	corev1 "k8s.io/api/core/v1"
// 	"k8s.io/client-go/informers"
// 	"k8s.io/client-go/tools/cache"

// 	topologyv1 "github.com/lyft/clutch/backend/api/config/service/topology/v1"
// 	k8sservice "github.com/lyft/clutch/backend/service/k8s"
// )

// type TopologyCache struct {
// 	cfg    *topologyv1.Cache
// 	db     *sql.DB
// 	k8sSvc k8sservice.Service
// }

// func New(cfg *topologyv1.Cache, db *sql.DB, svc k8sservice.Service) (*TopologyCache, error) {
// 	tc := &TopologyCache{
// 		cfg:    cfg,
// 		db:     db,
// 		k8sSvc: svc,
// 	}

// 	tc.PopulateCacheFromKubernetes()

// 	return tc, nil
// }

// // pretend we are the leader
// func (t *TopologyCache) PopulateCacheFromKubernetes() {
// 	log.Print("topology is enabled and starting k8s cache.")

// 	log.Print("populate cache")
// 	stop := make(chan struct{})

// 	for name, cs := range t.k8sSvc.GetClientSets() {
// 		log.Printf("starting informer for cluster: %s", name)
// 		t.startInformers(cs, stop)
// 	}
// }

// func (t *TopologyCache) startInformers(cs k8sservice.ContextClientset, stop chan struct{}) {
// 	factory := informers.NewSharedInformerFactoryWithOptions(cs, time.Minute*1)

// 	podInformer := factory.Core().V1().Pods().Informer()
// 	deploymentInformer := factory.Apps().V1().Deployments().Informer()

// 	informerHandlers := cache.ResourceEventHandlerFuncs{
// 		AddFunc:    t.informerAddHandler,
// 		UpdateFunc: t.informerUpdateHandler,
// 		DeleteFunc: t.informerDeleteHandler,
// 	}

// 	podInformer.AddEventHandler(informerHandlers)
// 	deploymentInformer.AddEventHandler(informerHandlers)

// 	go func() {
// 		podInformer.Run(stop)
// 	}()

// 	// go func() {
// 	// 	deploymentInformer.Run(stop)
// 	// }()
// }

// // switch to select type?
// func (t *TopologyCache) informerAddHandler(obj interface{}) {
// 	log.Print("Add Handler")
// 	// log.Printf("%v", obj.(runtime.Object))

// 	k8sObj := obj.(*corev1.Pod)
// 	// todo: make a switch for this choiceniss
// 	// k8sObj.GetObjectKind()
// 	b, _ := json.Marshal(k8sObj)

// 	t.upsertCache(k8sObj.Name, b, "pod")
// }

// func (t *TopologyCache) informerUpdateHandler(oldObj, newObj interface{}) {
// 	log.Print("Update Handler")

// 	k8sObj := newObj.(*corev1.Pod)
// 	b, _ := json.Marshal(k8sObj)
// 	t.upsertCache(k8sObj.Name, b, "pod")
// }

// func (t *TopologyCache) informerDeleteHandler(obj interface{}) {
// 	log.Print("Delete handler")
// 	// log.Printf("%v", obj.(runtime.Object))
// 	k8sObj := obj.(*corev1.Pod)

// 	t.deleteCache(k8sObj.Name)
// }

// func (t *TopologyCache) deleteCache(id string) {
// 	const deleteQuery = `
// 		DELETE FROM topology_cache WHERE id = $1
// 	`
// 	_, err := t.db.ExecContext(context.Background(), deleteQuery, id)
// 	if err != nil {
// 		log.Printf("%v", err)
// 		return
// 	}
// }

// func (t *TopologyCache) upsertCache(id string, data []byte, resolver_type_url string) {
// 	const upsertQuery = `
// 		INSERT INTO topology_cache (id, data, resolver_type_url)
// 		VALUES ($1, $2, $3)
// 		ON CONFLICT (id) DO UPDATE SET
// 			id = EXCLUDED.id,
// 			data = EXCLUDED.data,
// 			resolver_type_url = EXCLUDED.resolver_type_url
// 	`

// 	_, err := t.db.ExecContext(context.Background(), upsertQuery, id, data, resolver_type_url)
// 	if err != nil {
// 		log.Printf("%v", err)
// 		return
// 	}
// }
