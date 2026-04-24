package audit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func sampleFilterReports() []DiffReport {
	return []DiffReport{
		{Path: "secret/app/db", OnlyInA: []string{"password"}, OnlyInB: []string{}},
		{Path: "secret/app/api", OnlyInA: []string{}, OnlyInB: []string{}},
		{Path: "secret/infra/tls", OnlyInA: []string{"cert"}, OnlyInB: []string{"key"}},
		{Path: "config/feature", OnlyInA: []string{}, OnlyInB: []string{}},
	}
}

func TestFilterReports_OnlyDiffs(t *testing.T) {
	result := FilterReports(sampleFilterReports(), FilterOptions{OnlyDiffs: true})
	assert.Len(t, result, 2)
	assert.Equal(t, "secret/app/db", result[0].Path)
	assert.Equal(t, "secret/infra/tls", result[1].Path)
}

func TestFilterReports_PathPrefix(t *testing.T) {
	result := FilterReports(sampleFilterReports(), FilterOptions{PathPrefix: "secret/app"})
	assert.Len(t, result, 2)
}

func TestFilterReports_CombinedFilters(t *testing.T) {
	result := FilterReports(sampleFilterReports(), FilterOptions{OnlyDiffs: true, PathPrefix: "secret/app"})
	assert.Len(t, result, 1)
	assert.Equal(t, "secret/app/db", result[0].Path)
}

func TestFilterReports_NoMatch(t *testing.T) {
	result := FilterReports(sampleFilterReports(), FilterOptions{PathPrefix: "nonexistent/"})
	assert.Empty(t, result)
}

func TestFilterReports_NoOptions(t *testing.T) {
	result := FilterReports(sampleFilterReports(), FilterOptions{})
	assert.Len(t, result, 4)
}

func TestFilterReports_EmptyInput(t *testing.T) {
	result := FilterReports([]DiffReport{}, FilterOptions{OnlyDiffs: true, PathPrefix: "secret/"})
	assert.Empty(t, result)
}
