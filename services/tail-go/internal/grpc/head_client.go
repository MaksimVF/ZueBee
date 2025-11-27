



package grpc

import (
    "context"
    "log"

    pb "tail-service/tail-service/api/proto"

    "google.golang.org/grpc"
)

type HeadClient struct {
    client pb.LLMServiceClient
}

func NewHeadClient(addr string) *HeadClient {
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("failed to dial head: %v", err)
    }

    return &HeadClient{
        client: pb.NewLLMServiceClient(conn),
    }
}

func (c *HeadClient) Stream(ctx context.Context, req *pb.GenerateRequest) (pb.LLMService_GenerateClient, error) {
    return c.client.Generate(ctx, req)
}



