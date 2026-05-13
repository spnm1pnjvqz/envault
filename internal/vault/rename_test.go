package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cqroot/envault/internal/vault"
)

func setupRenameVault(t *testing.T, content string) (*vault.Vault, string) {
	t.Helper()
	dir := t.TempDir()
	pub, priv, err := vault.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	cfg := &vault.Config{
		PublicKey:     pub,
		PrivateKeyFile: filepath.Join(dir, "key.age"),
		EncryptedFile: filepath.Join(dir, ".env.age"),
		EnvFile:       filepath.Join(dir, ".env"),
	}
	if err := os.WriteFile(cfg.EnvFile, []byte(content), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := vault.SaveKeyPair(pub, priv, cfg.PrivateKeyFile); err != nil {
		t.Fatalf("SaveKeyPair: %v", err)
	}
	v, err := vault.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := v.Lock(); err != nil {
		t.Fatalf("Lock: %v", err)
	}
	return v, dir
}

func TestRenameKeyBasic(t *testing.T) {
	v, _ := setupRenameVault(t, "FOO=bar\nBAZ=qux\n")
	res, err := v.RenameKey("FOO", "FOO_RENAMED", false)
	if err != nil {
		t.Fatalf("RenameKey: %v", err)
	}
	if !res.Found {
		t.Fatal("expected key to be found")
	}
	entries, err := v.View()
	if err != nil {
		t.Fatalf("View: %v", err)
	}
	if _, ok := entries["FOO"]; ok {
		t.Error("old key FOO should not exist")
	}
	if _, ok := entries["FOO_RENAMED"]; !ok {
		t.Error("new key FOO_RENAMED should exist")
	}
}

func TestRenameKeyNotFound(t *testing.T) {
	v, _ := setupRenameVault(t, "FOO=bar\n")
	res, err := v.RenameKey("MISSING", "NEW", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Found {
		t.Error("expected Found=false for missing key")
	}
}

func TestRenameKeyConflictNoOverwrite(t *testing.T) {
	v, _ := setupRenameVault(t, "FOO=bar\nBAZ=qux\n")
	_, err := v.RenameKey("FOO", "BAZ", false)
	if err == nil {
		t.Fatal("expected error when overwrite=false and key exists")
	}
}

func TestRenameKeyConflictOverwrite(t *testing.T) {
	v, _ := setupRenameVault(t, "FOO=bar\nBAZ=qux\n")
	_, err := v.RenameKey("FOO", "BAZ", true)
	if err != nil {
		t.Fatalf("RenameKey overwrite: %v", err)
	}
	entries, err := v.View()
	if err != nil {
		t.Fatalf("View: %v", err)
	}
	if val, ok := entries["BAZ"]; !ok || val != "bar" {
		t.Errorf("expected BAZ=bar after overwrite, got %q", val)
	}
}

func TestRenameKeyEmptyArgs(t *testing.T) {
	v, _ := setupRenameVault(t, "FOO=bar\n")
	if _, err := v.RenameKey("", "NEW", false); err == nil {
		t.Error("expected error for empty old key")
	}
	if _, err := v.RenameKey("FOO", "", false); err == nil {
		t.Error("expected error for empty new key")
	}
	if _, err := v.RenameKey("FOO", "FOO", false); err == nil {
		t.Error("expected error for identical keys")
	}
}
