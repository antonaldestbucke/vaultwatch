package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// DigestEntry represents a hashed summary of a report set at a point in time.
type DigestEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Hash      string    `json:"hash"`
	PathCount int       `json:"path_count"`
	DriftCount int     `json:"drift_count"`
}

// BuildDigest produces a deterministic SHA-256 hash over a slice of CompareReports.
func BuildDigest(reports []CompareReport) (DigestEntry, error) {
	if len(reports) == 0 {
		return DigestEntry{}, fmt.Errorf("no reports provided")
	}

	// Sort for determinism
	sorted := make([]CompareReport, len(reports))
	copy(sorted, reports)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Path < sorted[j].Path
	})

	data, err := json.Marshal(sorted)
	if err != nil {
		return DigestEntry{}, fmt.Errorf("marshal error: %w", err)
	}

	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])

	driftCount := 0
	for _, r := range sorted {
		if hasDiff(r) {
			driftCount++
		}
	}

	return DigestEntry{
		Timestamp:  time.Now().UTC(),
		Hash:       hash,
		PathCount:  len(sorted),
		DriftCount: driftCount,
	}, nil
}

// DigestsMatch returns true if two DigestEntry values share the same hash.
func DigestsMatch(a, b DigestEntry) bool {
	return a.Hash == b.Hash
}
