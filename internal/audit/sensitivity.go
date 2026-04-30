package audit

import (
	"math"
	"sort"
	"strings"
)

// SensitivityResult holds the sensitivity classification for a single path.
type SensitivityResult struct {
	Path        string  `json:"path"`
	Env         string  `json:"env"`
	Score       float64 `json:"score"`
	Label       string  `json:"label"`
	MatchedKeys []string `json:"matched_keys"`
}

var sensitivePatterns = []string{
	"password", "passwd", "secret", "token", "apikey", "api_key",
	"private_key", "cert", "credential", "auth", "key",
}

// BuildSensitivity evaluates how many keys in each report match known
// sensitive naming patterns and returns a sorted list of results.
func BuildSensitivity(reports []CompareReport) []SensitivityResult {
	var results []SensitivityResult

	for _, r := range reports {
		allKeys := append(append([]string{}, r.OnlyInA...), r.OnlyInB...)
		matched := matchSensitiveKeys(allKeys)
		if len(allKeys) == 0 {
			continue
		}
		score := math.Round((float64(len(matched))/float64(len(allKeys)))*100) / 100
		results = append(results, SensitivityResult{
			Path:        r.Path,
			Env:         firstEnvFromScored(nil),
			Score:       score,
			Label:       classifySensitivity(score),
			MatchedKeys: matched,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Path < results[j].Path
	})
	return results
}

func matchSensitiveKeys(keys []string) []string {
	var matched []string
	seen := map[string]bool{}
	for _, k := range keys {
		lower := strings.ToLower(k)
		for _, p := range sensitivePatterns {
			if strings.Contains(lower, p) && !seen[k] {
				matched = append(matched, k)
				seen[k] = true
				break
			}
		}
	}
	return matched
}

func classifySensitivity(score float64) string {
	switch {
	case score >= 0.7:
		return "critical"
	case score >= 0.4:
		return "high"
	case score >= 0.1:
		return "medium"
	default:
		return "low"
	}
}
