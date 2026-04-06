# gobekli_api

<p align="center">
  <img src="assets/images/gobekli_api_art_01.png" alt="gobekli_api" width="400">
</p>

A production-oriented Go package for API key generation, parsing, hashing, verification, and storage abstraction.

[![GoDoc](https://godoc.org/github.com/package-foundry/gobekli_api?status.svg)](https://pkg.go.dev/github.com/package-foundry/gobekli_api)
[![Test](https://github.com/package-foundry/gobekli_api/actions/workflows/test.yml/badge.svg)](https://github.com/package-foundry/gobekli_api/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

`gobekli_api` provides a complete, secure, and reusable solution for generating and managing API keys in Go applications. It follows security best practices and is designed to be framework-agnostic, making it suitable for any Go service that needs to generate, store, and verify API keys.

The package is part of the [Package Foundry](https://github.com/package-foundry) ecosystem, designed for modular, production-ready Go libraries.

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Key Format](#key-format)
5. [Complete Usage Guide](#complete-usage-guide)
   - [Generating API Keys](#generating-api-keys)
   - [Hashing Keys for Storage](#hashing-keys-for-storage)
   - [Validating Keys](#validating-keys)
   - [Verifying Keys](#verifying-keys)
   - [Storing Key Metadata](#storing-key-metadata)
   - [Storage Interface](#storage-interface)
6. [API Reference](#api-reference)
   - [Constants](#constants)
   - [Error Variables](#error-variables)
   - [Generator](#generator)
   - [Parsing Functions](#parsing-functions)
   - [Hashing Functions](#hashing-functions)
   - [Verification Functions](#verification-functions)
   - [Storage Types](#storage-types)
7. [Complete Examples](#complete-examples)
   - [Basic Key Management](#basic-key-management)
   - [Full API Key Service](#full-api-key-service)
   - [Custom Store Implementation](#custom-store-implementation)
8. [Why Base58?](#why-base58)
9. [Security Considerations](#security-considerations)
10. [Design Decisions](#design-decisions)
11. [Testing](#testing)
12. [License](#license)

## Features

- **Cryptographically Secure**: Uses `crypto/rand` for secure random number generation
- **Base58 Encoding**: User-friendly encoding without ambiguous characters (0, O, I, l)
- **SHA-256 Hashing**: Secure one-way hashing for persistent storage
- **Constant-Time Verification**: Timing-safe comparison to prevent timing attacks
- **Storage Abstraction**: Pluggable `Store` interface for any database backend
- **Zero External Dependencies**: Pure Go standard library
- **Well-Tested**: Comprehensive unit test coverage

## Installation

```bash
go get github.com/package-foundry/gobekli_api
```

To use the package in your Go code:

```go
import (
    "github.com/package-foundry/gobekli_api"
)

// Most code uses the alias "apikey" for brevity
```

Verify the installation by building and testing:

```bash
go build ./...
go test -v ./...
```

For detailed installation instructions, see [Installation Guide](docs/docs/installation.md).

## Quick Start

The fastest way to get started with `gobekli_api`:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/package-foundry/gobekli_api"
)

func main() {
    // Step 1: Create a generator with your prefix and desired key size
    gen, err := apikey.NewGenerator("sk_live", 32)
    if err != nil {
        log.Fatal(err)  // Handle error: invalid prefix or key size
    }

    // Step 2: Generate a new API key
    key, err := gen.Generate()
    if err != nil {
        log.Fatal(err)  // Handle error: key generation failed
    }

    // Step 3: Hash the key for storage (NEVER store plaintext!)
    hash := apikey.HashKey(key)
    
    // Print results
    fmt.Println("Your API key:", key)
    fmt.Printf("Hash (hex): %x\n", hash)
}
```

Output:
```
Your API key: sk_live_4Er7pT9kQ2mX8vN6aBcD3
Hash (hex): a1b2c3d4e5f6789012345678901234567890123456789012345678901234abcd
```

**IMPORTANT**: Show the API key to the user only once at creation time! After that, only store the hash.

## Key Format

Keys are generated with the following structure:

```
<prefix>_<base58_random_part>
```

- **Prefix**: Identifies the key type (e.g., `sk_live`, `sk_test`, `pk_live`)
- **Separator**: Underscore character (`_`)
- **Random Part**: Cryptographically secure random bytes encoded in Base58

### Example Keys

```
sk_live_4Er7pT9kQ2mX8vN6aBcD3
sk_test_9Xf2qR8sT5uV1wY3zA7bC6d
pk_live_Nk4mL9pQ2rS6tU8vW0xY3z
pk_test_Qr5sT9uV2wX4yZ7aB0cD3e
```

### Common Prefixes

| Prefix | Type | Environment |
|--------|------|------------|
| `sk_live` | Secret Key | Production |
| `sk_test` | Secret Key | Testing |
| `pk_live` | Public Key | Production |
| `pk_test` | Public Key | Testing |

You can also use custom prefixes like `service_`, `internal_`, `partner_`, etc.

## Complete Usage Guide

### Generating API Keys

The `Generator` type creates cryptographically secure API keys. Each call to `Generate()` produces a unique key.

#### Creating a Generator

```go
// Create a generator for production secret keys with 32-character total length
gen, err := apikey.NewGenerator("sk_live", 32)
if err != nil {
    return fmt.Errorf("failed to create generator: %w", err)
}
```

```go
// Create a generator for testing keys with 64-character total length
gen, err := apikey.NewGenerator("sk_test", 64)
if err != nil {
    return fmt.Errorf("failed to create generator: %w", err)
}
```

```go
// Create a generator for public keys
gen, err := apikey.NewGenerator("pk_live", 32)
if err != nil {
    return fmt.Errorf("failed to create generator: %w", err)
}
```

```go
// Use custom prefixes for different services
gen, err := apikey.NewGenerator("myservice_api_key", 32)
if err != nil {
    return fmt.Errorf("failed to create generator: %w", err)
}
```

#### Generator Options

```go
// Set your own key size
gen, err := apikey.NewGenerator("sk_live", 48)  // 48 characters
if err != nil {
    return err
}

// Or use the default key size
gen, err := apikey.NewGenerator("sk_live", apikey.DefaultKeySize())  // 32
if err != nil {
    return err
}
```

```go
// Access generator properties
fmt.Println("Prefix:", gen.Prefix())    // "sk_live"
fmt.Println("Key Size:", gen.KeySize())  // 32
```

#### Generating Keys

```go
// Generate a single key
key, err := gen.Generate()
if err != nil {
    return fmt.Errorf("failed to generate key: %w", err)
}
fmt.Println("Generated key:", key)

// Each call produces a unique key
for i := 0; i < 5; i++ {
    key, _ := gen.Generate()
    fmt.Printf("Key %d: %s\n", i+1, key)
}
```

```go
// Example output:
// Key 1: sk_live_4Er7pT9kQ2mX8vN6aBcD3
// Key 2: sk_live_9Xf2qR8sT5uV1wY3zA7
// Key 3: sk_live_B8cD3eF6gH9iJ0kL2mN4
// Key 4: sk_live_O5pQ7rS9tU1vW3xY5zA8b
// Key 5: sk_live_C0dE2fG4hJ6kM8nO0pQ2
```

#### Error Handling

The generator can return errors for invalid configuration:

```go
// Empty prefix
_, err := apikey.NewGenerator("", 32)
if errors.Is(err, apikey.ErrPrefixEmpty) {
    fmt.Println("Prefix cannot be empty")
}

// Prefix with whitespace
_, err := apikey.NewGenerator("sk live", 32)
if errors.Is(err, apikey.ErrPrefixHasSpace) {
    fmt.Println("Prefix cannot contain whitespace")
}

// Key size too small
_, err := apikey.NewGenerator("sk_live", 5)
if errors.Is(err, apikey.ErrInvalidKeySize) {
    fmt.Println("Key size is too small")
}

// Prefix too long for key size
_, err := apikey.NewGenerator("sk_live", 9)
if errors.Is(err, apikey.ErrPrefixTooLong) {
    fmt.Println("Prefix is too long for the specified key size")
}
```

### Hashing Keys for Storage

**CRITICAL SECURITY NOTE**: Never store plaintext API keys. Always hash them before persistence.

#### Hash Functions

```go
// Generate a key
gen, _ := apikey.NewGenerator("sk_live", 32)
key, _ := gen.Generate()

// Hash the key - returns 32 bytes (SHA-256)
hash := apikey.HashKey(key)
fmt.Printf("Hash bytes: %d\n", len(hash))  // 32

// Hash as hex string - returns 64 characters
hashHex := apikey.HashKeyHex(key)
fmt.Printf("Hash hex: %s\n", hashHex)  // 64 character hex string

// Fingerprint for display - returns first 8 bytes as hex (16 chars)
fingerprint := apikey.Fingerprint(key)
fmt.Printf("Fingerprint: %s\n", fingerprint)  // 16 character hex

// Also available as bytes
fpBytes := apikey.FingerprintBytes(key)
fmt.Printf("Fingerprint bytes: %d\n", len(fpBytes))  // 8
```

```go
// Example output:
// Key:              sk_live_4Er7pT9kQ2mX8vN6aBcD3
// Hash bytes:       32
// Hash hex:        a1b2c3d4e5f6789012345678901234567890123456789012345678901234abcd
// Fingerprint:     a1b2c3d4e5f60718
// Fingerprint bytes: 8
```

#### Comparing Hashes

```go
// Hash two keys
hash1 := apikey.HashKey(key1)
hash2 := apikey.HashKey(key2)

// Compare hashes directly (not constant-time)
if string(hash1) == string(hash2) {
    fmt.Println("Same key")
}

// Or compare hex strings
if apikey.HashKeyHex(key1) == apikey.HashKeyHex(key2) {
    fmt.Println("Same key")
}
```

### Validating Keys

Validate the format of incoming keys before processing.

#### ValidateKey Function

```go
// Validate a key's format, prefix, and size
incomingKey := "sk_live_4Er7pT9kQ2mX8vN6aBcD3"

err := apikey.ValidateKey(incomingKey, "sk_live", 32)
if err != nil {
    return fmt.Errorf("invalid key: %w", err)
}

// This function checks:
// - Key is not empty
// - Has the correct prefix
// - Has the separator
// - Has the correct length (if expectedSize > 0)
```

```go
// Validation with size check
err := apikey.ValidateKey("sk_live_4Er7pT9kQ2mX8vN6aBcD3", "sk_live", 32)
if err != nil {
    fmt.Println("Invalid:", err)
}

// Validation without size check (0 allows any length)
err := apikey.ValidateKey("sk_live_4Er7pT9kQ2mX8vN6aBcD3", "sk_live", 0)
if err != nil {
    fmt.Println("Invalid:", err)
}
```

```go
// Validation error examples
err := apikey.ValidateKey("", "sk_live", 32)
fmt.Println(err)  // invalid API key

err := apikey.ValidateKey("wrong_4Er7pT9kQ2mX8vN6aBcD3", "sk_live", 32)
fmt.Println(err)  // invalid prefix

err := apikey.ValidateKey("sk_test_Ab", "sk_live", 32)
fmt.Println(err)  // key is too short

err := apikey.ValidateKey("sklive4Er7pT9kQ2mX8vN6aBcD3", "sk_live", 32)
fmt.Println(err)  // missing separator
```

#### Extract Prefix

```go
// Extract the prefix from a key
key := "sk_live_4Er7pT9kQ2mX8vN6aBcD3"
prefix, err := apikey.ExtractPrefix(key)
if err != nil {
    return err
}
fmt.Println("Prefix:", prefix)  // "sk_live"
```

#### Extract Payload

```go
// Extract the random payload from a key
key := "sk_live_4Er7pT9kQ2mX8vN6aBcD3"
payload, err := apikey.ExtractPayload(key)
if err != nil {
    return err
}
fmt.Println("Payload:", payload)  // "4Er7pT9kQ2mX8vN6aBcD3"
```

#### Has Valid Prefix

```go
// Check if key has one of several valid prefixes
key := "sk_live_4Er7pT9kQ2mX8vN6aBcD3"
validPrefixes := []string{"sk_live", "sk_test", "pk_live", "pk_test"}

if apikey.HasValidPrefix(key, validPrefixes) {
    fmt.Println("Valid prefix")
} else {
    fmt.Println("Unknown prefix")
}
```

### Verifying Keys

Verify a presented key against a stored hash using constant-time comparison.

#### VerifyKey Function

```go
// The key that will be verified
presentedKey := "sk_live_4Er7pT9kQ2mX8vN6aBcD3"

// The hash stored in the database (from when the key was created)
storedHash := apikey.HashKey(presentedKey)

// Verify the key - uses constant-time comparison
valid := apikey.VerifyKey(presentedKey, storedHash)
if valid {
    fmt.Println("Key is valid")
} else {
    fmt.Println("Key is invalid")
}
```

```go
// Verification returns false for invalid keys
invalidKey := "sk_live_wrongKey123456789"
valid := apikey.VerifyKey(invalidKey, storedHash)
fmt.Println("Valid:", valid)  // false
```

```go
// Verification with hex hashes
presentedKey := "sk_live_4Er7pT9kQ2mX8vN6aBcD3"
storedHashHex := apikey.HashKeyHex(presentedKey)

valid := apikey.VerifyKeyHex(presentedKey, storedHashHex)
if valid {
    fmt.Println("Key is valid")
}
```

#### Full Verification Flow

```go
// Complete verification example
func verifyAPIKey(keyToVerify string, storedHash []byte) (bool, error) {
    // First, validate the format
    if err := apikey.ValidateKey(keyToVerify, "sk_live", 32); err != nil {
        return false, fmt.Errorf("invalid format: %w", err)
    }
    
    // Then verify against the stored hash using constant-time comparison
    if !apikey.VerifyKey(keyToVerify, storedHash) {
        return false, nil  // Key doesn't match
    }
    
    return true, nil
}
```

### Storing Key Metadata

Use the `Record` struct to store key metadata along with the hash.

#### Record Structure

```go
import "time"

// Create a record for a new key
gen, _ := apikey.NewGenerator("sk_live", 32)
key, _ := gen.Generate()
hash := apikey.HashKey(key)

record := &apikey.Record{
    ID:          "key-001",                    // Unique identifier
    Prefix:      "sk_live",                 // Key prefix
    Hash:       hash,                     // SHA-256 hash (32 bytes)
    Fingerprint:  apikey.Fingerprint(key),     // Short fingerprint (8 bytes)
    CreatedAt:   time.Now(),             // Creation timestamp
    
    // Optional fields:
    // RevokedAt:   nil,                    // Revocation timestamp (set when revoked)
    // Disabled:  false,                 // Disabled flag
    // Metadata:  make(map[string]interface{}), // Custom metadata
}
```

```go
// Access record fields
fmt.Println("ID:", record.ID)
fmt.Println("Prefix:", record.Prefix)
fmt.Printf("Hash: %x\n", record.Hash)
fmt.Println("Fingerprint:", record.Fingerprint)
fmt.Println("Created:", record.CreatedAt)

// Optional fields
if record.RevokedAt != nil {
    fmt.Println("Revoked:", *record.RevokedAt)
}
fmt.Println("Disabled:", record.Disabled)
```

#### IsActive Method

```go
// Check if a key is active (not revoked and not disabled)
if record.IsActive() {
    fmt.Println("Key is active and can be used")
} else {
    fmt.Println("Key is not active")
}
```

```go
// IsActive returns false when disabled
record.Disabled = true
fmt.Println("Active:", record.IsActive())  // false

// IsActive returns false when revoked
record.Disabled = false
now := time.Now()
record.RevokedAt = &now
fmt.Println("Active:", record.IsActive())  // false
```

### Storage Interface

For production use, implement the `Store` interface for your database.

#### Store Interface Methods

```go
type Store interface {
    // Create stores a new key with its record
    Create(key string, record *Record) error
    
    // GetByID retrieves a record by its ID
    GetByID(id string) (*Record, error)
    
    // GetByHash retrieves a record by its hash
    GetByHash(hash []byte) (*Record, error)
    
    // Update updates an existing record
    Update(id string, record *Record) error
    
    // Delete removes a record by its ID
    Delete(id string) error
    
    // List returns all records
    List() ([]*Record, error)
    
    // Revoke marks a key as revoked
    Revoke(id string) error
    
    // Disable marks a key as disabled
    Disable(id string) error
}
```

#### Using MemStore (Testing)

```go
// Create an in-memory store for testing
store := apikey.NewMemStore()

// Create a key
gen, _ := apikey.NewGenerator("sk_live", 32)
key, _ := gen.Generate()

record := &apikey.Record{
    ID:     "key-001",
    Prefix: "sk_live",
}

// Store the key
err := store.Create(key, record)
if err != nil {
    return err
}
```

```go
// Get by ID
retrieved, err := store.GetByID("key-001")
if err != nil {
    return err
}
fmt.Println("Retrieved:", retrieved.ID)
```

```go
// Get by hash
hash := apikey.HashKey(key)
retrieved, err := store.GetByHash(hash)
if err != nil {
    return err
}
fmt.Println("Retrieved by hash:", retrieved.ID)
```

```go
// List all keys
allRecords, err := store.List()
for _, r := range allRecords {
    fmt.Println("Record:", r.ID, r.Prefix, r.Fingerprint)
}
```

```go
// Revoke a key
err := store.Revoke("key-001")
if err != nil {
    return err
}

// Now IsActive returns false
retrieved, _ := store.GetByID("key-001")
fmt.Println("Active:", retrieved.IsActive())  // false
```

```go
// Disable a key
err := store.Disable("key-001")
if err != nil {
    return err
}
```

```go
// Delete a key
err := store.Delete("key-001")
if err != nil {
    return err
}
```

## API Reference

### Constants

```go
// Separator is the character used between prefix and random part
const Separator = "_"

// DefaultKeySize is the recommended key size
const DefaultKeySize = 32
```

### Error Variables

```go
import "errors"

// Generator errors
var (
    ErrPrefixEmpty     = errors.New("prefix cannot be empty")
    ErrPrefixTooLong   = errors.New("prefix is too long for the specified key size")
    ErrPrefixHasSpace  = errors.New("prefix cannot contain whitespace")
    ErrInvalidKeySize = errors.New("key size must be greater than prefix + separator")
)

// Parsing errors
var (
    ErrInvalidKey       = errors.New("invalid API key")
    ErrInvalidPrefix   = errors.New("invalid prefix")
    ErrKeyTooShort     = errors.New("key is too short")
    ErrMissingSeparator = errors.New("missing prefix separator")
)

// Storage errors
var (
    ErrRecordNotFound = errors.New("record not found")
    ErrDuplicateKey   = errors.New("key already exists")
)
```

### Generator

```go
// NewGenerator creates a new key generator with the specified prefix and key size
func NewGenerator(prefix string, keySize int) (*Generator, error)

// DefaultKeySize returns the default key size (32)
func DefaultKeySize() int

// Generator holds the configuration for key generation
type Generator struct {
    prefix  string
    keySize int
}

// Generate creates a new API key
func (g *Generator) Generate() (string, error)

// Prefix returns the configured prefix
func (g *Generator) Prefix() string

// KeySize returns the configured key size
func (g *Generator) KeySize() int
```

### Parsing Functions

```go
// ValidateKey validates a key's format, prefix, and size
func ValidateKey(key string, expectedPrefix string, expectedSize int) error

// ExtractPrefix extracts the prefix from a key
func ExtractPrefix(key string) (string, error)

// ExtractPayload extracts the random payload from a key
func ExtractPayload(key string) (string, error)

// HasValidPrefix checks if a key has any of the valid prefixes
func HasValidPrefix(key string, validPrefixes []string) bool
```

### Hashing Functions

```go
// HashKey hashes a key using SHA-256, returns 32 bytes
func HashKey(key string) []byte

// HashKeyHex hashes a key and returns hex string (64 chars)
func HashKeyHex(key string) string

// Fingerprint returns first 8 bytes as hex (16 chars)
func Fingerprint(key string) string

// FingerprintBytes returns first 8 bytes
func FingerprintBytes(key string) []byte
```

### Verification Functions

```go
// VerifyKey compares a key against a stored hash using constant-time comparison
func VerifyKey(presented string, storedHash []byte) bool

// VerifyKeyHex compares a key against a stored hex hash using constant-time comparison
func VerifyKeyHex(presented string, storedHashHex string) bool
```

### Storage Types

```go
import "time"

// Record represents key metadata stored in the database
type Record struct {
    ID           string                 // Unique identifier
    Prefix       string                // Key prefix (sk_live, sk_test, etc.)
    Hash         []byte               // SHA-256 hash (32 bytes)
    Fingerprint  string              // Short fingerprint (8 bytes as hex)
    CreatedAt   time.Time          // Creation timestamp
    RevokedAt   *time.Time       // Revocation timestamp (nil if not revoked)
    Disabled    bool              // Disabled flag
    Metadata    map[string]interface{} // Custom metadata
}

// IsActive returns true if the key is valid
func (r *Record) IsActive() bool

// Store defines the interface for key storage
type Store interface {
    Create(key string, record *Record) error
    GetByID(id string) (*Record, error)
    GetByHash(hash []byte) (*Record, error)
    Update(id string, record *Record) error
    Delete(id string) error
    List() ([]*Record, error)
    Revoke(id string) error
    Disable(id string) error
}

// NewMemStore creates an in-memory store (for testing)
func NewMemStore() *MemStore
```

## Complete Examples

### Basic Key Management

This example demonstrates creating, storing, and verifying API keys:

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/package-foundry/gobekli_api"
)

func main() {
    // 1. Create an in-memory store
    store := apikey.NewMemStore()

    // 2. Create a generator
    gen, err := apikey.NewGenerator("sk_live", 32)
    if err != nil {
        log.Fatal(err)
    }

    // 3. Generate some keys
    for i := 0; i < 3; i++ {
        key, err := gen.Generate()
        if err != nil {
            log.Fatal(err)
        }

        // Hash for storage
        hash := apikey.HashKey(key)
        fingerprint := apikey.Fingerprint(key)

        // Create record
        record := &apikey.Record{
            ID:          fmt.Sprintf("key-%d", i+1),
            Prefix:      gen.Prefix(),
            Hash:       hash,
            Fingerprint: fingerprint,
            CreatedAt:   time.Now(),
        }

        // Store it
        if err := store.Create(key, record); err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Created: %s (Fingerprint: %s)\n", key, fingerprint)
    }

    // 4. List all keys
    fmt.Println("\nAll keys:")
    records, _ := store.List()
    for _, r := range records {
        status := "active"
        if !r.IsActive() {
            status = "inactive"
        }
        fmt.Printf("  - %s: %s [%s]\n", r.ID, r.Fingerprint, status)
    }

    // 5. Revoke the second key
    fmt.Println("\nRevoking key-2...")
    if err := store.Revoke("key-2"); err != nil {
        log.Fatal(err)
    }

    // 6. List keys again
    fmt.Println("\nAfter revocation:")
    records, _ = store.List()
    for _, r := range records {
        status := "active"
        if !r.IsActive() {
            status = "revoked"
        }
        fmt.Printf("  - %s: %s [%s]\n", r.ID, r.Fingerprint, status)
    }
}
```

### Full API Key Service

A complete API key service implementation:

```go
package main

import (
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/package-foundry/gobekli_api"
)

// APIKeyService manages the lifecycle of API keys
type APIKeyService struct {
    store apikey.Store
    gen   *apikey.Generator
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(store apikey.Store, prefix string) (*APIKeyService, error) {
    if store == nil {
        return nil, errors.New("store is required")
    }

    gen, err := apikey.NewGenerator(prefix, apikey.DefaultKeySize())
    if err != nil {
        return nil, err
    }

    return &APIKeyService{
        store: store,
        gen:   gen,
    }, nil
}

// CreateKey creates a new API key and returns it once
func (s *APIKeyService) CreateKey() (string, *apikey.Record, error) {
    key, err := s.gen.Generate()
    if err != nil {
        return "", nil, err
    }

    hash := apikey.HashKey(key)
    fingerprint := apikey.Fingerprint(key)

    record := &apikey.Record{
        ID:          fmt.Sprintf("key-%d", time.Now().UnixNano()),
        Prefix:      s.gen.Prefix(),
        Hash:       hash,
        Fingerprint: fingerprint,
        CreatedAt:  time.Now(),
    }

    if err := s.store.Create(key, record); err != nil {
        return "", nil, err
    }

    return key, record, nil
}

// VerifyKey verifies an API key and returns the record if valid
func (s *APIKeyService) VerifyKey(key string) (*apikey.Record, error) {
    // Validate format
    if err := apikey.ValidateKey(key, s.gen.Prefix(), s.gen.KeySize()); err != nil {
        return nil, fmt.Errorf("invalid key format: %w", err)
    }

    // Hash and look up
    hash := apikey.HashKey(key)
    record, err := s.store.GetByHash(hash)
    if err != nil {
        return nil, fmt.Errorf("key not found: %w", err)
    }

    // Constant-time verify
    if !apikey.VerifyKey(key, record.Hash) {
        return nil, errors.New("invalid key")
    }

    // Check if active
    if !record.IsActive() {
        return nil, errors.New("key is revoked or disabled")
    }

    return record, nil
}

// RevokeKey revokes an API key by ID
func (s *APIKeyService) RevokeKey(id string) error {
    record, err := s.store.GetByID(id)
    if err != nil {
        return err
    }

    if !record.IsActive() {
        return errors.New("key already revoked or disabled")
    }

    return s.store.Revoke(id)
}

// DisableKey disables an API key by ID
func (s *APIKeyService) DisableKey(id string) error {
    record, err := s.store.GetByID(id)
    if err != nil {
        return err
    }

    if !record.IsActive() {
        return errors.New("key already revoked or disabled")
    }

    return s.store.Disable(id)
}

// ListKeys returns all key records
func (s *APIKeyService) ListKeys() ([]*apikey.Record, error) {
    return s.store.List()
}

func main() {
    store := apikey.NewMemStore()
    service, err := NewAPIKeyService(store, "sk_live")
    if err != nil {
        log.Fatal(err)
    }

    // Create new keys
    fmt.Println("Creating keys...")
    for i := 0; i < 3; i++ {
        key, record, err := service.CreateKey()
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Created: %s (ID: %s)\n", key, record.ID)
    }

    // Verify a key
    fmt.Println("\nVerifying key...")
    key, _, _ := service.CreateKey()
    record, err := service.VerifyKey(key)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Verified: %s\n", record.ID)

    // Try to verify again after revocation
    fmt.Println("\nRevoking key...")
    if err := service.RevokeKey(record.ID); err != nil {
        log.Fatal(err)
    }

    // Should fail now
    _, err = service.VerifyKey(key)
    if err != nil {
        fmt.Printf("As expected, verification failed: %s\n", err)
    }
}
```

### Custom Store Implementation

Example of implementing a custom store for a SQL database:

```go
package myapp

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/package-foundry/gobekli_api"
)

// SQLStore implements apikey.Store for a SQL database
type SQLStore struct {
    db *sql.DB
}

// NewSQLStore creates a new SQL store
func NewSQLStore(db *sql.DB) *SQLStore {
    return &SQLStore{db: db}
}

// Create inserts a new key record
func (s *SQLStore) Create(key string, record *apikey.Record) error {
    query := `
        INSERT INTO api_keys (id, prefix, hash, fingerprint, created_at, revoked_at, disabled)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
    _, err := s.db.Exec(query,
        record.ID,
        record.Prefix,
        record.Hash,
        record.Fingerprint,
        record.CreatedAt,
        record.RevokedAt,
        record.Disabled,
    )
    return err
}

// GetByID retrieves a record by ID
func (s *SQLStore) GetByID(id string) (*apikey.Record, error) {
    query := `SELECT id, prefix, hash, fingerprint, created_at, revoked_at, disabled FROM api_keys WHERE id = ?`
    return s.scanRow(s.db.QueryRow(query, id))
}

// GetByHash retrieves a record by hash
func (s *SQLStore) GetByHash(hash []byte) (*apikey.Record, error) {
    query := `SELECT id, prefix, hash, fingerprint, created_at, revoked_at, disabled FROM api_keys WHERE hash = ?`
    return s.scanRow(s.db.QueryRow(query, hash))
}

// Update updates an existing record
func (s *SQLStore) Update(id string, record *apikey.Record) error {
    query := `
        UPDATE api_keys
        SET prefix = ?, hash = ?, fingerprint = ?, revoked_at = ?, disabled = ?
        WHERE id = ?
    `
    _, err := s.db.Exec(query,
        record.Prefix,
        record.Hash,
        record.Fingerprint,
        record.RevokedAt,
        record.Disabled,
        id,
    )
    return err
}

// Delete removes a record
func (s *SQLStore) Delete(id string) error {
    _, err := s.db.Exec("DELETE FROM api_keys WHERE id = ?", id)
    return err
}

// List returns all records
func (s *SQLStore) List() ([]*apikey.Record, error) {
    query := `SELECT id, prefix, hash, fingerprint, created_at, revoked_at, disabled FROM api_keys`
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var records []*apikey.Record
    for rows.Next() {
        r, err := s.scanRow(rows)
        if err != nil {
            return nil, err
        }
        records = append(records, r)
    }
    return records, nil
}

// Revoke marks a key as revoked
func (s *SQLStore) Revoke(id string) error {
    query := `UPDATE api_keys SET revoked_at = ? WHERE id = ?`
    _, err := s.db.Exec(query, time.Now(), id)
    return err
}

// Disable marks a key as disabled
func (s *SQLStore) Disable(id string) error {
    query := `UPDATE api_keys SET disabled = true WHERE id = ?`
    _, err := s.db.Exec(query, id)
    return err
}

// scanRow helper function (simplified)
func (s *SQLStore) scanRow(row *sql.Row) (*apikey.Record, error) {
    var r apikey.Record
    var createdAt, revokedAt sql.NullTime
    var disabled bool
    var hash []byte

    err := row.Scan(&r.ID, &r.Prefix, &hash, &r.Fingerprint, &createdAt, &revokedAt, &disabled)
    if err != nil {
        return nil, err
    }

    r.Hash = hash
    r.CreatedAt = createdAt.Time
    r.Disabled = disabled
    if revokedAt.Valid {
        r.RevokedAt = &revokedAt.Time
    }

    return &r, nil
}
```

## Why Base58?

Base58 encoding offers several important advantages over Base64 for API keys:

1. **No Ambiguous Characters**
   - Excludes: 0 (zero), O (capital o), I (capital i), l (lowercase L)
   - These characters look similar in many fonts and can cause confusion

2. **Compact Encoding**
   - No padding required (= signs)
   - Shorter keys for the same amount of entropy

3. **Proven Standard**
   - Used by Bitcoin, Ethereum, and other cryptocurrencies
   - Well-tested and battle-hardened

4. **URL-Safe**
   - Doesn't require URL encoding
   - Suitable for headers, query parameters, and more

Example comparison:
```
Base64: SkRvMzJEN2U5Rk0xcjNHN2U5Rk0xcjNHN2U5Rk0xcg==
Base58: SkRvMzJDN2U5Rk0xcjNHN2U5Rk0xcjNH
```

## Security Considerations

### 1. One-Time Display

```go
// RIGHT: Show the key once, then never again
key, record, err := service.CreateKey()
if err != nil {
    return err
}
displayToUser(key)  // Show only once!
fmt.Println("Key ID:", record.ID)  // Only show the ID after this
```

```go
// WRONG: Storing plaintext
// Never do this:
db.Exec("INSERT INTO keys (key) VALUES (?)", key)
```

### 2. Constant-Time Comparison

```go
// Uses crypto/subtle.ConstantTimeCompare to prevent timing attacks
valid := apikey.VerifyKey(presentedKey, storedHash)
```

### 3. Cryptographic Randomness

```go
// Uses crypto/rand (NOT math/rand)
randomBytes := make([]byte, 32)
if _, err := rand.Read(randomBytes); err != nil {
    return err
}
```

### 4. One-Way Hashing

```go
// SHA-256 is computationally infeasible to reverse
hash := apikey.HashKey(key)
// Cannot recover key from hash!
```

### 5. No Key Recovery

```go
// If a key is lost, generate a new one
// There is no way to recover the original key from the hash
```

## Design Decisions

### Adapted from Reference Python Package (cx-apikey)

1. **Factory Pattern**: The `Generator` struct provides a clean factory for creating keys with consistent prefix/size.

2. **Prefix Validation**: Strict validation of prefix format (no whitespace, reasonable length).

3. **Record Metadata**: The `Record` struct stores key metadata separately from the key itself.

4. **Validation Functions**: Clear separation between key generation and key validation.

### Deliberate Changes for Idiomatic Go

1. **No Global State**: All configuration is explicit and instance-based.

2. **Interface-Based Storage**: The `Store` interface allows any storage backend.

3. **Built-In Base58**: A minimal Base58 implementation is included.

4. **Explicit Error Handling**: Go errors are explicit, not exceptions.

5. **Separate Concerns**: Generation, parsing, hashing, and verification are in separate files.

## Testing

Run all tests:

```bash
go test -v ./...
```

Run with coverage:

```bash
go test -v -cover ./...
```

Run examples:

```bash
go test -v -run TestExamples ./...
```

### Test Categories

- **Unit Tests**: Test individual functions
- **Table-Driven Tests**: Test edge cases with various inputs
- **Uniqueness Tests**: Ensure generated keys are unique
- **Error Tests**: Test error conditions
- **Example Tests**: Test complete workflows

## License

MIT License - Copyright (c) 2026 Antonio Ognio

See [LICENSE](LICENSE) for full details.

## Support

- **Issues**: https://github.com/package-foundry/gobekli_api/issues
- **Discussions**: https://github.com/package-foundry/gobekli_api/discussions

---

Maintained with :heart_emoji: from Peru :peruvian_flag_emoji:. El Perú es clave. :key_emoji: