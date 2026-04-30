package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// RetentionRule defines how long scored reports for a path prefix should be kept.
type RetentionRule struct {
	PathPrefix string        `json:"path_prefix"`
	MaxAge     time.Duration `json:"max_age"`
}

// RetentionStore holds a list of retention rules.
type RetentionStore struct {
	Rules []RetentionRule `json:"rules"`
}

// RetentionResult describes whether a scored report was retained or pruned.
type RetentionResult struct {
	Path    string
	EnvA    string
	EnvB    string
	Pruned  bool
	Reason  string
}

// SaveRetention writes a RetentionStore to the given file path as JSON.
func SaveRetention(path string, store RetentionStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal retention: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadRetention reads a RetentionStore from the given file path.
func LoadRetention(path string) (RetentionStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return RetentionStore{}, fmt.Errorf("read retention file: %w", err)
	}
	var store RetentionStore
	if err := json.Unmarshal(data, &store); err != nil {
		return RetentionStore{}, fmt.Errorf("parse retention file: %w", err)
	}
	return store, nil
}

// EvaluateRetention checks each scored report against retention rules.
// Reports older than the matching rule's MaxAge are marked as pruned.
func EvaluateRetention(reports []ScoredReport, store RetentionStore, now time.Time) []RetentionResult {
	var results []RetentionResult
	for _, r := range reports {
		result := RetentionResult{
			Path: r.Path,
			EnvA: r.EnvA,
			EnvB: r.EnvB,
			Pruned: false,
		}
		for _, rule := range store.Rules {
			if hasPrefix(r.Path, rule.PathPrefix) {
				age := now.Sub(r.ScannedAt)
				if age > rule.MaxAge {
					result.Pruned = true
					result.Reason = fmt.Sprintf("age %s exceeds max_age %s for prefix %q", age.Round(time.Second), rule.MaxAge, rule.PathPrefix)
				}
				break
			}
		}
		results = append(results, result)
	}
	return results
}
