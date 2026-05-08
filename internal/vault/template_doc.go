// Package vault — template rendering.
//
// The RenderTemplate method substitutes environment variable placeholders
// in a text template file using secrets stored in the vault's decrypted
// .env file. Both $VAR_NAME and ${VAR_NAME} syntax are supported.
//
// Usage:
//
//	v, _ := vault.New(cfg, "")
//	result, err := v.RenderTemplate("nginx.conf.tmpl")
//	if err != nil { ... }
//	fmt.Println(result.Output)
//	if len(result.Missing) > 0 {
//	    log.Printf("unresolved: %v", result.Missing)
//	}
//
// The TemplateResult struct exposes:
//   - Output:      the fully rendered string
//   - Missing:     variable names that had no matching key in the .env file
//   - Substituted: count of successful replacements
//
// Unresolved placeholders are preserved verbatim in the output so that
// partial renders can be inspected or retried after adding missing keys.
package vault
