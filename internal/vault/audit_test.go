package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppendAndLoadAuditEvent(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.json")

	event := AuditEvent{
		Operation: "lock",
		VaultFile: ".env.age",
		Success:   true,
	}

	if err := AppendAuditEvent(logPath, event); err != nil {
		t.Fatalf("AppendAuditEvent: %v", err)
	}

	log, err := LoadAuditLog(logPath)
	if err != nil {
		t.Fatalf("LoadAuditLog: %v", err)
	}

	if len(log.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(log.Events))
	}
	if log.Events[0].Operation != "lock" {
		t.Errorf("expected operation 'lock', got %q", log.Events[0].Operation)
	}
	if log.Events[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestAppendMultipleEvents(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.json")

	ops := []string{"lock", "unlock", "view"}
	for _, op := range ops {
		if err := AppendAuditEvent(logPath, AuditEvent{Operation: op, Success: true}); err != nil {
			t.Fatalf("AppendAuditEvent(%s): %v", op, err)
		}
	}

	log, err := LoadAuditLog(logPath)
	if err != nil {
		t.Fatalf("LoadAuditLog: %v", err)
	}
	if len(log.Events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(log.Events))
	}
	for i, op := range ops {
		if log.Events[i].Operation != op {
			t.Errorf("event %d: expected %q, got %q", i, op, log.Events[i].Operation)
		}
	}
}

func TestAppendAuditEventEmptyPath(t *testing.T) {
	err := AppendAuditEvent("", AuditEvent{Operation: "lock"})
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestLoadAuditLogNotFound(t *testing.T) {
	_, err := LoadAuditLog("/nonexistent/path/audit.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadAuditLogEmptyPath(t *testing.T) {
	_, err := LoadAuditLog("")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestDefaultAuditLogPath(t *testing.T) {
	dir := t.TempDir()
	path := DefaultAuditLogPath(dir)
	if path == "" {
		t.Error("expected non-empty audit log path")
	}
	_ = os.MkdirAll(filepath.Dir(path), 0700)
}
