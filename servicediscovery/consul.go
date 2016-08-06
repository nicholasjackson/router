package servicediscovery

import "github.com/hashicorp/consul/api"

// Consul represents a backened which can connect to a Consul server
type Consul struct {
	client  *api.Client
	catalog *api.Catalog
	options *api.QueryOptions
}

// NewConsul returns a new instance of the Consul backend
func NewConsul(client *api.Client) *Consul {
	return &Consul{
		client:  client,
		catalog: client.Catalog(),
	}
}

//{"syslog-514":["dev"]}
//[{"Node":"6fc799543d4f","Address":"172.20.0.3","ServiceID":"4fca1506802a:minkeufdwfxnb6wyvkorb_syslog_1:514","ServiceName":"syslog-514","ServiceTags":["project2","dev"],"ServiceAddress":"172.20.0.5","ServicePort":514}]
func (c *Consul) Services() map[string][]*Service {
	options := api.QueryOptions{}
	s, _, _ := c.catalog.Services(&options)

	services := map[string][]*Service{}

	for key := range s {
		services[key] = c.getService(key)
	}

	return services
}

func (c *Consul) getService(service string) []*Service {
	services, _, _ := c.catalog.Service(service, "", nil)

	returns := []*Service{}

	for _, service := range services {
		myservice := Service{
			Address:        service.Address,
			ServiceAddress: service.ServiceAddress,
			ServicePort:    service.ServicePort,
			ServiceTags:    service.ServiceTags,
		}

		returns = append(returns, &myservice)
	}

	return returns
}
