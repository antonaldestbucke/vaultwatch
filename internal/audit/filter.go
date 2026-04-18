package audit

import "strings"

// FilterOptions controls which report entries are included.
type FilterOptions struct {
	OnlyDiffs  bool
	PathPrefix string
}

// FilterReports returns a filtered slice of DiffReport based on the given options.
func FilterReports(reports []DiffReport, opts FilterOptions) []DiffReport {
	var result []DiffReport
	for _, r := range reports {
		if opts.PathPrefix != "" && !strings.HasPrefix(r.Path, opts.PathPrefix) {
			continue
		}
		if opts.OnlyDiffs && !hasDiff(r) {
			continue
		}
		result = append(result, r)
	}
	return result
}

func hasDiff(r DiffReport) bool {
	return len(r.OnlyInA) > 0 || len(r.OnlyInB) > 0
}
