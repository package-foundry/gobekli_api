# CLAUDE.md - AI Agent Instructions

You are an expert Go developer working with the `gobekli_api` package. This file provides specific instructions for AI agents like Claude Code.

## Package Overview

`gobekli_api` provides API key generation, parsing, hashing, verification, and storage abstraction for Go applications.

## Package Location

- Module: `github.com/package-foundry/gobekli_api`
- Location: `/Users/gnrfan/code/libraries/golang/gobekli_api`

## Key Commands

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Build the package
go build ./...

# Run examples
go test -v -run TestExamples ./...

# Run specific test
go test -v -run TestGenerator_NewGenerator
```

## Package Structure

```
gobekli_api/
├── generator.go     # Key generation with crypto/rand + Base58
├── parse.go        # Key validation/parsing
├── hash.go        # SHA-256 hashing
├── verify.go     # Constant-time verification
├── store.go      # Storage abstraction interface + MemStore
├── base58.go     # Base58 encoding (internal)
├── apikey_test.go # Unit tests
├── examples_test.go # Usage examples
└── SKILL.md      # AI skill documentation
```

## Important Packages

The package uses these Go standard library packages:
- `crypto/rand` - Cryptographic randomness
- `crypto/sha256` - SHA-256 hashing
- `crypto/subtle` - Constant-time comparison
- `encoding/hex` - Hex encoding
- `strings` - String utilities
- `time` - Timestamps

## Security Guidelines

1. Never log or display plaintext API keys after creation
2. Always hash keys before storage (never store plaintext)
3. Use `VerifyKey()` with constant-time comparison
4. Implement the `Store` interface for production databases

## Usage Patterns

### Generate and Hash a Key
```go
gen, _ := apikey.NewGenerator("sk_live", 32)
key, _ := gen.Generate()
hash := apikey.HashKey(key) // Store this, not the key
```

### Verify a Key
```go
valid := apikey.VerifyKey(presentedKey, storedHash)
```

### Common Errors

When adding features:
- Use table-driven tests for edge cases
- Test uniqueness of generated keys
- Test error conditions
- Use constant-time comparison for verification

## Testing Requirements

Run tests before any code submission:
```bash
go test -v ./...
```

All 14 tests must pass.