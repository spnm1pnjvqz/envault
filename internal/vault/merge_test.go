package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envault/internal/env"
)

func writeMergeEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeMergeEnv: %v", err)
	}
	return p
}

func TestMergeAddsNewKeys(t *testing.T) {
	dir := t.TempDir()
	dst := writeMergeEnv(t, dir, "dst.env", "FOO=bar\n")
	src := writeMergeEnv(t, dir, "src.env", "BAZ=qux\n")

	res, err := MergeEnvFiles(dst, src, MergeStrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "BAZ" {
		t.Errorf("expected Added=[BAZ], got %v", res.Added)
	}
	entries, _ := env.ReadFile(dst)
	m := env.ToMap(entries)
	if m["BAZ"] != "qux" {
		t.Errorf("BAZ not written: %v", m)
	}
}

func TestMergeStrategyOursSkipsConflict(t *testing.T) {
	dir := t.TempDir()
	dst := writeMergeEnv(t, dir, "dst.env", "FOO=original\n")
	src := writeMergeEnv(t, dir, "src.env", "FOO=new\n")

	res, err := MergeEnvFiles(dst, src, MergeStrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "FOO" {
		t.Errorf("expected Skipped=[FOO], got %v", res.Skipped)
	}
	entries, _ := env.ReadFile(dst)
	m := env.ToMap(entries)
	if m["FOO"] != "original" {
		t.Errorf("FOO should remain 'original', got %q", m["FOO"])
	}
}

func TestMergeStrategyTheirsOverwrites(t *testing.T) {
	dir := t.TempDir()
	dst := writeMergeEnv(t, dir, "dst.env", "FOO=original\n")
	src := writeMergeEnv(t, dir, "src.env", "FOO=new\n")

	res, err := MergeEnvFiles(dst, src, MergeStrategyTheirs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 || res.Overwritten[0] != "FOO" {
		t.Errorf("expected Overwritten=[FOO], got %v", res.Overwritten)
	}
	entries, _ := env.ReadFile(dst)
	m := env.ToMap(entries)
	if m["FOO"] != "new" {
		t.Errorf("FOO should be 'new', got %q", m["FOO"])
	}
}

func TestMergeStrategyErrorOnConflict(t *testing.T) {
	dir := t.TempDir()
	dst := writeMergeEnv(t, dir, "dst.env", "FOO=bar\n")
	src := writeMergeEnv(t, dir, "src.env", "FOO=baz\n")

	_, err := MergeEnvFiles(dst, src, MergeStrategyError)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}

func TestMergeSrcNotFound(t *testing.T) {
	dir := t.TempDir()
	dst := writeMergeEnv(t, dir, "dst.env", "FOO=bar\n")
	_, err := MergeEnvFiles(dst, filepath.Join(dir, "missing.env"), MergeStrategyOurs)
	if err == nil {
		t.Fatal("expected error for missing src, got nil")
	}
}
