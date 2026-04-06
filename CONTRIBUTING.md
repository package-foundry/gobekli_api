# Contributing to gobekli_api

Thank you for your interest in contributing to `gobekli_api`!

## Quick Summary

**Agentic modifications are welcome.** This package is designed to be enhanced by AI agents and human developers alike. We're interested in improving:

- Key generation improvements
- New storage backends
- Additional hashing algorithms
- Performance optimizations
- Bug fixes
- Documentation improvements

## How to Contribute

### 1. Fork and Clone

```bash
git fork https://github.com/package-foundry/gobekli_api
git clone https://github.com/<your-username>/gobekli_api
cd gobekli_api
```

### 2. Make Changes

Follow the package structure:
- `generator.go` - Key generation logic
- `parse.go` - Key validation
- `hash.go` - Hashing functions
- `verify.go` - Verification logic
- `store.go` - Storage interface + implementations
- `base58.go` - Base58 encoding

### 3. Run Tests

```bash
go test -v ./...
```

All tests must pass before submitting.

### 4. Test Requirements

When adding features, include:
- Unit tests for new functions
- Table-driven tests for edge cases
- Tests for uniqueness (if applicable)
- Tests for error conditions

Example test structure:

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"normal case", "input", "expected", false},
        {"error case", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := YourFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("got = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 5. Submit Changes

```bash
git checkout -b feature/your-feature
git add .
git commit -m "Add: your feature description"
git push origin feature/your-feature
```

Then open a pull request.

## What Makes a Good Contribution

### Agentic Modification Guidelines

AI agents (and humans) should:

1. **Run existing tests first** - Don't break existing functionality
2. **Add tests** - New features need test coverage
3. **Update documentation** - Keep README and SKILL.md in sync
4. **Keep it simple** - Avoid over-engineering
5. **Follow Go conventions** - Use standard library where possible

### Code Style

- Use `go fmt` for formatting
- Use meaningful variable names
- Add GoDoc comments for exported symbols
- Keep dependencies minimal (standard library preferred)

### Security Considerations

- Never log or display API keys
- Use `crypto/rand` for randomness
- Use constant-time comparison for verification
- Hash keys before storage

## Types of Contributions

### Welcome:

- **New storage backends** - Implement `Store` interface for databases
- **Key formats** - Different prefix/size combinations
- **Hash algorithms** - Additional hashing options
- **Tests** - More edge case coverage
- **Documentation** - Improvements and translations
- **Examples** - Usage patterns for different languages

### Not Welcome:

- Introducing external dependencies without discussion
- Removing existing functionality
- Breaking changes without v2 discussion

## Getting Help

- Open an issue: https://github.com/package-foundry/gobekli_api/issues
- Discussions: https://github.com/package-foundry/gobekli_api/discussions

## License

By contributing, you agree that your contributions will be licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

**Made with :heart_emoji: from Peru**

El Perú es clave. :key_emoji: