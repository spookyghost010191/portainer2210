package proxy

import (
	"fmt"
	"net/http"

	"github.com/portainer/portainer/api/dataservices"
	dockerclient "github.com/portainer/portainer/api/docker/client"
	"github.com/portainer/portainer/api/http/proxy/factory/kubernetes"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/portainer/portainer/api/kubernetes/cli"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/http/proxy/factory"
)

type (
	// Manager represents a service used to manage proxies to environments (endpoints) and extensions.
	Manager struct {
		proxyFactory     *factory.ProxyFactory
		endpointProxies  cmap.ConcurrentMap
		k8sClientFactory *cli.ClientFactory
	}
)

// NewManager initializes a new proxy Service
func NewManager(kubernetesClientFactory *cli.ClientFactory) *Manager {
	return &Manager{
		endpointProxies:  cmap.New(),
		k8sClientFactory: kubernetesClientFactory,
	}
}

func (manager *Manager) NewProxyFactory(dataStore dataservices.DataStore, signatureService portainer.DigitalSignatureService, tunnelService portainer.ReverseTunnelService, clientFactory *dockerclient.ClientFactory, kubernetesClientFactory *cli.ClientFactory, kubernetesTokenCacheManager *kubernetes.TokenCacheManager, gitService portainer.GitService, snapshotService portainer.SnapshotService) {
	manager.proxyFactory = factory.NewProxyFactory(dataStore, signatureService, tunnelService, clientFactory, kubernetesClientFactory, kubernetesTokenCacheManager, gitService, snapshotService)
}

// CreateAndRegisterEndpointProxy creates a new HTTP reverse proxy based on environment(endpoint) properties and and adds it to the registered proxies.
// It can also be used to create a new HTTP reverse proxy and replace an already registered proxy.
func (manager *Manager) CreateAndRegisterEndpointProxy(endpoint *portainer.Endpoint) (http.Handler, error) {
	if manager.proxyFactory == nil {
		return nil, fmt.Errorf("proxy factory not init")
	}

	proxy, err := manager.proxyFactory.NewEndpointProxy(endpoint)
	if err != nil {
		return nil, err
	}

	manager.endpointProxies.Set(fmt.Sprint(endpoint.ID), proxy)
	return proxy, nil
}

// CreateAgentProxyServer creates a new HTTP reverse proxy based on environment(endpoint) properties and and adds it to the registered proxies.
// It can also be used to create a new HTTP reverse proxy and replace an already registered proxy.
func (manager *Manager) CreateAgentProxyServer(endpoint *portainer.Endpoint) (*factory.ProxyServer, error) {
	if manager.proxyFactory == nil {
		return nil, fmt.Errorf("proxy factory not init")
	}
	return manager.proxyFactory.NewAgentProxy(endpoint)
}

// GetEndpointProxy returns the proxy associated to a key
func (manager *Manager) GetEndpointProxy(endpoint *portainer.Endpoint) http.Handler {
	proxy, ok := manager.endpointProxies.Get(fmt.Sprint(endpoint.ID))
	if !ok {
		return nil
	}

	return proxy.(http.Handler)
}

// DeleteEndpointProxy deletes the proxy associated to a key
// and cleans the k8s environment(endpoint) client cache. DeleteEndpointProxy
// is currently only called for edge connection clean up and when endpoint is updated
func (manager *Manager) DeleteEndpointProxy(endpointID portainer.EndpointID) {
	manager.endpointProxies.Remove(fmt.Sprint(endpointID))

	if manager.k8sClientFactory != nil {
		manager.k8sClientFactory.RemoveKubeClient(endpointID)
	}
}

// CreateGitlabProxy creates a new HTTP reverse proxy that can be used to send requests to the Gitlab API
func (manager *Manager) CreateGitlabProxy(url string) (http.Handler, error) {
	if manager.proxyFactory == nil {
		return nil, fmt.Errorf("proxy factory not init")
	}
	return manager.proxyFactory.NewGitlabProxy(url)
}
