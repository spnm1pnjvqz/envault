package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"

	"github.com/user/envault/internal/vault"
)

func generateTestKeys(t *testing.T) (pub, priv string) {
	t.Helper()
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("failed to generate age key: %v", err)
	}
	return identity.Recipient().String(), identity.String()
}

func TestLockUnlockRoundtrip(t *testing.T) {
	pub, priv := generateTestKeys(t)

	v, err := vault.New(pub, priv)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, ".vault")
	outPath := filepath.Join(dir, "decrypted.env")

	original := "APP_ENV=production\nSECRET_KEY=supersecret\n"
	if err := os.WriteFile(envPath, []byte(original), 0o600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	if err := v.Lock(envPath, vaultPath); err != nil {
		t.Fatalf("Lock: %v", err)
	}

	if _, err := os.Stat(vaultPath); err != nil {
		t.Fatalf("vault file not created: %v", err)
	}

	if err := v.Unlock(vaultPath, outPath); err != nil {
		t.Fatalf("Unlock: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	if string(got) != original {
		t.Errorf("roundtrip mismatch:\ngot:  %q\nwant: %q", got, original)
	}
}

func TestView(t *testing.T) {
	pub, priv := generateTestKeys(t)

	v, err := vault.New(pub, priv)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, ".vault")

	if err := os.WriteFile(envPath, []byte("FOO=bar\nBAZ=qux\n"), 0o600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	if err := v.Lock(envPath, vaultPath); err != nil {
		t.Fatalf("Lock: %v", err)
	}

	entries, err := v.View(vaultPath)
	if err != nil {
		t.Fatalf("View: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestLockMissingEnvFile(t *testing.T) {
	pub, priv := generateTestKeys(t)

	v, err := vault.New(pub, priv)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := v.Lock("/nonexistent/.env", ""); err == nil {
		t.Error("expected error for missing env file, got nil")
	}
}
