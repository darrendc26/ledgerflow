package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	registry          *prometheus.Registry
	AccountsCreated   prometheus.Counter
	PaymentsProcessed prometheus.Counter
	PaymentsFailed    prometheus.Counter
	PaymentsLatency   prometheus.Histogram
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

		AccountsCreated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "accounts_created_total",
				Help: "Total number of accounts created",
			},
		),
	}

	registry.MustRegister(p.PaymentsProcessed)
	registry.MustRegister(p.PaymentsFailed)
	registry.MustRegister(p.AccountsCreated)

	p.PaymentsLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "payment_processing_seconds",
			Help:    "Time taken to process payments",
			Buckets: prometheus.DefBuckets,
		},
	)

	registry.MustRegister(p.PaymentsLatency)
	return p
}

func (p *Prometheus) Start(port string) {

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{}))

	go func() {
		log.Println("Prometheus metrics exposed on", port+"/metrics")
		log.Println(http.ListenAndServe(port, mux))
	}()
}
