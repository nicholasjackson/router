package servicediscovery

import (
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert"
)

var mockBackend = MockBackend{}

func TestMain(m *testing.M) {
	services := map[string][]*Service{}
	services["TestService1"] = []*Service{
		&Service{
			Address:     "127.0.0.1",
			ServiceName: "TestService1",
			ServiceTags: []string{"project1"},
		},
		&Service{
			Address: "127.0.0.1",

			ServiceName: "TestService1",
			ServiceTags: []string{"project2"},
		},
	}
	services["TestService2"] = []*Service{
		&Service{
			Address:     "127.0.0.1",
			ServiceName: "TestService2",
			ServiceTags: []string{"project1"},
		},
	}

	mockBackend.Mock.On("Services").Return(services)

	os.Exit(m.Run())
}

func Test_NewServiceDiscovery_returns_valid_instance(t *testing.T) {
	discovery := NewServiceDiscovery(&mockBackend)

	assert.Equal(t, discovery.refreshInterval, (30 * time.Second))
	assert.Equal(t, discovery.backend, &mockBackend)
}

func Test_ServiceDiscovery_returns_services_with_correct_tag(t *testing.T) {
	discovery := NewServiceDiscovery(&mockBackend)
	services := discovery.ServicesByTag("TestService1", "project1")

	assert.Equal(t, len(services), 1)
}

func Test_ServiceDiscovery_returns_empty_when_service_not_exist(t *testing.T) {
	discovery := NewServiceDiscovery(&mockBackend)
	services := discovery.ServicesByTag("TestService3", "project1")

	assert.Equal(t, len(services), 0)
}

func Test_ServiceDiscovery_returns_empty_when_project_not_exist(t *testing.T) {
	discovery := NewServiceDiscovery(&mockBackend)
	services := discovery.ServicesByTag("TestService1", "project3")

	assert.Equal(t, len(services), 0)
}
