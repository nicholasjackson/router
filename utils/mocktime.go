package utils

import "time"

// MockTimer implements CurrentTime and is used for mocking time instances
type MockTimer struct {
	Time time.Time
}

// CurrentTime returns the current value stored in the mock
func (m *MockTimer) CurrentTime() time.Time {
	return m.Time
}
