package audit

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ExportFormat defines the output format for reports.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
)

// ExportReport writes the given reports to w in the specified format.
func ExportReport(w io.Writer, reports []PathReport, format ExportFormat) error {
	switch format {
	case FormatJSON:
		return exportJSON(w, reports)
	case FormatCSV:
		return exportCSV(w, reports)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

type jsonExport struct {
	GeneratedAt time.Time    `json:"generated_at"`
	Reports     []PathReport `json:"reports"`
}

func exportJSON(w io.Writer, reports []PathReport) error {
	payload := jsonExport{
		GeneratedAt: time.Now().UTC(),
		Reports:     reports,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func exportCSV(w io.Writer, reports []PathReport) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"path", "only_in_a", "only_in_b", "env_a", "env_b"}); err != nil {
		return err
	}
	for _, r := range reports {
		row := []string{
			r.Path,
			fmt.Sprintf("%v", r.Diff.OnlyInA),
			fmt.Sprintf("%v", r.Diff.OnlyInB),
			r.EnvA,
			r.EnvB,
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
