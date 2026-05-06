package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envault/internal/env"
)

func TestRotateKeys(t *testing.T) {
	dir := t.TempDir()

	recipient, identity, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}

	keyPath := filepath.Join(dir, "identity.txt")
	cfgPath := filepath.Join(dir, ".envault.toml")
	vaultPath := filepath.Join(dir, ".env.age")
	envPath := filepath.Join(dir, ".env")

	if err := SaveKeyPair(recipient, identity, keyPath); err != nil {
		t.Fatalf("save key pair: %v", err)
	}

	cfg := DefaultConfig()
	cfg.PublicKey = recipient.String()
	cfg.PrivateKeyPath = keyPath
	cfg.VaultFile = vaultPath
	cfg.EnvFile = envPath

	if err := SaveConfig(cfg, cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	// Write initial .env file and lock it
	initialEntries := []env.Entry{
		{Key: "SECRET", Value: "mysecret"},
		{Key: "DB_URL", Value: "postgres://localhost/dev"},
	}

	v, err := New(cfgPath)
	if err != nil {
		t.Fatalf("new vault: %v", err)
	}
	if err := v.Lock(initialEntries); err != nil {
		t.Fatalf("lock: %v", err)
	}

	// Rotate keys
	if err := v.RotateKeys(); err != nil {
		t.Fatalf("rotate keys: %v", err)
	}

	// Verify old key was archived
	matches, err := filepath.Glob(keyPath + ".*.bak")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(matches) == 0 {
		t.Error("expected archived key file, found none")
	}

	// Verify vault can still be unlocked with new keys
	v2, err := New(cfgPath)
	if err != nil {
		t.Fatalf("new vault after rotate: %v", err)
	}
	entries, err := v2.View()
	if err != nil {
		t.Fatalf("view after rotate: %v", err)
	}

	entryMap := make(map[string]string)
	for _, e := range entries {
		entryMap[e.Key] = e.Value
	}

	if entryMap["SECRET"] != "mysecret" {
		t.Errorf("expected SECRET=mysecret, got %q", entryMap["SECRET"])
	}
	if entryMap["DB_URL"] != "postgres://localhost/dev" {
		t.Errorf("expected DB_URL=postgres://localhost/dev, got %q", entryMap["DB_URL"])
	}
}

func TestRotateKeysUpdatesConfigPublicKey(t *testing.T) {
	dir := t.TempDir()

	recipient, identity, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}

	keyPath := filepath.Join(dir, "identity.txt")
	cfgPath := filepath.Join(dir, ".envault.toml")
	vaultPath := filepath.Join(dir, ".env.age")
	envPath := filepath.Join(dir, ".env")

	if err := SaveKeyPair(recipient, identity, keyPath); err != nil {
		t.Fatalf("save key pair: %v", err)
	}

	cfg := DefaultConfig()
	cfg.PublicKey = recipient.String()
	cfg.PrivateKeyPath = keyPath
	cfg.VaultFile = vaultPath
	cfg.EnvFile = envPath

	if err := SaveConfig(cfg, cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	v, err := New(cfgPath)
	if err != nil {
		t.Fatalf("new vault: %v", err)
	}
	if err := v.Lock([]env.Entry{{Key: "FOO", Value: "bar"}}); err != nil {
		t.Fatalf("lock: %v", err)
	}

	originalPublicKey := cfg.PublicKey

	if err := v.RotateKeys(); err != nil {
		t.Fatalf("rotate keys: %v", err)
	}

	// Re-load config from disk and verify the public key was updated
	newCfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("load config after rotate: %v", err)
	}
	if newCfg.PublicKey == originalPublicKey {
		t.Error("expected public key to change after rotation, but it remained the same")
	}
	if newCfg.PublicKey == "" {
		t.Error("expected non-empty public key after rotation")
	}
}

func TestArchiveFileNotExist(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "nonexistent.txt")
	dst := filepath.Join(dir, "archive", "nonexistent.bak")

	// Should not return an error if source doesn't exist
	if err := archiveFile(src, dst); err != nil {
		t.Errorf("expected no error for missing src, got: %v", err)
	}

	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Error("expected dst to not exist")
	}
}
