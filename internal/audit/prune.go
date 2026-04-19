package audit

import (
	"time"
)

// PruneOptions controls what gets removed from a snapshot or report set.
type PruneOptions struct {
	OlderThan time.Duration // remove entries last seen older than this
	PathPrefix string       // only prune entries matching this prefix
	DryRun     bool
}

// PruneResult summarises what was (or would be) removed.
type PruneResult struct {
	Removed []string
	Kept    []string
}

// PruneSnapshots removes snapshot keys that are stale based on options.
// The snapshot map is keyed by path; lastSeen maps path -> timestamp.
func PruneSnapshots(snapshot map[string][]string, lastSeen map[string]time.Time, opts PruneOptions) (map[string][]string, PruneResult) {
	result := PruneResult{}
	pruned := make(map[string][]string)
	cutoff := time.Now().Add(-opts.OlderThan)

	for path, keys := range snapshot {
		if opts.PathPrefix != "" && !hasPrefix(path, opts.PathPrefix) {
			pruned[path] = keys
			result.Kept = append(result.Kept, path)
			continue
		}
		ts, ok := lastSeen[path]
		if ok && ts.Before(cutoff) {
			result.Removed = append(result.Removed, path)
			if opts.DryRun {
				pruned[path] = keys
			}
		} else {
			pruned[path] = keys
			result.Kept = append(result.Kept, path)
		}
	}
	return pruned, result
}

func hasPrefix(s, prefix string) bool {
	if len(prefix) == 0 {
		return true
	}
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
