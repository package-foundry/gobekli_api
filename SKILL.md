# gobekli_api Skill

This skill enables you to write code that uses the `gobekli_api` package for API key generation, parsing, hashing, verification, and storage abstraction in Go applications.

## Package Import

```go
import "github.com/package-foundry/gobekli_api"
```

For brevity in code examples, reference it as `apikey`.

## Package Purpose

`gobekli_api` provides a production-oriented Go package for:
- API key generation with cryptographically secure randomness
- Key validation and parsing
- SHA-256 hashing for secure storage
- Constant-time verification
- Storage abstraction via a Store interface

## Core Concepts

### 1. Key Generation

Generate cryptographically secure API keys:

```go
// Create generator with prefix "sk_live" and total key size 32
gen, err := apikey.NewGenerator("sk_live", 32)
if err != nil {
    return err
}

// Generate key
key, err := gen.Generate() // e.g., "sk_live_4Er7pT9kQ2mX8vN6aBcD3"
if err != nil {
    return err
}
```

Key format: `<prefix>_<base58_random>`

Valid prefixes: `sk_live`, `sk_test`, `pk_live`, `pk_test`, or custom
Valid key sizes: 8+ characters (typically 32)

### 2. Key Validation

Validate incoming keys:

```go
// Validate key format and prefix
err := apikey.ValidateKey(incomingKey, "sk_live", 32)
if err != nil {
    // Invalid key
    return err
}

// Extract prefix from key
prefix, err := apikey.ExtractPrefix(key)

// Extract payload (random part)
payload, err := apikey.ExtractPayload(key)
```

### 3. Secure Hashing

**IMPORTANT**: Never store plaintext API keys. Always hash before persistence:

```go
// Hash key (32 bytes SHA-256)
hash := apikey.HashKey(key)       // []byte

// Or get hex string
hashHex := apikey.HashKeyHex(key) // "a1b2c3..."

// Short fingerprint for display/logging (8 bytes)
fingerprint := apikey.Fingerprint(key) // "a1b2c3d4e5f60718"
```

### 4. Verification

Verify presented key against stored hash (constant-time):

```go
valid := apikey.VerifyKey(presentedKey, storedHash)
if !valid {
    return errors.New("invalid key")
}

// Or with hex hash
valid := apikey.VerifyKeyHex(presentedKey, storedHashHex)
```

Uses `crypto/subtle.ConstantTimeCompare` to prevent timing attacks.

### 5. Storage

Store key metadata (not plaintext):

```go
// Record structure
record := &apikey.Record{
    ID:          "key-001",
    Prefix:      "sk_live",
    Hash:        hash,
    Fingerprint: fingerprint,
    CreatedAt:   time.Now(),
    // RevokedAt:   nil (optional)
    // Disabled:   false (optional)
    // Metadata:   map[string]interface{}{} (optional)
}

// Check if active
if record.IsActive() {
    // Key is valid and not revoked/disabled
}
```

### 6. Storage Interface

Implement your own storage:

```go
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
```

In-memory example for testing:

```go
store := apikey.NewMemStore()
if err := store.Create(key, record); err != nil {
    return err
}

// Retrieve
retrieved, err := store.GetByHash(hash)

// Revoke
store.Revoke(keyID)
```

## Common Patterns

### Pattern 1: Create New API Key

```go
func CreateAPIKey(prefix string) (string, *apikey.Record, error) {
    gen, err := apikey.NewGenerator(prefix, 32)
    if err != nil {
        return "", nil, err
    }

    key, err := gen.Generate()
    if err != nil {
        return "", nil, err
    }

    hash := apikey.HashKey(key)
    fingerprint := apikey.Fingerprint(key)

    record := &apikey.Record{
        ID:          uuid.New().String(),
        Prefix:      prefix,
        Hash:       hash,
        Fingerprint: fingerprint,
        CreatedAt:  time.Now(),
    }

    return key, record, nil
}
```

### Pattern 2: Verify Incoming Key

```go
func VerifyAPIKey(store apikey.Store, presentedKey string) (*apikey.Record, error) {
    // Hash the presented key
    hash := apikey.HashKey(presentedKey)

    // Look up by hash
    record, err := store.GetByHash(hash)
    if err != nil {
        return nil, err // Not found
    }

    // Double-check with constant-time comparison
    if !apikey.VerifyKey(presentedKey, record.Hash) {
        return nil, errors.New("invalid key")
    }

    // Check if active
    if !record.IsActive() {
        return nil, errors.New("key revoked or disabled")
    }

    return record, nil
}
```

### Pattern 3: List Active Keys

```go
func ListActiveKeys(store apikey.Store) ([]*apikey.Record, error) {
    records, err := store.List()
    if err != nil {
        return nil, err
    }

    var active []*apikey.Record
    for _, r := range records {
        if r.IsActive() {
            active = append(active, r)
        }
    }

    return active, nil
}
```

### Pattern 4: Revoke Key

```go
func RevokeAPIKey(store apikey.Store, keyID string) error {
    record, err := store.GetByID(keyID)
    if err != nil {
        return err
    }

    if !record.IsActive() {
        return errors.New("key already revoked or disabled")
    }

    return store.Revoke(keyID)
}
```

## Error Handling

```go
import "github.com/package-foundry/gobekli_api"

var (
    ErrPrefixEmpty     = apikey.ErrPrefixEmpty
    ErrPrefixTooLong   = apikey.ErrPrefixTooLong
    ErrPrefixHasSpace  = apikey.ErrPrefixHasSpace
    ErrInvalidKeySize = apikey.ErrPrefixEmpty // note: ErrPrefixEmpty is reused
    
    ErrInvalidKey       = errors.New("invalid API key")
    ErrInvalidPrefix   = errors.New("invalid prefix")
    ErrKeyTooShort     = errors.New("key is too short")
    ErrMissingSeparator = errors.New("missing prefix separator")
    
    ErrRecordNotFound = errors.New("record not found")
    ErrDuplicateKey   = errors.New("key already exists")
)
```

## Constants

```go
const apikey.Separator = "_"  // Between prefix and random part
const apikey.DefaultKeySize() int // 32
```

## Best Practices

1. **Show key once**: Display the plaintext key to the user only at creation time
2. **Hash immediately**: Hash the key before storing; never store plaintext
3. **Use fingerprints**: Store the fingerprint for logging/display (not reversible)
4. **Track metadata**: Store prefix, created time, revoked time
5. **Implement Store**: Create a real database implementation for production

## Testing

The package includes an in-memory `MemStore` for testing:

```go
store := apikey.NewMemStore()
gen, _ := apikey.NewGenerator("sk_test", 32)
key, _ := gen.Generate()

record := &apikey.Record{
    ID:      "test-key-001",
    Prefix:  "sk_test",
}

if err := store.Create(key, record); err != nil {
    t.Fatal(err)
}
```

## File Structure

- `generator.go` - Key generation
- `parse.go` - Key validation/parsing
- `hash.go` - SHA-256 hashing
- `verify.go` - Constant-time verification
- `store.go` - Storage abstraction + in-memory impl
- `base58.go` - Base58 encoding (internal)

## GoDoc Reference

For complete API documentation, see [GoDoc](https://pkg.go.dev/github.com/package-foundry/gobekli_api).

## Trigger Keywords

Use this skill when:
- Working with API keys
- Generating keys
- Hashing for storage
- Verifying keys
- Managing key lifecycle (create, revoke, disable)
- Implementing key storage

Examples:
- "Generate an API key"
- "Hash the key for storage"
- "Verify the API key"
- "Create a key management service"
- "Implement key storage"
- "Revoke an API key"