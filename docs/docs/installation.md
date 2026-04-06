# Installation

## Requirements

- Go 1.21 or later

## Install

```bash
go get github.com/package-foundry/gobekli_api
```

## Verify

```bash
go build ./...
```

## Import

```go
import "github.com/package-foundry/gobekli_api"
```

The package uses the import path `github.com/package-foundry/gobekli_api` (or just `gobekli_api` if using Go modules with the module name). All public symbols are prefixed with `apikey.` (lowercase for package, though in code you'll use `apikey.` prefix).

## Dependencies

This package has no external dependencies outside the Go standard library:

- `crypto/rand` - Cryptographic randomness
- `crypto/sha256` - SHA-256 hashing
- `crypto/subtle` - Constant-time comparison
- `encoding/hex` - Hex encoding
- `strings` - String utilities
- `time` - Timestamps