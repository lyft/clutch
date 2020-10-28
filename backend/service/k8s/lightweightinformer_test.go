package k8s

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
	fakecontroller "k8s.io/client-go/tools/cache/testing"
)

func TestLightweightInformer(t *testing.T) {
	cs := fake.NewSimpleClientset()
	csm := &managerImpl{
		clientsets: map[string]*ctxClientsetImpl{"foo": &ctxClientsetImpl{
			Interface: cs,
			namespace: "default",
			cluster:   "core-testing",
		}},
	}

	numAddActions := 0
	numUpdateActions := 0
	numDeleteActions := 0

	informerHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			numAddActions++
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			numUpdateActions++
		},
		DeleteFunc: func(obj interface{}) {
			numDeleteActions++
		},
	}

	// This fake controller is a fake source for listing and watching resources
	fc := fakecontroller.NewFakeControllerSource()

	stop := make(chan struct{})
	podInformer := NewLightweightInformer(
		csm.clientsets["foo"],
		fc,
		&v1.Pod{},
		time.Minute,
		informerHandlers,
	)

	go podInformer.Run(stop)

	// We sleep below to give some time for the informer to get the event and hit a handler
	fc.Add(&v1.Pod{})
	time.Sleep(time.Millisecond * 10)
	assert.Greater(t, numAddActions, 0)

	fc.Modify(&v1.Pod{})
	time.Sleep(time.Millisecond * 10)
	assert.Greater(t, numUpdateActions, 0)

	fc.Delete(&v1.Pod{})
	time.Sleep(time.Millisecond * 10)
	assert.Greater(t, numDeleteActions, 0)
}
