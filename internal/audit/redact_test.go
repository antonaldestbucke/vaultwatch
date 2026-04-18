package audit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func sampleRedactReports() []CompareReport {
	return []CompareReport{
		{
			Path:    "secret/app",
			OnlyInA: []string{"password", "username", "api_key"},
			OnlyInB: []string{"token", "host"},
		},
	}
}

func TestRedactReports_MatchingKeys(t *testing.T) {
	reports := sampleRedactReports()
	opts := RedactOptions{KeyPatterns: []string{"password", "token", "key"}}
	result := RedactReports(reports, opts)

	assert.Equal(t, "***REDACTED***", result[0].OnlyInA[0]) // password
	assert.Equal(t, "username", result[0].OnlyInA[1])        // not matched
	assert.Equal(t, "***REDACTED***", result[0].OnlyInA[2]) // api_key
	assert.Equal(t, "***REDACTED***", result[0].OnlyInB[0]) // token
	assert.Equal(t, "host", result[0].OnlyInB[1])            // not matched
}

func TestRedactReports_NoPatterns(t *testing.T) {
	reports := sampleRedactReports()
	opts := RedactOptions{}
	result := RedactReports(reports, opts)

	assert.Equal(t, reports[0].OnlyInA, result[0].OnlyInA)
	assert.Equal(t, reports[0].OnlyInB, result[0].OnlyInB)
}

func TestRedactReports_EmptyKeys(t *testing.T) {
	reports := []CompareReport{{Path: "secret/empty", OnlyInA: nil, OnlyInB: nil}}
	opts := RedactOptions{KeyPatterns: []string{"password"}}
	result := RedactReports(reports, opts)

	assert.Nil(t, result[0].OnlyInA)
	assert.Nil(t, result[0].OnlyInB)
}

func TestRedactReports_CaseInsensitive(t *testing.T) {
	reports := []CompareReport{
		{Path: "secret/app", OnlyInA: []string{"PASSWORD", "Secret_Token"}, OnlyInB: []string{}},
	}
	opts := RedactOptions{KeyPatterns: []string{"password", "token"}}
	result := RedactReports(reports, opts)

	assert.Equal(t, "***REDACTED***", result[0].OnlyInA[0])
	assert.Equal(t, "***REDACTED***", result[0].OnlyInA[1])
}
