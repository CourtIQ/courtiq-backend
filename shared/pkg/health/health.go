package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Status represents the health status of a service
type Status string

const (
	// StatusUp indicates the service is healthy
	StatusUp Status = "UP"
	
	// StatusDown indicates the service is unhealthy
	StatusDown Status = "DOWN"
	
	// StatusDegraded indicates the service is partially healthy
	StatusDegraded Status = "DEGRADED"
)

// HealthCheck defines a function that performs a health check
type HealthCheck func(ctx context.Context) (Status, error)

// Component represents a component of the system
type Component struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
}

// HealthResponse represents the response of a health check
type HealthResponse struct {
	Status     Status      `json:"status"`
	Components []Component `json:"components"`
	Timestamp  time.Time   `json:"timestamp"`
}

// Handler is a health check handler
type Handler struct {
	mu        sync.RWMutex
	checks    map[string]HealthCheck
	timeout   time.Duration
	readiness bool
}

// NewHandler creates a new health check handler
func NewHandler(timeout time.Duration) *Handler {
	return &Handler{
		checks:    make(map[string]HealthCheck),
		timeout:   timeout,
		readiness: false,
	}
}

// AddCheck adds a health check
func (h *Handler) AddCheck(name string, check HealthCheck) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = check
}

// AddMongoDBCheck adds a MongoDB health check
func (h *Handler) AddMongoDBCheck(name string, client *mongo.Client) {
	h.AddCheck(name, func(ctx context.Context) (Status, error) {
		err := client.Ping(ctx, nil)
		if err != nil {
			return StatusDown, err
		}
		return StatusUp, nil
	})
}

// SetReady sets the readiness status
func (h *Handler) SetReady(ready bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.readiness = ready
}

// HandleLiveness is an HTTP handler for liveness checks
func (h *Handler) HandleLiveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Liveness checks are simple - just return 200 OK
	resp := HealthResponse{
		Status:    StatusUp,
		Timestamp: time.Now(),
	}
	
	json.NewEncoder(w).Encode(resp)
}

// HandleReadiness is an HTTP handler for readiness checks
func (h *Handler) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	readiness := h.readiness
	checks := make(map[string]HealthCheck, len(h.checks))
	for k, v := range h.checks {
		checks[k] = v
	}
	h.mu.RUnlock()
	
	// If not ready yet, return 503 Service Unavailable
	if !readiness {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(HealthResponse{
			Status:    StatusDown,
			Timestamp: time.Now(),
		})
		return
	}
	
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()
	
	// Run all health checks
	components := make([]Component, 0, len(checks))
	overallStatus := StatusUp
	
	for name, check := range checks {
		status, err := check(ctx)
		
		component := Component{
			Name:   name,
			Status: status,
		}
		
		if err != nil {
			component.Error = err.Error()
		}
		
		components = append(components, component)
		
		// If any component is down, overall status is down
		if status == StatusDown {
			overallStatus = StatusDown
		} else if status == StatusDegraded && overallStatus != StatusDown {
			// If any component is degraded and none are down, overall status is degraded
			overallStatus = StatusDegraded
		}
	}
	
	resp := HealthResponse{
		Status:     overallStatus,
		Components: components,
		Timestamp:  time.Now(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	if overallStatus != StatusUp {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	
	json.NewEncoder(w).Encode(resp)
}