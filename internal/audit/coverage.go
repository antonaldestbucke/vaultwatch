package audit

import "sort"

// CoverageReport holds coverage stats for a single secret path.
type CoverageReport struct {
	Path        string
	Envs        []string
	PresentIn   []string
	MissingFrom []string
	CoveragePct float64
}

// BuildCoverage computes how many environments each path is present in
// relative to the full set of environments seen across all reports.
func BuildCoverage(reports []CompareReport) []CoverageReport {
	if len(reports) == 0 {
		return nil
	}

	// Collect all environments seen.
	envSet := map[string]struct{}{}
	for _, r := range reports {
		for _, e := range r.Envs {
			envSet[e] = struct{}{}
		}
	}
	allEnvs := make([]string, 0, len(envSet))
	for e := range envSet {
		allEnvs = append(allEnvs, e)
	}
	sort.Strings(allEnvs)

	// Group by path.
	pathEnvs := map[string]map[string]struct{}{}
	for _, r := range reports {
		if _, ok := pathEnvs[r.Path]; !ok {
			pathEnvs[r.Path] = map[string]struct{}{}
		}
		for _, e := range r.Envs {
			pathEnvs[r.Path][e] = struct{}{}
		}
	}

	paths := make([]string, 0, len(pathEnvs))
	for p := range pathEnvs {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	result := make([]CoverageReport, 0, len(paths))
	for _, p := range paths {
		present := []string{}
		missing := []string{}
		for _, e := range allEnvs {
			if _, ok := pathEnvs[p][e]; ok {
				present = append(present, e)
			} else {
				missing = append(missing, e)
			}
		}
		pct := float64(len(present)) / float64(len(allEnvs)) * 100.0
		result = append(result, CoverageReport{
			Path:        p,
			Envs:        allEnvs,
			PresentIn:   present,
			MissingFrom: missing,
			CoveragePct: pct,
		})
	}
	return result
}
