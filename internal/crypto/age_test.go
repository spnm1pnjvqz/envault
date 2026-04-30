package crypto

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
)

func generateTestKeys(t *testing.T) (pub, priv string) {
	t.Helper()
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generating age identity: %v", err)
	}
	return identity.Recipient().String(), identity.String()
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	pub, priv := generateTestKeys(t)

	enc, err := NewEncryptor(pub, priv)
	if err != nil {
		t.Fatalf("NewEncryptor: %v", err)
	}

	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, ".env")
	encFile := filepath.Join(tmpDir, ".env.age")
	dstFile := filepath.Join(tmpDir, ".env.decrypted")

	original := []byte("API_KEY=secret123\nDB_PASS=hunter2\n")
	if err := os.WriteFile(srcFile, original, 0600); err != nil {
		t.Fatalf("writing source file: %v", err)
	}

	if err := enc.EncryptFile(srcFile, encFile); err != nil {
		t.Fatalf("EncryptFile: %v", err)
	}

	encContent, _ := os.ReadFile(encFile)
	if string(encContent) == string(original) {
		t.Error("encrypted content should differ from plaintext")
	}

	if err := enc.DecryptFile(encFile, dstFile); err != nil {
		t.Fatalf("DecryptFile: %v", err)
	}

	decrypted, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("reading decrypted file: %v", err)
	}

	if string(decrypted) != string(original) {
		t.Errorf("expected %q, got %q", original, decrypted)
	}
}

func TestNewEncryptorInvalidKeys(t *testing.T) {
	_, err := NewEncryptor("not-a-key", "not-a-key")
	if err == nil {
		t.Error("expected error for invalid keys, got nil")
	}
}

func TestDecryptFileNotFound(t *testing.T) {
	pub, priv := generateTestKeys(t)
	enc, _ := NewEncryptor(pub, priv)

	err := enc.DecryptFile("/nonexistent/path.age", "/tmp/out.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
