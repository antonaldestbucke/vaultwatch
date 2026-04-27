package audit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func sampleConsensusReports() []ScoredReport {
	return []ScoredReport{
		{
			Environment: "prod",
			Score:       90,
			Reports: []DiffReport{
				{Path: "secret/app", OnlyInA: []string{"db_pass", "api_key"}},
			},
		},
		{
			Environment: "staging",
			Score:       90,
			Reports: []DiffReport{
				{Path: "secret/app", OnlyInA: []string{"db_pass", "api_key"}},
			},
		},
		{
			Environment: "dev",
			Score:       60,
			Reports: []DiffReport{
				{Path: "secret/app", OnlyInA: []string{"db_pass"}},
			},
		},
	}
}

func TestBuildConsensus_MajorityAgrees(t *testing.T) {
	reports := sampleConsensusReports()
	results := BuildConsensus(reports, 60.0)
	assert.Len(t, results, 1)
	assert.Equal(t, "secret/app", results[0].Path)
	assert.InDelta(t, 66.6, results[0].AgreementPct, 1.0)
	assert.True(t, results[0].Consensus)
}

func TestBuildConsensus_DissentEnvIdentified(t *testing.T) {
	reports := sampleConsensusReports()
	results := BuildConsensus(reports, 60.0)
	assert.Len(t, results[0].DissentEnvs, 1)
	assert.Equal(t, "dev", results[0].DissentEnvs[0])
}

func TestBuildConsensus_BelowThreshold(t *testing.T) {
	reports := sampleConsensusReports()
	results := BuildConsensus(reports, 90.0)
	assert.Len(t, results, 1)
	assert.False(t, results[0].Consensus)
}

func TestBuildConsensus_Empty(t *testing.T) {
	results := BuildConsensus(nil, 75.0)
	assert.Nil(t, results)
}

func TestBuildConsensus_AllAgree(t *testing.T) {
	reports := []ScoredReport{
		{Environment: "prod", Reports: []DiffReport{{Path: "secret/db", OnlyInA: []string{"pass"}}}},
		{Environment: "staging", Reports: []DiffReport{{Path: "secret/db", OnlyInA: []string{"pass"}}}},
	}
	results := BuildConsensus(reports, 100.0)
	assert.Len(t, results, 1)
	assert.True(t, results[0].Consensus)
	assert.Empty(t, results[0].DissentEnvs)
	assert.InDelta(t, 100.0, results[0].AgreementPct, 0.01)
}
