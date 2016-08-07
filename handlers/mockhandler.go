package handlers

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockHandler implements a mock instance of a HTTPHandler
type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = m.Mock.Called(w, r)
}
