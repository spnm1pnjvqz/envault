package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func writeLintEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writeLintEnv: %v", err)
	}
	return path
}

func TestLintCleanFile(t *testing.T) {
	path := writeLintEnv(t, "APP_NAME=envault\nDEBUG=false\n")
	result, err := LintEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d: %+v", len(result.Issues), result.Issues)
	}
}

func TestLintEmptyValue(t *testing.T) {
	path := writeLintEnv(t, "APP_NAME=\n")
	result, err := LintEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) != 1 || result.Issues[0].Severity != "warn" {
		t.Errorf("expected 1 warn for empty value, got %+v", result.Issues)
	}
}

func TestLintDuplicateKey(t *testing.T) {
	path := writeLintEnv(t, "FOO=bar\nFOO=baz\n")
	result, err := LintEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasErrors() {
		t.Error("expected HasErrors() to be true for duplicate key")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "FOO" && issue.Severity == "error" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected duplicate key error for FOO, got %+v", result.Issues)
	}
}

func TestLintNamingConvention(t *testing.T) {
	path := writeLintEnv(t, "appName=envault\n")
	result, err := LintEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "appName" && issue.Message == "key is not UPPER_SNAKE_CASE" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected naming convention warn, got %+v", result.Issues)
	}
}

func TestLintSensitiveKey(t *testing.T) {
	path := writeLintEnv(t, "DB_PASSWORD=supersecret\n")
	result, err := LintEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "DB_PASSWORD" && issue.Severity == "warn" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected sensitive key warn, got %+v", result.Issues)
	}
}

func TestLintFileNotFound(t *testing.T) {
	_, err := LintEnvFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
