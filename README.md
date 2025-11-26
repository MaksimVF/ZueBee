
# ZueBee Project

## Project Overview

ZueBee is a platform consisting of two interconnected services: Head and Tail. The Head-Service (Go) acts as the front-end proxy layer for accessing LLM (Large Language Models).

## Head-Service Purpose

Head-Service is the front-end proxy layer of your LLM access platform. It accepts client requests (HTTP/HTTPS, WebSocket streaming), performs authentication, authorization, quota/limit checks, and forwards requests downstream via gRPC to Tail-Service (lower proxy), which then accesses local inference, OpenAI-compatible APIs, or multi-backends.

## Main Roles of Head-Service

### 1. Authentication and Security

- Client API key verification (via Billing-Service)
- Subscription, token limit, and plan checks
- Request origin verification (IP allowlist, signature)
- Rate limiting + priority limits
- Web Application Firewall (primary)

### 2. Priorities and Billing

- Billing-Service request: token reservation
- Request priority assignment
- Request metadata tagging for Tail-Service queue
- Upon request completion: commit/release tokens

### 3. Traffic Distribution

- Request proxying via gRPC to Tail-Service
- Support for multiple Tail backends: local LLM, vLLM, OpenAI, Anthropic
- Selection of best Tail node by availability, load, and plan

### 4. WebSocket Streaming Gateway

- WebSocket → gRPC streaming conversion
- SSE support (optional)
- Support for pause, stop, and additional commands

## Distinctive Features (in the context of your architecture)

### 1. Multi-level Authentication Model

- API key → verification in Billing-Service
- JWT → for multi-user portal
- Request signature → for enterprise clients

### 2. Unified Metrics Endpoint

Head-Service aggregates:
- Its own metrics
- Tail-Service metrics (via gRPC health/metrics)
- Billing-Service metrics (quotas, tokens, burn rate)
- Priority request processing statistics

This addresses your need to collect head and tail metrics in one place.

## Tasks Performed by Head-Service

### 1. Pre-processing

- API key → verification
- Billing.reserveTokens
- User priority extraction
- Load queue for Tail-Service

### 2. During Processing

- WebSocket → stream chunks translation
- Heartbeat support
- Timeout tracking
- Inferencer error forwarding

### 3. Post-processing

- Billing.commitTokens
- Logging
- Metrics: latency, tokens_in/out, model, priority

## Technologies

- Go 1.22+
- gRPC + protobuf
- OpenTelemetry (Tracing + Metrics)
- Prometheus + Grafana
- JWT + API Keys + HMAC-Signature
- Redis (limit cache + rate limit + lua)
- PostgreSQL / Prisma (billing-service)

## Why Head-Service in Go?

- High performance and low latency
- Easy gRPC streaming implementation
- Better compatibility with high-load architecture
- Native WebSocket support
- Safe concurrent structures (goroutine/mutex)
- Low memory overhead - crucial for gateway services

## Architecture Position

```
Client → Head-Service → Tail-Service → LLM backend
  ↑       ↑               ↑
  │       │               └→ Metrics, Health
  │       └→ Billing-Service (reserve/commit tokens)
  └→ Redis (ratelimiting, caching)
```

## Development Status

Head-Service is a critical point requiring:
- Security
- Fault tolerance
- Strict SLAs
- Proper billing integration

## Full File Set

- main.go
- server.go
- routes.go
- middleware/auth.go
- middleware/rate_limit.go
- billing_client.go (gRPC)
- tail_client.go (gRPC)
- handlers/ws.go
- handlers/chat_completion.go
- metrics.go
- config.go
- proto/*.proto
- Makefile
- go.mod
- Dockerfile
