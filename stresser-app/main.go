package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/marcoarc01/aws-stresser-observability/stresser-app/handlers"
	"github.com/marcoarc01/aws-stresser-observability/stresser-app/metrics"
	"github.com/marcoarc01/aws-stresser-observability/stresser-app/stress"
	"github.com/marcoarc01/aws-stresser-observability/stresser-app/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	// Inicializa OpenTelemetry
	shutdown := tracing.InitTracer()
	defer shutdown()

	// Inicializa o motor de stress
	engine := stress.NewEngine()

	// Configura as rotas
	mux := http.NewServeMux()

	// Rotas da API
	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("GET /api/state", handlers.StateHandler(engine))
	mux.HandleFunc("POST /api/stress", handlers.StressHandler(engine))
	mux.Handle("GET /metrics", promhttp.Handler())

	// Arquivos estaticos (CSS, JS)
	mux.Handle("GET /static/", handlers.StaticHandler())

	// UI estatica
	mux.HandleFunc("GET /", handlers.UIHandler)

	// Porta do servidor
	port := os.Getenv("STRESSER_PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Stresser App rodando em http://localhost%s", addr)
	log.Printf("Metricas em http://localhost%s/metrics", addr)

	// Aplica middlewares: OTel (tracing) + Metrics
	handler := otelhttp.NewHandler(metrics.MetricsMiddleware(mux), "stresser-app")

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
