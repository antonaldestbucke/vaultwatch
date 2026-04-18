package audit

import "fmt"

// RiskLevel represents the severity of drift detected.
type RiskLevel string

const (
	RiskNone   RiskLevel = "none"
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

// ScoreResult holds the drift score and risk level for a set of reports.
type ScoreResult struct {
	TotalPaths  int
	DriftedPaths int
	Score       float64
	Risk        RiskLevel
}

// String returns a human-readable summary of the score result.
func (s ScoreResult) String() string {
	return fmt.Sprintf("Drift Score: %.1f%% (%d/%d paths drifted) — Risk: %s",
		s.Score, s.DriftedPaths, s.TotalPaths, s.Risk)
}

// ScoreReports calculates a drift score from a slice of CompareReport.
// Score is the percentage of paths with at least one diff.
func ScoreReports(reports []CompareReport) ScoreResult {
	if len(reports) == 0 {
		return ScoreResult{Risk: RiskNone}
	}

	drifted := 0
	for _, r := range reports {
		if hasDiff(r) {
			drifted++
		}
	}

	score := float64(drifted) / float64(len(reports)) * 100
	risk := classifyRisk(score)

	return ScoreResult{
		TotalPaths:   len(reports),
		DriftedPaths: drifted,
		Score:        score,
		Risk:         risk,
	}
}

func classifyRisk(score float64) RiskLevel {
	switch {
	case score == 0:
		return RiskNone
	case score < 25:
		return RiskLow
	case score < 60:
		return RiskMedium
	default:
		return RiskHigh
	}
}
