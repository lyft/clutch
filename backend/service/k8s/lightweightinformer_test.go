package k8s

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	fakecontroller "k8s.io/client-go/tools/cache/testing"
)

func TestLightweightInformer(t *testing.T) {
	addActionChan := make(chan int)
	updateActionChan := make(chan int)
	deleteActionChan := make(chan int)

	informerHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			addActionChan <- 1
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			updateActionChan <- 1
		},
		DeleteFunc: func(obj interface{}) {
			deleteActionChan <- 1
		},
	}

	// This fake controller is a fake source for listing and watching resources
	fc := fakecontroller.NewFakeControllerSource()

	stop := make(chan struct{})
	defer close(stop)

	podInformer := NewLightweightInformer(
		fc,
		&v1.Pod{},
		time.Minute,
		informerHandlers,
	)

	go podInformer.Run(stop)

	fc.Add(&v1.Pod{})
	assert.Greater(t, <-addActionChan, 0)

	fc.Modify(&v1.Pod{})
	assert.Greater(t, <-updateActionChan, 0)

	fc.Delete(&v1.Pod{})
	assert.Greater(t, <-deleteActionChan, 0)
}
