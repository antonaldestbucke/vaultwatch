package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// TagStore maps secret paths to a list of string tags.
type TagStore map[string][]string

// SaveTags writes the tag store to a JSON file.
func SaveTags(path string, store TagStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tags: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadTags reads a tag store from a JSON file.
func LoadTags(path string) (TagStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read tags: %w", err)
	}
	var store TagStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("unmarshal tags: %w", err)
	}
	return store, nil
}

// ApplyTags attaches tags to each CompareReport whose path matches an entry in
// the store. Tags are stored in the report's annotation note field via a
// dedicated Tags slice returned alongside the reports.
func ApplyTags(reports []CompareReport, store TagStore) map[string][]string {
	result := make(map[string][]string, len(reports))
	for _, r := range reports {
		if tags, ok := store[r.Path]; ok {
			result[r.Path] = tags
		}
	}
	return result
}

// FilterByTag returns only the reports whose path has at least one of the given
// tags in the store.
func FilterByTag(reports []CompareReport, store TagStore, tag string) []CompareReport {
	tag = strings.ToLower(tag)
	var out []CompareReport
	for _, r := range reports {
		for _, t := range store[r.Path] {
			if strings.ToLower(t) == tag {
				out = append(out, r)
				break
			}
		}
	}
	return out
}
