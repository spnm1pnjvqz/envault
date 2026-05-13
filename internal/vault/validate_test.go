package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envault/internal/vault"
)

func writeValidateEnv(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	return p
}

func newValidateVault(t *testing.T, envPath string) *Vault {
	t.Helper()
	v := &Vault{Config: &Config{EnvFile: envPath}}
	return v
}

func TestValidateCleanFile(t *testing.T) {
	dir := t.TempDir()
	envPath := writeValidateEnv(t, dir, "APP_NAME=envault\nDEBUG=false\n")
	v := newValidateVault(t, envPath)
	issues, err := ValidateEnvFile(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestValidateEmptyValue(t *testing.T) {
	dir := t.TempDir()
	envPath := writeValidateEnv(t, dir, "APP_NAME=\n")
	v := newValidateVault(t, envPath)
	issues, err := ValidateEnvFile(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, issue := range issues {
		if issue.Key == "APP_NAME" && issue.Rule == "no-empty-value" {
			found = true
		}
	}
	if !found {
		t.Error("expected no-empty-value issue for APP_NAME")
	}
}

func TestValidateLowercaseKey(t *testing.T) {
	dir := t.TempDir()
	envPath := writeValidateEnv(t, dir, "app_name=envault\n")
	v := newValidateVault(t, envPath)
	issues, err := ValidateEnvFile(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rules := map[string]bool{}
	for _, issue := range issues {
		if issue.Key == "app_name" {
			rules[issue.Rule] = true
		}
	}
	if !rules["key-uppercase"] {
		t.Error("expected key-uppercase issue for app_name")
	}
	if !rules["key-valid-chars"] {
		t.Error("expected key-valid-chars issue for app_name")
	}
}

func TestValidateWhitespaceValue(t *testing.T) {
	dir := t.TempDir()
	envPath := writeValidateEnv(t, dir, "APP_NAME= envault \n")
	v := newValidateVault(t, envPath)
	issues, err := ValidateEnvFile(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, issue := range issues {
		if issue.Key == "APP_NAME" && issue.Rule == "no-whitespace-value" {
			found = true
		}
	}
	if !found {
		t.Error("expected no-whitespace-value issue for APP_NAME")
	}
}

func TestFormatValidationIssuesEmpty(t *testing.T) {
	out := FormatValidationIssues(nil)
	if out != "No validation issues found." {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatValidationIssuesNonEmpty(t *testing.T) {
	issues := []ValidationIssue{
		{Key: "foo", Rule: "key-uppercase", Message: "key should be uppercase"},
	}
	out := FormatValidationIssues(issues)
	if out == "" {
		t.Error("expected non-empty output")
	}
}
