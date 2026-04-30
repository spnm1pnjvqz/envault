package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadWriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	content := "# test\nKEY=value\nFOO=bar\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	entries, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	outPath := filepath.Join(dir, ".env.out")
	if err := WriteFile(outPath, entries); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	if string(data) != content {
		t.Errorf("content mismatch:\nwant: %q\ngot:  %q", content, string(data))
	}
}

func TestReadFileNotFound(t *testing.T) {
	_, err := ReadFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestToMap(t *testing.T) {
	entries := []Entry{
		{Comment: "# comment"},
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	m := ToMap(entries)
	if len(m) != 2 {
		t.Fatalf("expected 2 map entries, got %d", len(m))
	}
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected map values: %v", m)
	}
}

func TestMaskValue(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"secret123", "se*******"},
		{"ab", "**"},
		{"a", "*"},
		{"", ""},
	}
	for _, c := range cases {
		got := MaskValue(c.input)
		if got != c.want {
			t.Errorf("MaskValue(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
