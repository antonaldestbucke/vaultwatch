package audit

import "sort"

// ConsensusResult represents agreement analysis across environments for a single path.
type ConsensusResult struct {
	Path        string
	Envs        []string
	MajorityKeys []string
	DissentEnvs  []string
	AgreementPct float64
	Consensus    bool
}

// BuildConsensus analyzes scored reports and determines which paths have
// majority agreement on their key sets across environments.
func BuildConsensus(reports []ScoredReport, threshold float64) []ConsensusResult {
	if len(reports) == 0 {
		return nil
	}

	type envKeys struct {
		env  string
		keys []string
	}

	byPath := map[string][]envKeys{}
	for _, r := range reports {
		env := firstEnvFromScored(r)
		for _, rep := range r.Reports {
			byPath[rep.Path] = append(byPath[rep.Path], envKeys{env: env, keys: rep.OnlyInA})
		}
	}

	var results []ConsensusResult
	for path, entries := range byPath {
		if len(entries) == 0 {
			continue
		}

		// Count key-set frequency
		freq := map[string]int{}
		for _, e := range entries {
			sig := keySignature(e.keys)
			freq[sig]++
		}

		// Find majority signature
		var majorSig string
		maxCount := 0
		for sig, cnt := range freq {
			if cnt > maxCount {
				maxCount = cnt
				majorSig = sig
			}
		}

		agreementPct := float64(maxCount) / float64(len(entries)) * 100.0

		var majorKeys []string
		var dissentEnvs []string
		var envNames []string
		for _, e := range entries {
			envNames = append(envNames, e.env)
			if keySignature(e.keys) == majorSig {
				majorKeys = e.keys
			} else {
				dissentEnvs = append(dissentEnvs, e.env)
			}
		}

		sort.Strings(envNames)
		sort.Strings(dissentEnvs)

		results = append(results, ConsensusResult{
			Path:         path,
			Envs:         envNames,
			MajorityKeys: majorKeys,
			DissentEnvs:  dissentEnvs,
			AgreementPct: agreementPct,
			Consensus:    agreementPct >= threshold,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})
	return results
}

func keySignature(keys []string) string {
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)
	sig := ""
	for _, k := range sorted {
		sig += k + "|"
	}
	return sig
}
