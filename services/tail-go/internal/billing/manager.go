




package billing

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
)

type BillingManager struct {
    BaseURL string
}

func NewBilling(base string) *BillingManager {
    return &BillingManager{BaseURL: base}
}

type ReserveRequest struct {
    ClientID   string `json:"client_id"`
    Amount     int    `json:"amount"`
    RequestID  string `json:"request_id"`
}

type ReserveResponse struct {
    Approved bool `json:"approved"`
}

func (b *BillingManager) ReserveTokens(clientID string, reqID string, tokens int) bool {
    body, _ := json.Marshal(ReserveRequest{
        ClientID:  clientID,
        Amount:    tokens,
        RequestID: reqID,
    })

    url := fmt.Sprintf("%s/reserve", b.BaseURL)
    resp, err := http.Post(url, "application/json", strings.NewReader(string(body)))
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    var data ReserveResponse
    json.NewDecoder(resp.Body).Decode(&data)

    return data.Approved
}




