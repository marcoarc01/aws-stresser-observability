package stress

import (
	"context"
	"log"
	"runtime"
	"sync"

	"github.com/marcoarc01/aws-stresser-observability/stresser-app/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("stresser-app")

// Engine controla o nível de stress da aplicação
type Engine struct {
	mu      sync.Mutex
	level   int
	workers int
	cancel  []chan struct{}
}

// NewEngine cria um novo motor de stress
func NewEngine() *Engine {
	return &Engine{
		level:   0,
		workers: 0,
		cancel:  make([]chan struct{}, 0),
	}
}

// GetLevel retorna o nível de stress atual
func (e *Engine) GetLevel() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.level
}

// GetWorkers retorna a quantidade de goroutines ativas
func (e *Engine) GetWorkers() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.workers
}

// SetLevel define o nível de stress e ajusta os workers
func (e *Engine) SetLevel(level int) {
	ctx := context.Background()
	_, span := tracer.Start(ctx, "engine.SetLevel")
	defer span.End()

	e.mu.Lock()
	defer e.mu.Unlock()

	oldLevel := e.level
	e.level = level

	span.SetAttributes(
		attribute.Int("stress.level.old", oldLevel),
		attribute.Int("stress.level.new", level),
	)

	// Para todos os workers atuais
	for _, ch := range e.cancel {
		close(ch)
	}
	e.cancel = make([]chan struct{}, 0)

	// Calcula quantos workers criar baseado no level
	maxCPUs := runtime.NumCPU()
	desiredWorkers := (level * maxCPUs) / 100

	// Inicia novos workers
	for i := 0; i < desiredWorkers; i++ {
		ch := make(chan struct{})
		e.cancel = append(e.cancel, ch)
		go cpuWorker(ch)
	}

	e.workers = desiredWorkers

	span.SetAttributes(
		attribute.Int("stress.workers.count", desiredWorkers),
		attribute.Int("stress.cpu.max", maxCPUs),
	)

	// MÉTRICAS — atualiza Prometheus
	metrics.StressLevel.Set(float64(level))
	metrics.StressCPUWorkers.Set(float64(desiredWorkers))
	metrics.StressChangesTotal.Inc()
	metrics.EstimatedCostUSD.Set(float64(level) * 0.001)

	log.Printf(`{"level":"info","component":"engine","stress":%d,"workers":%d,"old_level":%d,"msg":"stress updated"}`,
		level, desiredWorkers, oldLevel)
}

// cpuWorker é uma goroutine que consome CPU até receber sinal de parada
func cpuWorker(stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
			_ = 1 + 1
		}
	}
}
