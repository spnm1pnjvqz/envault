package vault

import (
	"fmt"
	"sort"
	"strings"
)

// DiffEntry represents a single change between two env states.
type DiffEntry struct {
	Key    string
	Old    string
	New    string
	Action string // "added", "removed", "changed"
}

// Diff compares two maps of env vars and returns a list of changes.
func Diff(oldEnv, newEnv map[string]string) []DiffEntry {
	var entries []DiffEntry

	for key, oldVal := range oldEnv {
		if newVal, ok := newEnv[key]; !ok {
			entries = append(entries, DiffEntry{
				Key:    key,
				Old:    oldVal,
				New:    "",
				Action: "removed",
			})
		} else if oldVal != newVal {
			entries = append(entries, DiffEntry{
				Key:    key,
				Old:    oldVal,
				New:    newVal,
				Action: "changed",
			})
		}
	}

	for key, newVal := range newEnv {
		if _, ok := oldEnv[key]; !ok {
			entries = append(entries, DiffEntry{
				Key:    key,
				Old:    "",
				New:    newVal,
				Action: "added",
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return entries
}

// FormatDiff renders a human-readable diff string.
// If maskSecrets is true, values are replaced with "***".
func FormatDiff(entries []DiffEntry, maskSecrets bool) string {
	if len(entries) == 0 {
		return "No changes detected."
	}

	var sb strings.Builder
	for _, e := range entries {
		oldVal := e.Old
		newVal := e.New
		if maskSecrets {
			if oldVal != "" {
				oldVal = "***"
			}
			if newVal != "" {
				newVal = "***"
			}
		}
		switch e.Action {
		case "added":
			sb.WriteString(fmt.Sprintf("+ %s=%s\n", e.Key, newVal))
		case "removed":
			sb.WriteString(fmt.Sprintf("- %s=%s\n", e.Key, oldVal))
		case "changed":
			sb.WriteString(fmt.Sprintf("~ %s: %s -> %s\n", e.Key, oldVal, newVal))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
