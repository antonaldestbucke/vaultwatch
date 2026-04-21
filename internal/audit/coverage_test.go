package audit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func sampleCoverageReports() []CompareReport {
	return []CompareReport{
		{Path: "secret/db", Envs: []string{"dev", "prod"}, OnlyInA: []string{}, OnlyInB: []string{}},
		{Path: "secret/api", Envs: []string{"dev"}, OnlyInA: []string{"key"}, OnlyInB: []string{}},
		{Path: "secret/cache", Envs: []string{"prod"}, OnlyInA: []string{}, OnlyInB: []string{"ttl"}},
	}
}

func TestBuildCoverage_AllEnvsPresent(t *testing.T) {
	reports := sampleCoverageReports()
	cov := BuildCoverage(reports)

	var db *CoverageReport
	for i := range cov {
		if cov[i].Path == "secret/db" {
			db = &cov[i]
		}
	}
	assert.NotNil(t, db)
	assert.Equal(t, 100.0, db.CoveragePct)
	assert.Empty(t, db.MissingFrom)
	assert.ElementsMatch(t, []string{"dev", "prod"}, db.PresentIn)
}

func TestBuildCoverage_PartialCoverage(t *testing.T) {
	reports := sampleCoverageReports()
	cov := BuildCoverage(reports)

	var api *CoverageReport
	for i := range cov {
		if cov[i].Path == "secret/api" {
			api = &cov[i]
		}
	}
	assert.NotNil(t, api)
	assert.Equal(t, 50.0, api.CoveragePct)
	assert.Equal(t, []string{"prod"}, api.MissingFrom)
	assert.Equal(t, []string{"dev"}, api.PresentIn)
}

func TestBuildCoverage_Empty(t *testing.T) {
	cov := BuildCoverage(nil)
	assert.Nil(t, cov)
}

func TestBuildCoverage_SortedPaths(t *testing.T) {
	reports := sampleCoverageReports()
	cov := BuildCoverage(reports)
	paths := make([]string, len(cov))
	for i, c := range cov {
		paths[i] = c.Path
	}
	assert.Equal(t, []string{"secret/api", "secret/cache", "secret/db"}, paths)
}

func TestBuildCoverage_EnvsAreSorted(t *testing.T) {
	reports := sampleCoverageReports()
	cov := BuildCoverage(reports)
	for _, c := range cov {
		assert.Equal(t, []string{"dev", "prod"}, c.Envs)
	}
}
