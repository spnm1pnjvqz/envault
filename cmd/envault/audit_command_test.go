package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envault/internal/vault"
)

func TestAuditCommandNoLog(t *testing.T) {
	dir := t.TempDir()
	cfg := vault.DefaultConfig(dir)
	cfgPath := filepath.Join(dir, "config.json")
	if err := vault.SaveConfig(cfgPath, cfg); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	origDefault := defaultConfigPath
	defaultConfigPath = cfgPath
	defer func() { defaultConfigPath = origDefault }()

	out, err := executeCommand(rootCmd, "audit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No audit log found") && !strings.Contains(out, "empty") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestAuditCommandWithEntries(t *testing.T) {
	dir := t.TempDir()
	cfg := vault.DefaultConfig(dir)
	cfgPath := filepath.Join(dir, "config.json")
	if err := vault.SaveConfig(cfgPath, cfg); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	logPath := vault.DefaultAuditLogPath(dir)
	if err := vault.AppendAuditEvent(logPath, vault.AuditEvent{
		Operation: "lock",
		VaultFile: ".env.age",
		Success:   true,
	}); err != nil {
		t.Fatalf("AppendAuditEvent: %v", err)
	}

	origDefault := defaultConfigPath
	defaultConfigPath = cfgPath
	defer func() { defaultConfigPath = origDefault }()

	out, err := executeCommand(rootCmd, "audit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "lock") {
		t.Errorf("expected 'lock' in output, got: %q", out)
	}
	if !strings.Contains(out, ".env.age") {
		t.Errorf("expected vault file in output, got: %q", out)
	}
}

func TestAuditCommandMissingConfig(t *testing.T) {
	origDefault := defaultConfigPath
	defaultConfigPath = "/nonexistent/config.json"
	defer func() { defaultConfigPath = origDefault }()

	_, err := executeCommand(rootCmd, "audit")
	if err == nil {
		t.Error("expected error for missing config")
	}
}

func TestAuditCommandLastFlag(t *testing.T) {
	dir := t.TempDir()
	cfg := vault.DefaultConfig(dir)
	cfgPath := filepath.Join(dir, "config.json")
	_ = vault.SaveConfig(cfgPath, cfg)

	logPath := vault.DefaultAuditLogPath(dir)
	for _, op := range []string{"lock", "unlock", "view"} {
		_ = vault.AppendAuditEvent(logPath, vault.AuditEvent{Operation: op, Success: true})
	}

	origDefault := defaultConfigPath
	defaultConfigPath = cfgPath
	defer func() { defaultConfigPath = origDefault }()

	out, err := executeCommand(rootCmd, "audit", "--last", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "lock") {
		t.Errorf("expected only last entry, but got 'lock': %q", out)
	}
	if !strings.Contains(out, "view") {
		t.Errorf("expected 'view' in last entry output: %q", out)
	}
	_ = os.RemoveAll(dir)
}
