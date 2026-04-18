package audit

import (
	"fmt"
	"time"
)

// Baseline represents a named snapshot used as a reference point for future comparisons.
type Baseline struct {
	Name      string                 `json:"name"`
	CreatedAt time.Time              `json:"created_at"`
	Path      string                 `json:"path"`
	Data      map[string][]string    `json:"data"` // env -> keys
}

// SaveBaseline persists a baseline to disk using the snapshot mechanism.
func SaveBaseline(dir, name, path string, reports []CompareReport) error {
	data := make(map[string][]string)
	for _, r := range reports {
		data[r.EnvA] = r.OnlyInA
		data[r.EnvB] = r.OnlyInB
	}
	b := Baseline{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Path:      path,
		Data:      data,
	}
	filePath := fmt.Sprintf("%s/%s.baseline.json", dir, name)
	return SaveSnapshot(filePath, b)
}

// LoadBaseline reads a named baseline from disk.
func LoadBaseline(dir, name string) (*Baseline, error) {
	filePath := fmt.Sprintf("%s/%s.baseline.json", dir, name)
	var b Baseline
	if err := LoadSnapshot(filePath, &b); err != nil {
		return nil, fmt.Errorf("load baseline %q: %w", name, err)
	}
	return &b, nil
}

// DiffAgainstBaseline compares current reports against a saved baseline.
// Returns paths that have changed since the baseline was taken.
func DiffAgainstBaseline(baseline *Baseline, reports []CompareReport) []string {
	changed := []string{}
	for _, r := range reports {
		baseA := baseline.Data[r.EnvA]
		baseB := baseline.Data[r.EnvB]
		if !equalSlices(baseA, r.OnlyInA) || !equalSlices(baseB, r.OnlyInB) {
			changed = append(changed, r.Path)
		}
	}
	return changed
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	set := toSet(a)
	for _, v := range b {
		if !set[v] {
			return false
		}
	}
	return true
}
