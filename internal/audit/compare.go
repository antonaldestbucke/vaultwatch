package audit

import (
	"context"
	"fmt"

	"vaultwatch/internal/vault"
)

// CompareReport holds the diff result for a single path across two environments.
type CompareReport struct {
	Path    string
	EnvA    string
	EnvB    string
	OnlyInA []string
	OnlyInB []string
}

// ComparePathAcrossEnvs lists keys at path for each client pair and diffs them.
func ComparePathAcrossEnvs(clients map[string]*vault.Client, path string) ([]CompareReport, error) {
	envs := make([]string, 0, len(clients))
	for env := range clients {
		envs = append(envs, env)
	}

	if len(envs) < 2 {
		return nil, fmt.Errorf("at least 2 environments required for comparison")
	}

	ctx := context.Background()
	keyMap := make(map[string][]string)
	for _, env := range envs {
		keys, err := clients[env].ListSecrets(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("list secrets for env %q: %w", env, err)
		}
		keyMap[env] = keys
	}

	var reports []CompareReport
	for i := 0; i < len(envs)-1; i++ {
		for j := i + 1; j < len(envs); j++ {
			a, b := envs[i], envs[j]
			diff := DiffKeys(path, keyMap[a], keyMap[b])
			reports = append(reports, CompareReport{
				Path:    path,
				EnvA:    a,
				EnvB:    b,
				OnlyInA: diff.OnlyInA,
				OnlyInB: diff.OnlyInB,
			})
		}
	}
	return reports, nil
}
