package audit

import (
	"math"
	"sort"
	"time"
)

// VelocityEntry represents the drift velocity for a single secret path.
type VelocityEntry struct {
	Path          string
	Env           string
	DriftRate     float64 // drifts per day
	TotalDrifts   int
	SpanDays      float64
	VelocityLabel string // "accelerating", "stable", "decelerating"
}

// BuildVelocity computes how quickly each path is drifting over time
// by analysing a sorted slice of TrendPoints.
func BuildVelocity(trends []TrendPoint, minPoints int) []VelocityEntry {
	if minPoints < 2 {
		minPoints = 2
	}

	// Group trend points by path+env key.
	type key struct{ path, env string }
	groups := make(map[key][]TrendPoint)
	for _, tp := range trends {
		k := key{tp.Path, tp.Env}
		groups[k] = append(groups[k], tp)
	}

	var entries []VelocityEntry
	for k, pts := range groups {
		if len(pts) < minPoints {
			continue
		}
		sort.Slice(pts, func(i, j int) bool {
			return pts[i].Timestamp.Before(pts[j].Timestamp)
		})

		spanDays := pts[len(pts)-1].Timestamp.Sub(pts[0].Timestamp).Hours() / 24
		if spanDays < 0.001 {
			spanDays = 0.001
		}

		totalDrifts := countDrifts(pts)
		driftRate := float64(totalDrifts) / spanDays

		// Compare first-half rate vs second-half rate to classify velocity.
		mid := len(pts) / 2
		firstRate := driftRateForSlice(pts[:mid])
		secondRate := driftRateForSlice(pts[mid:])
		label := classifyVelocity(firstRate, secondRate)

		entries = append(entries, VelocityEntry{
			Path:          k.path,
			Env:           k.env,
			DriftRate:     math.Round(driftRate*100) / 100,
			TotalDrifts:   totalDrifts,
			SpanDays:      math.Round(spanDays*10) / 10,
			VelocityLabel: label,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].DriftRate > entries[j].DriftRate
	})
	return entries
}

func countDrifts(pts []TrendPoint) int {
	count := 0
	for _, p := range pts {
		if p.Drifted {
			count++
		}
	}
	return count
}

func driftRateForSlice(pts []TrendPoint) float64 {
	if len(pts) < 2 {
		return 0
	}
	span := pts[len(pts)-1].Timestamp.Sub(pts[0].Timestamp).Hours() / 24
	if span < 0.001 {
		span = 0.001
	}
	return float64(countDrifts(pts)) / span
}

func classifyVelocity(first, second float64) string {
	switch {
	case second > first*1.2:
		return "accelerating"
	case second < first*0.8:
		return "decelerating"
	default:
		return "stable"
	}
}

// TrendPoint is a minimal time-stamped drift observation used by velocity.
// It mirrors what BuildTrend produces so callers can reuse the same data.
type TrendPoint struct {
	Path      string
	Env       string
	Timestamp time.Time
	Score     float64
	Drifted   bool
}
