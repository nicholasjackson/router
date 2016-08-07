package logging

import "github.com/stretchr/testify/mock"

// MockStatsD is a mock statsd client to be used for testing
type MockStatsD struct {
	mock.Mock
}

// Increment increments a StatsD bucket by 1
func (m *MockStatsD) Increment(label string) {
	_ = m.Mock.Called(label)
}
