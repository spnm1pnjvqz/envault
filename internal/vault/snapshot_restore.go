package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RestoreSnapshot restores a previously created snapshot back to the active .env file.
// The current .env file is archived before the restore takes place.
func RestoreSnapshot(v *Vault, snapshotName string) error {
	snapshotDir := SnapshotDir(v.Config.EnvFile)
	snapshotPath := filepath.Join(snapshotDir, snapshotName)

	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
		return fmt.Errorf("snapshot not found: %s", snapshotName)
	}

	// Archive the current .env file before overwriting
	if _, err := os.Stat(v.Config.EnvFile); err == nil {
		ts := time.Now().UTC().Format("20060102T150405Z")
		archiveName := fmt.Sprintf(".env.pre-restore.%s", ts)
		archivePath := filepath.Join(filepath.Dir(v.Config.EnvFile), archiveName)
		if err := archiveFile(v.Config.EnvFile, archivePath); err != nil {
			return fmt.Errorf("archiving current env file: %w", err)
		}
	}

	src, err := os.ReadFile(snapshotPath)
	if err != nil {
		return fmt.Errorf("reading snapshot: %w", err)
	}

	if err := os.WriteFile(v.Config.EnvFile, src, 0600); err != nil {
		return fmt.Errorf("writing restored env file: %w", err)
	}

	return nil
}

// SnapshotInfo holds metadata about a snapshot parsed from its filename.
type SnapshotInfo struct {
	Name      string
	Timestamp time.Time
	Label     string
}

// ParseSnapshotName parses a snapshot filename into a SnapshotInfo.
// Expected format: .env.snapshot.<timestamp>[.<label>]
func ParseSnapshotName(name string) (SnapshotInfo, error) {
	parts := strings.SplitN(name, ".", 4)
	// parts: ["", "env", "snapshot", "<timestamp>[.<label>]"]
	if len(parts) < 4 || parts[0] != "" || parts[1] != "env" || parts[2] != "snapshot" {
		return SnapshotInfo{}, fmt.Errorf("invalid snapshot name format: %s", name)
	}

	remainder := parts[3]
	var tsStr, label string

	idx := strings.Index(remainder, ".")
	if idx >= 0 {
		tsStr = remainder[:idx]
		label = remainder[idx+1:]
	} else {
		tsStr = remainder
	}

	ts, err := time.Parse("20060102T150405Z", tsStr)
	if err != nil {
		return SnapshotInfo{}, fmt.Errorf("parsing snapshot timestamp %q: %w", tsStr, err)
	}

	return SnapshotInfo{Name: name, Timestamp: ts, Label: label}, nil
}
