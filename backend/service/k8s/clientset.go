package k8s

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	k8s "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	k8sconfigv1 "github.com/lyft/clutch/backend/api/config/service/k8s/v1"
)

const (
	inCluster = "in-cluster"
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

func newClientsetManager(rules *clientcmd.ClientConfigLoadingRules, restClientConfig *k8sconfigv1.RestClientConfig, logger *zap.Logger) (ClientsetManager, error) {
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	apiConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load apiconfig: %w", err)
	}

	lookup := make(map[string]*ctxClientsetImpl, len(apiConfig.Contexts))
	for name, ctxInfo := range apiConfig.Contexts {
		contextConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			rules,
			&clientcmd.ConfigOverrides{CurrentContext: name},
		)

		restConfig, err := contextConfig.ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("could not load restconfig: %w", err)
		}

		if err := ApplyRestClientConfig(restConfig, restClientConfig); err != nil {
			return nil, err
		}

		clientset, err := k8s.NewForConfig(restConfig)
		if err != nil {
			return nil, fmt.Errorf("could not create k8s clientset from config: %w", err)
		}

		ns, _, err := contextConfig.Namespace()
		if err != nil {
			return nil, err
		}
		lookup[name] = &ctxClientsetImpl{Interface: clientset, namespace: ns, cluster: ctxInfo.Cluster}
	}

	// If there is no configured cluster produced fallback to InClusterConfig
	if len(lookup) == 0 {
		logger.Info("no kubeconfig was found, falling back to InClusterConfig")

		restConfig, err := restclient.InClusterConfig()
		if err := ApplyRestClientConfig(restConfig, restClientConfig); err != nil {
			return nil, err
		}

		switch err {
		case restclient.ErrNotInCluster:
			logger.Warn("not in a kubernetes cluster, unable to configure kube clientset")
		case nil:
			clientset, err := k8s.NewForConfig(restConfig)
			if err != nil {
				return nil, fmt.Errorf("could not create k8s InClusterConfig: %w", err)
			}
			lookup[inCluster] = &ctxClientsetImpl{Interface: clientset, namespace: "default", cluster: inCluster}
		default:
			return nil, err
		}
	}

	return &managerImpl{clientsets: lookup}, nil
}

func ApplyRestClientConfig(restConfig *restclient.Config, restClientConfig *k8sconfigv1.RestClientConfig) error {
	if restClientConfig == nil {
		return nil
	}

	if restClientConfig.Burst != 0 {
		restConfig.Burst = int(restClientConfig.Burst)
	}

	if restClientConfig.Qps >= 0 {
		restConfig.QPS = restClientConfig.Qps
	}

	if restClientConfig.Timeout != nil {
		timeout, err := ptypes.Duration(restClientConfig.Timeout)
		if err != nil {
			return err
		}
		restConfig.Timeout = timeout
	}
	return nil
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
