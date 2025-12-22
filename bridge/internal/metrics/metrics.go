package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bridge_requests_total",
			Help: "Total requests by service",
		}, []string{"service"})

	RateLimit = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bridge_rate_limited_total",
			Help: "Total rate-limited requests",
		}, []string{"service"})

	HealthStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bridge_service_healthy",
			Help: "Service health status",
		}, []string{"service", "addr"})
)

func Init() {
	prometheus.MustRegister(Requests, RateLimit, HealthStatus)
}

func StartHTTP(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(addr, nil)
}
