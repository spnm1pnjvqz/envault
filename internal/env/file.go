package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// ReadFile parses a .env file from the given path.
func ReadFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer f.Close()

	entries, err := Parse(f)
	if err != nil {
		return nil, fmt.Errorf("parsing env file %q: %w", path, err)
	}
	return entries, nil
}

// WriteFile serializes entries and writes them to the given path.
func WriteFile(path string, entries []Entry) error {
	content := Serialize(entries)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return fmt.Errorf("writing env file %q: %w", path, err)
	}
	return nil
}

// ToMap converts a slice of entries into a key→value map.
// Comment-only entries are skipped.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.Key != "" {
			m[e.Key] = e.Value
		}
	}
	return m
}

// FromMap converts a map into a sorted slice of entries.
// Keys are sorted alphabetically to produce deterministic output.
func FromMap(m map[string]string) []Entry {
	entries := make([]Entry, 0, len(m))
	for k, v := range m {
		entries = append(entries, Entry{Key: k, Value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// MaskValue replaces all but the first two characters of a secret value
// with asterisks, useful for display purposes.
func MaskValue(value string) string {
	if len(value) <= 2 {
		return strings.Repeat("*", len(value))
	}
	return value[:2] + strings.Repeat("*", len(value)-2)
}
