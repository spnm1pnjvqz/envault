package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envault/internal/vault"
)

func writeTemplateTestEnv(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func writeTemplateTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func setupTemplateConfig(t *testing.T, dir string) string {
	t.Helper()
	cfgPath := filepath.Join(dir, ".envault.json")
	cfg := &vault.Config{EnvFile: filepath.Join(dir, ".env")}
	if err := vault.SaveConfig(cfgPath, cfg); err != nil {
		t.Fatal(err)
	}
	return cfgPath
}

func TestTemplateCommandStdout(t *testing.T) {
	dir := t.TempDir()
	writeTemplateTestEnv(t, dir, "GREETING=hello\nNAME=world\n")
	tmpl := writeTemplateTestFile(t, dir, "msg.tmpl", "${GREETING}, ${NAME}!\n")
	cfgPath := setupTemplateConfig(t, dir)

	out, err := executeCommand(rootCmd, "--config", cfgPath, "template", tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "hello, world!") {
		t.Errorf("expected rendered output, got: %s", out)
	}
}

func TestTemplateCommandOutputFile(t *testing.T) {
	dir := t.TempDir()
	writeTemplateTestEnv(t, dir, "HOST=localhost\n")
	tmpl := writeTemplateTestFile(t, dir, "db.conf", "host=${HOST}\n")
	outFile := filepath.Join(dir, "db.rendered")
	cfgPath := setupTemplateConfig(t, dir)

	_, err := executeCommand(rootCmd, "--config", cfgPath, "template", "--output", outFile, tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(outFile)
	if !strings.Contains(string(data), "host=localhost") {
		t.Errorf("expected rendered file content, got: %s", string(data))
	}
}

func TestTemplateCommandMissingConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "--config", "/nonexistent/.envault.json", "template", "any.tmpl")
	if err == nil {
		t.Error("expected error for missing config")
	}
}

func TestTemplateCommandMissingTemplate(t *testing.T) {
	dir := t.TempDir()
	writeTemplateTestEnv(t, dir, "KEY=val\n")
	cfgPath := setupTemplateConfig(t, dir)

	_, err := executeCommand(rootCmd, "--config", cfgPath, "template", filepath.Join(dir, "ghost.tmpl"))
	if err == nil {
		t.Error("expected error for missing template")
	}
}
