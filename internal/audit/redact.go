package audit

import "strings"

// RedactOptions controls which keys are redacted in reports.
type RedactOptions struct {
	KeyPatterns []string // substrings to match against key names
}

// RedactReports returns a copy of reports with matching keys replaced by "***REDACTED***".
func RedactReports(reports []CompareReport, opts RedactOptions) []CompareReport {
	result := make([]CompareReport, len(reports))
	for i, r := range reports {
		result[i] = CompareReport{
			Path:   r.Path,
			OnlyInA: redactKeys(r.OnlyInA, opts),
			OnlyInB: redactKeys(r.OnlyInB, opts),
		}
	}
	return result
}

func redactKeys(keys []string, opts RedactOptions) []string {
	if len(keys) == 0 {
		return keys
	}
	out := make([]string, len(keys))
	for i, k := range keys {
		if matchesAny(k, opts.KeyPatterns) {
			out[i] = "***REDACTED***"
		} else {
			out[i] = k
		}
	}
	return out
}

func matchesAny(key string, patterns []string) bool {
	for _, p := range patterns {
		if strings.Contains(strings.ToLower(key), strings.ToLower(p)) {
			return true
		}
	}
	return false
}
