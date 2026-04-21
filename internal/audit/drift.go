package audit

import (
	"fmt"
	"sort"
	"strings"
)

// DriftSummary holds aggregated drift statistics across all scored reports.
type DriftSummary struct {
	TotalPaths   int
	DriftedPaths int
	CleanPaths   int
	DriftRate    float64
	TopDrifted   []string
	EnvBreakdown map[string]int
}

// BuildDriftSummary computes a high-level drift summary from a slice of ScoredReports.
func BuildDriftSummary(reports []ScoredReport, topN int) DriftSummary {
	if len(reports) == 0 {
		return DriftSummary{EnvBreakdown: map[string]int{}}
	}

	envCounts := map[string]int{}
	type pathScore struct {
		path  string
		score float64
	}
	var drifted []pathScore

	for _, r := range reports {
		if r.DriftScore > 0 {
			drifted = append(drifted, pathScore{path: r.Path, score: r.DriftScore})
			for _, env := range extractEnvs(r) {
				envCounts[env]++
			}
		}
	}

	sort.Slice(drifted, func(i, j int) bool {
		return drifted[i].score > drifted[j].score
	})

	top := []string{}
	for i, d := range drifted {
		if i >= topN {
			break
		}
		top = append(top, fmt.Sprintf("%s (%.2f)", d.path, d.score))
	}

	total := len(reports)
	driftedCount := len(drifted)
	rate := 0.0
	if total > 0 {
		rate = float64(driftedCount) / float64(total) * 100
	}

	return DriftSummary{
		TotalPaths:   total,
		DriftedPaths: driftedCount,
		CleanPaths:   total - driftedCount,
		DriftRate:    rate,
		TopDrifted:   top,
		EnvBreakdown: envCounts,
	}
}

// extractEnvs pulls environment names from a ScoredReport's diff entries.
func extractEnvs(r ScoredReport) []string {
	seen := map[string]bool{}
	for _, d := range r.Diffs {
		for _, part := range strings.SplitN(d, ":", 2) {
			part = strings.TrimSpace(part)
			if part != "" && !strings.Contains(part, " ") {
				if !seen[part] {
					seen[part] = true
				}
			}
		}
	}
	envs := make([]string, 0, len(seen))
	for e := range seen {
		envs = append(envs, e)
	}
	return envs
}
