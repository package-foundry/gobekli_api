# API Reference

## Generator

Functions and types for generating API keys.

### Variables

```go
var (
    ErrPrefixEmpty     = errors.New("prefix cannot be empty")
    ErrPrefixTooLong   = errors.New("prefix is too long for the specified key size")
    ErrPrefixHasSpace  = errors.New("prefix cannot contain whitespace")
    ErrInvalidKeySize  = errors.New("key size must be greater than prefix + separator")
)
```

### Constants

```go
const Separator = "_"
```

### Types

```go
type Generator struct {
    prefix  string
    keySize int
}
```

### Functions

```go
func NewGenerator(prefix string, keySize int) (*Generator, error)
```

Creates a new key generator with the specified prefix and key size.

**Parameters:**
- `prefix` - The key prefix (e.g., "sk_live", "sk_test")
- `keySize` - Total key size in characters

**Errors:**
- `ErrPrefixEmpty` - If prefix is empty
- `ErrPrefixHasSpace` - If prefix contains whitespace
- `ErrInvalidKeySize` - If key size is too small
- `ErrPrefixTooLong` - If prefix + separator exceeds key size

### Methods

```go
func (g *Generator) Generate() (string, error)
```

Generates a new API key.

**Returns:** A key string with format `<prefix>_<base58_random>`

```go
func (g *Generator) Prefix() string
```

Returns the configured prefix.

```go
func (g *Generator) KeySize() int
```

Returns the configured key size.

```go
func DefaultKeySize() int
```

Returns the default key size (32).

## Usage

```go
gen, err := apikey.NewGenerator("sk_live", 32)
if err != nil {
    log.Fatal(err)
}

key, err := gen.Generate()
if err != nil {
    log.Fatal(err)
}

fmt.Println(key) // sk_live_4Er7pT9kQ2mX8vN6aBcD3
```