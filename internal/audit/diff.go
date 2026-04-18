package audit

// DiffResult holds the keys unique to each side and the path compared.
type DiffResult struct {
	Path    string
	OnlyInA []string
	OnlyInB []string
}

// DiffKeys compares two slices of keys and returns a DiffResult.
func DiffKeys(path string, a, b []string) DiffResult {
	setA := toSet(a)
	setB := toSet(b)

	var onlyA, onlyB []string
	for k := range setA {
		if !setB[k] {
			onlyA = append(onlyA, k)
		}
	}
	for k := range setB {
		if !setA[k] {
			onlyB = append(onlyB, k)
		}
	}
	return DiffResult{Path: path, OnlyInA: onlyA, OnlyInB: onlyB}
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
