package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func testNodeClientset() *fake.Clientset {
	nodes := &corev1.NodeList{
		Items: []corev1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-node",
				},
				Spec: corev1.NodeSpec{
					Unschedulable: false,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "already-cordoned-node",
				},
				Spec: corev1.NodeSpec{
					Unschedulable: true,
				},
			},
		},
	}

	fakeClient := fake.NewSimpleClientset([]runtime.Object{nodes}...)
	return fakeClient
}

func testListNodesClientSet() *fake.Clientset {
	nodes := &corev1.NodeList{
		Items: []corev1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-node",
				},
				Spec: corev1.NodeSpec{
					Unschedulable: false,
				},
			},
		},
	}

	fakeClient := fake.NewSimpleClientset([]runtime.Object{nodes}...)
	fakeClient.AddReactor("list", "nodes",
		func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, nodes, nil
		})
	return fakeClient
}

func TestDescribeNode(t *testing.T) {
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{
				"test": {
					Interface: testListNodesClientSet(),
					namespace: metav1.NamespaceAll,
					cluster:   "boba9000",
				},
				"test-with-duplicates": {
					Interface: testNodeClientset(),
					namespace: metav1.NamespaceAll,
					cluster:   "boba9000",
				},
			},
		},
	}

	// Found 1 node
	node, err := s.DescribeNode(context.Background(), "test", "boba9000", "test-node")
	assert.NoError(t, err)
	assert.Equal(t, "test-node", node.Name)
	assert.False(t, node.Unschedulable)

	// Found >1 node
	node, err = s.DescribeNode(context.Background(), "test-with-duplicates", "boba9000", "test-node")
	assert.Error(t, err)
	assert.Nil(t, node)

	// Node not found
	node, err = s.DescribeNode(context.Background(), "", "boba9000", "")
	assert.Error(t, err)
	assert.Nil(t, node)
}

func TestCordonNode(t *testing.T) {
	cs := testNodeClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"test": {
				Interface: cs,
				namespace: metav1.NamespaceAll,
				cluster:   "boba9000",
			}},
		},
	}

	err := s.UpdateNode(context.Background(), "test", "boba9000", "test-node", true)
	assert.NoError(t, err)

	node, _ := cs.CoreV1().Nodes().Get(context.Background(), "test-node", metav1.GetOptions{})
	assert.Equal(t, true, node.Spec.Unschedulable)

	err = s.UpdateNode(context.Background(), "test", "boba9000", "already-cordoned-node", true)
	assert.NoError(t, err)

	node, _ = cs.CoreV1().Nodes().Get(context.Background(), "already-cordoned-node", metav1.GetOptions{})
	assert.Equal(t, true, node.Spec.Unschedulable)
}
