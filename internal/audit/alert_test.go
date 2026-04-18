package audit

import (
	"bytes"
	"strings"
	"testing"
)

func sampleScoredForAlert() []ScoredReport {
	return []ScoredReport{
		{Path: "secret/prod/db", Score: 80.0, Risk: "high"},
		{Path: "secret/prod/api", Score: 40.0, Risk: "medium"},
		{Path: "secret/staging/db", Score: 0.0, Risk: "none"},
	}
}

func TestEvaluateAlerts_MinScore(t *testing.T) {
	reports := sampleScoredForAlert()
	alerts := EvaluateAlerts(reports, AlertRule{MinRiskScore: 50.0})
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Path != "secret/prod/db" {
		t.Errorf("unexpected path: %s", alerts[0].Path)
	}
}

func TestEvaluateAlerts_OnlyDrifted(t *testing.T) {
	reports := sampleScoredForAlert()
	alerts := EvaluateAlerts(reports, AlertRule{OnlyDrifted: true})
	for _, a := range alerts {
		if a.Risk == "none" {
			t.Errorf("expected no 'none' risk alerts, got %s", a.Path)
		}
	}
	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
}

func TestEvaluateAlerts_PathPrefix(t *testing.T) {
	reports := sampleScoredForAlert()
	alerts := EvaluateAlerts(reports, AlertRule{PathPrefix: "secret/staging"})
	if len(alerts) != 1 || alerts[0].Path != "secret/staging/db" {
		t.Errorf("unexpected alerts: %+v", alerts)
	}
}

func TestEvaluateAlerts_NoMatch(t *testing.T) {
	reports := sampleScoredForAlert()
	alerts := EvaluateAlerts(reports, AlertRule{MinRiskScore: 99.0})
	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestPrintAlerts_WithAlerts(t *testing.T) {
	alerts := []Alert{{Path: "secret/prod/db", Score: 80.0, Risk: "high", Message: "path \"secret/prod/db\" has risk \"high\" with score 80.00"}}
	var buf bytes.Buffer
	PrintAlerts(&buf, alerts)
	if !strings.Contains(buf.String(), "[HIGH]") {
		t.Errorf("expected [HIGH] in output, got: %s", buf.String())
	}
}

func TestPrintAlerts_Empty(t *testing.T) {
	var buf bytes.Buffer
	PrintAlerts(&buf, nil)
	if !strings.Contains(buf.String(), "No alerts") {
		t.Errorf("expected no-alerts message, got: %s", buf.String())
	}
}
