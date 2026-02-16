package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/marcoarc01/aws-stresser-observability/stresser-app/stress"
)

// StateResponse é o JSON retornado pelo GET /api/state
type StateResponse struct {
	StressLevel int `json:"stress_level"`
	CPUWorkers  int `json:"cpu_workers"`
}

// StateHandler retorna o estado atual do stresser
// GET /api/state → { "stress_level": 50, "cpu_workers": 4 }
func StateHandler(engine *stress.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(StateResponse{
			StressLevel: engine.GetLevel(),
			CPUWorkers:  engine.GetWorkers(),
		})
	}
}