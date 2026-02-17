package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/marcoarc01/aws-stresser-observability/stresser-app/metrics"
)

// HealthHandler responde com o status da aplicação
// GET /health → { "status": "ok" }
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	metrics.HTTPRequestsTotal.WithLabelValues("GET", "/health", "200").Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}