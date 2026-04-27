package audit

import (
	"math"
	"sort"
	"strings"
)

// EntropyResult holds the Shannon entropy score for a secret path.
type EntropyResult struct {
	Path    string
	EnvA    string
	EnvB    string
	Keys    []string
	Entropy float64
	Risk    string
}

// BuildEntropy computes Shannon entropy over the union of diffed keys
// across environments, surfacing paths with high key-space disorder.
func BuildEntropy(reports []CompareReport) []EntropyResult {
	results := make([]EntropyResult, 0, len(reports))

	for _, r := range reports {
		if len(r.Diff.OnlyInA) == 0 && len(r.Diff.OnlyInB) == 0 {
			continue
		}

		allKeys := append(r.Diff.OnlyInA, r.Diff.OnlyInB...)
		e := shannonEntropy(allKeys)

		results = append(results, EntropyResult{
			Path:    r.Path,
			EnvA:    r.EnvA,
			EnvB:    r.EnvB,
			Keys:    allKeys,
			Entropy: e,
			Risk:    classifyEntropy(e),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Entropy > results[j].Entropy
	})

	return results
}

// shannonEntropy computes H(X) = -sum(p * log2(p)) over normalised key prefixes.
func shannonEntropy(keys []string) float64 {
	if len(keys) == 0 {
		return 0
	}

	freq := make(map[string]int)
	for _, k := range keys {
		parts := strings.SplitN(k, "_", 2)
		freq[parts[0]]++
	}

	total := float64(len(keys))
	var h float64
	for _, count := range freq {
		p := float64(count) / total
		h -= p * math.Log2(p)
	}
	return math.Round(h*1000) / 1000
}

// classifyEntropy maps an entropy score to a risk label.
func classifyEntropy(e float64) string {
	switch {
	case e >= 3.0:
		return "critical"
	case e >= 2.0:
		return "high"
	case e >= 1.0:
		return "medium"
	default:
		return "low"
	}
}
