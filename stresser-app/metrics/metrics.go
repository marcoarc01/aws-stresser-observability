package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)


// Nível atual do slider (0-100)
var CPUStressLevel = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "stresser_cpu_stress_level",
	Help: "Current CPU stress level percentage (0-100)",
})

// Workers ativos agora
var CPUWorkersActive = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "stresser_cpu_workers_active",
	Help: "Number of active CPU stress workers",
})


// Total de requests HTTP recebidos (por método e path)
var HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "stresser_http_requests_total",
	Help: "Total HTTP requests received",
}, []string{"method", "path", "status"})

// Quantas vezes o stress foi alterado
var StressChangesTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "stresser_stress_changes_total",
	Help: "Total number of stress level changes",
})