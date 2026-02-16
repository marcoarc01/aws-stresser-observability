package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthHandler responde com o status da aplicação
// GET /health → { "status": "ok" }
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}