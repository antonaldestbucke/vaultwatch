package audit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func sampleImpactReports() []ScoredReport {
	return []ScoredReport{
		{
			Score: 40,
			Risk:  "high",
			Report: CompareReport{
				Path:    "secret/app/db",
				Envs:    []string{"prod", "staging", "dev"},
				OnlyInA: []string{"password", "host", "port", "user", "ssl"},
				OnlyInB: []string{},
			},
		},
		{
			Score: 70,
			Risk:  "medium",
			Report: CompareReport{
				Path:    "secret/app/api",
				Envs:    []string{"prod", "staging"},
				OnlyInA: []string{"key"},
				OnlyInB: []string{"token"},
			},
		},
		{
			Score: 100,
			Risk:  "low",
			Report: CompareReport{
				Path:    "secret/app/cache",
				Envs:    []string{"prod", "staging"},
				OnlyInA: []string{},
				OnlyInB: []string{},
			},
		},
	}
}

func TestBuildImpact_ExcludesNoDrift(t *testing.T) {
	reports := sampleImpactReports()
	summary := BuildImpact(reports)
	assert.Equal(t, 2, summary.Total)
}

func TestBuildImpact_HighCountCorrect(t *testing.T) {
	reports := sampleImpactReports()
	summary := BuildImpact(reports)
	assert.Equal(t, 1, summary.HighCount)
}

func TestBuildImpact_SortedByWeight(t *testing.T) {
	reports := sampleImpactReports()
	summary := BuildImpact(reports)
	assert.Equal(t, ImpactHigh, summary.Results[0].Level)
	assert.Equal(t, ImpactMedium, summary.Results[1].Level)
}

func TestBuildImpact_Empty(t *testing.T) {
	summary := BuildImpact([]ScoredReport{})
	assert.Equal(t, 0, summary.Total)
	assert.Equal(t, 0, summary.HighCount)
}

func TestClassifyImpact_Thresholds(t *testing.T) {
	assert.Equal(t, ImpactHigh, classifyImpact(5, 1))
	assert.Equal(t, ImpactHigh, classifyImpact(1, 3))
	assert.Equal(t, ImpactMedium, classifyImpact(2, 1))
	assert.Equal(t, ImpactMedium, classifyImpact(1, 2))
	assert.Equal(t, ImpactLow, classifyImpact(1, 1))
}
