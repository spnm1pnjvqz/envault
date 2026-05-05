// Package main is the entry point for the envault CLI tool.
//
// envault provides commands to manage encrypted .env files:
//
//	envault init    - Initialize a new vault and generate an age key pair
//	envault lock    - Encrypt the .env file into .env.age
//	envault unlock  - Decrypt .env.age back into .env
//	envault view    - View secrets in-memory without writing to disk
//
// Configuration is stored in .envault.toml which should be committed
// to version control. The private key is stored separately (default:
// ~/.config/envault/key.txt) and must never be committed.
//
// Example workflow:
//
//	$ envault init
//	$ echo 'API_KEY=supersecret' > .env
//	$ envault lock
//	$ git add .env.age .envault.toml
//	$ envault unlock
package main
