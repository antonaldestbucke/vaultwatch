package audit

import (
	"context"
	"fmt"

	"github.com/vaultwatch/internal/vault"
)

// EnvSecrets holds the secret keys found at a path for one environment.
type EnvSecrets struct {
	Env  string
	Keys []string
}

// CompareResult holds the diff result between two environments.
type CompareResult struct {
	Path    string
	EnvA    string
	EnvB    string
	DiffResult DiffResult
}

// ComparePathAcrossEnvs lists keys at the given path in two Vault clients
// and returns a CompareResult containing the diff.
func ComparePathAcrossEnvs(ctx context.Context, path string, envA, envB string, clientA, clientB *vault.Client) (*CompareResult, error) {
	keysA, err := clientA.ListSecrets(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("listing secrets for env %q at %q: %w", envA, path, err)
	}

	keysB, err := clientB.ListSecrets(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("listing secrets for env %q at %q: %w", envB, path, err)
	}

	diff := DiffKeys(path, keysA, keysB)

	return &CompareResult{
		Path:       path,
		EnvA:       envA,
		EnvB:       envB,
		DiffResult: diff,
	}, nil
}
