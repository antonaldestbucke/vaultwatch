package audit

import "sort"

// SecretKeys represents a set of secret keys at a given path.
type SecretKeys map[string]struct{}

// DiffResult holds the result of comparing secrets between two environments.
type DiffResult struct {
	Path       string
	OnlyInA    []string
	OnlyInB    []string
	InBoth     []string
}

// DiffKeys compares two sets of secret keys and returns a DiffResult.
func DiffKeys(path string, a, b []string) DiffResult {
	aSet := toSet(a)
	bSet := toSet(b)

	result := DiffResult{Path: path}

	for k := range aSet {
		if _, ok := bSet[k]; ok {
			result.InBoth = append(result.InBoth, k)
		} else {
			result.OnlyInA = append(result.OnlyInA, k)
		}
	}

	for k := range bSet {
		if _, ok := aSet[k]; !ok {
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.InBoth)

	return result
}

// HasDiff returns true if there are any keys present in only one of the two sets.
func (d DiffResult) HasDiff() bool {
	return len(d.OnlyInA) > 0 || len(d.OnlyInB) > 0
}

func toSet(keys []string) SecretKeys {
	s := make(SecretKeys, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
