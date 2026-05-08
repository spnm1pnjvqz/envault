package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemplateEnv(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func writeTemplateFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func newTemplateVault(t *testing.T, dir string) *Vault {
	t.Helper()
	cfg := &Config{EnvFile: filepath.Join(dir, ".env")}
	v, err := New(cfg, "")
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestRenderTemplateBasic(t *testing.T) {
	dir := t.TempDir()
	writeTemplateEnv(t, dir, "APP_NAME=envault\nPORT=8080\n")
	tmpl := writeTemplateFile(t, dir, "app.conf", "name=${APP_NAME}\nport=$PORT\n")
	v := newTemplateVault(t, dir)

	res, err := v.RenderTemplate(tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Substituted != 2 {
		t.Errorf("expected 2 substitutions, got %d", res.Substituted)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing vars, got %v", res.Missing)
	}
	if !strings.Contains(res.Output, "name=envault") {
		t.Errorf("expected APP_NAME substituted, got: %s", res.Output)
	}
}

func TestRenderTemplateMissingVars(t *testing.T) {
	dir := t.TempDir()
	writeTemplateEnv(t, dir, "APP_NAME=envault\n")
	tmpl := writeTemplateFile(t, dir, "app.conf", "name=${APP_NAME}\ndb=${DB_HOST}\n")
	v := newTemplateVault(t, dir)

	res, err := v.RenderTemplate(tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST missing, got %v", res.Missing)
	}
	if !strings.Contains(res.Output, "${DB_HOST}") {
		t.Errorf("expected unresolved placeholder preserved")
	}
}

func TestRenderTemplateFileNotFound(t *testing.T) {
	dir := t.TempDir()
	writeTemplateEnv(t, dir, "KEY=val\n")
	v := newTemplateVault(t, dir)
	_, err := v.RenderTemplate(filepath.Join(dir, "nonexistent.tmpl"))
	if err == nil {
		t.Error("expected error for missing template file")
	}
}

func TestRenderTemplateEnvNotFound(t *testing.T) {
	dir := t.TempDir()
	tmpl := writeTemplateFile(t, dir, "app.conf", "val=${KEY}\n")
	cfg := &Config{EnvFile: filepath.Join(dir, "missing.env")}
	v := &Vault{cfg: cfg}
	_, err := v.RenderTemplate(tmpl)
	if err == nil {
		t.Error("expected error when env file missing")
	}
}
