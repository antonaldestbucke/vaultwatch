package audit

import (
	"sort"
	"time"
)

// TrendPoint represents drift score at a point in time.
type TrendPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score"`
	Drifted   int       `json:"drifted"`
	Total     int       `json:"total"`
}

// TrendReport holds a series of trend points for analysis.
type TrendReport struct {
	Points []TrendPoint `json:"points"`
}

// BuildTrend constructs a TrendReport from a slice of scored snapshots.
func BuildTrend(snapshots []ScoredReport) TrendReport {
	points := make([]TrendPoint, 0, len(snapshots))
	for _, s := range snapshots {
		points = append(points, TrendPoint{
			Timestamp: s.Timestamp,
			Score:     s.Score,
			Drifted:   s.Drifted,
			Total:     s.Total,
		})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Timestamp.Before(points[j].Timestamp)
	})
	return TrendReport{Points: points}
}

// AverageScore returns the mean drift score across all trend points.
func (t TrendReport) AverageScore() float64 {
	if len(t.Points) == 0 {
		return 0
	}
	var sum float64
	for _, p := range t.Points {
		sum += p.Score
	}
	return sum / float64(len(t.Points))
}

// WorstPoint returns the TrendPoint with the highest drift score.
func (t TrendReport) WorstPoint() (TrendPoint, bool) {
	if len(t.Points) == 0 {
		return TrendPoint{}, false
	}
	worst := t.Points[0]
	for _, p := range t.Points[1:] {
		if p.Score > worst.Score {
			worst = p
		}
	}
	return worst, true
}
