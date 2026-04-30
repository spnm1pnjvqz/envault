package env

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
}

// Parse reads .env formatted content from r and returns a slice of entries.
// It preserves comments and blank lines are skipped.
func Parse(r io.Reader) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			entries = append(entries, Entry{Comment: line})
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = stripQuotes(value)

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}

		entries = append(entries, Entry{Key: key, Value: value})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return entries, nil
}

// Serialize converts a slice of entries back to .env file format.
func Serialize(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		if e.Comment != "" {
			sb.WriteString(e.Comment + "\n")
		} else {
			sb.WriteString(e.Key + "=" + e.Value + "\n")
		}
	}
	return sb.String()
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
