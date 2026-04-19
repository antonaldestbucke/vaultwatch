package audit

import (
	"fmt"
	"strings"
)

// RemediationAction describes a suggested fix for a drifted secret path.
type RemediationAction struct {
	Path        string `json:"path"`
	Environment string `json:"environment"`
	Action      string `json:"action"`
	Detail      string `json:"detail"`
}

// RemediationPlan holds all suggested actions for a set of reports.
type RemediationPlan struct {
	Actions []RemediationAction `json:"actions"`
}

// BuildRemediationPlan generates remediation suggestions from scored reports.
func BuildRemediationPlan(reports []ScoredReport) RemediationPlan {
	plan := RemediationPlan{}
	for _, r := range reports {
		for _, diff := range r.Report.Diffs {
			for _, key := range diff.OnlyInA {
				plan.Actions = append(plan.Actions, RemediationAction{
					Path:        diff.Path,
					Environment: secondEnv(r.Report.Environments),
					Action:      "add_key",
					Detail:      fmt.Sprintf("Key %q is missing; add it to match source environment", key),
				})
			}
			for _, key := range diff.OnlyInB {
				plan.Actions = append(plan.Actions, RemediationAction{
					Path:        diff.Path,
					Environment: firstEnv(r.Report.Environments),
					Action:      "add_key",
					Detail:      fmt.Sprintf("Key %q is missing; add it to match target environment", key),
				})
			}
		}
	}
	return plan
}

// PrintRemediationPlan writes the plan to stdout in human-readable form.
func PrintRemediationPlan(plan RemediationPlan) {
	if len(plan.Actions) == 0 {
		fmt.Println("No remediation actions required.")
		return
	}
	fmt.Printf("Remediation Plan (%d action(s)):\n", len(plan.Actions))
	fmt.Println(strings.Repeat("-", 60))
	for _, a := range plan.Actions {
		fmt.Printf("[%s] %s (env: %s)\n  -> %s\n", a.Action, a.Path, a.Environment, a.Detail)
	}
}

func firstEnv(envs []string) string {
	if len(envs) > 0 {
		return envs[0]
	}
	return "unknown"
}

func secondEnv(envs []string) string {
	if len(envs) > 1 {
		return envs[1]
	}
	return "unknown"
}
