package audit

import (
	"math"
	"sort"
	"time"
)

// ForecastPoint represents a predicted drift score at a future timestamp.
type ForecastPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score"`
}

// ForecastResult holds the linear regression forecast for a path.
type ForecastResult struct {
	Path        string           `json:"path"`
	Slope       float64          `json:"slope"`
	Intercept   float64          `json:"intercept"`
	Forecast    []ForecastPoint  `json:"forecast"`
	Trending    string           `json:"trending"` // "up", "down", "stable"
}

// BuildForecast performs a simple linear regression over scored trend points
// and projects drift scores forward by the given number of steps (each step = interval).
func BuildForecast(trend []ScoredTrendPoint, steps int, interval time.Duration) []ForecastResult {
	if len(trend) == 0 || steps <= 0 {
		return nil
	}

	// Group by path
	byPath := map[string][]ScoredTrendPoint{}
	for _, pt := range trend {
		byPath[pt.Path] = append(byPath[pt.Path], pt)
	}

	paths := make([]string, 0, len(byPath))
	for p := range byPath {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	results := make([]ForecastResult, 0, len(paths))
	for _, path := range paths {
		pts := byPath[path]
		sort.Slice(pts, func(i, j int) bool {
			return pts[i].Timestamp.Before(pts[j].Timestamp)
		})

		slope, intercept := linearRegression(pts)
		base := pts[len(pts)-1].Timestamp

		forecast := make([]ForecastPoint, steps)
		for i := 1; i <= steps; i++ {
			t := base.Add(interval * time.Duration(i))
			x := float64(t.Unix())
			score := math.Max(0, math.Min(100, slope*x+intercept))
			forecast[i-1] = ForecastPoint{Timestamp: t, Score: math.Round(score*100) / 100}
		}

		trending := classifyTrend(slope)
		results = append(results, ForecastResult{
			Path:      path,
			Slope:     math.Round(slope*1e9) / 1e9,
			Intercept: math.Round(intercept*100) / 100,
			Forecast:  forecast,
			Trending:  trending,
		})
	}
	return results
}

func linearRegression(pts []ScoredTrendPoint) (slope, intercept float64) {
	n := float64(len(pts))
	if n < 2 {
		if n == 1 {
			return 0, pts[0].Score
		}
		return 0, 0
	}
	var sumX, sumY, sumXY, sumX2 float64
	for _, pt := range pts {
		x := float64(pt.Timestamp.Unix())
		y := pt.Score
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	denom := n*sumX2 - sumX*sumX
	if denom == 0 {
		return 0, sumY / n
	}
	slope = (n*sumXY - sumX*sumY) / denom
	intercept = (sumY - slope*sumX) / n
	return slope, intercept
}

func classifyTrend(slope float64) string {
	const threshold = 1e-7
	switch {
	case slope > threshold:
		return "up"
	case slope < -threshold:
		return "down"
	default:
		return "stable"
	}
}
