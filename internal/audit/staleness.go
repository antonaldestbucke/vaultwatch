package audit

import (
	"sort"
	"time"
)

// StalenessEntry represents a path that has not changed within the expected window.
type StalenessEntry struct {
	Path        string        `json:"path"`
	Environment string        `json:"environment"`
	LastSeen    time.Time     `json:"last_seen"`
	Age         time.Duration `json:"age"`
	Stale       bool          `json:"stale"`
}

// StalenessReport is the result of a staleness evaluation.
type StalenessReport struct {
	GeneratedAt time.Time        `json:"generated_at"`
	Threshold   time.Duration    `json:"threshold"`
	Entries     []StalenessEntry `json:"entries"`
	StaleCount  int              `json:"stale_count"`
}

// BuildStaleness evaluates lineage history and flags paths whose most recent
// entry is older than the given threshold.
func BuildStaleness(store LineageStore, threshold time.Duration, now time.Time) StalenessReport {
	report := StalenessReport{
		GeneratedAt: now,
		Threshold:   threshold,
	}

	// Collect the latest timestamp per (path, env) pair.
	type key struct{ path, env string }
	latest := make(map[key]time.Time)

	for _, entry := range store.Entries {
		k := key{entry.Path, entry.Environment}
		if t, ok := latest[k]; !ok || entry.Timestamp.After(t) {
			latest[k] = entry.Timestamp
		}
	}

	for k, lastSeen := range latest {
		age := now.Sub(lastSeen)
		stale := age > threshold
		report.Entries = append(report.Entries, StalenessEntry{
			Path:        k.path,
			Environment: k.env,
			LastSeen:    lastSeen,
			Age:         age,
			Stale:       stale,
		})
		if stale {
			report.StaleCount++
		}
	}

	sort.Slice(report.Entries, func(i, j int) bool {
		if report.Entries[i].Path != report.Entries[j].Path {
			return report.Entries[i].Path < report.Entries[j].Path
		}
		return report.Entries[i].Environment < report.Entries[j].Environment
	})

	return report
}
