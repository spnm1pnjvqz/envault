package vault

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nicholasgasior/envault/internal/env"
)

// ValidationRule defines a rule applied to env entries.
type ValidationRule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// ValidationIssue represents a single failed rule for a key.
type ValidationIssue struct {
	Key     string
	Rule    string
	Message string
}

var defaultRules = []ValidationRule{
	{
		Name:    "no-empty-value",
		Message: "value is empty",
		Check:   func(_, value string) bool { return value == "" },
	},
	{
		Name:    "key-uppercase",
		Message: "key should be uppercase",
		Check:   func(key, _ string) bool { return key != strings.ToUpper(key) },
	},
	{
		Name:    "no-whitespace-value",
		Message: "value contains leading or trailing whitespace",
		Check:   func(_, value string) bool { return value != strings.TrimSpace(value) },
	},
	{
		Name:    "key-valid-chars",
		Message: "key contains invalid characters (only A-Z, 0-9, _ allowed)",
		Check: func(key, _ string) bool {
			matched, _ := regexp.MatchString(`^[A-Z][A-Z0-9_]*$`, key)
			return !matched
		},
	},
}

// ValidateEnvFile reads the decrypted env file and applies all validation
// rules, returning a list of issues found.
func ValidateEnvFile(v *Vault) ([]ValidationIssue, error) {
	entries, err := env.ReadFile(v.Config.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}

	var issues []ValidationIssue
	for _, entry := range entries {
		for _, rule := range defaultRules {
			if rule.Check(entry.Key, entry.Value) {
				issues = append(issues, ValidationIssue{
					Key:     entry.Key,
					Rule:    rule.Name,
					Message: rule.Message,
				})
			}
		}
	}
	return issues, nil
}

// FormatValidationIssues returns a human-readable summary of issues.
func FormatValidationIssues(issues []ValidationIssue) string {
	if len(issues) == 0 {
		return "No validation issues found."
	}
	var sb strings.Builder
	for _, issue := range issues {
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", issue.Rule, issue.Key, issue.Message))
	}
	return strings.TrimRight(sb.String(), "\n")
}
