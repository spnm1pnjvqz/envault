package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func setupSnapshotVault(t *testing.T) (*Vault, string) {
	t.Helper()
	dir := t.TempDir()

	pub, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("APP=hello\nSECRET=world\n"), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	cfg := &Config{
		EnvFile:    envFile,
		PublicKey:  pub,
		Encrypted:  false,
	}
	configPath := filepath.Join(dir, "envault.json")
	if err := SaveConfig(cfg, configPath); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	v, err := New(cfg, priv)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return v, configPath
}

func TestCreateSnapshot(t *testing.T) {
	v, configPath := setupSnapshotVault(t)

	snap, err := CreateSnapshot(v, configPath, "before-deploy")
	if err != nil {
		t.Fatalf("CreateSnapshot: %v", err)
	}
	if snap.ID == "" {
		t.Error("expected non-empty snapshot ID")
	}
	if snap.Label != "before-deploy" {
		t.Errorf("label = %q, want %q", snap.Label, "before-deploy")
	}
	if snap.Entries["APP"] != "hello" {
		t.Errorf("APP = %q, want %q", snap.Entries["APP"], "hello")
	}
}

func TestListSnapshots(t *testing.T) {
	v, configPath := setupSnapshotVault(t)

	if _, err := CreateSnapshot(v, configPath, "snap1"); err != nil {
		t.Fatalf("CreateSnapshot 1: %v", err)
	}
	if _, err := CreateSnapshot(v, configPath, "snap2"); err != nil {
		t.Fatalf("CreateSnapshot 2: %v", err)
	}

	snaps, err := ListSnapshots(configPath)
	if err != nil {
		t.Fatalf("ListSnapshots: %v", err)
	}
	if len(snaps) != 2 {
		t.Errorf("len(snaps) = %d, want 2", len(snaps))
	}
}

func TestListSnapshotsEmpty(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "envault.json")

	snaps, err := ListSnapshots(configPath)
	if err != nil {
		t.Fatalf("ListSnapshots on missing dir: %v", err)
	}
	if snaps != nil {
		t.Errorf("expected nil slice, got %v", snaps)
	}
}

func TestSnapshotDir(t *testing.T) {
	configPath := "/home/user/project/envault.json"
	got := SnapshotDir(configPath)
	want := "/home/user/project/.snapshots"
	if got != want {
		t.Errorf("SnapshotDir = %q, want %q", got, want)
	}
}
