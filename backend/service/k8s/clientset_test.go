package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestClientsetManager(t *testing.T) {
	fcs := fake.NewSimpleClientset()
	m := managerImpl{clientsets: map[string]*ctxClientsetImpl{
		"core-0":  NewContextClientset("core-namespace", "core-cluster-0", fcs).(*ctxClientsetImpl),
		"core-1a": NewContextClientset("core-namespace", "core-cluster-1", fcs).(*ctxClientsetImpl),
		"core-1b": NewContextClientset("core-namespace", "core-cluster-1", fcs).(*ctxClientsetImpl),
	}}

	// No match found.
	cs, err := m.GetK8sClientset(context.Background(), "foo", "foo", "foo")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, cs)

	// Clientset found but cluster does not match.
	cs, err = m.GetK8sClientset(context.Background(), "core-0", "foo", "foo")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match clientset")
	assert.Nil(t, cs)

	// Clientset match and check namespace.
	cs, err = m.GetK8sClientset(context.Background(), "core-0", "core-cluster-0", "foo")
	assert.NoError(t, err)
	assert.NotNil(t, cs)
	assert.Equal(t, "foo", cs.Namespace())
	assert.Equal(t, "core-cluster-0", cs.Cluster())

	// Cluster match.
	cs, err = m.GetK8sClientset(context.Background(), "undefined", "core-cluster-0", "foo")
	assert.NoError(t, err)
	assert.NotNil(t, cs)

	// Multiple cluster match error.
	cs, err = m.GetK8sClientset(context.Background(), "undefined", "core-cluster-1", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "impossible to determine")
	assert.Nil(t, cs)
}
