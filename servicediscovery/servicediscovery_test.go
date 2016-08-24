package servicediscovery

import (
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/nicholasjackson/router/logging"
	"github.com/nicholasjackson/router/utils"
	"github.com/stretchr/testify/mock"
)

var mockBackend MockBackend
var mockStatsD logging.MockStatsD
var discovery *ServiceDiscovery
var mockTimer = utils.MockTimer{}
var backendError error

func getBackendError() error {
	return backendError
}

func setupSDTests() {
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

	mockStatsD = logging.MockStatsD{}
	mockStatsD.Mock.On("Increment", mock.Anything)

	mockBackend = MockBackend{}
	mockBackend.Mock.On("Services").Return(services, getBackendError)

	discovery = NewServiceDiscovery(&mockBackend, &mockTimer, &mockStatsD)

	backendError = nil
}

func Test_NewServiceDiscovery_returns_valid_instance(t *testing.T) {
	setupSDTests()

	assert.Equal(t, discovery.refreshInterval, (30 * time.Second))
	assert.Equal(t, discovery.backend, &mockBackend)
	assert.Equal(t, discovery.timer, &mockTimer)
	assert.Equal(t, discovery.statsd, &mockStatsD)
}

func Test_NewServiceDiscovery_sets_up_polling_interval(t *testing.T) {
	setupSDTests()
	discovery.StartPolling()

	assert.Equal(t, discovery.refreshInterval, mockTimer.Duration)
}

func Test_NewServiceDiscovery_sets_up_polling_callback(t *testing.T) {
	setupSDTests()
	discovery.StartPolling()

	assert.NotNil(t, mockTimer.Func)
}

func Test_NewServiceDiscovery_polls_for_changes(t *testing.T) {
	setupSDTests()
	discovery.StartPolling()
	mockTimer.Func()

	mockBackend.AssertNumberOfCalls(t, "Services", 1)
}

func Test_NewServiceDiscovery_polls_for_changes_at_30s_interval(t *testing.T) {
	setupSDTests()
	discovery.StartPolling()
	mockTimer.Func()
	mockTimer.Func()

	mockBackend.AssertNumberOfCalls(t, "Services", 2)
}

func Test_NewServiceDiscovery_polling_success_calls_StatsD(t *testing.T) {
	setupSDTests()
	discovery.StartPolling()
	mockTimer.Func()

	mockStatsD.AssertCalled(t, "Increment", pollingBackendSuccess)
}

func Test_NewServiceDiscovery_polling_error_calls_StatsD(t *testing.T) {
	setupSDTests()
	backendError = fmt.Errorf("Unable to reach server")

	discovery.StartPolling()
	mockTimer.Func()

	mockStatsD.AssertCalled(t, "Increment", pollingBackendError)
}

func Test_ServiceDiscovery_returns_services_with_correct_tag(t *testing.T) {
	setupSDTests()
	discovery.StartPolling()
	mockTimer.Func()

	services := discovery.ServicesByTag("TestService1", "project1")

	assert.Equal(t, len(services), 1)
}

func Test_ServiceDiscovery_returns_empty_when_service_not_exist(t *testing.T) {
	setupSDTests()
	discovery.StartPolling()
	mockTimer.Func()

	services := discovery.ServicesByTag("TestService3", "project1")

	assert.Equal(t, len(services), 0)
}

func Test_ServiceDiscovery_returns_empty_when_project_not_exist(t *testing.T) {
	setupSDTests()
	discovery.StartPolling()
	mockTimer.Func()
	services := discovery.ServicesByTag("TestService1", "project3")

	assert.Equal(t, len(services), 0)
}
