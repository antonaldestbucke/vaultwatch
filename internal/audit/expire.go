package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ExpiryRule defines a TTL-based expiry rule for a secret path pattern.
type ExpiryRule struct {
	Pattern   string        `json:"pattern"`
	TTL       time.Duration `json:"ttl"`
	NotifyAt  time.Duration `json:"notify_at"` // warn when this close to expiry
}

// ExpiryStore holds all expiry rules and per-path last-seen timestamps.
type ExpiryStore struct {
	Rules    []ExpiryRule         `json:"rules"`
	LastSeen map[string]time.Time `json:"last_seen"`
}

// ExpiryResult describes the expiry status of a single path.
type ExpiryResult struct {
	Path      string
	Env       string
	Status    string // "ok", "warning", "expired"
	ExpiresAt time.Time
	Message   string
}

// SaveExpiry persists an ExpiryStore to disk.
func SaveExpiry(path string, store ExpiryStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal expiry store: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadExpiry reads an ExpiryStore from disk.
func LoadExpiry(path string) (ExpiryStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ExpiryStore{LastSeen: make(map[string]time.Time)}, nil
		}
		return ExpiryStore{}, fmt.Errorf("read expiry store: %w", err)
	}
	var store ExpiryStore
	if err := json.Unmarshal(data, &store); err != nil {
		return ExpiryStore{}, fmt.Errorf("unmarshal expiry store: %w", err)
	}
	if store.LastSeen == nil {
		store.LastSeen = make(map[string]time.Time)
	}
	return store, nil
}

// EvaluateExpiry checks each ScoredReport against expiry rules and returns results.
func EvaluateExpiry(reports []ScoredReport, store ExpiryStore, now time.Time) []ExpiryResult {
	var results []ExpiryResult
	for _, r := range reports {
		for _, rule := range store.Rules {
			if !matchesAny(r.Path, []string{rule.Pattern}) {
				continue
			}
			key := r.Env + ":" + r.Path
			lastSeen, ok := store.LastSeen[key]
			if !ok {
				lastSeen = now
			}
			expiresAt := lastSeen.Add(rule.TTL)
			result := ExpiryResult{
				Path:      r.Path,
				Env:       r.Env,
				ExpiresAt: expiresAt,
			}
			switch {
			case now.After(expiresAt):
				result.Status = "expired"
				result.Message = fmt.Sprintf("expired %s ago", now.Sub(expiresAt).Round(time.Second))
			case now.Add(rule.NotifyAt).After(expiresAt):
				result.Status = "warning"
				result.Message = fmt.Sprintf("expires in %s", expiresAt.Sub(now).Round(time.Second))
			default:
				result.Status = "ok"
				result.Message = fmt.Sprintf("expires in %s", expiresAt.Sub(now).Round(time.Second))
			}
			results = append(results, result)
			break
		}
	}
	return results
}
