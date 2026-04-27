package audit

import (
	"testing"
	"time"
)

func sampleTrendPoints() []TrendPoint {
	now := time.Now()
	return []TrendPoint{
		{Path: "secret/app", Env: "prod", Timestamp: now.Add(-6 * 24 * time.Hour), Score: 80, Drifted: false},
		{Path: "secret/app", Env: "prod", Timestamp: now.Add(-4 * 24 * time.Hour), Score: 60, Drifted: true},
		{Path: "secret/app", Env: "prod", Timestamp: now.Add(-2 * 24 * time.Hour), Score: 50, Drifted: true},
		{Path: "secret/app", Env: "prod", Timestamp: now, Score: 40, Drifted: true},
		{Path: "secret/db", Env: "staging", Timestamp: now.Add(-5 * 24 * time.Hour), Score: 90, Drifted: false},
		{Path: "secret/db", Env: "staging", Timestamp: now.Add(-1 * 24 * time.Hour), Score: 85, Drifted: false},
	}
}

func TestBuildVelocity_ReturnsBothPaths(t *testing.T) {
	pts := sampleTrendPoints()
	result := BuildVelocity(pts, 2)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestBuildVelocity_SortedByDriftRate(t *testing.T) {
	pts := sampleTrendPoints()
	result := BuildVelocity(pts, 2)
	if result[0].DriftRate < result[1].DriftRate {
		t.Errorf("expected descending drift rate order")
	}
}

func TestBuildVelocity_DriftRatePositive(t *testing.T) {
	pts := sampleTrendPoints()
	result := BuildVelocity(pts, 2)
	for _, e := range result {
		if e.DriftRate < 0 {
			t.Errorf("drift rate should not be negative for path %s", e.Path)
		}
	}
}

func TestBuildVelocity_NoDriftPath(t *testing.T) {
	pts := sampleTrendPoints()
	result := BuildVelocity(pts, 2)
	var dbEntry *VelocityEntry
	for i := range result {
		if result[i].Path == "secret/db" {
			dbEntry = &result[i]
		}
	}
	if dbEntry == nil {
		t.Fatal("expected entry for secret/db")
	}
	if dbEntry.TotalDrifts != 0 {
		t.Errorf("expected 0 drifts for secret/db, got %d", dbEntry.TotalDrifts)
	}
}

func TestBuildVelocity_InsufficientPoints(t *testing.T) {
	now := time.Now()
	pts := []TrendPoint{
		{Path: "secret/lone", Env: "prod", Timestamp: now, Score: 50, Drifted: true},
	}
	result := BuildVelocity(pts, 2)
	if len(result) != 0 {
		t.Errorf("expected no entries for path with fewer than minPoints, got %d", len(result))
	}
}

func TestBuildVelocity_AcceleratingLabel(t *testing.T) {
	now := time.Now()
	pts := []TrendPoint{
		{Path: "secret/fast", Env: "prod", Timestamp: now.Add(-8 * 24 * time.Hour), Score: 70, Drifted: false},
		{Path: "secret/fast", Env: "prod", Timestamp: now.Add(-6 * 24 * time.Hour), Score: 65, Drifted: false},
		{Path: "secret/fast", Env: "prod", Timestamp: now.Add(-2 * 24 * time.Hour), Score: 50, Drifted: true},
		{Path: "secret/fast", Env: "prod", Timestamp: now, Score: 30, Drifted: true},
	}
	result := BuildVelocity(pts, 2)
	if len(result) == 0 {
		t.Fatal("expected at least one entry")
	}
	if result[0].VelocityLabel != "accelerating" {
		t.Errorf("expected accelerating, got %s", result[0].VelocityLabel)
	}
}

func TestBuildVelocity_Empty(t *testing.T) {
	result := BuildVelocity(nil, 2)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil input")
	}
}
