package audit

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

// ReportFormat defines the output format for reports.
type ReportFormat string

const (
	FormatText ReportFormat = "text"
	FormatJSON ReportFormat = "json"
)

// PathReport holds the diff result for a single secret path.
type PathReport struct {
	Path    string
	EnvA    string
	EnvB    string
	Diff    DiffResult
}

// DiffResult mirrors the output of DiffKeys.
type DiffResult struct {
	OnlyInA []string
	OnlyInB []string
}

// PrintTextReport writes a human-readable diff report to w.
func PrintTextReport(w io.Writer, reports []PathReport) {
	for _, r := range reports {
		if len(r.Diff.OnlyInA) == 0 && len(r.Diff.OnlyInB) == 0 {
			fmt.Fprintf(w, "[OK] %s — no differences\n", r.Path)
			continue
		}
		fmt.Fprintf(w, "[DIFF] %s\n", r.Path)
		for _, k := range r.Diff.OnlyInA {
			color.New(color.FgRed).Fprintf(w, "  - [%s only] %s\n", r.EnvA, k)
		}
		for _, k := range r.Diff.OnlyInB {
			color.New(color.FgGreen).Fprintf(w, "  + [%s only] %s\n", r.EnvB, k)
		}
	}
}

// Summary returns a one-line summary string.
func Summary(reports []PathReport) string {
	total := len(reports)
	diffCount := 0
	for _, r := range reports {
		if len(r.Diff.OnlyInA) > 0 || len(r.Diff.OnlyInB) > 0 {
			diffCount++
		}
	}
	return strings.TrimSpace(fmt.Sprintf("%d path(s) checked, %d with differences", total, diffCount))
}
