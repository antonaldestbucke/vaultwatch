package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// QuotaRule defines a maximum allowed number of drifted keys for a path prefix.
type QuotaRule struct {
	PathPrefix string `json:"path_prefix"`
	MaxDrifted int    `json:"max_drifted"`
}

// QuotaStore holds a collection of quota rules.
type QuotaStore struct {
	Rules []QuotaRule `json:"rules"`
}

// QuotaViolation describes a path that has exceeded its quota.
type QuotaViolation struct {
	Path       string
	Drifted    int
	MaxAllowed int
	Rule       string
}

// SaveQuota writes a QuotaStore to disk as JSON.
func SaveQuota(path string, store QuotaStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal quota: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadQuota reads a QuotaStore from a JSON file.
func LoadQuota(path string) (QuotaStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return QuotaStore{}, fmt.Errorf("read quota file: %w", err)
	}
	var store QuotaStore
	if err := json.Unmarshal(data, &store); err != nil {
		return QuotaStore{}, fmt.Errorf("parse quota file: %w", err)
	}
	return store, nil
}

// EvaluateQuota checks scored reports against quota rules and returns violations.
func EvaluateQuota(reports []ScoredReport, store QuotaStore) []QuotaViolation {
	var violations []QuotaViolation

	for _, report := range reports {
		drifted := len(report.Report.OnlyInA) + len(report.Report.OnlyInB)
		if drifted == 0 {
			continue
		}
		for _, rule := range store.Rules {
			if rule.PathPrefix == "" || hasPrefix(report.Report.Path, rule.PathPrefix) {
				if drifted > rule.MaxDrifted {
					violations = append(violations, QuotaViolation{
						Path:       report.Report.Path,
						Drifted:    drifted,
						MaxAllowed: rule.MaxDrifted,
						Rule:       rule.PathPrefix,
					})
				}
				break
			}
		}
	}

	sort.Slice(violations, func(i, j int) bool {
		return violations[i].Path < violations[j].Path
	})
	return violations
}
