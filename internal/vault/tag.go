package vault

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nicholasgasior/envault/internal/env"
)

// TagResult holds a key and its associated tags.
type TagResult struct {
	Key  string
	Tags []string
}

// SetTags writes or replaces the tags comment for a key in the env file.
// Tags are stored as inline comments in the form: KEY=VALUE # tags:tag1,tag2
func SetTags(v *Vault, key string, tags []string) error {
	entries, err := env.ReadFile(v.Config.EnvFile)
	if err != nil {
		return fmt.Errorf("read env file: %w", err)
	}

	found := false
	for i, e := range entries {
		if e.Key == key {
			found = true
			if len(tags) == 0 {
				entries[i].Comment = stripTagsFromComment(e.Comment)
			} else {
				sorted := make([]string, len(tags))
				copy(sorted, tags)
				sort.Strings(sorted)
				tagStr := "tags:" + strings.Join(sorted, ",")
				base := stripTagsFromComment(e.Comment)
				if base != "" {
					entries[i].Comment = base + " " + tagStr
				} else {
					entries[i].Comment = tagStr
				}
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("key %q not found in env file", key)
	}

	return env.WriteFile(v.Config.EnvFile, entries)
}

// GetTags returns the tags for a given key.
func GetTags(v *Vault, key string) ([]string, error) {
	entries, err := env.ReadFile(v.Config.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}
	for _, e := range entries {
		if e.Key == key {
			return parseTagsFromComment(e.Comment), nil
		}
	}
	return nil, fmt.Errorf("key %q not found in env file", key)
}

// ListByTag returns all keys that carry the given tag.
func ListByTag(v *Vault, tag string) ([]TagResult, error) {
	entries, err := env.ReadFile(v.Config.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}
	var results []TagResult
	for _, e := range entries {
		tags := parseTagsFromComment(e.Comment)
		for _, t := range tags {
			if t == tag {
				results = append(results, TagResult{Key: e.Key, Tags: tags})
				break
			}
		}
	}
	return results, nil
}

func parseTagsFromComment(comment string) []string {
	for _, part := range strings.Fields(comment) {
		if strings.HasPrefix(part, "tags:") {
			raw := strings.TrimPrefix(part, "tags:")
			if raw == "" {
				return nil
			}
			return strings.Split(raw, ",")
		}
	}
	return nil
}

func stripTagsFromComment(comment string) string {
	parts := strings.Fields(comment)
	var kept []string
	for _, p := range parts {
		if !strings.HasPrefix(p, "tags:") {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, " ")
}
