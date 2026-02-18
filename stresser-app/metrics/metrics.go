package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Nível atual do slider (0-100)
var StressLevel = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "stress_level",
	Help: "Current CPU stress level percentage (0-100)",
})

// Workers ativos agora
var StressCPUWorkers = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "stress_cpu_workers",
	Help: "Number of active CPU stress workers",
})

// Total de requests HTTP recebidos
var HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total HTTP requests received",
}, []string{"method", "path", "status"})

// Quantas vezes o stress foi alterado
var StressChangesTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "stress_changes_total",
	Help: "Total number of stress level changes",
})

// Histogram de latencia HTTP
var HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_request_duration_seconds",
	Help: "Duration of HTTP requests in seconds",
}, []string{"path", "method"})

// Emails enviados com sucesso
var EmailsSentTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "emails_sent_total",
	Help: "Total emails sent successfully",
})

// Erros ao enviar email
var EmailSendErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "email_send_errors_total",
	Help: "Total email send errors",
})

// Uploads S3 com sucesso
var S3UploadsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "s3_uploads_total",
	Help: "Total S3 uploads completed",
})

// Erros de upload S3
var S3UploadErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "s3_upload_errors_total",
	Help: "Total S3 upload errors",
})

// Custo estimado simulado
var EstimatedCostUSD = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "estimated_cost_usd",
	Help: "Estimated cost in USD (simulated)",
})