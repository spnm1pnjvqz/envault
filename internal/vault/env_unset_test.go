package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/subtlepseudonym/envault/internal/env"
)

func setupUnsetVault(t *testing.T, content string) *Vault {
	t.Helper()
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(content), 0600); err != nil {
		t.Fatalf("write env file: %v", err)
	}
	return &Vault{cfg: &Config{EnvFile: envFile}}
}

func TestUnsetKeyBasic(t *testing.T) {
	v := setupUnsetVault(t, "FOO=bar\nBAZ=qux\n")
	res, err := v.UnsetKeys([]string{"FOO"}, UnsetOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "FOO" {
		t.Errorf("expected FOO removed, got %v", res.Removed)
	}
	entries, _ := env.ReadFile(v.cfg.EnvFile)
	for _, e := range entries {
		if e.Key == "FOO" {
			t.Error("FOO should have been removed")
		}
	}
}

func TestUnsetKeyNotFound(t *testing.T) {
	v := setupUnsetVault(t, "FOO=bar\n")
	res, err := v.UnsetKeys([]string{"MISSING"}, UnsetOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.NotFound) != 1 || res.NotFound[0] != "MISSING" {
		t.Errorf("expected MISSING in NotFound, got %v", res.NotFound)
	}
}

func TestUnsetKeyMustExistError(t *testing.T) {
	v := setupUnsetVault(t, "FOO=bar\n")
	_, err := v.UnsetKeys([]string{"MISSING"}, UnsetOptions{MustExist: true})
	if err == nil {
		t.Fatal("expected error for missing key with MustExist")
	}
}

func TestUnsetMultipleKeys(t *testing.T) {
	v := setupUnsetVault(t, "A=1\nB=2\nC=3\n")
	res, err := v.UnsetKeys([]string{"A", "C"}, UnsetOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
	entries, _ := env.ReadFile(v.cfg.EnvFile)
	if len(entries) != 1 || entries[0].Key != "B" {
		t.Errorf("expected only B remaining, got %v", entries)
	}
}

func TestUnsetEmptyKey(t *testing.T) {
	v := setupUnsetVault(t, "FOO=bar\n")
	_, err := v.UnsetKeys([]string{""}, UnsetOptions{})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestUnsetNoKeys(t *testing.T) {
	v := setupUnsetVault(t, "FOO=bar\n")
	_, err := v.UnsetKeys([]string{}, UnsetOptions{})
	if err == nil {
		t.Fatal("expected error when no keys provided")
	}
}
