package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envault/internal/env"
)

func writeCopyEnv(t *testing.T, dir, name string, entries []env.Entry) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := env.WriteFile(p, entries); err != nil {
		t.Fatalf("writeCopyEnv: %v", err)
	}
	return p
}

func TestCopyCommandBasic(t *testing.T) {
	dir := t.TempDir()
	src := writeCopyEnv(t, dir, "src.env", []env.Entry{{Key: "FOO", Value: "bar"}})
	dst := writeCopyEnv(t, dir, "dst.env", []env.Entry{})

	out, err := executeCommand(rootCmd, "copy", src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Copied:    1") {
		t.Errorf("expected copied count in output, got: %s", out)
	}
}

func TestCopyCommandSkips(t *testing.T) {
	dir := t.TempDir()
	src := writeCopyEnv(t, dir, "src.env", []env.Entry{{Key: "FOO", Value: "new"}})
	dst := writeCopyEnv(t, dir, "dst.env", []env.Entry{{Key: "FOO", Value: "old"}})

	out, err := executeCommand(rootCmd, "copy", src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Skipped:   1") {
		t.Errorf("expected skipped count in output, got: %s", out)
	}
	if !strings.Contains(out, "--overwrite") {
		t.Errorf("expected overwrite hint in output, got: %s", out)
	}
}

func TestCopyCommandOverwrite(t *testing.T) {
	dir := t.TempDir()
	src := writeCopyEnv(t, dir, "src.env", []env.Entry{{Key: "FOO", Value: "new"}})
	dst := writeCopyEnv(t, dir, "dst.env", []env.Entry{{Key: "FOO", Value: "old"}})

	out, err := executeCommand(rootCmd, "copy", "--overwrite", src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Overwrote: 1") {
		t.Errorf("expected overwrote count in output, got: %s", out)
	}
}

func TestCopyCommandMissingArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "copy", "only-one-arg")
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestCopyCommandMissingSource(t *testing.T) {
	dir := t.TempDir()
	dst := writeCopyEnv(t, dir, "dst.env", []env.Entry{})
	_, err := executeCommand(rootCmd, "copy", filepath.Join(dir, "missing.env"), dst)
	if err == nil {
		t.Error("expected error for missing source file")
	}
}
