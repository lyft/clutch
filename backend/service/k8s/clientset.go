package k8s

import (
	"errors"
	"fmt"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	types "k8s.io/client-go/tools/clientcmd/api"
)

type ClientsetManager interface {
	Clientsets() map[string]ContextClientset
	GetK8sClientset(clientset, cluster, namespace string) (ContextClientset, error)
}

type ContextClientset interface {
	k8s.Interface
	Namespace() string
	Cluster() string
}

func NewContextClientset(namespace string, cluster string, clientset k8s.Interface) ContextClientset {
	return &ctxClientsetImpl{
		Interface: clientset,
		namespace: namespace,
		cluster:   cluster,
	}
}

type ctxClientsetImpl struct {
	k8s.Interface
	namespace string
	cluster   string
}

func (c *ctxClientsetImpl) Namespace() string { return c.namespace }
func (c *ctxClientsetImpl) Cluster() string   { return c.cluster }

func newClientsetManager(rules *clientcmd.ClientConfigLoadingRules) (ClientsetManager, error) {
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	apiConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load apiconfig: %w", err)
	}

	lookup := make(map[string]*ctxClientsetImpl, len(apiConfig.Contexts))
	for name, ctxInfo := range apiConfig.Contexts {
		clientset, err := createClientsetImpl(name, ctxInfo, rules)
		if err != nil {
			return nil, err
		}
		lookup[name] = clientset
	}

	if len(lookup) == 0 {
		restConfig, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		clientset, err := k8s.NewForConfig(restConfig)
		if err != nil {
			return nil, fmt.Errorf("could not create k8s clientset from config: %w", err)
		}

		lookup["local"] = &ctxClientsetImpl{Interface: clientset, namespace: "default", cluster: "local"}
	}

	return &managerImpl{clientsets: lookup}, nil
}

func createClientsetImpl(name string, ctxInfo *types.Context, rules *clientcmd.ClientConfigLoadingRules) (*ctxClientsetImpl, error) {
	contextConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		rules,
		&clientcmd.ConfigOverrides{CurrentContext: name},
	)

	restConfig, err := contextConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load restconfig: %w", err)
	}

	clientset, err := k8s.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create k8s clientset from config: %w", err)
	}

	ns, _, err := contextConfig.Namespace()
	if err != nil {
		return nil, err
	}
	return &ctxClientsetImpl{Interface: clientset, namespace: ns, cluster: ctxInfo.Cluster}, nil
}

type managerImpl struct {
	clientsets map[string]*ctxClientsetImpl
}

func (m *managerImpl) Clientsets() map[string]ContextClientset {
	ret := make(map[string]ContextClientset)
	for k, v := range m.clientsets {
		ret[k] = v
	}
	return ret
}

func (m *managerImpl) GetK8sClientset(clientset, cluster, namespace string) (ContextClientset, error) {
	cs, ok := m.clientsets[clientset]
	if !ok {
		return nil, errors.New("not found")
	}

	if cluster != "" && cluster != cs.cluster {
		return nil, errors.New("specified cluster does not match clientset")
	}

	if namespace == "" {
		// Use the clients' default namespace.
		return cs, nil
	}

	// Shallow copy and update namespace.
	ret := *cs
	ret.namespace = namespace
	return &ret, nil
}
