package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRestoreSnapshot(t *testing.T) {
	v := setupSnapshotVault(t)

	// Write initial content and create a snapshot
	initial := "DB_HOST=localhost\nDB_PORT=5432\n"
	if err := os.WriteFile(v.Config.EnvFile, []byte(initial), 0600); err != nil {
		t.Fatal(err)
	}

	snapName, err := CreateSnapshot(v, "before-change")
	if err != nil {
		t.Fatalf("CreateSnapshot: %v", err)
	}

	// Overwrite .env with new content
	updated := "DB_HOST=prod-server\nDB_PORT=5432\nDB_NAME=mydb\n"
	if err := os.WriteFile(v.Config.EnvFile, []byte(updated), 0600); err != nil {
		t.Fatal(err)
	}

	// Restore the snapshot
	if err := RestoreSnapshot(v, snapName); err != nil {
		t.Fatalf("RestoreSnapshot: %v", err)
	}

	restored, err := os.ReadFile(v.Config.EnvFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(restored) != initial {
		t.Errorf("expected restored content %q, got %q", initial, string(restored))
	}
}

func TestRestoreSnapshotNotFound(t *testing.T) {
	v := setupSnapshotVault(t)
	err := RestoreSnapshot(v, ".env.snapshot.19991231T235959Z")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
	if !strings.Contains(err.Error(), "snapshot not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRestoreSnapshotArchivesCurrentFile(t *testing.T) {
	v := setupSnapshotVault(t)

	original := "KEY=original\n"
	if err := os.WriteFile(v.Config.EnvFile, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}

	snapName, err := CreateSnapshot(v, "")
	if err != nil {
		t.Fatalf("CreateSnapshot: %v", err)
	}

	if err := os.WriteFile(v.Config.EnvFile, []byte("KEY=changed\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := RestoreSnapshot(v, snapName); err != nil {
		t.Fatalf("RestoreSnapshot: %v", err)
	}

	dir := filepath.Dir(v.Config.EnvFile)
	entries, _ := os.ReadDir(dir)
	found := false
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".env.pre-restore.") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected pre-restore archive file to exist")
	}
}

func TestParseSnapshotName(t *testing.T) {
	info, err := ParseSnapshotName(".env.snapshot.20240101T120000Z.my-label")
	if err != nil {
		t.Fatalf("ParseSnapshotName: %v", err)
	}
	if info.Label != "my-label" {
		t.Errorf("expected label 'my-label', got %q", info.Label)
	}
	if info.Timestamp.Year() != 2024 {
		t.Errorf("unexpected year: %d", info.Timestamp.Year())
	}
}

func TestParseSnapshotNameNoLabel(t *testing.T) {
	info, err := ParseSnapshotName(".env.snapshot.20230615T083000Z")
	if err != nil {
		t.Fatalf("ParseSnapshotName: %v", err)
	}
	if info.Label != "" {
		t.Errorf("expected empty label, got %q", info.Label)
	}
}

func TestParseSnapshotNameInvalid(t *testing.T) {
	_, err := ParseSnapshotName("not-a-snapshot")
	if err == nil {
		t.Fatal("expected error for invalid snapshot name")
	}
}
