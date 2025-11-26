





package clients

import (
    "context"
    "time"

    "github.com/yourorg/head/gen"
    "google.golang.org/grpc"
)

type TailClient struct {
    addr string
    conn *grpc.ClientConn
    stub gen.ChatServiceClient
}

func NewTailClient(addr string) *TailClient { return &TailClient{addr: addr} }

func (c *TailClient) Init(ctx context.Context) error {
    conn, err := grpc.DialContext(ctx, c.addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
    if err != nil { return err }
    c.conn = conn
    c.stub = gen.NewChatServiceClient(conn)
    return nil
}

func (c *TailClient) Close() { if c.conn != nil { c.conn.Close() } }





