package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/vaultwatch/internal/vault"
)

// WatchConfig defines polling behavior for watch mode.
type WatchConfig struct {
	Interval  time.Duration
	Path      string
	Clients   map[string]*vault.Client
	OnChange  func(results []CompareResult)
}

// CompareResult holds the diff result for a single environment pair.
type CompareResult struct {
	EnvA    string
	EnvB    string
	Diff    DiffResult
	CheckedAt time.Time
}

// Watch polls Vault at the given interval and calls OnChange when diffs are detected.
func Watch(ctx context.Context, cfg WatchConfig) error {
	if cfg.Interval <= 0 {
		return fmt.Errorf("watch interval must be positive")
	}
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	var prev []CompareResult

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			results, err := pollOnce(cfg)
			if err != nil {
				return fmt.Errorf("poll error: %w", err)
			}
			if hasChanged(prev, results) {
				cfg.OnChange(results)
			}
			prev = results
		}
	}
}

func pollOnce(cfg WatchConfig) ([]CompareResult, error) {
	envNames := make([]string, 0, len(cfg.Clients))
	for name := range cfg.Clients {
		envNames = append(envNames, name)
	}
	if len(envNames) < 2 {
		return nil, fmt.Errorf("need at least two environments to compare")
	}
	var results []CompareResult
	for i := 0; i < len(envNames)-1; i++ {
		a, b := envNames[i], envNames[i+1]
		diff, err := ComparePathAcrossEnvs(cfg.Path, map[string]*vault.Client{a: cfg.Clients[a], b: cfg.Clients[b]})
		if err != nil {
			return nil, err
		}
		for _, d := range diff {
			results = append(results, CompareResult{EnvA: a, EnvB: b, Diff: d, CheckedAt: time.Now()})
		}
	}
	return results, nil
}

func hasChanged(prev, next []CompareResult) bool {
	if len(prev) != len(next) {
		return true
	}
	for i := range next {
		if next[i].Diff != prev[i].Diff {
			return true
		}
	}
	return false
}
