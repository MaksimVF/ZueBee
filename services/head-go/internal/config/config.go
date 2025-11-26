

package config

import (
    "encoding/json"
    "os"
    "strings"
    "time"

    "github.com/joho/godotenv"
)

type MTLSConfig struct {
    Enabled  bool
    CertFile string
    KeyFile  string
    CAFile   string
}

type Config struct {
    GRPCAddr    string
    MetricsPort int
    ProviderKeys map[string]string
    Fallbacks    map[string][]string
    CacheAddr    string
    CacheEnabled bool
    MTLS         MTLSConfig
    ProviderFailThreshold int
    ProviderDisableMinutes int
    ProviderConcurrency int
    CacheTTL int
}

func Load() *Config {
    // optional .env
    _ = godotenv.Load()

    cfg := &Config{
        GRPCAddr: ":50055",
        MetricsPort: 9001,
        CacheAddr: "cache:50053",
        CacheEnabled: false,
        ProviderFailThreshold: 5,
        ProviderDisableMinutes: 5,
        ProviderConcurrency: 8,
        CacheTTL: 300,
    }

    if v := os.Getenv("CHAT_ADDR"); v != "" { cfg.GRPCAddr = v }
    if v := os.Getenv("METRICS_PORT"); v != "" { /* parse */ }
    if v := os.Getenv("CACHE_ADDR"); v != "" { cfg.CacheAddr = v }
    if v := os.Getenv("CACHE_ENABLED"); strings.ToLower(v) == "true" { cfg.CacheEnabled = true }

    // parse JSON envs safely (strip quotes/newlines)
    parseJSON := func(env string, out interface{}) {
        raw := os.Getenv(env)
        if raw == "" { return }
        s := strings.TrimSpace(raw)
        if (strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) || (strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")) {
            s = s[1:len(s)-1]
        }
        s = strings.ReplaceAll(s, "\\n", "")
        _ = json.Unmarshal([]byte(s), out)
    }

    parseJSON("PROVIDER_KEYS", &cfg.ProviderKeys)
    if cfg.ProviderKeys == nil {
        cfg.ProviderKeys = map[string]string{}
    }
    parseJSON("FALLBACKS", &cfg.Fallbacks)
    if cfg.Fallbacks == nil {
        cfg.Fallbacks = map[string][]string{}
    }

    // mTLS envs
    if os.Getenv("MTLS_ENABLED") == "1" {
        cfg.MTLS.Enabled = true
        cfg.MTLS.CertFile = os.Getenv("MTLS_CERT_FILE")
        cfg.MTLS.KeyFile = os.Getenv("MTLS_KEY_FILE")
        cfg.MTLS.CAFile = os.Getenv("MTLS_CA_FILE")
    }

    if v := os.Getenv("PROVIDER_FAIL_THRESHOLD"); v != "" { /* parse int */ }

    return cfg
}

