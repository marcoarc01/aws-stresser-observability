package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/marcoarc01/aws-stresser-observability/stresser-app/stress"
)

// StressRequest é o JSON esperado no POST /api/stress
type StressRequest struct {
	Level int `json:"level"`
}

// StressResponse é o JSON retornado após alterar o stress
type StressResponse struct {
	Message     string `json:"message"`
	StressLevel int    `json:"stress_level"`
	CPUWorkers  int    `json:"cpu_workers"`
}

// StressHandler altera o nível de stress
// POST /api/stress → body: { "level": 0..100 }
func StressHandler(engine *stress.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req StressRequest

		// Decodifica o JSON do body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "JSON inválido"}`, http.StatusBadRequest)
			return
		}

		// Valida o range 0-100
		if req.Level < 0 || req.Level > 100 {
			http.Error(w, `{"error": "level deve ser entre 0 e 100"}`, http.StatusBadRequest)
			return
		}

		// Aplica o novo nível de stress
		engine.SetLevel(req.Level)

		log.Printf(" Stress alterado para %d%% (%d workers)", req.Level, engine.GetWorkers())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(StressResponse{
			Message:     "Stress level atualizado",
			StressLevel: engine.GetLevel(),
			CPUWorkers:  engine.GetWorkers(),
		})
	}
}