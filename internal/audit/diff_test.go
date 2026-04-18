package audit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffKeys_NoChanges(t *testing.T) {
	a := []string{"foo", "bar", "baz"}
	b := []string{"foo", "bar", "baz"}

	result := DiffKeys("secret/app", a, b)

	assert.False(t, result.HasDiff())
	assert.ElementsMatch(t, []string{"bar", "baz", "foo"}, result.InBoth)
	assert.Empty(t, result.OnlyInA)
	assert.Empty(t, result.OnlyInB)
}

func TestDiffKeys_OnlyInA(t *testing.T) {
	a := []string{"foo", "extra"}
	b := []string{"foo"}

	result := DiffKeys("secret/app", a, b)

	assert.True(t, result.HasDiff())
	assert.Equal(t, []string{"extra"}, result.OnlyInA)
	assert.Empty(t, result.OnlyInB)
}

func TestDiffKeys_OnlyInB(t *testing.T) {
	a := []string{"foo"}
	b := []string{"foo", "new-key"}

	result := DiffKeys("secret/app", a, b)

	assert.True(t, result.HasDiff())
	assert.Empty(t, result.OnlyInA)
	assert.Equal(t, []string{"new-key"}, result.OnlyInB)
}

func TestDiffKeys_DisjointSets(t *testing.T) {
	a := []string{"alpha", "beta"}
	b := []string{"gamma", "delta"}

	result := DiffKeys("secret/svc", a, b)

	assert.True(t, result.HasDiff())
	assert.Equal(t, []string{"alpha", "beta"}, result.OnlyInA)
	assert.Equal(t, []string{"delta", "gamma"}, result.OnlyInB)
	assert.Empty(t, result.InBoth)
}

func TestDiffKeys_PathIsPreserved(t *testing.T) {
	result := DiffKeys("secret/mypath", []string{}, []string{})
	assert.Equal(t, "secret/mypath", result.Path)
}
