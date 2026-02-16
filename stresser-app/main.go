package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/marcoarc01/aws-stresser-observability/stresser-app/handlers"
	"github.com/marcoarc01/aws-stresser-observability/stresser-app/stress"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Inicializa o motor de stress
	engine := stress.NewEngine()

	// Configura as rotas
	mux := http.NewServeMux()

	// Rotas da API
	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("GET /api/state", handlers.StateHandler(engine))
	mux.HandleFunc("POST /api/stress", handlers.StressHandler(engine))
	mux.Handle("GET /metrics", promhttp.Handler())

	// UI estática
	mux.HandleFunc("GET /", handlers.UIHandler)

	// Porta do servidor
	port := os.Getenv("STRESSER_PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Stresser App rodando em http://localhost%s", addr)
	log.Printf("Métricas em http://localhost%s/metrics", addr)

	// Inicia o servidor HTTP
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}