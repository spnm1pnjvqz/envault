package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	home := "/home/testuser"
	cfg := DefaultConfig(home)

	if cfg.EnvFile != DefaultEnvFile {
		t.Errorf("expected EnvFile %q, got %q", DefaultEnvFile, cfg.EnvFile)
	}
	if cfg.VaultFile != DefaultVaultFile {
		t.Errorf("expected VaultFile %q, got %q", DefaultVaultFile, cfg.VaultFile)
	}
	expectedIdentity := filepath.Join(home, DefaultConfigDir, DefaultKeyFile)
	if cfg.IdentityPath != expectedIdentity {
		t.Errorf("expected IdentityPath %q, got %q", expectedIdentity, cfg.IdentityPath)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")

	orig := &Config{
		IdentityPath: filepath.Join(dir, "identity.age"),
		EnvFile:      ".env",
		VaultFile:    ".env.age",
	}

	if err := SaveConfig(cfgPath, orig); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	info, err := os.Stat(cfgPath)
	if err != nil {
		t.Fatalf("config file not created: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected file perm 0600, got %04o", perm)
	}

	loaded, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loaded.IdentityPath != orig.IdentityPath {
		t.Errorf("IdentityPath mismatch: got %q, want %q", loaded.IdentityPath, orig.IdentityPath)
	}
	if loaded.EnvFile != orig.EnvFile {
		t.Errorf("EnvFile mismatch: got %q, want %q", loaded.EnvFile, orig.EnvFile)
	}
	if loaded.VaultFile != orig.VaultFile {
		t.Errorf("VaultFile mismatch: got %q, want %q", loaded.VaultFile, orig.VaultFile)
	}
}

func TestSaveConfigEmptyPath(t *testing.T) {
	err := SaveConfig("", &Config{})
	if err == nil {
		t.Error("expected error for empty path, got nil")
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing config file, got nil")
	}
}

func TestLoadConfigEmptyPath(t *testing.T) {
	_, err := LoadConfig("")
	if err == nil {
		t.Error("expected error for empty path, got nil")
	}
}
