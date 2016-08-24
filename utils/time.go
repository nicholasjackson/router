package utils

import "time"

// Timer is an interface which returns the current time
type Timer interface {
	CurrentTime() time.Time
	AfterFunc(d time.Duration, f func()) *time.Timer
}

// RealTime implements Time and returns the current clock time
type RealTime struct{}

// CurrentTime returns the current clock time
func (t *RealTime) CurrentTime() time.Time {
	return time.Now()
}

// AfterFunc calls a function after the given duration
func (t *RealTime) AfterFunc(d time.Duration, f func()) *time.Timer {
	return time.AfterFunc(d, f)
}
