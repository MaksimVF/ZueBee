


package cache

import (
    "context"
    "time"
    "log"

    "github.com/yourorg/head/gen"
    "google.golang.org/grpc"
)

type Client struct {
    addr string
    conn *grpc.ClientConn
    stub gen.CacheServiceClient
    enabled bool
}

func NewClient(addr string, enabled bool) *Client {
    return &Client{ addr: addr, enabled: enabled }
}

func (c *Client) Init(ctx context.Context) error {
    if !c.enabled { return nil }
    conn, err := grpc.DialContext(ctx, c.addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
    if err!=nil { return err }
    c.conn = conn
    c.stub = gen.NewCacheServiceClient(conn)
    return nil
}

func (c *Client) Close() {
    if c.conn != nil { c.conn.Close() }
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
    if !c.enabled { return "", nil }
    resp, err := c.stub.Get(ctx, &gen.CacheGetRequest{Key:key})
    if err!=nil { return "", err }
    if !resp.Hit { return "", nil }
    return resp.ValueJson, nil
}

func (c *Client) Set(ctx context.Context, key string, value string, ttl int) error {
    if !c.enabled { return nil }
    _, err := c.stub.Set(ctx, &gen.CacheSetRequest{Key:key, ValueJson:value, TtlSeconds:int32(ttl)})
    return err
}


