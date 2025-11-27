


package config

import (
    "os"
)

type Config struct {
    Addr         string
    HeadGRPCAddr string
    BillingURL   string
}

func Load() Config {
    return Config{
        Addr:         get("TAIL_ADDR", ":8080"),
        HeadGRPCAddr: get("HEAD_GRPC_ADDR", "head:9000"),
        BillingURL:   get("BILLING_URL", "http://billing:8000"),
    }
}

func get(k, d string) string {
    v := os.Getenv(k)
    if v == "" {
        return d
    }
    return v
}


