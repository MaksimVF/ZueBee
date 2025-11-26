





package clients

import (
    "context"
    "time"

    "github.com/yourorg/head/gen"
    "google.golang.org/grpc"
)

type BillingClient struct {
    addr string
    conn *grpc.ClientConn
    stub gen.BillingServiceClient
}

func NewBillingClient(addr string) *BillingClient { return &BillingClient{addr: addr} }

func (c *BillingClient) Init(ctx context.Context) error {
    conn, err := grpc.DialContext(ctx, c.addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
    if err != nil { return err }
    c.conn = conn
    c.stub = gen.NewBillingServiceClient(conn)
    return nil
}

func (c *BillingClient) Close() { if c.conn != nil { c.conn.Close() } }

func (c *BillingClient) Reserve(ctx context.Context, userID, requestID string, tokensEstimate int32, costEstimate float64) (*gen.ReserveResponse, error) {
    return c.stub.Reserve(ctx, &gen.ReserveRequest{UserId:userID, RequestId:requestID, TokensEstimate:tokensEstimate, CostEstimate:costEstimate})
}

func (c *BillingClient) Commit(ctx context.Context, reservationID string, tokensActual int32, costActual float64) (*gen.CommitResponse, error) {
    return c.stub.Commit(ctx, &gen.CommitRequest{ReservationId:reservationID, TokensActual:tokensActual, CostActual:costActual})
}




