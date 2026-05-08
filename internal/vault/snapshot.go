package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of vault env entries.
type Snapshot struct {
	ID        string            `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label,omitempty"`
	Entries   map[string]string `json:"entries"`
}

// SnapshotDir returns the directory used to store snapshots relative to a vault config path.
func SnapshotDir(configPath string) string {
	return filepath.Join(filepath.Dir(configPath), ".snapshots")
}

// CreateSnapshot decrypts the current vault and saves a named snapshot to disk.
func CreateSnapshot(v *Vault, configPath, label string) (*Snapshot, error) {
	entries, err := v.View()
	if err != nil {
		return nil, fmt.Errorf("snapshot: failed to read vault: %w", err)
	}

	envMap := make(map[string]string, len(entries))
	for _, e := range entries {
		envMap[e.Key] = e.Value
	}

	snap := &Snapshot{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Entries:   envMap,
	}

	dir := SnapshotDir(configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("snapshot: mkdir: %w", err)
	}

	fileName := snap.ID + ".json"
	if label != "" {
		fileName = label + "_" + snap.ID + ".json"
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir, fileName), data, 0600); err != nil {
		return nil, fmt.Errorf("snapshot: write: %w", err)
	}

	return snap, nil
}

// ListSnapshots returns all snapshots stored in the snapshot directory.
func ListSnapshots(configPath string) ([]Snapshot, error) {
	dir := SnapshotDir(configPath)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot: readdir: %w", err)
	}

	var snaps []Snapshot
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("snapshot: read %s: %w", entry.Name(), err)
		}
		var s Snapshot
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, fmt.Errorf("snapshot: parse %s: %w", entry.Name(), err)
		}
		snaps = append(snaps, s)
	}
	return snaps, nil
}
