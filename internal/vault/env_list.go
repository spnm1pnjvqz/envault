package vault

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yourusername/envault/internal/env"
)

// ListOptions controls how keys are listed.
type ListOptions struct {
	// SortBy can be "key" (default) or "none".
	SortBy string
	// FilterPrefix limits results to keys with this prefix.
	FilterPrefix string
	// MaskValues replaces values with asterisks.
	MaskValues bool
	// ShowTags includes inline tag annotations.
	ShowTags bool
}

// ListResult holds a single key/value entry from the env file.
type ListResult struct {
	Key   string
	Value string
	Tags  []string
}

// ListKeys reads the decrypted env file and returns all key/value pairs
// according to the provided options.
func (v *Vault) ListKeys(opts ListOptions) ([]ListResult, error) {
	entries, err := env.ReadFile(v.config.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("list: read env file: %w", err)
	}

	results := make([]ListResult, 0, len(entries))
	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		if opts.FilterPrefix != "" && !strings.HasPrefix(e.Key, opts.FilterPrefix) {
			continue
		}
		val := e.Value
		if opts.MaskValues {
			val = env.MaskValue(val)
		}
		r := ListResult{Key: e.Key, Value: val}
		if opts.ShowTags {
			r.Tags, _ = GetTags(v.config.EnvFile, e.Key)
		}
		results = append(results, r)
	}

	if opts.SortBy != "none" {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Key < results[j].Key
		})
	}
	return results, nil
}

// FormatList renders a slice of ListResult into a human-readable string.
func FormatList(results []ListResult, maskValues bool) string {
	if len(results) == 0 {
		return "(no keys found)"
	}
	var sb strings.Builder
	for _, r := range results {
		line := fmt.Sprintf("%s=%s", r.Key, r.Value)
		if len(r.Tags) > 0 {
			line += fmt.Sprintf("  # tags:%s", strings.Join(r.Tags, ","))
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}
