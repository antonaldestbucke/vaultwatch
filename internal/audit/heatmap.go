package audit

import (
	"fmt"
	"sort"
	"strings"
)

// HeatmapEntry represents drift frequency for a single path.
type HeatmapEntry struct {
	Path       string
	DriftCount int
	Envs       []string
}

// Heatmap holds a ranked list of paths by drift frequency.
type Heatmap struct {
	Entries []HeatmapEntry
	Total   int
}

// BuildHeatmap aggregates drift counts per path across scored reports.
func BuildHeatmap(reports []ScoredReport) Heatmap {
	type key struct{ path string }
	counts := make(map[string]int)
	envSet := make(map[string]map[string]struct{})

	for _, r := range reports {
		if r.Report.OnlyInA == nil && r.Report.OnlyInB == nil {
			continue
		}
		if len(r.Report.OnlyInA)+len(r.Report.OnlyInB) == 0 {
			continue
		}
		counts[r.Report.Path]++
		if envSet[r.Report.Path] == nil {
			envSet[r.Report.Path] = make(map[string]struct{})
		}
		envSet[r.Report.Path][r.Report.EnvA] = struct{}{}
		envSet[r.Report.Path][r.Report.EnvB] = struct{}{}
	}

	entries := make([]HeatmapEntry, 0, len(counts))
	for path, count := range counts {
		envs := make([]string, 0, len(envSet[path]))
		for e := range envSet[path] {
			envs = append(envs, e)
		}
		sort.Strings(envs)
		entries = append(entries, HeatmapEntry{Path: path, DriftCount: count, Envs: envs})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].DriftCount != entries[j].DriftCount {
			return entries[i].DriftCount > entries[j].DriftCount
		}
		return entries[i].Path < entries[j].Path
	})

	return Heatmap{Entries: entries, Total: len(reports)}
}

// PrintHeatmap writes a text heatmap table to stdout.
func PrintHeatmap(h Heatmap) {
	if len(h.Entries) == 0 {
		fmt.Println("No drift detected across reports.")
		return
	}
	fmt.Printf("%-50s %10s  %s\n", "PATH", "DRIFT_COUNT", "ENVS")
	fmt.Println(strings.Repeat("-", 80))
	for _, e := range h.Entries {
		fmt.Printf("%-50s %10d  %s\n", e.Path, e.DriftCount, strings.Join(e.Envs, ","))
	}
	fmt.Printf("\nTotal reports analysed: %d\n", h.Total)
}
