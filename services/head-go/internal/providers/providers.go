

package providers

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "sync"
    "time"

    "github.com/yourorg/head/internal/config"
)

type Message struct {
    Role string `json:"role"`
    Content string `json:"content"`
}

type ProviderManager struct {
    cfg *config.Config
    mu sync.Mutex
    state map[string]*providerState
    client *http.Client
}

type providerState struct {
    failCount int
    disabledUntil time.Time
}

func NewManager(cfg *config.Config) *ProviderManager {
    return &ProviderManager{
        cfg: cfg,
        state: make(map[string]*providerState),
        client: &http.Client{ Timeout: 30 * time.Second },
    }
}

func (m *ProviderManager) isAvailable(provider string) bool {
    m.mu.Lock(); defer m.mu.Unlock()
    st, ok := m.state[provider]
    if !ok { return true }
    if st.disabledUntil.IsZero() { return true }
    return time.Now().After(st.disabledUntil)
}

func (m *ProviderManager) markFailure(provider string) {
    m.mu.Lock(); defer m.mu.Unlock()
    st := m.state[provider]
    if st == nil { st = &providerState{}; m.state[provider] = st }
    st.failCount++
    if st.failCount >= m.cfg.ProviderFailThreshold {
        st.disabledUntil = time.Now().Add(time.Duration(m.cfg.ProviderDisableMinutes) * time.Minute)
    }
}

func (m *ProviderManager) markSuccess(provider string) {
    m.mu.Lock(); defer m.mu.Unlock()
    st := m.state[provider]
    if st==nil { return }
    st.failCount = 0
    st.disabledUntil = time.Time{}
}

// Call providers according to fallback list or provider keys order.
// returns provider name and response text + tokensUsed (best-effort)
func (m *ProviderManager) Call(ctx context.Context, model string, messages []Message, temperature float32, maxTokens int, stream bool) (string, string, int, error) {
    // define providers order
    pl := m.cfg.Fallbacks[model]
    if len(pl) == 0 {
        for k := range m.cfg.ProviderKeys { pl = append(pl, k) }
    }
    var lastErr error
    for _, p := range pl {
        if !m.isAvailable(p) { continue }
        // pick implementation
        switch p {
        case "openai":
            txt, tokens, err := m.callOpenAI(ctx, model, messages, temperature, maxTokens)
            if err != nil { m.markFailure(p); lastErr = err; continue }
            m.markSuccess(p)
            return p, txt, tokens, nil
        case "local":
            txt, tokens, err := m.localEcho(model, messages)
            if err != nil { m.markFailure(p); lastErr = err; continue }
            m.markSuccess(p)
            return p, txt, tokens, nil
        default:
            // unknown provider: skip
            continue
        }
    }
    if lastErr!=nil { return "", "", 0, lastErr }
    return "", "", 0, errors.New("no provider available")
}

// Example simple OpenAI REST call (non-streaming). Adapt payload per your provider.
func (m *ProviderManager) callOpenAI(ctx context.Context, model string, messages []Message, temperature float32, maxTokens int) (string, int, error) {
    apiKey := m.cfg.ProviderKeys["openai"]
    if apiKey=="" { return "", 0, fmt.Errorf("openai key not provided") }

    payload := map[string]interface{}{
        "model": model,
        "messages": messages,
        "temperature": temperature,
        "max_tokens": maxTokens,
    }
    b, _ := json.Marshal(payload)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")
    resp, err := m.client.Do(req)
    if err!=nil { return "", 0, err }
    defer resp.Body.Close()
    if resp.StatusCode>=400 {
        body, _ := io.ReadAll(resp.Body)
        return "", 0, fmt.Errorf("openai error %d: %s", resp.StatusCode, string(body))
    }
    var out struct {
        Choices []struct{
            Message struct { Content string `json:"content"` } `json:"message"`
        } `json:"choices"`
        Usage struct { TotalTokens int `json:"total_tokens"` } `json:"usage"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&out); err!=nil { return "", 0, err }
    text := ""
    for _, c := range out.Choices { text += c.Message.Content }
    return text, out.Usage.TotalTokens, nil
}

func (m *ProviderManager) localEcho(model string, messages []Message) (string, int, error) {
    // demo: concat messages
    s := ""
    for _, m := range messages { s += m.Content + " " }
    if len(s)>0 { s = s[:len(s)-1] }
    tokens := len(s)/4 + 1
    return "Echo: " + s, tokens, nil
}

