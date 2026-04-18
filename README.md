# vaultwatch

> CLI tool to audit and diff HashiCorp Vault secret paths across environments

---

## Installation

```bash
go install github.com/yourusername/vaultwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultwatch.git
cd vaultwatch && go build -o vaultwatch .
```

---

## Usage

```bash
# Audit secret paths in a Vault environment
vaultwatch audit --addr https://vault.prod.example.com --path secret/myapp

# Diff secret paths between two environments
vaultwatch diff \
  --src https://vault.staging.example.com \
  --dst https://vault.prod.example.com \
  --path secret/myapp
```

**Common flags:**

| Flag | Description |
|------|-------------|
| `--addr` | Vault server address |
| `--path` | Secret path to audit |
| `--src` / `--dst` | Source and destination for diff |
| `--token` | Vault token (or set `VAULT_TOKEN`) |
| `--output` | Output format: `text` (default) or `json` |
| `--recursive` | Recursively traverse sub-paths |

---

## Requirements

- Go 1.21+
- HashiCorp Vault 1.x

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)
