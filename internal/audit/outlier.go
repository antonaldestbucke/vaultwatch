package audit

import (
	"math"
	"sort"
)

// OutlierResult represents a path whose drift score deviates significantly
// from the mean across all scored reports.
type OutlierResult struct {
	Path    string
	Envs    []string
	Score   float64
	ZScore  float64
	IsOutlier bool
}

// BuildOutliers identifies paths whose drift scores are statistical outliers
// based on z-score with the given threshold (e.g. 1.5 or 2.0).
func BuildOutliers(reports []ScoredReport, zThreshold float64) []OutlierResult {
	if len(reports) == 0 {
		return nil
	}

	mean := averageScoreFor(reports)
	stddev := stddevScore(reports, mean)

	results := make([]OutlierResult, 0, len(reports))
	for _, r := range reports {
		z := 0.0
		if stddev > 0 {
			z = (r.Score - mean) / stddev
		}
		results = append(results, OutlierResult{
			Path:      r.Path,
			Envs:      r.Envs,
			Score:     r.Score,
			ZScore:    math.Abs(z),
			IsOutlier: math.Abs(z) >= zThreshold,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ZScore > results[j].ZScore
	})

	return results
}

func averageScoreFor(reports []ScoredReport) float64 {
	if len(reports) == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range reports {
		sum += r.Score
	}
	return sum / float64(len(reports))
}

func stddevScore(reports []ScoredReport, mean float64) float64 {
	if len(reports) == 0 {
		return 0
	}
	variance := 0.0
	for _, r := range reports {
		diff := r.Score - mean
		variance += diff * diff
	}
	return math.Sqrt(variance / float64(len(reports)))
}
