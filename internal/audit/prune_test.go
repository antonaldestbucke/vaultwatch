package audit

import (
	"testing"
	"time"
)

func sampleSnapshotForPrune() (map[string][]string, map[string]time.Time) {
	now := time.Now()
	snapshot := map[string][]string{
		"secret/old":    {"key1"},
		"secret/recent": {"key2"},
		"other/old":     {"key3"},
	}
	lastSeen := map[string]time.Time{
		"secret/old":    now.Add(-48 * time.Hour),
		"secret/recent": now.Add(-1 * time.Hour),
		"other/old":     now.Add(-72 * time.Hour),
	}
	return snapshot, lastSeen
}

func TestPruneSnapshots_RemovesStale(t *testing.T) {
	snap, lastSeen := sampleSnapshotForPrune()
	opts := PruneOptions{OlderThan: 24 * time.Hour}
	result, pr := PruneSnapshots(snap, lastSeen, opts)

	if _, ok := result["secret/old"]; ok {
		t.Error("expected secret/old to be pruned")
	}
	if _, ok := result["secret/recent"]; !ok {
		t.Error("expected secret/recent to be kept")
	}
	if len(pr.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(pr.Removed))
	}
}

func TestPruneSnapshots_DryRunKeepsAll(t *testing.T) {
	snap, lastSeen := sampleSnapshotForPrune()
	opts := PruneOptions{OlderThan: 1 * time.Hour, DryRun: true}
	result, pr := PruneSnapshots(snap, lastSeen, opts)

	if len(result) != len(snap) {
		t.Errorf("dry run should keep all entries, got %d", len(result))
	}
	if len(pr.Removed) == 0 {
		t.Error("expected removed list to be populated in dry run")
	}
}

func TestPruneSnapshots_PathPrefixFilter(t *testing.T) {
	snap, lastSeen := sampleSnapshotForPrune()
	opts := PruneOptions{OlderThan: 24 * time.Hour, PathPrefix: "secret/"}
	result, pr := PruneSnapshots(snap, lastSeen, opts)

	if _, ok := result["other/old"]; !ok {
		t.Error("expected other/old to be kept (outside prefix)")
	}
	if _, ok := result["secret/old"]; ok {
		t.Error("expected secret/old to be pruned")
	}
	if len(pr.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(pr.Removed))
	}
}

func TestPruneSnapshots_NothingStale(t *testing.T) {
	snap := map[string][]string{"secret/a": {"k"}}
	lastSeen := map[string]time.Time{"secret/a": time.Now()}
	opts := PruneOptions{OlderThan: 24 * time.Hour}
	result, pr := PruneSnapshots(snap, lastSeen, opts)

	if len(result) != 1 {
		t.Error("expected all entries kept")
	}
	if len(pr.Removed) != 0 {
		t.Error("expected no removals")
	}
}
