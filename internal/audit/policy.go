package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// PolicyRule defines a rule that a secret path must satisfy.
type PolicyRule struct {
	PathPrefix  string   `json:"path_prefix"`
	RequiredKeys []string `json:"required_keys"`
	ForbiddenKeys []string `json:"forbidden_keys,omitempty"`
}

// PolicyStore holds a collection of policy rules.
type PolicyStore struct {
	Rules []PolicyRule `json:"rules"`
}

// PolicyViolation describes a single rule violation.
type PolicyViolation struct {
	Path    string
	Rule    PolicyRule
	Message string
}

// SavePolicy writes a PolicyStore to a JSON file.
func SavePolicy(path string, store PolicyStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal policy: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPolicy reads a PolicyStore from a JSON file.
func LoadPolicy(path string) (PolicyStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PolicyStore{}, fmt.Errorf("read policy: %w", err)
	}
	var store PolicyStore
	if err := json.Unmarshal(data, &store); err != nil {
		return PolicyStore{}, fmt.Errorf("unmarshal policy: %w", err)
	}
	return store, nil
}

// EvaluatePolicy checks reports against policy rules and returns violations.
func EvaluatePolicy(reports []CompareReport, store PolicyStore) []PolicyViolation {
	var violations []PolicyViolation
	for _, r := range reports {
		for _, rule := range store.Rules {
			if !strings.HasPrefix(r.Path, rule.PathPrefix) {
				continue
			}
			allKeys := append(r.OnlyInA, r.OnlyInB...)
			keySet := make(map[string]struct{}, len(allKeys))
			for _, k := range allKeys {
				keySet[strings.ToLower(k)] = struct{}{}
			}
			for _, req := range rule.RequiredKeys {
				if _, ok := keySet[strings.ToLower(req)]; !ok {
					violations = append(violations, PolicyViolation{
						Path:    r.Path,
						Rule:    rule,
						Message: fmt.Sprintf("required key %q missing", req),
					})
				}
			}
			for _, forbidden := range rule.ForbiddenKeys {
				if _, ok := keySet[strings.ToLower(forbidden)]; ok {
					violations = append(violations, PolicyViolation{
						Path:    r.Path,
						Rule:    rule,
						Message: fmt.Sprintf("forbidden key %q present", forbidden),
					})
				}
			}
		}
	}
	return violations
}
