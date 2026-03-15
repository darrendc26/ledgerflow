package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	registry *prometheus.Registry

	PaymentsProcessed prometheus.Counter
	PaymentsFailed    prometheus.Counter
}

func NewPrometheus() *Prometheus {

	registry := prometheus.NewRegistry()

	p := &Prometheus{
		registry: registry,

		PaymentsProcessed: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "payments_processed_total",
				Help: "Total number of payments processed",
			},
		),

		PaymentsFailed: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "payments_failed_total",
				Help: "Total number of payments failed",
			},
		),
	}

	registry.MustRegister(p.PaymentsProcessed)
	registry.MustRegister(p.PaymentsFailed)

	return p
}

func (p *Prometheus) Start() {

	http.Handle("/metrics", promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{}))

	go func() {
		log.Println("Prometheus metrics exposed on :2112/metrics")
		log.Println(http.ListenAndServe(":2112", nil))
	}()
}
