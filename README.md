# envault

> Lightweight secrets manager that encrypts `.env` files using [age](https://github.com/FiloSottile/age) encryption for local dev workflows.

---

## Installation

```bash
go install github.com/yourusername/envault@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envault/releases).

---

## Usage

**Encrypt a `.env` file:**

```bash
envault encrypt .env --output .env.age
```

**Decrypt and inject into your shell session:**

```bash
envault decrypt .env.age --output .env
```

**Run a command with decrypted secrets loaded into the environment:**

```bash
envault run --secret .env.age -- go run main.go
```

On first use, envault generates an age key pair stored at `~/.config/envault/keys.txt`. Share the public key with teammates so they can encrypt secrets for you, and keep your private key out of version control.

**Typical workflow:**

```bash
# Add .env and .env.age to .gitignore, commit only .env.age
echo ".env" >> .gitignore
git add .env.age .gitignore
git commit -m "chore: add encrypted secrets"
```

---

## Requirements

- Go 1.21+
- [age](https://github.com/FiloSottile/age) (bundled, no separate install needed)

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)