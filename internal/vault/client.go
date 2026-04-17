package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with environment metadata.
type Client struct {
	api  *vaultapi.Client
	Env  string
	Addr string
}

// NewClient creates a new Vault client for the given address and token.
func NewClient(env, addr, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client for env %q: %w", env, err)
	}

	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	api.SetToken(token)

	return &Client{
		api:  api,
		Env:  env,
		Addr: addr,
	}, nil
}

// ListSecrets returns all secret keys under the given KV v2 path.
func (c *Client) ListSecrets(mountPath, secretPath string) ([]string, error) {
	logical := c.api.Logical()
	listPath := fmt.Sprintf("%s/metadata/%s", mountPath, secretPath)

	secret, err := logical.List(listPath)
	if err != nil {
		return nil, fmt.Errorf("listing secrets at %q: %w", listPath, err)
	}
	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	rawKeys, ok := secret.Data["keys"]
	if !ok {
		return []string{}, nil
	}

	ifaces, ok := rawKeys.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for keys at %q", listPath)
	}

	keys := make([]string, 0, len(ifaces))
	for _, k := range ifaces {
		if s, ok := k.(string); ok {
			keys = append(keys, s)
		}
	}
	return keys, nil
}

// ReadSecret reads key-value pairs from a KV v2 secret path.
func (c *Client) ReadSecret(mountPath, secretPath string) (map[string]interface{}, error) {
	readPath := fmt.Sprintf("%s/data/%s", mountPath, secretPath)
	secret, err := c.api.Logical().Read(readPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", readPath, err)
	}
	if secret == nil || secret.Data == nil {
		return map[string]interface{}{}, nil
	}
	data, _ := secret.Data["data"].(map[string]interface{})
	return data, nil
}
