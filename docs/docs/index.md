# gobekli_api

A production-oriented Go package for API key generation, parsing, hashing, verification, and storage abstraction.

## Overview

`gobekli_api` provides a complete solution for generating and managing API keys in Go applications. It follows security best practices and is designed to be framework-agnostic, making it suitable for any Go service.

## Key Features

- **Cryptographically Secure**: Uses `crypto/rand` for secure random number generation
- **Base58 Encoding**: User-friendly encoding without ambiguous characters (0, O, I, l)
- **Hash Before Storage**: SHA-256 hashing for secure persistent storage
- **Constant-Time Verification**: Timing-safe comparison to prevent timing attacks
- **Storage Abstraction**: Pluggable storage interface for any database backend

## Key Format

```
<prefix>_<base58_random_part>
```

Example:
```
sk_live_4Er7pT9kQ2mX8vN6aBcD3fGh
```

## Quick Example

```go
package main

import (
    "log"
    "github.com/package-foundry/gobekli_api"
)

func main() {
    gen, err := apikey.NewGenerator("sk_live", 32)
    if err != nil {
        log.Fatal(err)
    }

    key, err := gen.Generate()
    if err != nil {
        log.Fatal(err)
    }

    // IMPORTANT: Show key to user once, then only store the hash
    hash := apikey.HashKey(key)
    
    // Store hash in database, never the plaintext key
    _ = hash
}
```

## Design Principles

1. **Separation of Concerns**: Generation, parsing, hashing, and verification are separate modules
2. **No Global State**: All configuration is explicit and instance-based
3. **Storage Agnostic**: Use any database by implementing the `Store` interface
4. **Security First**: Cryptographic best practices built-in by default

## Why Base58?

Base58 was chosen over Base64 because:
- Excludes visually ambiguous characters: 0, O, I, l
- More user-friendly for manual transcription  
- Compact encoding without padding
- Proven by Bitcoin and other cryptocurrencies

## Installation

```bash
go get github.com/package-foundry/gobekli_api
```

## Related Packages

This package is part of the Package Foundry ecosystem. See [github.com/package-foundry](https://github.com/package-foundry) for more packages.