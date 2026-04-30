package audit

import (
	"testing"
)

func sampleSensitivityReports() []CompareReport {
	return []CompareReport{
		{
			Path:    "secret/app/prod",
			OnlyInA: []string{"db_password", "api_token", "host"},
			OnlyInB: []string{"db_password", "api_token"},
		},
		{
			Path:    "secret/app/staging",
			OnlyInA: []string{"host", "port", "region"},
			OnlyInB: []string{"host"},
		},
		{
			Path:    "secret/infra/certs",
			OnlyInA: []string{"cert_pem", "private_key", "ca_cert"},
			OnlyInB: []string{"cert_pem", "private_key", "ca_cert"},
		},
	}
}

func TestBuildSensitivity_CriticalPath(t *testing.T) {
	reports := sampleSensitivityReports()
	results := BuildSensitivity(reports)

	var certResult *SensitivityResult
	for i := range results {
		if results[i].Path == "secret/infra/certs" {
			certResult = &results[i]
			break
		}
	}
	if certResult == nil {
		t.Fatal("expected result for secret/infra/certs")
	}
	if certResult.Label != "critical" {
		t.Errorf("expected critical, got %s", certResult.Label)
	}
}

func TestBuildSensitivity_LowSensitivity(t *testing.T) {
	reports := []CompareReport{
		{
			Path:    "config/app",
			OnlyInA: []string{"region", "zone", "replica_count"},
			OnlyInB: []string{"region"},
		},
	}
	results := BuildSensitivity(reports)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Label != "low" {
		t.Errorf("expected low, got %s", results[0].Label)
	}
}

func TestBuildSensitivity_Empty(t *testing.T) {
	results := BuildSensitivity([]CompareReport{})
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestBuildSensitivity_SortedByScore(t *testing.T) {
	reports := sampleSensitivityReports()
	results := BuildSensitivity(reports)
	for i := 1; i < len(results); i++ {
		if results[i].Score > results[i-1].Score {
			t.Errorf("results not sorted by score descending at index %d", i)
		}
	}
}

func TestBuildSensitivity_MatchedKeysPopulated(t *testing.T) {
	reports := []CompareReport{
		{
			Path:    "secret/db",
			OnlyInA: []string{"db_password", "db_user"},
			OnlyInB: []string{},
		},
	}
	results := BuildSensitivity(reports)
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	if len(results[0].MatchedKeys) == 0 {
		t.Error("expected matched keys to be populated")
	}
}

func TestBuildSensitivity_SkipsEmptyKeyPaths(t *testing.T) {
	reports := []CompareReport{
		{Path: "empty/path", OnlyInA: []string{}, OnlyInB: []string{}},
	}
	results := BuildSensitivity(reports)
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty key path, got %d", len(results))
	}
}
