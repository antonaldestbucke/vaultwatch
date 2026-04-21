package audit

import (
	"math"
	"sort"
	"strings"
)

// SimilarityResult holds the similarity score between two paths across environments.
type SimilarityResult struct {
	Path       string
	EnvA       string
	EnvB       string
	ScoreA     []string
	ScoreB     []string
	Jaccard    float64
	Similarity string // "high", "medium", "low"
}

// ComputeSimilarity calculates Jaccard similarity between key sets of two environments
// for each report path and classifies the result.
func ComputeSimilarity(reports []DiffReport, envA, envB string) []SimilarityResult {
	results := make([]SimilarityResult, 0, len(reports))

	for _, r := range reports {
		setA := toStringSet(r.OnlyInA)
		setB := toStringSet(r.OnlyInB)
		all := unionSets(setA, setB)

		// keys present in both = total - only in A - only in B
		totalUnique := len(all)
		if totalUnique == 0 {
			// identical empty sets — perfect similarity
			results = append(results, SimilarityResult{
				Path:       r.Path,
				EnvA:       envA,
				EnvB:       envB,
				ScoreA:     r.OnlyInA,
				ScoreB:     r.OnlyInB,
				Jaccard:    1.0,
				Similarity: "high",
			})
			continue
		}

		intersection := totalUnique - len(setA) - len(setB)
		if intersection < 0 {
			intersection = 0
		}
		jaccard := float64(intersection) / float64(totalUnique)
		jaccard = math.Round(jaccard*100) / 100

		results = append(results, SimilarityResult{
			Path:       r.Path,
			EnvA:       envA,
			EnvB:       envB,
			ScoreA:     r.OnlyInA,
			ScoreB:     r.OnlyInB,
			Jaccard:    jaccard,
			Similarity: classifySimilarity(jaccard),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Jaccard < results[j].Jaccard
	})

	return results
}

func classifySimilarity(score float64) string {
	switch {
	case score >= 0.8:
		return "high"
	case score >= 0.5:
		return "medium"
	default:
		return "low"
	}
}

func toStringSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[strings.ToLower(k)] = struct{}{}
	}
	return s
}

func unionSets(a, b map[string]struct{}) map[string]struct{} {
	union := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		union[k] = struct{}{}
	}
	for k := range b {
		union[k] = struct{}{}
	}
	return union
}
