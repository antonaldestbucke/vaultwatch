package audit

import (
	"math"
	"sort"
)

// MaturityLevel represents the assessed maturity of a secret path.
type MaturityLevel string

const (
	MaturityExemplary MaturityLevel = "exemplary"
	MaturityMature    MaturityLevel = "mature"
	MaturityDeveloping MaturityLevel = "developing"
	MaturityImmature  MaturityLevel = "immature"
)

// MaturityResult holds the maturity assessment for a single path.
type MaturityResult struct {
	Path          string
	EnvCoverage   float64
	AvgDriftScore float64
	Level         MaturityLevel
	Notes         []string
}

// BuildMaturity evaluates secret path maturity based on env coverage,
// drift score stability, and consistency across scored reports.
func BuildMaturity(scored []ScoredReport, allEnvs []string) []MaturityResult {
	if len(scored) == 0 || len(allEnvs) == 0 {
		return nil
	}

	type pathData struct {
		envsSeen map[string]bool
		scores   []float64
	}

	index := map[string]*pathData{}
	for _, r := range scored {
		if _, ok := index[r.Path]; !ok {
			index[r.Path] = &pathData{envsSeen: map[string]bool{}}
		}
		pd := index[r.Path]
		pd.envsSeen[r.Env] = true
		pd.scores = append(pd.scores, r.Score)
	}

	results := make([]MaturityResult, 0, len(index))
	for path, pd := range index {
		coverage := float64(len(pd.envsSeen)) / float64(len(allEnvs))
		avgScore := averageMaturityScore(pd.scores)
		level, notes := classifyMaturity(coverage, avgScore)
		results = append(results, MaturityResult{
			Path:          path,
			EnvCoverage:   math.Round(coverage*100) / 100,
			AvgDriftScore: math.Round(avgScore*100) / 100,
			Level:         level,
			Notes:         notes,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})
	return results
}

func averageMaturityScore(scores []float64) float64 {
	if len(scores) == 0 {
		return 0
	}
	sum := 0.0
	for _, s := range scores {
		sum += s
	}
	return sum / float64(len(scores))
}

func classifyMaturity(coverage, avgScore float64) (MaturityLevel, []string) {
	var notes []string
	if coverage < 1.0 {
		notes = append(notes, "not present in all environments")
	}
	if avgScore < 70 {
		notes = append(notes, "high drift detected")
	}
	switch {
	case coverage >= 1.0 && avgScore >= 90:
		return MaturityExemplary, notes
	case coverage >= 0.75 && avgScore >= 70:
		return MaturityMature, notes
	case coverage >= 0.5 && avgScore >= 50:
		return MaturityDeveloping, notes
	default:
		return MaturityImmature, notes
	}
}
