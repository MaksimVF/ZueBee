




package main

import (
    "log"
    "net/http"

    "tail-service/internal/config"
    "tail-service/internal/server"
)

func main() {
    cfg := config.Load()

    srv := server.NewServer(cfg)

    http.Handle("/ws", srv.WebSocketHandler())
    http.Handle("/metrics", srv.MetricsHandler())

    log.Printf("Starting tail service on %s", cfg.Addr)
    log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}




