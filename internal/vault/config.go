package vault

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// EnvConfig holds connection details for a single Vault environment.
type EnvConfig struct {
	Name      string `yaml:"name"`
	Addr      string `yaml:"addr"`
	Token     string `yaml:"token"`
	MountPath string `yaml:"mount_path"`
}

// Config is the top-level vaultwatch configuration.
type Config struct {
	Environments []EnvConfig `yaml:"environments"`
}

// LoadConfig reads and parses a YAML config file from the given path.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decoding config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Environments) == 0 {
		return fmt.Errorf("config must define at least one environment")
	}
	seen := map[string]bool{}
	for _, e := range c.Environments {
		if e.Name == "" {
			return fmt.Errorf("environment entry missing 'name'")
		}
		if e.Addr == "" {
			return fmt.Errorf("environment %q missing 'addr'", e.Name)
		}
		if seen[e.Name] {
			return fmt.Errorf("duplicate environment name %q", e.Name)
		}
		seen[e.Name] = true
	}
	return nil
}

// ClientsFromConfig builds a Vault Client for every configured environment.
func ClientsFromConfig(cfg *Config) ([]*Client, error) {
	clients := make([]*Client, 0, len(cfg.Environments))
	for _, e := range cfg.Environments {
		c, err := NewClient(e.Name, e.Addr, e.Token)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}
