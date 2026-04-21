package audit

import (
	"testing"
)

func sampleSimilarityReports() []DiffReport {
	return []DiffReport{
		{
			Path:    "secret/app/db",
			OnlyInA: []string{"password", "host"},
			OnlyInB: []string{"port"},
		},
		{
			Path:    "secret/app/cache",
			OnlyInA: []string{},
			OnlyInB: []string{},
		},
		{
			Path:    "secret/app/api",
			OnlyInA: []string{"key", "secret", "endpoint"},
			OnlyInB: []string{"token", "url", "version", "timeout"},
		},
	}
}

func TestComputeSimilarity_SortedByJaccard(t *testing.T) {
	reports := sampleSimilarityReports()
	results := ComputeSimilarity(reports, "staging", "production")

	for i := 1; i < len(results); i++ {
		if results[i].Jaccard < results[i-1].Jaccard {
			t.Errorf("expected results sorted ascending by Jaccard, got %v before %v",
				results[i-1].Jaccard, results[i].Jaccard)
		}
	}
}

func TestComputeSimilarity_EmptyKeys_IsHigh(t *testing.T) {
	reports := []DiffReport{
		{Path: "secret/app/cache", OnlyInA: []string{}, OnlyInB: []string{}},
	}
	results := ComputeSimilarity(reports, "dev", "prod")

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Jaccard != 1.0 {
		t.Errorf("expected Jaccard 1.0 for empty sets, got %f", results[0].Jaccard)
	}
	if results[0].Similarity != "high" {
		t.Errorf("expected similarity 'high', got %s", results[0].Similarity)
	}
}

func TestComputeSimilarity_LowSimilarity(t *testing.T) {
	reports := []DiffReport{
		{
			Path:    "secret/app/api",
			OnlyInA: []string{"key", "secret", "endpoint"},
			OnlyInB: []string{"token", "url", "version", "timeout"},
		},
	}
	results := ComputeSimilarity(reports, "staging", "production")

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Similarity != "low" {
		t.Errorf("expected similarity 'low', got %s", results[0].Similarity)
	}
	if results[0].Jaccard >= 0.5 {
		t.Errorf("expected Jaccard < 0.5, got %f", results[0].Jaccard)
	}
}

func TestComputeSimilarity_EnvFieldsSet(t *testing.T) {
	reports := []DiffReport{
		{Path: "secret/x", OnlyInA: []string{"a"}, OnlyInB: []string{"b"}},
	}
	results := ComputeSimilarity(reports, "dev", "prod")

	if results[0].EnvA != "dev" || results[0].EnvB != "prod" {
		t.Errorf("expected envA=dev envB=prod, got %s %s", results[0].EnvA, results[0].EnvB)
	}
}

func TestClassifySimilarity(t *testing.T) {
	cases := []struct {
		score    float64
		want     string
	}{
		{1.0, "high"},
		{0.8, "high"},
		{0.79, "medium"},
		{0.5, "medium"},
		{0.49, "low"},
		{0.0, "low"},
	}
	for _, tc := range cases {
		got := classifySimilarity(tc.score)
		if got != tc.want {
			t.Errorf("classifySimilarity(%v) = %s, want %s", tc.score, got, tc.want)
		}
	}
}
