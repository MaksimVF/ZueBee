





package server

import (
    "context"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
    pb "tail-service/tail-service/api/proto"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type WebSocketRequest struct {
    Model      string `json:"model"`
    Prompt     string `json:"prompt"`
    ClientID  string `json:"client_id"`
    RequestID string `json:"request_id"`
    MaxTokens int32   `json:"max_tokens"`
}

func (s *Server) WebSocketHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Printf("Failed to upgrade to websocket: %v", err)
            return
        }
        defer conn.Close()

        var req WebSocketRequest
        if err := conn.ReadJSON(&req); err != nil {
            log.Printf("Failed to read request: %v", err)
            return
        }

        // Reserve tokens
        if !s.billing.ReserveTokens(req.ClientID, req.RequestID, int(req.MaxTokens)) {
            conn.WriteJSON(map[string]string{"error": "insufficient tokens"})
            return
        }

        // Create gRPC request
        grpcReq := &pb.GenerateRequest{
            Model:     req.Model,
            Prompt:    req.Prompt,
            ClientId:  req.ClientID,
            RequestId: req.RequestID,
            MaxTokens: req.MaxTokens,
        }

        // Stream from head service
        stream, err := s.headClient.Stream(context.Background(), grpcReq)
        if err != nil {
            log.Printf("Failed to stream from head: %v", err)
            conn.WriteJSON(map[string]string{"error": "stream failed"})
            return
        }

        // Stream responses to websocket
        for {
            resp, err := stream.Recv()
            if err != nil {
                log.Printf("Stream error: %v", err)
                break
            }

            if err := conn.WriteJSON(resp); err != nil {
                log.Printf("Failed to write to websocket: %v", err)
                break
            }

            if resp.IsEnd {
                break
            }
        }
    })
}





