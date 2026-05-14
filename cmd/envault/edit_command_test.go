package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joeshaw/envault/internal/vault"
	"github.com/joeshaw/envault/internal/env"
)

func setupEditConfig(t *testing.T) (cfgPath string, v *vault.Vault) {
	t.Helper()
	dir := t.TempDir()
	pub, priv, err := vault.GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate keys: %v", err)
	}
	keyPath := filepath.Join(dir, "key.txt")
	if err := vault.SaveKeyPair(pub, priv, keyPath); err != nil {
		t.Fatalf("save key pair: %v", err)
	}
	cfg := vault.DefaultConfig()
	cfg.EnvFile = filepath.Join(dir, ".env.enc")
	cfg.PublicKey = pub
	cfg.PrivateKeyPath = keyPath
	cfgPath = filepath.Join(dir, ".envault.toml")
	if err := vault.SaveConfig(cfg, cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}
	v2, err := vault.New(cfg)
	if err != nil {
		t.Fatalf("new vault: %v", err)
	}
	if err := v2.Lock(env.WriteFile, []env.Pair{{Key: "FOO", Value: "bar"}}); err != nil {
		t.Fatalf("lock: %v", err)
	}
	return cfgPath, v2
}

func TestEditCommandNoChange(t *testing.T) {
	cfgPath, _ := setupEditConfig(t)
	// 'true' editor makes no changes
	out, err := executeCommand(rootCmd, "edit", "--config", cfgPath, "--editor", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		// stderr output is fine; we just care there was no crash
	}
}

func TestEditCommandWithChange(t *testing.T) {
	cfgPath, _ := setupEditConfig(t)
	dir := filepath.Dir(cfgPath)

	script := filepath.Join(dir, "editor.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho 'NEW=val' >> \"$1\"\n"), 0755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	out, err := executeCommand(rootCmd, "edit", "--config", cfgPath, "--editor", script)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Log("(output may be on stderr)")
	}
}

func TestEditCommandMissingConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "edit", "--config", "/nonexistent/.envault.toml")
	if err == nil {
		t.Error("expected error for missing config")
	}
}
