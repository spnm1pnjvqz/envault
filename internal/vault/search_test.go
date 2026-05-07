package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envault/internal/env"
)

func setupSearchVault(t *testing.T) (*Vault, string) {
	t.Helper()
	dir := t.TempDir()

	recipient, identity := generateTestKeys(t)

	cfg := &Config{
		EnvFile:    filepath.Join(dir, ".env"),
		VaultFile:  filepath.Join(dir, ".env.age"),
		PublicKey:  recipient,
		PrivateKey: filepath.Join(dir, "key.txt"),
	}

	entries := []env.Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/mydb"},
		{Key: "API_KEY", Value: "supersecret"},
		{Key: "DEBUG", Value: "true"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	if err := env.WriteFile(cfg.EnvFile, entries); err != nil {
		t.Fatal(err)
	}

	v, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := v.Lock(); err != nil {
		t.Fatal(err)
	}
	os.Remove(cfg.EnvFile)

	_ = identity
	// Store identity so decrypt works
	if err := os.WriteFile(cfg.PrivateKey, []byte(identity), 0600); err != nil {
		t.Fatal(err)
	}
	return v, dir
}

func TestSearchByKeyPrefix(t *testing.T) {
	v, _ := setupSearchVault(t)
	results, err := v.Search(SearchOptions{Query: "DB", CaseSensitive: false, KeyOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchCaseSensitive(t *testing.T) {
	v, _ := setupSearchVault(t)
	results, err := v.Search(SearchOptions{Query: "db", CaseSensitive: true, KeyOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results for case-sensitive 'db', got %d", len(results))
	}
}

func TestSearchMaskValues(t *testing.T) {
	v, _ := setupSearchVault(t)
	results, err := v.Search(SearchOptions{Query: "API_KEY", MaskValues: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Value == "supersecret" {
		t.Error("expected value to be masked")
	}
	if !results[0].Masked {
		t.Error("expected Masked flag to be true")
	}
}

func TestSearchNoMatch(t *testing.T) {
	v, _ := setupSearchVault(t)
	results, err := v.Search(SearchOptions{Query: "NONEXISTENT"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
