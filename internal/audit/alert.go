package audit

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// AlertRule defines a condition that triggers an alert.
type AlertRule struct {
	MinRiskScore float64
	PathPrefix   string
	OnlyDrifted  bool
}

// Alert represents a triggered alert for a scored report.
type Alert struct {
	Path    string
	Score   float64
	Risk    string
	Message string
}

// EvaluateAlerts checks scored reports against rules and returns triggered alerts.
func EvaluateAlerts(reports []ScoredReport, rule AlertRule) []Alert {
	var alerts []Alert
	for _, r := range reports {
		if rule.PathPrefix != "" && !strings.HasPrefix(r.Path, rule.PathPrefix) {
			continue
		}
		if rule.OnlyDrifted && r.Risk == "none" {
			continue
		}
		if r.Score < rule.MinRiskScore {
			continue
		}
		alerts = append(alerts, Alert{
			Path:    r.Path,
			Score:   r.Score,
			Risk:    r.Risk,
			Message: fmt.Sprintf("path %q has risk %q with score %.2f", r.Path, r.Risk, r.Score),
		})
	}
	return alerts
}

// PrintAlerts writes alerts to the given writer.
func PrintAlerts(w io.Writer, alerts []Alert) {
	if len(alerts) == 0 {
		fmt.Fprintln(w, "No alerts triggered.")
		return
	}
	fmt.Fprintf(w, "%d alert(s) triggered:\n", len(alerts))
	for _, a := range alerts {
		fmt.Fprintf(w, "  [%s] %s\n", strings.ToUpper(a.Risk), a.Message)
	}
}

// PrintAlertsToStdout is a convenience wrapper.
func PrintAlertsToStdout(alerts []Alert) {
	PrintAlerts(os.Stdout, alerts)
}
