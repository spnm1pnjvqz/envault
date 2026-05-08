package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envault/internal/env"
)

func writeTempEnv(t *testing.T, dir, name string, entries []env.Entry) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := env.WriteFile(p, entries); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestCopyKeysBasic(t *testing.T) {
	dir := t.TempDir()
	src := writeTempEnv(t, dir, "src.env", []env.Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}})
	dst := writeTempEnv(t, dir, "dst.env", []env.Entry{{Key: "EXISTING", Value: "yes"}})

	res, err := CopyKeys(src, dst, CopyOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(res.Skipped))
	}

	entries, _ := env.ReadFile(dst)
	m := env.ToMap(entries)
	if m["FOO"] != "bar" || m["BAZ"] != "qux" || m["EXISTING"] != "yes" {
		t.Errorf("destination map unexpected: %v", m)
	}
}

func TestCopyKeysSkipExisting(t *testing.T) {
	dir := t.TempDir()
	src := writeTempEnv(t, dir, "src.env", []env.Entry{{Key: "FOO", Value: "new"}})
	dst := writeTempEnv(t, dir, "dst.env", []env.Entry{{Key: "FOO", Value: "old"}})

	res, err := CopyKeys(src, dst, CopyOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}

	entries, _ := env.ReadFile(dst)
	if env.ToMap(entries)["FOO"] != "old" {
		t.Error("FOO should not have been overwritten")
	}
}

func TestCopyKeysOverwrite(t *testing.T) {
	dir := t.TempDir()
	src := writeTempEnv(t, dir, "src.env", []env.Entry{{Key: "FOO", Value: "new"}})
	dst := writeTempEnv(t, dir, "dst.env", []env.Entry{{Key: "FOO", Value: "old"}})

	res, err := CopyKeys(src, dst, CopyOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwrote) != 1 {
		t.Errorf("expected 1 overwrote, got %d", len(res.Overwrote))
	}

	entries, _ := env.ReadFile(dst)
	if env.ToMap(entries)["FOO"] != "new" {
		t.Error("FOO should have been overwritten")
	}
}

func TestCopyKeysAllowlist(t *testing.T) {
	dir := t.TempDir()
	src := writeTempEnv(t, dir, "src.env", []env.Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}, {Key: "C", Value: "3"}})
	dst := writeTempEnv(t, dir, "dst.env", []env.Entry{})

	res, err := CopyKeys(src, dst, CopyOptions{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}

	entries, _ := env.ReadFile(dst)
	m := env.ToMap(entries)
	if _, ok := m["B"]; ok {
		t.Error("B should not have been copied")
	}
}

func TestCopyKeysMissingSource(t *testing.T) {
	dir := t.TempDir()
	dst := writeTempEnv(t, dir, "dst.env", []env.Entry{})
	_, err := CopyKeys(filepath.Join(dir, "missing.env"), dst, CopyOptions{})
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestCopyKeysMissingDest(t *testing.T) {
	dir := t.TempDir()
	src := writeTempEnv(t, dir, "src.env", []env.Entry{{Key: "X", Value: "1"}})
	_, err := CopyKeys(src, filepath.Join(dir, "missing.env"), CopyOptions{})
	if err == nil {
		t.Error("expected error for missing destination")
	}
	_ = os.Remove(src)
}
