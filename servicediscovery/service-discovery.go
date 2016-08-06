package servicediscovery

type Service struct {
	Node           string
	Address        string
	ServiceID      string
	ServiceName    string
	ServiceAddress string
	ServiceTags    []string
	ServicePort    int
}

type Backend interface {
	Services() []*Service
}
