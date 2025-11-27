




package server

import (
    "net/http"

    "tail-service/internal/config"
    "tail-service/internal/grpc"
    "tail-service/internal/billing"
    "tail-service/internal/priority"
    "tail-service/internal/metrics"
)

type Server struct {
    cfg          config.Config
    headClient   *grpc.HeadClient
    billing      *billing.BillingManager
    priorityQueue *priority.PriorityQueue
}

func NewServer(cfg config.Config) *Server {
    return &Server{
        cfg:          cfg,
        headClient:   grpc.NewHeadClient(cfg.HeadGRPCAddr),
        billing:      billing.NewBilling(cfg.BillingURL),
        priorityQueue: priority.New(),
    }
}

func (s *Server) MetricsHandler() http.Handler {
    return metrics.MetricsHandler()
}




