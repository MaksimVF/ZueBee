
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/yourorg/head/internal/config"
    "github.com/yourorg/head/internal/metrics"
    "github.com/yourorg/head/internal/server"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

func main() {
    cfg := config.Load()

    // Start Prometheus metrics http server
    go metrics.Start(cfg.MetricsPort)

    // Initialize server (gRPC)
    srv := server.NewServer(cfg)

    // Build gRPC server options: mTLS if enabled
    var opts []grpc.ServerOption
    if cfg.MTLS.Enabled {
        creds, err := credentials.NewServerTLSFromFile(cfg.MTLS.CertFile, cfg.MTLS.KeyFile)
        if err != nil {
            log.Fatalf("failed to load server TLS cert/key: %v", err)
        }
        opts = append(opts, grpc.Creds(creds))
        log.Println("mTLS: enabled for gRPC server")
    }

    grpcSrv := srv.NewGRPCServer(opts...)

    // Listen and serve
    go func() {
        if err := srv.Serve(grpcSrv); err != nil {
            log.Fatalf("gRPC serve error: %v", err)
        }
    }()

    // Graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop
    log.Println("Shutting down head-service...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.GracefulStop(ctx, grpcSrv); err != nil {
        log.Printf("Graceful stop error: %v", err)
    } else {
        log.Println("Stopped cleanly")
    }
}
