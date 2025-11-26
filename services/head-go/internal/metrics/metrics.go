


package metrics

import (
    "net/http"
    "log"
    "fmt"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    Requests = prometheus.NewCounterVec(prometheus.CounterOpts{
        Name: "head_requests_total",
        Help: "Total requests to head",
    }, []string{"model", "status"})

    Latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name: "head_request_latency_seconds",
        Help: "Request latency seconds",
    }, []string{"model"})
)

func Start(port int) {
    prometheus.MustRegister(Requests, Latency)
    addr := fmt.Sprintf(":%d", port)
    http.Handle("/metrics", promhttp.Handler())
    log.Printf("metrics: starting on %s", addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}


