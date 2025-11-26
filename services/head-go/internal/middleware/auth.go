







package middleware

import (
    "context"
    "fmt"

    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

// Unary interceptor checks internal metadata token if present
func UnaryAuthInterceptor(internalSecret string) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        if internalSecret == "" {
            return handler(ctx, req)
        }
        md, _ := metadata.FromIncomingContext(ctx)
        auth := ""
        if v := md.Get("authorization"); len(v) > 0 {
            auth = v[0]
        }
        if auth != fmt.Sprintf("Bearer %s", internalSecret) {
            return nil, grpc.Errorf(grpc.Code(grpc.Unauthenticated), "unauthorized")
        }
        return handler(ctx, req)
    }
}







