package utils

import "time"

// MockTimer implements CurrentTime and is used for mocking time instances
type MockTimer struct {
	Time     time.Time
	Duration time.Duration
	Func     func()
}

// CurrentTime returns the current value stored in the mock
func (m *MockTimer) CurrentTime() time.Time {
	return m.Time
}

// AddInterval adds a duration to the current mock time
func (m *MockTimer) AddInterval(d time.Duration) {
	m.Time = m.Time.Add(d)
}

// AfterFunc calls a function after the given duration
func (t *MockTimer) AfterFunc(d time.Duration, f func()) *time.Timer {
	t.Func = f
	t.Duration = d
	return &time.Timer{}
}
