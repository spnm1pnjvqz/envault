// Package vault provides high-level vault operations for envault.
//
// It ties together the crypto and env packages to offer a simple API for
// locking (encrypting) and unlocking (decrypting) .env files:
//
//	v, err := vault.New(publicKey, privateKey)
//
//	// Encrypt .env → .vault
//	err = v.Lock(".env", ".env.vault")
//
//	// Decrypt .vault → .env
//	err = v.Unlock(".env.vault", ".env")
//
//	// Inspect without writing to disk
//	entries, err := v.View(".env.vault")
//
// Key material is never written to disk by this package; callers are
// responsible for sourcing keys from a secure location (e.g. environment
// variables or a key file with restricted permissions).
package vault
