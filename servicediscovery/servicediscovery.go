package servicediscovery

import (
	"time"

	"github.com/nicholasjackson/router/logging"
	"github.com/nicholasjackson/router/utils"
)

const (
	pollingBackendSuccess string = "servicerouter.Backend.Polling.Success"
	pollingBackendError   string = "servicerouter.Backend.Polling.Error"
)

// Service represents a backend service which is registered with the dynamic service registry
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
	Services() (map[string][]*Service, error)
}

// Discoverable defines an interface for ServiceDiscovery
type Discoverable interface {
	//ServicesByTag returns an array of Service for the given service name and tag
	ServicesByTag(name string, tag string) []*Service
	// StartPolling starts polling the backend for changes
	StartPolling()
	//StopPolling stops polling the backend for changes
	StopPolling()
}

// ServiceDiscovery implements the Discoverable interface and is the main class for finding details for backends
type ServiceDiscovery struct {
	backend         Backend
	refreshInterval time.Duration
	timer           utils.Timer
	statsd          logging.StatsD
	backendServices map[string][]*Service
}

// NewServiceDiscovery creates and returns a ServiceDiscovery type with the given Backend
func NewServiceDiscovery(backend Backend, timer utils.Timer, statsd logging.StatsD) *ServiceDiscovery {
	return &ServiceDiscovery{
		backend:         backend,
		timer:           timer,
		statsd:          statsd,
		refreshInterval: 30 * time.Second}
}

// StartPolling starts to poll the backend for changes
func (s *ServiceDiscovery) StartPolling() {
	go s.internalPolling()
	time.Sleep(100 * time.Millisecond)
}

func (s *ServiceDiscovery) internalPolling() {
	s.timer.AfterFunc(s.refreshInterval, func() {
		var err error
		s.backendServices, err = s.backend.Services()

		if err != nil {
			s.statsd.Increment(pollingBackendError)
		} else {
			s.statsd.Increment(pollingBackendSuccess)
		}
	})
}

//ServicesByTag returns an array of Service for the given service name and tag
func (s *ServiceDiscovery) ServicesByTag(name string, tag string) []*Service {
	services := []*Service{}

	for _, s := range s.backendServices[name] {
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
