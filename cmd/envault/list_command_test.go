package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/envault/internal/vault"
)

func setupListConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(content), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	cfg := vault.DefaultConfig()
	cfg.EnvFile = envFile
	cfgPath := filepath.Join(dir, "envault.json")
	if err := vault.SaveConfig(cfgPath, cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	return cfgPath
}

func TestListCommandBasic(t *testing.T) {
	cfgPath := setupListConfig(t, "FOO=bar\nBAZ=qux\n")
	out, err := executeCommand(rootCmd, "list", "--config", cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", out)
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got: %s", out)
	}
}

func TestListCommandMask(t *testing.T) {
	cfgPath := setupListConfig(t, "SECRET=topsecret\n")
	out, err := executeCommand(rootCmd, "list", "--config", cfgPath, "--mask")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "topsecret") {
		t.Errorf("expected value to be masked, got: %s", out)
	}
}

func TestListCommandPrefix(t *testing.T) {
	cfgPath := setupListConfig(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=test\n")
	out, err := executeCommand(rootCmd, "list", "--config", cfgPath, "--prefix", "DB_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "APP_NAME") {
		t.Errorf("APP_NAME should be filtered out, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
}

func TestListCommandMissingConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "list", "--config", "/nonexistent/envault.json")
	if err == nil {
		t.Error("expected error for missing config")
	}
}

func TestListCommandEmpty(t *testing.T) {
	cfgPath := setupListConfig(t, "")
	out, err := executeCommand(rootCmd, "list", "--config", cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no keys found") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
