



package metrics

import (
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    RequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "tail_requests_total",
            Help: "Total number of requests processed by tail service",
        },
        []string{"status"},
    )

    RequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "tail_request_duration_seconds",
            Help:    "Duration of requests in tail service",
            Buckets: prometheus.DefBuckets,
        },
        []string{"status"},
    )
)

func init() {
    prometheus.MustRegister(RequestsTotal, RequestDuration)
}

func MetricsHandler() http.Handler {
    return promhttp.Handler()
}



