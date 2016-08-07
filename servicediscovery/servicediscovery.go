package servicediscovery

import "time"

// Service represents a backend service which is registered with the dynamci service registry
type Service struct {
	Node           string
	Address        string
	ServiceID      string
	ServiceName    string
	ServiceAddress string
	ServiceTags    []string
	ServicePort    int
}

// Backend represents a dynamic service registry like Consul or etcd
type Backend interface {
	//Services returns a map of Service from the Backend
	// The key to the map is the name of the service
	// Services returns an array as there may be more than one endpoint registered for a service name
	Services() map[string][]*Service
}

// Discoverable defines an interface for ServiceDiscovery
type Discoverable interface {
	//ServicesByTag returns an array of Service for the given service name and tag
	ServicesByTag(name string, tag string) []*Service
}

// ServiceDiscovery implements the Discoverable interface and is the main class for finding details for backends
type ServiceDiscovery struct {
	backend         Backend
	refreshInterval time.Duration
}

// NewServiceDiscovery creates and returns a ServiceDiscovery type with the given Backend
func NewServiceDiscovery(backend Backend) *ServiceDiscovery {
	return &ServiceDiscovery{backend: backend, refreshInterval: 30 * time.Second}
}

//ServicesByTag returns an array of Service for the given service name and tag
func (s *ServiceDiscovery) ServicesByTag(name string, tag string) []*Service {
	services := []*Service{}
	backendServices := s.backend.Services()

	for _, s := range backendServices[name] {
		if containsTag(s, tag) {
			services = append(services, s)
		}
	}

	return services
}

func containsTag(service *Service, tag string) bool {
	for _, t := range service.ServiceTags {
		if t == tag {
			return true
		}
	}

	return false
}
