# AGENTS.md - Guidelines for AI Agents

This file provides guidance for AI agents interacting with this codebase.

## Package Purpose

`gobekli_api` is a Go package for API key generation, parsing, hashing, verification, and storage abstraction.

## Key Commands

```bash
# Run tests
go test -v ./...

# Build
go build ./...

# Run examples
go test -v -run TestExamples ./...
```

## Code Structure

```
gobekli_api/
├── generator.go    # Key generation
├── parse.go       # Key validation/parsing
├── hash.go        # SHA-256 hashing
├── verify.go     # Constant-time verification
├── store.go      # Storage abstraction
├── base58.go    # Base58 encoding
├── apikey_test.go # Unit tests
├── examples_test.go # Usage examples
├── LICENSE       # MIT License
├── README.md    # Package documentation
├── go.mod       # Go module
└── AGENTS.md    # This file
```

## Usage Pattern

```go
// 1. Generate key
gen, _ := apikey.NewGenerator("sk_live", 32)
key, _ := gen.Generate() // "sk_live_4Er7pT9kQ2mX8vN6aBcD3"

// 2. Hash for storage (never store plaintext)
hash := apikey.HashKey(key)

// 3. Hash is 32 bytes (SHA-256)
fmt.Printf("%x", hash)

// 4. Verify using constant-time comparison
valid := apikey.VerifyKey(presentedKey, storedHash)
```

## Security Considerations

1. Never log or display plaintext API keys after initial creation
2. Always hash keys before storing
3. Use `VerifyKey()` with constant-time comparison
4. Implement the `Store` interface for production databases

## Testing Requirements

When modifying this package:
- Ensure all tests pass: `go test -v ./...`
- Test key uniqueness across multiple generations
- Test error cases for validation
- Verify constant-time property of verification

## Dependencies

This package has no external dependencies outside Go standard library:
- `crypto/rand` - Cryptographic randomness
- `crypto/sha256` - SHA-256 hashing
- `crypto/subtle` - Constant-time comparison
- `encoding/hex` - Hex encoding
- `strings` - String utilities
- `time` - Timestamps