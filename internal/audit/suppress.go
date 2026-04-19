package audit

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

// SuppressRule defines a rule to suppress alerts for a specific path and reason.
type SuppressRule struct {
	Path      string    `json:"path"`
	Reason    string    `json:"reason"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SuppressStore holds a list of suppression rules.
type SuppressStore struct {
	Rules []SuppressRule `json:"rules"`
}

// SaveSuppressions writes the store to a JSON file.
func SaveSuppressions(path string, store SuppressStore) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(store)
}

// LoadSuppressions reads a SuppressStore from a JSON file.
func LoadSuppressions(path string) (SuppressStore, error) {
	f, err := os.Open(path)
	if err != nil {
		return SuppressStore{}, err
	}
	defer f.Close()
	var store SuppressStore
	if err := json.NewDecoder(f).Decode(&store); err != nil {
		return SuppressStore{}, err
	}
	return store, nil
}

// IsSuppressed returns true if the given path matches an active (non-expired) suppression rule.
func IsSuppressed(store SuppressStore, path string) bool {
	now := time.Now()
	for _, rule := range store.Rules {
		if strings.HasPrefix(path, rule.Path) && now.Before(rule.ExpiresAt) {
			return true
		}
	}
	return false
}

// ApplySuppressions filters out scored reports whose paths are suppressed.
func ApplySuppressions(reports []ScoredReport, store SuppressStore) []ScoredReport {
	var out []ScoredReport
	for _, r := range reports {
		if !IsSuppressed(store, r.Path) {
			out = append(out, r)
		}
	}
	return out
}
