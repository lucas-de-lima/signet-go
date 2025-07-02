package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lucas-de-lima/signet-go/signet"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RecorderPrometheus implementa MetricsRecorder usando Prometheus
type RecorderPrometheus struct {
	success prometheus.Counter
	errors  *prometheus.CounterVec
}

func NewRecorderPrometheus() *RecorderPrometheus {
	r := &RecorderPrometheus{
		success: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "signet_token_validation_success_total",
			Help: "Total de tokens validados com sucesso",
		}),
		errors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "signet_token_validation_errors_total",
				Help: "Total de falhas na validação de tokens por motivo",
			},
			[]string{"reason"},
		),
	}
	prometheus.MustRegister(r.success, r.errors)
	return r
}

func (r *RecorderPrometheus) IncrementTokenValidation(ctx context.Context, success bool, failureReason string) {
	if success {
		r.success.Inc()
	} else {
		r.errors.WithLabelValues(failureReason).Inc()
	}
}

func main() {
	recorder := NewRecorderPrometheus()
	// Exemplo de uso manual:
	recorder.IncrementTokenValidation(context.Background(), true, signet.ReasonSuccess)
	recorder.IncrementTokenValidation(context.Background(), false, signet.ReasonTokenExpired)
	recorder.IncrementTokenValidation(context.Background(), false, signet.ReasonInvalidSignature)

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Prometheus metrics em http://localhost:2112/metrics")
	http.ListenAndServe(":2112", nil)
}
