package audit

import (
	"testing"
	"time"
)

var baseTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func sampleLineageForStaleness() LineageStore {
	return LineageStore{
		Entries: []LineageEntry{
			{Path: "secret/app", Environment: "prod", Timestamp: baseTime.Add(-10 * 24 * time.Hour)},
			{Path: "secret/app", Environment: "prod", Timestamp: baseTime.Add(-5 * 24 * time.Hour)},
			{Path: "secret/db", Environment: "prod", Timestamp: baseTime.Add(-40 * 24 * time.Hour)},
			{Path: "secret/app", Environment: "staging", Timestamp: baseTime.Add(-2 * 24 * time.Hour)},
		},
	}
}

func TestBuildStaleness_FlagsOldEntries(t *testing.T) {
	store := sampleLineageForStaleness()
	threshold := 30 * 24 * time.Hour
	report := BuildStaleness(store, threshold, baseTime)

	if report.StaleCount != 1 {
		t.Errorf("expected 1 stale entry, got %d", report.StaleCount)
	}
}

func TestBuildStaleness_NonStaleEntry(t *testing.T) {
	store := sampleLineageForStaleness()
	threshold := 30 * 24 * time.Hour
	report := BuildStaleness(store, threshold, baseTime)

	for _, e := range report.Entries {
		if e.Path == "secret/app" && e.Environment == "prod" && e.Stale {
			t.Errorf("secret/app prod should not be stale (last seen 5 days ago)")
		}
	}
}

func TestBuildStaleness_UsesLatestTimestamp(t *testing.T) {
	store := sampleLineageForStaleness()
	threshold := 7 * 24 * time.Hour
	report := BuildStaleness(store, threshold, baseTime)

	for _, e := range report.Entries {
		if e.Path == "secret/app" && e.Environment == "prod" {
			// latest entry is 5 days ago — within 7 day threshold
			if e.Stale {
				t.Errorf("expected secret/app prod to be fresh; age=%v", e.Age)
			}
		}
	}
}

func TestBuildStaleness_Empty(t *testing.T) {
	store := LineageStore{}
	report := BuildStaleness(store, 24*time.Hour, baseTime)

	if len(report.Entries) != 0 {
		t.Errorf("expected no entries, got %d", len(report.Entries))
	}
	if report.StaleCount != 0 {
		t.Errorf("expected stale count 0, got %d", report.StaleCount)
	}
}

func TestBuildStaleness_SortedOutput(t *testing.T) {
	store := sampleLineageForStaleness()
	report := BuildStaleness(store, 30*24*time.Hour, baseTime)

	for i := 1; i < len(report.Entries); i++ {
		prev := report.Entries[i-1]
		curr := report.Entries[i]
		if prev.Path > curr.Path {
			t.Errorf("entries not sorted by path: %s > %s", prev.Path, curr.Path)
		}
	}
}
