package audit

import (
	"math"
	"sort"
)

// ClusterResult groups paths that share similar drift patterns.
type ClusterResult struct {
	Centroid string   `json:"centroid"`
	Paths    []string `json:"paths"`
	AvgScore float64  `json:"avg_score"`
}

// ClusterReports groups ScoredReports into clusters based on drift score proximity.
// threshold defines the max score distance to belong to the same cluster.
func ClusterReports(reports []ScoredReport, threshold float64) []ClusterResult {
	if len(reports) == 0 {
		return nil
	}

	sorted := make([]ScoredReport, len(reports))
	copy(sorted, reports)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score < sorted[j].Score
	})

	var clusters []ClusterResult
	var current []ScoredReport

	for _, r := range sorted {
		if len(current) == 0 || math.Abs(r.Score-current[0].Score) <= threshold {
			current = append(current, r)
		} else {
			clusters = append(clusters, buildCluster(current))
			current = []ScoredReport{r}
		}
	}
	if len(current) > 0 {
		clusters = append(clusters, buildCluster(current))
	}
	return clusters
}

func buildCluster(reports []ScoredReport) ClusterResult {
	var total float64
	paths := make([]string, 0, len(reports))
	for _, r := range reports {
		total += r.Score
		paths = append(paths, r.Path)
	}
	avg := total / float64(len(reports))
	return ClusterResult{
		Centroid: reports[len(reports)/2].Path,
		Paths:    paths,
		AvgScore: math.Round(avg*100) / 100,
	}
}
