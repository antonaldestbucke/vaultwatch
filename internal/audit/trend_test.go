package audit

import (
	"testing"
	"time"
)

func sampleScoredReports(scores []float64) []ScoredReport {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	var out []ScoredReport
	for i, s := range scores {
		out = append(out, ScoredReport{
			Timestamp: base.Add(time.Duration(i) * time.Hour),
			Score:     s,
			Drifted:   int(s * 10),
			Total:     10,
		})
	}
	return out
}

func TestBuildTrend_SortsChronologically(t *testing.T) {
	scores := []float64{0.8, 0.2, 0.5}
	snaps := sampleScoredReports(scores)
	// Reverse order
	snaps[0], snaps[2] = snaps[2], snaps[0]
	tr := BuildTrend(snaps)
	for i := 1; i < len(tr.Points); i++ {
		if tr.Points[i].Timestamp.Before(tr.Points[i-1].Timestamp) {
			t.Errorf("points not sorted at index %d", i)
		}
	}
}

func TestBuildTrend_Empty(t *testing.T) {
	tr := BuildTrend(nil)
	if len(tr.Points) != 0 {
		t.Errorf("expected empty trend")
	}
}

func TestAverageScore(t *testing.T) {
	tr := BuildTrend(sampleScoredReports([]float64{0.4, 0.6, 0.8}))
	avg := tr.AverageScore()
	if avg < 0.59 || avg > 0.61 {
		t.Errorf("expected ~0.6, got %f", avg)
	}
}

func TestAverageScore_Empty(t *testing.T) {
	tr := TrendReport{}
	if tr.AverageScore() != 0 {
		t.Errorf("expected 0 for empty trend")
	}
}

func TestWorstPoint(t *testing.T) {
	tr := BuildTrend(sampleScoredReports([]float64{0.2, 0.9, 0.4}))
	w, ok := tr.WorstPoint()
	if !ok {
		t.Fatal("expected a worst point")
	}
	if w.Score != 0.9 {
		t.Errorf("expected worst score 0.9, got %f", w.Score)
	}
}

func TestWorstPoint_Empty(t *testing.T) {
	tr := TrendReport{}
	_, ok := tr.WorstPoint()
	if ok {
		t.Error("expected no worst point for empty trend")
	}
}
