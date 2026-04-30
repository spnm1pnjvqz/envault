package crypto

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"filippo.io/age"
	"filippo.io/age/armor"
)

// Encryptor handles age-based encryption and decryption of .env files.
type Encryptor struct {
	recipient age.Recipient
	identity  age.Identity
}

// NewEncryptor creates an Encryptor from an age public/private key pair.
func NewEncryptor(publicKey, privateKey string) (*Encryptor, error) {
	recipient, err := age.ParseX25519Recipient(publicKey)
	if err != nil {
		return nil, fmt.Errorf("invalid public key: %w", err)
	}

	identity, err := age.ParseX25519Identity(privateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	return &Encryptor{recipient: recipient, identity: identity}, nil
}

// EncryptFile reads a plaintext .env file and writes an armored encrypted file.
func (e *Encryptor) EncryptFile(src, dst string) error {
	plaintext, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading source file: %w", err)
	}

	var buf bytes.Buffer
	armorWriter := armor.NewWriter(&buf)

	w, err := age.Encrypt(armorWriter, e.recipient)
	if err != nil {
		return fmt.Errorf("initializing encryption: %w", err)
	}
	if _, err := w.Write(plaintext); err != nil {
		return fmt.Errorf("encrypting data: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("finalizing encryption: %w", err)
	}
	if err := armorWriter.Close(); err != nil {
		return fmt.Errorf("closing armor writer: %w", err)
	}

	if err := os.WriteFile(dst, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("writing encrypted file: %w", err)
	}
	return nil
}

// DecryptFile reads an armored encrypted file and writes the plaintext .env file.
func (e *Encryptor) DecryptFile(src, dst string) error {
	ciphertext, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading encrypted file: %w", err)
	}

	armorReader := armor.NewReader(bytes.NewReader(ciphertext))
	r, err := age.Decrypt(armorReader, e.identity)
	if err != nil {
		return fmt.Errorf("decrypting data: %w", err)
	}

	plaintext, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading decrypted stream: %w", err)
	}

	if err := os.WriteFile(dst, plaintext, 0600); err != nil {
		return fmt.Errorf("writing decrypted file: %w", err)
	}
	return nil
}
