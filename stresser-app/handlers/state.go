package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/marcoarc01/aws-stresser-observability/stresser-app/stress"
	"go.opentelemetry.io/otel/attribute"
)

// StateResponse é o JSON retornado pelo GET /api/state
type StateResponse struct {
	StressLevel int `json:"stress_level"`
	CPUWorkers  int `json:"cpu_workers"`
}

// StateHandler retorna o estado atual do stresser
func StateHandler(engine *stress.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, span := tracer.Start(ctx, "GET /api/state")
		defer span.End()

		level := engine.GetLevel()
		workers := engine.GetWorkers()

		span.SetAttributes(
			attribute.Int("stress.level", level),
			attribute.Int("stress.workers", workers),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(StateResponse{
			StressLevel: level,
			CPUWorkers:  workers,
		})
	}
}
