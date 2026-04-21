package audit

import (
	"fmt"
	"strings"
)

// PromotionAction describes a single secret key promotion between environments.
type PromotionAction struct {
	Path    string
	Key     string
	FromEnv string
	ToEnv   string
	Reason  string
}

// PromotionPlan holds all actions required to promote secrets from one env to another.
type PromotionPlan struct {
	FromEnv string
	ToEnv   string
	Actions []PromotionAction
}

// BuildPromotionPlan generates a plan to promote keys that exist in fromEnv but
// are missing in toEnv, based on scored diff reports.
func BuildPromotionPlan(reports []ScoredReport, fromEnv, toEnv string) PromotionPlan {
	plan := PromotionPlan{FromEnv: fromEnv, ToEnv: toEnv}

	for _, r := range reports {
		for _, key := range r.Report.OnlyInA {
			if strings.EqualFold(r.Report.EnvA, fromEnv) {
				plan.Actions = append(plan.Actions, PromotionAction{
					Path:    r.Report.Path,
					Key:     key,
					FromEnv: fromEnv,
					ToEnv:   toEnv,
					Reason:  "missing in target environment",
				})
			}
		}
		for _, key := range r.Report.OnlyInB {
			if strings.EqualFold(r.Report.EnvB, fromEnv) {
				plan.Actions = append(plan.Actions, PromotionAction{
					Path:    r.Report.Path,
					Key:     key,
					FromEnv: fromEnv,
					ToEnv:   toEnv,
					Reason:  "missing in target environment",
				})
			}
		}
	}

	return plan
}

// PrintPromotionPlan writes a human-readable promotion plan to stdout.
func PrintPromotionPlan(plan PromotionPlan) {
	if len(plan.Actions) == 0 {
		fmt.Printf("No promotion actions needed from %s → %s\n", plan.FromEnv, plan.ToEnv)
		return
	}

	fmt.Printf("Promotion Plan: %s → %s (%d action(s))\n", plan.FromEnv, plan.ToEnv, len(plan.Actions))
	fmt.Println(strings.Repeat("-", 60))
	for _, a := range plan.Actions {
		fmt.Printf("  [PROMOTE] %s  key=%s  (%s)\n", a.Path, a.Key, a.Reason)
	}
}
