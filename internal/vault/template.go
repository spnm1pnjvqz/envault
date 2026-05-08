package vault

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nicholasgasior/envault/internal/env"
)

// TemplateResult holds the output of a template rendering operation.
type TemplateResult struct {
	Output    string
	Missing   []string
	Substituted int
}

// templateVarRe matches ${VAR_NAME} and $VAR_NAME patterns.
var templateVarRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// RenderTemplate reads a template file, substitutes variables from the
// decrypted vault env file, and returns the rendered result.
// Missing variables are collected rather than causing an error.
func (v *Vault) RenderTemplate(templatePath string) (*TemplateResult, error) {
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("read template: %w", err)
	}

	pairs, err := env.ReadFile(v.cfg.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}

	vars := env.ToMap(pairs)
	result := &TemplateResult{}
	missingSet := map[string]struct{}{}

	output := templateVarRe.ReplaceAllStringFunc(string(tmplBytes), func(match string) string {
		key := extractKey(match)
		if val, ok := vars[key]; ok {
			result.Substituted++
			return val
		}
		missingSet[key] = struct{}{}
		return match
	})

	for k := range missingSet {
		result.Missing = append(result.Missing, k)
	}
	result.Output = output
	return result, nil
}

// extractKey pulls the variable name from a $VAR or ${VAR} match.
func extractKey(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}
