package audit

import "fmt"

// DiffResult holds keys present in only one of the two compared sets.
type DiffResult struct {
	Path    string
	OnlyInA []string
	OnlyInB []string
}

// HasDiff returns true when there is at least one key difference.
func (d DiffResult) HasDiff() bool {
	return len(d.OnlyInA) > 0 || len(d.OnlyInB) > 0
}

// DiffKeys compares two slices of secret keys rooted at path and returns a DiffResult.
func DiffKeys(path string, keysA, keysB []string) DiffResult {
	setA := toSet(keysA)
	setB := toSet(keysB)

	result := DiffResult{Path: path}

	for k := range setA {
		if !setB[k] {
			result.OnlyInA = append(result.OnlyInA, fmt.Sprintf("%s/%s", path, k))
		}
	}
	for k := range setB {
		if !setA[k] {
			result.OnlyInB = append(result.OnlyInB, fmt.Sprintf("%s/%s", path, k))
		}
	}

	return result
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
