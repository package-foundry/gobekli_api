# Quick Start

This guide gets you started with `gobekli_api` in 5 minutes.

## Step 1: Install

```bash
go get github.com/package-foundry/gobekli_api
```

## Step 2: Generate a Key

```go
package main

import (
    "fmt"
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

    fmt.Println("Your API key:", key)
}
```

Output:
```
Your API key: sk_live_4Er7pT9kQ2mX8vN6aBcD3
```

**IMPORTANT**: Show this key to the user once! After this, only the hash should be stored.

## Step 3: Hash for Storage

```go
hash := apikey.HashKey(key)
// hash is 32 bytes (SHA-256)

hashHex := apikey.HashKeyHex(key)
// hashHex is "a1b2c3d4..."

fingerprint := apikey.Fingerprint(key)
// fingerprint is "a1b2c3d4e5f60718" (for logging)
```

## Step 4: Store Metadata

```go
record := &apikey.Record{
    ID:          "key-001",
    Prefix:      "sk_live",
    Hash:       hash,
    Fingerprint: fingerprint,
    CreatedAt:  time.Now(),
}

store := apikey.NewMemStore()
if err := store.Create(key, record); err != nil {
    log.Fatal(err)
}
```

## Step 5: Verify an Incoming Key

```go
// Incoming key from request
incomingKey := "sk_live_4Er7pT9kQ2mX8vN6aBcD3"

// Hash it
presentedHash := apikey.HashKey(incomingKey)

// Look up
record, err := store.GetByHash(presentedHash)
if err != nil {
    // Key not found in store
    return errors.New("invalid key")
}

// Constant-time verify (defense in depth)
if !apikey.VerifyKey(incomingKey, record.Hash) {
    return errors.New("invalid key")
}

// Check if active
if !record.IsActive() {
    return errors.New("key revoked or disabled")
}

// Key is valid!
```

## Complete Example

```go
package main

import (
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/package-foundry/gobekli_api"
)

func main() {
    // 1. Create a new key
    key, record, err := createAPIKey("sk_live")
    if err != nil {
        log.Fatal(err)
    }

    // Show to user ONCE
    fmt.Println("API Key:", key)
    fmt.Println("Fingerprint:", record.Fingerprint)

    // 2. Store only the metadata (not plaintext)
    store := apikey.NewMemStore()
    if err := store.Create(key, record); err != nil {
        log.Fatal(err)
    }

    // 3. Verify a key
    validRecord, err := verifyKey(store, key)
    if err != nil {
        log.Printf("Verification failed: %v", err)
    } else {
        fmt.Println("Verified! ID:", validRecord.ID)
    }

    // 4. Revoke the key
    if err := store.Revoke(record.ID); err != nil {
        log.Fatal(err)
    }

    // 5. Try to verify again
    _, err = verifyKey(store, key)
    if err != nil {
        fmt.Println("As expected, key is revoked:", err)
    }
}

func createAPIKey(prefix string) (string, *apikey.Record, error) {
    gen, err := apikey.NewGenerator(prefix, 32)
    if err != nil {
        return "", nil, err
    }

    key, err := gen.Generate()
    if err != nil {
        return "", nil, err
    }

    hash := apikey.HashKey(key)
    record := &apikey.Record{
        ID:          "key-" + time.Now().Format("20060102150405"),
        Prefix:      prefix,
        Hash:       hash,
        Fingerprint: apikey.Fingerprint(key),
        CreatedAt:  time.Now(),
    }

    return key, record, nil
}

func verifyKey(store apikey.Store, key string) (*apikey.Record, error) {
    hash := apikey.HashKey(key)
    record, err := store.GetByHash(hash)
    if err != nil {
        return nil, err
    }

    if !apikey.VerifyKey(key, record.Hash) {
        return nil, errors.New("key mismatch")
    }

    if !record.IsActive() {
        return nil, errors.New("key not active")
    }

    return record, nil
}
```

## Next Steps

- Read the [Security Guide](security.md) for production deployment
- See [API Reference](api/generator.md) for complete documentation
- Implement the [Storage Interface](api/store.md) for your database