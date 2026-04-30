// Package env provides utilities for reading, writing, and manipulating
// .env files used by envault.
//
// A .env file consists of lines in the form:
//
//	KEY=value
//
// Lines beginning with '#' are treated as comments and are preserved
// during serialization. Quoted values (single or double) are automatically
// unquoted during parsing.
//
// Typical usage:
//
//	entries, err := env.ReadFile(".env")
//	if err != nil { ... }
//
//	m := env.ToMap(entries)
//	fmt.Println(m["DB_HOST"])
package env
