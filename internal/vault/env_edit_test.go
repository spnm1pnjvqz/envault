package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joeshaw/envault/internal/env"
)

func setupEditVault(t *testing.T) (*Vault, string) {
	t.Helper()
	dir := t.TempDir()
	pub, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate keys: %v", err)
	}
	cfg := &Config{
		EnvFile: filepath.Join(dir, ".env.enc"),
		PublicKey: pub,
		PrivateKeyPath: filepath.Join(dir, "key.txt"),
	}
	if err := SaveKeyPair(pub, priv, cfg.PrivateKeyPath); err != nil {
		t.Fatalf("save key pair: %v", err)
	}
	v, err := New(cfg)
	if err != nil {
		t.Fatalf("new vault: %v", err)
	}
	pairs := []env.Pair{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	if err := v.Lock(env.WriteFile, pairs); err != nil {
		t.Fatalf("lock: %v", err)
	}
	return v, dir
}

func TestEditNoChange(t *testing.T) {
	v, dir := setupEditVault(t)
	_ = dir

	// Use a no-op editor (cat just re-writes the file identically on most systems;
	// use 'true' which does nothing — file stays the same).
	result, err := v.EditEnvFile(EditOptions{Editor: "true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Modified {
		t.Error("expected Modified=false when editor makes no changes")
	}
}

func TestEditAddsKey(t *testing.T) {
	v, dir := setupEditVault(t)

	// Write a helper script that appends a new key.
	script := filepath.Join(dir, "editor.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho 'NEW_KEY=hello' >> \"$1\"\n"), 0755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	result, err := v.EditEnvFile(EditOptions{Editor: script})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Modified {
		t.Fatal("expected Modified=true")
	}
	if len(result.KeysAdded) != 1 || result.KeysAdded[0] != "NEW_KEY" {
		t.Errorf("expected KeysAdded=[NEW_KEY], got %v", result.KeysAdded)
	}
}

func TestEditRemovesKey(t *testing.T) {
	v, dir := setupEditVault(t)

	// Script that keeps only the first line.
	script := filepath.Join(dir, "editor.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nhead -1 \"$1\" > /tmp/_ev_tmp && mv /tmp/_ev_tmp \"$1\"\n"), 0755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	result, err := v.EditEnvFile(EditOptions{Editor: script})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Modified {
		t.Fatal("expected Modified=true")
	}
	if len(result.KeysRemoved) != 1 {
		t.Errorf("expected 1 removed key, got %v", result.KeysRemoved)
	}
}

func TestBuildEditResult(t *testing.T) {
	before := []env.Pair{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	after := []env.Pair{{Key: "A", Value: "changed"}, {Key: "C", Value: "3"}}
	r := buildEditResult(before, after)
	if !r.Modified {
		t.Error("expected Modified=true")
	}
	if len(r.KeysAdded) != 1 || r.KeysAdded[0] != "C" {
		t.Errorf("KeysAdded: got %v", r.KeysAdded)
	}
	if len(r.KeysRemoved) != 1 || r.KeysRemoved[0] != "B" {
		t.Errorf("KeysRemoved: got %v", r.KeysRemoved)
	}
	if len(r.KeysChanged) != 1 || r.KeysChanged[0] != "A" {
		t.Errorf("KeysChanged: got %v", r.KeysChanged)
	}
}
