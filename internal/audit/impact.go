package audit

import "sort"

// ImpactLevel categorizes how broadly a drift affects environments.
type ImpactLevel string

const (
	ImpactLow    ImpactLevel = "low"
	ImpactMedium ImpactLevel = "medium"
	ImpactHigh   ImpactLevel = "high"
)

// ImpactResult holds the computed impact for a single secret path.
type ImpactResult struct {
	Path         string      `json:"path"`
	AffectedEnvs []string    `json:"affected_envs"`
	DriftedKeys  int         `json:"drifted_keys"`
	Level        ImpactLevel `json:"level"`
}

// ImpactSummary holds all impact results sorted by severity.
type ImpactSummary struct {
	Results []ImpactResult `json:"results"`
	Total   int            `json:"total"`
	HighCount int          `json:"high_count"`
}

// BuildImpact computes an impact assessment from scored reports.
func BuildImpact(reports []ScoredReport) ImpactSummary {
	var results []ImpactResult

	for _, r := range reports {
		if len(r.Report.OnlyInA) == 0 && len(r.Report.OnlyInB) == 0 {
			continue
		}
		drifted := len(r.Report.OnlyInA) + len(r.Report.OnlyInB)
		envs := collectEnvs(r)
		level := classifyImpact(drifted, len(envs))
		results = append(results, ImpactResult{
			Path:         r.Report.Path,
			AffectedEnvs: envs,
			DriftedKeys:  drifted,
			Level:        level,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return impactWeight(results[i].Level) > impactWeight(results[j].Level)
	})

	high := 0
	for _, r := range results {
		if r.Level == ImpactHigh {
			high++
		}
	}

	return ImpactSummary{Results: results, Total: len(results), HighCount: high}
}

func collectEnvs(r ScoredReport) []string {
	seen := map[string]struct{}{}
	for _, e := range r.Report.Envs {
		seen[e] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for e := range seen {
		out = append(out, e)
	}
	sort.Strings(out)
	return out
}

func classifyImpact(driftedKeys, envCount int) ImpactLevel {
	switch {
	case driftedKeys >= 5 || envCount >= 3:
		return ImpactHigh
	case driftedKeys >= 2 || envCount >= 2:
		return ImpactMedium
	default:
		return ImpactLow
	}
}

func impactWeight(l ImpactLevel) int {
	switch l {
	case ImpactHigh:
		return 3
	case ImpactMedium:
		return 2
	default:
		return 1
	}
}
