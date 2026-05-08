package vault

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/envault/internal/env"
)

// LintIssue represents a single linting warning or error for an env entry.
type LintIssue struct {
	Key      string
	Severity string // "warn" or "error"
	Message  string
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
}

// HasErrors returns true if any issue has severity "error".
func (r *LintResult) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

// LintEnvFile reads a plaintext .env file and returns a LintResult.
func LintEnvFile(path string) (*LintResult, error) {
	entries, err := env.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lint: reading file: %w", err)
	}

	result := &LintResult{}
	seen := make(map[string]bool)

	for _, entry := range entries {
		key := entry.Key

		// Duplicate key check
		if seen[key] {
			result.Issues = append(result.Issues, LintIssue{
				Key:      key,
				Severity: "error",
				Message:  "duplicate key",
			})
		}
		seen[key] = true

		// Empty value warning
		if entry.Value == "" {
			result.Issues = append(result.Issues, LintIssue{
				Key:      key,
				Severity: "warn",
				Message:  "empty value",
			})
		}

		// Naming convention: should be UPPER_SNAKE_CASE
		if key != strings.ToUpper(key) {
			result.Issues = append(result.Issues, LintIssue{
				Key:      key,
				Severity: "warn",
				Message:  "key is not UPPER_SNAKE_CASE",
			})
		}

		// Detect potential secrets stored as plain values (heuristic)
		lower := strings.ToLower(key)
		if (strings.Contains(lower, "password") || strings.Contains(lower, "secret") || strings.Contains(lower, "token")) && len(entry.Value) > 0 {
			result.Issues = append(result.Issues, LintIssue{
				Key:      key,
				Severity: "warn",
				Message:  "sensitive key has a plaintext value; ensure this file is encrypted",
			})
		}
	}

	return result, nil
}
