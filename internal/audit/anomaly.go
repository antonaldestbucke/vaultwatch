package audit

import "math"

// AnomalyResult holds the result of an anomaly detection check for a single path.
type AnomalyResult struct {
	Path        string
	Env         string
	Score       float64
	MeanScore   float64
	StdDev      float64
	ZScore      float64
	IsAnomaly   bool
	Severity    string
}

// BuildAnomalies detects paths whose drift score deviates significantly from
// the mean across all scored reports, using a configurable z-score threshold.
func BuildAnomalies(reports []ScoredReport, zThreshold float64) []AnomalyResult {
	if len(reports) == 0 {
		return nil
	}

	scores := make([]float64, len(reports))
	for i, r := range reports {
		scores[i] = r.Score
	}

	mean := meanFloat(scores)
	std := stddevFloat(scores, mean)

	var results []AnomalyResult
	for _, r := range reports {
		z := 0.0
		if std > 0 {
			z = (r.Score - mean) / std
		}
		isAnomaly := math.Abs(z) >= zThreshold
		severity := classifyAnomalySeverity(math.Abs(z))
		results = append(results, AnomalyResult{
			Path:      r.Path,
			Env:       firstEnvFromScored(r),
			Score:     r.Score,
			MeanScore: mean,
			StdDev:    std,
			ZScore:    z,
			IsAnomaly: isAnomaly,
			Severity:  severity,
		})
	}
	return results
}

func classifyAnomalySeverity(absZ float64) string {
	switch {
	case absZ >= 3.0:
		return "critical"
	case absZ >= 2.0:
		return "high"
	case absZ >= 1.0:
		return "medium"
	default:
		return "low"
	}
}

func meanFloat(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func stddevFloat(vals []float64, mean float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		diff := v - mean
		sum += diff * diff
	}
	return math.Sqrt(sum / float64(len(vals)))
}
