package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envault/internal/vault"
)

func setupTagConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(content), 0o600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	cfg := vault.DefaultConfig(dir)
	cfg.EnvFile = envFile
	cfgPath := filepath.Join(dir, "envault.json")
	if err := vault.SaveConfig(cfgPath, cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	return cfgPath
}

func TestTagSetCommand(t *testing.T) {
	cfgPath := setupTagConfig(t, "DB_HOST=localhost\n")
	out, err := executeCommand(rootCmd, "--config", cfgPath, "tag", "set", "DB_HOST", "infra,database")
	if err != nil {
		t.Fatalf("tag set: %v", err)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected key in output, got: %s", out)
	}
}

func TestTagGetCommand(t *testing.T) {
	cfgPath := setupTagConfig(t, "API_KEY=secret\n")
	_, _ = executeCommand(rootCmd, "--config", cfgPath, "tag", "set", "API_KEY", "auth")
	out, err := executeCommand(rootCmd, "--config", cfgPath, "tag", "get", "API_KEY")
	if err != nil {
		t.Fatalf("tag get: %v", err)
	}
	if !strings.Contains(out, "auth") {
		t.Errorf("expected 'auth' in output, got: %s", out)
	}
}

func TestTagGetNoTags(t *testing.T) {
	cfgPath := setupTagConfig(t, "FOO=bar\n")
	out, err := executeCommand(rootCmd, "--config", cfgPath, "tag", "get", "FOO")
	if err != nil {
		t.Fatalf("tag get: %v", err)
	}
	if !strings.Contains(out, "no tags") {
		t.Errorf("expected 'no tags' in output, got: %s", out)
	}
}

func TestTagListCommand(t *testing.T) {
	cfgPath := setupTagConfig(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	_, _ = executeCommand(rootCmd, "--config", cfgPath, "tag", "set", "DB_HOST", "database")
	_, _ = executeCommand(rootCmd, "--config", cfgPath, "tag", "set", "DB_PORT", "database")
	out, err := executeCommand(rootCmd, "--config", cfgPath, "tag", "list", "database")
	if err != nil {
		t.Fatalf("tag list: %v", err)
	}
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "DB_PORT") {
		t.Errorf("expected both keys in output, got: %s", out)
	}
}

func TestTagListNoMatch(t *testing.T) {
	cfgPath := setupTagConfig(t, "FOO=bar\n")
	out, err := executeCommand(rootCmd, "--config", cfgPath, "tag", "list", "ghost")
	if err != nil {
		t.Fatalf("tag list: %v", err)
	}
	if !strings.Contains(out, "no keys tagged") {
		t.Errorf("expected 'no keys tagged' in output, got: %s", out)
	}
}

func TestTagCommandMissingConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "--config", "/nonexistent/envault.json", "tag", "list", "foo")
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}
