package servicediscovery

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

var server *httptest.Server
var address string

func setupServer() {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/v1/catalog/services" {
			fmt.Fprintln(w, "{\"syslog-514\":[\"dev\"]}")
		} else {
			fmt.Fprintln(w, "[{\"Node\":\"6fc799543d4f\",\"Address\":\"172.20.0.3\",\"ServiceID\":\"4fca1506802a:minkeufdwfxnb6wyvkorb_syslog_1:514\",\"ServiceName\":\"syslog-514\",\"ServiceTags\":[\"project1\",\"dev\"],\"ServiceAddress\":\"172.20.0.5\",\"ServicePort\":514},{\"Node\":\"6fc799543d4f\",\"Address\":\"172.20.0.3\",\"ServiceID\":\"4fca1506802a:minkeufdwfxnb6wyvkorb_syslog_1:514\",\"ServiceName\":\"syslog-514\",\"ServiceTags\":[\"project2\",\"dev\"],\"ServiceAddress\":\"172.20.0.5\",\"ServicePort\":514}]")
		}
	}))

	address = strings.Replace(server.URL, "http://", "", -1)
}

func createConsul(address string) *Consul {
	config := api.DefaultConfig()
	config.Address = address

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	return NewConsul(client)
}

func Test_NewConsul_returns_a_new_consul_backend(t *testing.T) {
	consul := createConsul("127.0.0.1")
	assert.NotNil(t, consul)
}

func Test_Services_returns_2_services_from_the_api(t *testing.T) {
	setupServer()
	defer server.Close()

	services := createConsul(address).Services()
	assert.Equal(t, 2, len(services["syslog-514"]))
}

func Test_Services_decodes_service_Address(t *testing.T) {
	setupServer()
	defer server.Close()

	services := createConsul(address).Services()
	assert.Equal(t, "172.20.0.3", services["syslog-514"][0].Address)
}

func Test_Services_decodes_service_ServiceAddress(t *testing.T) {
	setupServer()
	defer server.Close()

	services := createConsul(address).Services()
	assert.Equal(t, "172.20.0.5", services["syslog-514"][0].ServiceAddress)
}

func Test_Services_decodes_service_ServicePort(t *testing.T) {
	setupServer()
	defer server.Close()

	services := createConsul(address).Services()
	assert.Equal(t, 514, services["syslog-514"][0].ServicePort)
}

func Test_Services_decodes_service_ServiceTags(t *testing.T) {
	setupServer()
	defer server.Close()

	services := createConsul(address).Services()
	assert.Equal(t, []string{"project1", "dev"}, services["syslog-514"][0].ServiceTags)
}
