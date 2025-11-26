



package server

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/yourorg/head/gen"
    "github.com/yourorg/head/internal/cache"
    "github.com/yourorg/head/internal/config"
    "github.com/yourorg/head/internal/metrics"
    "github.com/yourorg/head/internal/providers"

    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

type Server struct {
    cfg *config.Config
    pm *providers.ProviderManager
    cache *cache.Client
}

func NewServer(cfg *config.Config) *Server {
    pm := providers.NewManager(cfg)
    return &Server{ cfg: cfg, pm: pm, cache: cache.NewClient(cfg.CacheAddr, cfg.CacheEnabled) }
}

func (s *Server) NewGRPCServer(opts ...grpc.ServerOption) *grpc.Server {
    grpcServer := grpc.NewServer(opts...)
    gen.RegisterChatServiceServer(grpcServer, s)
    reflection.Register(grpcServer)
    return grpcServer
}

func (s *Server) Serve(grpcServer *grpc.Server) error {
    // init cache client
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := s.cache.Init(ctx); err != nil {
        log.Printf("cache init error: %v (continuing)", err)
    }

    lis, err := net.Listen("tcp", s.cfg.GRPCAddr)
    if err != nil {
        return fmt.Errorf("listen: %w", err)
    }
    log.Printf("gRPC listening on %s", s.cfg.GRPCAddr)
    return grpcServer.Serve(lis)
}

func (s *Server) GracefulStop(ctx context.Context, grpcSrv *grpc.Server) error {
    done := make(chan struct{})
    go func() {
        grpcSrv.GracefulStop()
        close(done)
    }()
    select {
    case <-done:
        s.cache.Close()
        return nil
    case <-ctx.Done():
        grpcSrv.Stop()
        s.cache.Close()
        return ctx.Err()
    }
}

// --- Implement ChatService methods ---

func (s *Server) ChatCompletion(ctx context.Context, req *gen.ChatRequest) (*gen.ChatResponse, error) {
    model := req.Model
    metrics.Requests.WithLabelValues(model, "received").Inc()
    start := time.Now()

    // cache key
    keyBuf, _ := json.Marshal(req)
    cacheKey := fmt.Sprintf("head:cache:%x", sha256.Sum256(keyBuf))

    if s.cfg.CacheEnabled {
        if v, err := s.cache.Get(ctx, cacheKey); err == nil && v!="" {
            metrics.Requests.WithLabelValues(model, "cache_hit").Inc()
            var obj map[string]interface{}
            _ = json.Unmarshal([]byte(v), &obj)
            resp := &gen.ChatResponse{
                RequestId: req.RequestId,
                FullText: fmt.Sprintf("%v", obj["full_text"]),
                Model: fmt.Sprintf("%v", obj["model"]),
                Provider: fmt.Sprintf("%v", obj["provider"]),
                TokensUsed: int32(int(obj["tokens_used"].(float64))),
            }
            metrics.Latency.WithLabelValues(model).Observe(time.Since(start).Seconds())
            metrics.Requests.WithLabelValues(model, "ok").Inc()
            return resp, nil
        }
    }

    // convert messages
    msgs := make([]providers.Message, 0, len(req.Messages))
    for _, m := range req.Messages {
        msgs = append(msgs, providers.Message{ Role: m.Role, Content: m.Content })
    }

    provider, text, tokens, err := s.pm.Call(ctx, model, msgs, float32(req.Temperature), int(req.MaxTokens), false)
    if err!=nil {
        metrics.Requests.WithLabelValues(model, "error").Inc()
        return nil, err
    }

    // cache set
    if s.cfg.CacheEnabled {
        payload := map[string]interface{}{"full_text": text, "model": model, "provider": provider, "tokens_used": tokens}
        b, _ := json.Marshal(payload)
        _ = s.cache.Set(ctx, cacheKey, string(b), s.cfg.CacheTTL)
    }

    metrics.Latency.WithLabelValues(model).Observe(time.Since(start).Seconds())
    metrics.Requests.WithLabelValues(model, "ok").Inc()
    return &gen.ChatResponse{
        RequestId: req.RequestId,
        FullText: text,
        Model: model,
        Provider: provider,
        TokensUsed: int32(tokens),
    }, nil
}

func (s *Server) ChatCompletionStream(req *gen.ChatRequest, stream gen.ChatService_ChatCompletionStreamServer) error {
    model := req.Model
    metrics.Requests.WithLabelValues(model, "stream_received").Inc()
    start := time.Now()

    msgs := make([]providers.Message, 0, len(req.Messages))
    for _, m := range req.Messages {
        msgs = append(msgs, providers.Message{ Role: m.Role, Content: m.Content })
    }

    // call manager with stream flag; here we use local provider stream simulation
    provider, text, tokens, err := s.pm.Call(stream.Context(), model, msgs, float32(req.Temperature), int(req.MaxTokens), true)
    if err!=nil {
        metrics.Requests.WithLabelValues(model, "stream_error").Inc()
        return err
    }

    // simple chunking: split by space — replace with provider streaming for real provider
    words := strings.Fields(text)
    for i, w := range words {
        if err := stream.Send(&gen.ChatResponseChunk{
            RequestId: req.RequestId,
            Chunk: w,
            IsFinal: false,
            Provider: provider,
            TokensUsed: int32(len(w)/4 + 1),
        }); err != nil { return err }
        time.Sleep(10 * time.Millisecond)
        if i==len(words)-1 {
            // final
            _ = stream.Send(&gen.ChatResponseChunk{
                RequestId: req.RequestId, Chunk: "", IsFinal: true, Provider: provider,
            })
        }
    }

    metrics.Latency.WithLabelValues(model).Observe(time.Since(start).Seconds())
    metrics.Requests.WithLabelValues(model, "stream_ok").Inc()
    return nil
}



