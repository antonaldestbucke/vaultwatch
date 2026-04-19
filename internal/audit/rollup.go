package audit

import (
	"fmt"
	"sort"
	"strings"
)

// RollupEntry summarizes drift across a path prefix.
type RollupEntry struct {
	Prefix      string
	TotalPaths  int
	DriftedPaths int
	DriftRate   float64
}

// RollupReport holds all rollup entries for a given depth.
type RollupReport struct {
	Depth   int
	Entries []RollupEntry
}

// BuildRollup aggregates ScoredReport results by path prefix at the given depth.
func BuildRollup(reports []ScoredReport, depth int) RollupReport {
	type bucket struct {
		total   int
		drifted int
	}
	buckets := map[string]*bucket{}

	for _, r := range reports {
		prefix := prefixAtDepth(r.Path, depth)
		if _, ok := buckets[prefix]; !ok {
			buckets[prefix] = &bucket{}
		}
		buckets[prefix].total++
		if r.Risk == "high" || r.Risk == "medium" {
			buckets[prefix].drifted++
		}
	}

	entries := make([]RollupEntry, 0, len(buckets))
	for prefix, b := range buckets {
		rate := 0.0
		if b.total > 0 {
			rate = float64(b.drifted) / float64(b.total) * 100
		}
		entries = append(entries, RollupEntry{
			Prefix:       prefix,
			TotalPaths:   b.total,
			DriftedPaths: b.drifted,
			DriftRate:    rate,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Prefix < entries[j].Prefix
	})

	return RollupReport{Depth: depth, Entries: entries}
}

// PrintRollup writes a rollup report to stdout.
func PrintRollup(r RollupReport) {
	fmt.Printf("Rollup at depth %d:\n", r.Depth)
	fmt.Printf("%-40s %8s %8s %10s\n", "Prefix", "Total", "Drifted", "DriftRate")
	fmt.Println(strings.Repeat("-", 70))
	for _, e := range r.Entries {
		fmt.Printf("%-40s %8d %8d %9.1f%%\n", e.Prefix, e.TotalPaths, e.DriftedPaths, e.DriftRate)
	}
}

func prefixAtDepth(path string, depth int) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if depth <= 0 || depth >= len(parts) {
		return path
	}
	return "/" + strings.Join(parts[:depth], "/")
}
