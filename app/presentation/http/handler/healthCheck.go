package handler

import (
	"net/http"
)

type HealthCheck struct {
	message string
}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

func (h *HealthCheck) WithMessage(message string) {
	h.message = message
}

func (h *HealthCheck) CheckServerHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(h.message))
}
