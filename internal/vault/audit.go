package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuditEvent represents a single recorded vault operation.
type AuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	VaultFile string    `json:"vault_file"`
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
}

// AuditLog holds a list of recorded events.
type AuditLog struct {
	Events []AuditEvent `json:"events"`
}

// DefaultAuditLogPath returns the default path for the audit log.
func DefaultAuditLogPath(dir string) string {
	return filepath.Join(dir, ".envault_audit.json")
}

// AppendAuditEvent loads the existing audit log (if any), appends the event, and saves it.
func AppendAuditEvent(logPath string, event AuditEvent) error {
	if logPath == "" {
		return fmt.Errorf("audit log path must not be empty")
	}

	log, err := LoadAuditLog(logPath)
	if err != nil {
		log = &AuditLog{}
	}

	event.Timestamp = time.Now().UTC()
	log.Events = append(log.Events, event)

	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal audit log: %w", err)
	}

	if err := os.WriteFile(logPath, data, 0600); err != nil {
		return fmt.Errorf("write audit log: %w", err)
	}
	return nil
}

// LoadAuditLog reads and parses the audit log from disk.
func LoadAuditLog(logPath string) (*AuditLog, error) {
	if logPath == "" {
		return nil, fmt.Errorf("audit log path must not be empty")
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		return nil, fmt.Errorf("read audit log: %w", err)
	}

	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("parse audit log: %w", err)
	}
	return &log, nil
}
