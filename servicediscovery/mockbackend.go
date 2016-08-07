package servicediscovery

import (
	"github.com/stretchr/testify/mock"
)

// MockBackend is a mock implementation of a servicediscovery.Backend
// used for testing
type MockBackend struct {
	mock.Mock
}

// Services returns service from the backend
func (m *MockBackend) Services() map[string][]*Service {
	args := m.Mock.Called()
	return args[0].(map[string][]*Service)
}
