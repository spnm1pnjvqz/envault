package env

import (
	"strings"
	"testing"
)

func TestParseBasic(t *testing.T) {
	input := `# database config
DB_HOST=localhost
DB_PORT=5432
DB_NAME="mydb"
`
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	if entries[0].Comment != "# database config" {
		t.Errorf("expected comment, got %q", entries[0].Comment)
	}

	if entries[1].Key != "DB_HOST" || entries[1].Value != "localhost" {
		t.Errorf("unexpected entry: %+v", entries[1])
	}

	if entries[3].Value != "mydb" {
		t.Errorf("expected quotes stripped, got %q", entries[3].Value)
	}
}

func TestParseInvalidLine(t *testing.T) {
	input := "INVALID_LINE_NO_EQUALS\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid line")
	}
}

func TestParseEmptyKey(t *testing.T) {
	input := "=value\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestSerializeRoundtrip(t *testing.T) {
	original := "# comment\nKEY=value\nFOO=bar\n"
	entries, err := Parse(strings.NewReader(original))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result := Serialize(entries)
	if result != original {
		t.Errorf("roundtrip mismatch:\nwant: %q\ngot:  %q", original, result)
	}
}

func TestParseSingleQuotes(t *testing.T) {
	input := "SECRET='my secret value'\n"
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "my secret value" {
		t.Errorf("expected single quotes stripped, got %q", entries[0].Value)
	}
}
