package vault_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourorg/vaultwatch/internal/vault"
)

func mockVaultServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func TestNewClient_Success(t *testing.T) {
	srv := mockVaultServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c, err := vault.NewClient("dev", srv.URL, "test-token")
	require.NoError(t, err)
	assert.Equal(t, "dev", c.Env)
	assert.Equal(t, srv.URL, c.Addr)
}

func TestReadSecret_ReturnsData(t *testing.T) {
	body := `{"data":{"data":{"username":"admin","password":"s3cr3t"}}}`
	srv := mockVaultServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	})

	c, err := vault.NewClient("dev", srv.URL, "test-token")
	require.NoError(t, err)

	data, err := c.ReadSecret("secret", "myapp/config")
	require.NoError(t, err)
	assert.Equal(t, "admin", data["username"])
	assert.Equal(t, "s3cr3t", data["password"])
}

func TestListSecrets_ReturnsKeys(t *testing.T) {
	body := `{"data":{"keys":["db","api/"]}}`
	srv := mockVaultServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	})

	c, err := vault.NewClient("staging", srv.URL, "test-token")
	require.NoError(t, err)

	keys, err := c.ListSecrets("secret", "myapp")
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"db", "api/"}, keys)
}
