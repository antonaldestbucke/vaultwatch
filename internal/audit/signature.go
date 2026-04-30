package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// SignatureEntry records a signed snapshot of scored reports at a point in time.
type SignatureEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Signature string    `json:"signature"`
	PathCount int       `json:"path_count"`
	DriftCount int     `json:"drift_count"`
}

// SignatureStore holds a list of signature entries.
type SignatureStore struct {
	Entries []SignatureEntry `json:"entries"`
}

// SignReports computes a deterministic SHA-256 signature over the given scored reports.
// The reports are sorted by path before hashing to ensure stability.
func SignReports(reports []ScoredReport) (string, error) {
	if len(reports) == 0 {
		return "", fmt.Errorf("no reports to sign")
	}

	sorted := make([]ScoredReport, len(reports))
	copy(sorted, reports)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Path < sorted[j].Path
	})

	data, err := json.Marshal(sorted)
	if err != nil {
		return "", fmt.Errorf("marshal reports: %w", err)
	}

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// RecordSignature appends a new SignatureEntry to the store based on the given reports.
func RecordSignature(store *SignatureStore, reports []ScoredReport) error {
	sig, err := SignReports(reports)
	if err != nil {
		return err
	}

	driftCount := 0
	for _, r := range reports {
		if r.Drifted {
			driftCount++
		}
	}

	entry := SignatureEntry{
		Timestamp:  time.Now().UTC(),
		Signature:  sig,
		PathCount:  len(reports),
		DriftCount: driftCount,
	}
	store.Entries = append(store.Entries, entry)
	return nil
}

// VerifySignature checks whether the given reports produce the expected signature.
func VerifySignature(reports []ScoredReport, expected string) (bool, error) {
	actual, err := SignReports(reports)
	if err != nil {
		return false, err
	}
	return actual == expected, nil
}
