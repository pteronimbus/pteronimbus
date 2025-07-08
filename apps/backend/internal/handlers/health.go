package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// HealthHandler provides health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Routes sets up the health check routes
func (h *HealthHandler) Routes() http.Handler {
	r := mux.NewRouter()

	// Health check endpoints
	r.HandleFunc("/health", h.Health).Methods("GET")
	r.HandleFunc("/healthz", h.Health).Methods("GET") // Kubernetes style
	r.HandleFunc("/ready", h.Ready).Methods("GET")    // Readiness probe
	r.HandleFunc("/live", h.Live).Methods("GET")      // Liveness probe

	return r
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version,omitempty"`
}

// Health handles general health checks
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Service:   "pteronimbus-backend",
		Version:   "0.1.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Ready handles readiness probe (when service is ready to receive traffic)
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now().UTC(),
		Service:   "pteronimbus-backend",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Live handles liveness probe (when service is alive)
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "alive",
		Timestamp: time.Now().UTC(),
		Service:   "pteronimbus-backend",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
