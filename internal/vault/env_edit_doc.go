// Package vault — edit support.
//
// # Interactive Editing
//
// EditEnvFile decrypts the vault to a secure temporary file, launches the
// user's preferred editor, and re-encrypts the result once the editor exits.
//
// Editor resolution order:
//  1. EditOptions.Editor field (non-empty)
//  2. $EDITOR environment variable
//  3. Fallback: "vi"
//
// The temporary file is always removed after the operation, whether or not
// the edit succeeded. If the editor exits without making any changes the
// vault file is left untouched and EditResult.Modified is false.
//
// # EditResult
//
// EditResult summarises the diff between the pre- and post-edit state:
//
//	KeysAdded   — keys present in the new file but not in the original
//	KeysRemoved — keys present in the original but removed in the new file
//	KeysChanged — keys present in both but whose values differ
//
// This information is surfaced by the CLI `edit` command and can also be
// consumed programmatically by callers that embed the vault package.
package vault
