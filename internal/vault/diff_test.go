package vault

import (
	"strings"
	"testing"
)

func TestDiffNoChanges(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	entries := Diff(env, env)
	if len(entries) != 0 {
		t.Fatalf("expected no diff entries, got %d", len(entries))
	}
}

func TestDiffAdded(t *testing.T) {
	oldEnv := map[string]string{"FOO": "bar"}
	newEnv := map[string]string{"FOO": "bar", "NEW_KEY": "value"}
	entries := Diff(oldEnv, newEnv)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != "added" || entries[0].Key != "NEW_KEY" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestDiffRemoved(t *testing.T) {
	oldEnv := map[string]string{"FOO": "bar", "OLD_KEY": "val"}
	newEnv := map[string]string{"FOO": "bar"}
	entries := Diff(oldEnv, newEnv)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != "removed" || entries[0].Key != "OLD_KEY" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestDiffChanged(t *testing.T) {
	oldEnv := map[string]string{"FOO": "bar"}
	newEnv := map[string]string{"FOO": "baz"}
	entries := Diff(oldEnv, newEnv)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Action != "changed" || e.Old != "bar" || e.New != "baz" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestFormatDiffNoChanges(t *testing.T) {
	out := FormatDiff(nil, false)
	if out != "No changes detected." {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatDiffMaskSecrets(t *testing.T) {
	entries := []DiffEntry{
		{Key: "SECRET", Old: "hunter2", New: "newpass", Action: "changed"},
	}
	out := FormatDiff(entries, true)
	if strings.Contains(out, "hunter2") || strings.Contains(out, "newpass") {
		t.Errorf("secrets not masked in output: %q", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected masked values in output: %q", out)
	}
}

func TestFormatDiffSymbols(t *testing.T) {
	entries := []DiffEntry{
		{Key: "A", Old: "", New: "1", Action: "added"},
		{Key: "B", Old: "2", New: "", Action: "removed"},
		{Key: "C", Old: "x", New: "y", Action: "changed"},
	}
	out := FormatDiff(entries, false)
	if !strings.Contains(out, "+ A=1") {
		t.Errorf("missing added line: %q", out)
	}
	if !strings.Contains(out, "- B=2") {
		t.Errorf("missing removed line: %q", out)
	}
	if !strings.Contains(out, "~ C: x -> y") {
		t.Errorf("missing changed line: %q", out)
	}
}
