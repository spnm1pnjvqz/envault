package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envault/internal/env"
)

func setupDeleteVault(t *testing.T, content string) *Vault {
	t.Helper()
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(content), 0600); err != nil {
		t.Fatalf("setup: write env: %v", err)
	}
	cfg := &Config{EnvFile: envFile}
	v, err := New(cfg)
	if err != nil {
		t.Fatalf("setup: new vault: %v", err)
	}
	return v
}

func TestDeleteKeyBasic(t *testing.T) {
	v := setupDeleteVault(t, "FOO=bar\nBAZ=qux\n")
	res, err := v.DeleteKeys(DeleteOptions{Keys: []string{"FOO"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Deleted) != 1 || res.Deleted[0] != "FOO" {
		t.Errorf("expected Deleted=[FOO], got %v", res.Deleted)
	}
	entries, _ := env.ReadFile(v.cfg.EnvFile)
	if len(entries) != 1 || entries[0].Key != "BAZ" {
		t.Errorf("expected only BAZ remaining, got %v", entries)
	}
}

func TestDeleteKeyNotFound(t *testing.T) {
	v := setupDeleteVault(t, "FOO=bar\n")
	res, err := v.DeleteKeys(DeleteOptions{Keys: []string{"MISSING"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.NotFound) != 1 || res.NotFound[0] != "MISSING" {
		t.Errorf("expected NotFound=[MISSING], got %v", res.NotFound)
	}
}

func TestDeleteKeyMustExistError(t *testing.T) {
	v := setupDeleteVault(t, "FOO=bar\n")
	_, err := v.DeleteKeys(DeleteOptions{Keys: []string{"MISSING"}, MustExist: true})
	if err == nil {
		t.Fatal("expected error for missing key with MustExist=true")
	}
}

func TestDeleteMultipleKeys(t *testing.T) {
	v := setupDeleteVault(t, "A=1\nB=2\nC=3\n")
	res, err := v.DeleteKeys(DeleteOptions{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Deleted) != 2 {
		t.Errorf("expected 2 deleted, got %d", len(res.Deleted))
	}
	entries, _ := env.ReadFile(v.cfg.EnvFile)
	if len(entries) != 1 || entries[0].Key != "B" {
		t.Errorf("expected only B remaining, got %v", entries)
	}
}

func TestDeleteNoKeysError(t *testing.T) {
	v := setupDeleteVault(t, "FOO=bar\n")
	_, err := v.DeleteKeys(DeleteOptions{Keys: []string{}})
	if err == nil {
		t.Fatal("expected error when no keys specified")
	}
}
