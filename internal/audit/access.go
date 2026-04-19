package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// AccessRule defines who can access a secret path.
type AccessRule struct {
	Path    string    `json:"path"`
	Owner   string    `json:"owner"`
	Team    string    `json:"team"`
	Expires time.Time `json:"expires,omitempty"`
}

// AccessStore holds all access rules.
type AccessStore struct {
	Rules []AccessRule `json:"rules"`
}

// SaveAccess writes the access store to a file.
func SaveAccess(path string, store AccessStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal access store: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadAccess reads the access store from a file.
func LoadAccess(path string) (AccessStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return AccessStore{}, nil
		}
		return AccessStore{}, fmt.Errorf("read access store: %w", err)
	}
	var store AccessStore
	if err := json.Unmarshal(data, &store); err != nil {
		return AccessStore{}, fmt.Errorf("unmarshal access store: %w", err)
	}
	return store, nil
}

// ApplyAccess annotates reports with owner/team from matching access rules.
func ApplyAccess(reports []CompareReport, store AccessStore) []CompareReport {
	for i, r := range reports {
		for _, rule := range store.Rules {
			if strings.HasPrefix(r.Path, rule.Path) {
				if reports[i].Annotations == nil {
					reports[i].Annotations = map[string]string{}
				}
				reports[i].Annotations["owner"] = rule.Owner
				reports[i].Annotations["team"] = rule.Team
				break
			}
		}
	}
	return reports
}

// LookupAccess returns the rule for a given path, if any.
func LookupAccess(store AccessStore, path string) (AccessRule, bool) {
	for _, rule := range store.Rules {
		if strings.HasPrefix(path, rule.Path) {
			return rule, true
		}
	}
	return AccessRule{}, false
}
