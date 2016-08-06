package servicediscovery

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

	//ServicesByTag returns an array of Service for the given service name and tag
	ServicesByTag(service string, tag string) []*Service
}
