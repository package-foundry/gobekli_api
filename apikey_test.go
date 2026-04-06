package apikey

import (
	"testing"
	"time"
)

func TestGenerator_NewGenerator(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		keySize int
		wantErr bool
	}{
		{"valid", "sk_live", 32, false},
		{"valid_long_prefix", "sk_live_verylongprefix", 64, false},
		{"empty_prefix", "", 32, true},
		{"whitespace_prefix", "sk live", 32, true},
		{"too_small_size", "sk_live", 5, true},
		{"prefix_too_long", "sk_live", 8, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := NewGenerator(tt.prefix, tt.keySize)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && gen != nil {
				if gen.Prefix() != tt.prefix {
					t.Errorf("Prefix() = %v, want %v", gen.Prefix(), tt.prefix)
				}
				if gen.KeySize() != tt.keySize {
					t.Errorf("KeySize() = %v, want %v", gen.KeySize(), tt.keySize)
				}
			}
		})
	}
}

func TestGenerator_Generate(t *testing.T) {
	gen, err := NewGenerator("sk_live", 32)
	if err != nil {
		t.Fatalf("NewGenerator() error = %v", err)
	}

	key, err := gen.Generate()
	if err != nil {
		t.Errorf("Generate() error = %v", err)
		return
	}

	if len(key) != 32 {
		t.Errorf("key length = %v, want 32", len(key))
	}

	if len(key) < len(gen.Prefix())+1 {
		t.Errorf("key too short: %s", key)
	}

	if key[:len(gen.Prefix())] != gen.Prefix() {
		t.Errorf("key prefix = %v, want %v", key[:len(gen.Prefix())], gen.Prefix())
	}

	if key[len(gen.Prefix())] != '_' {
		t.Errorf("separator missing in key: %s", key)
	}
}

func TestGenerator_GenerateUniqueness(t *testing.T) {
	gen, _ := NewGenerator("sk_live", 32)
	keys := make(map[string]bool)

	for i := 0; i < 100; i++ {
		key, err := gen.Generate()
		if err != nil {
			t.Errorf("Generate() error = %v", err)
			return
		}
		if keys[key] {
			t.Errorf("duplicate key generated: %s", key)
		}
		keys[key] = true
	}
}

func TestValidateKey(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		expectedPrefix string
		expectedSize   int
		wantErr        bool
	}{
		{"valid", "live_AbCdEfGh12345", "live", 18, false},
		{"wrong_prefix", "test_AbCdEfGh12345", "live", 18, true},
		{"too_short", "live_Ab", "live", 18, true},
		{"wrong_size", "sk_live_AbCdEfGh123456789012345", "sk_live", 64, true},
		{"empty", "", "sk_live", 32, true},
		{"missing_separator", "skliveAbCdEfGh123456789012345678", "sk_live", 32, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateKey(tt.key, tt.expectedPrefix, tt.expectedSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtractPrefix(t *testing.T) {
	tests := []struct {
		key     string
		want    string
		wantErr bool
	}{
		{"prefix_randompayload", "prefix", false},
		{"testkey_123", "testkey", false},
		{"nokey", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			prefix, err := ExtractPrefix(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPrefix() error = %v, wantErr %v", err, tt.wantErr)
			}
			if prefix != tt.want {
				t.Errorf("ExtractPrefix() = %v, want %v", prefix, tt.want)
			}
		})
	}
}

func TestHashKey(t *testing.T) {
	key := "sk_live_AbCdEfGh123456789012345678"

	hash1 := HashKey(key)
	hash2 := HashKey(key)

	if string(hash1) != string(hash2) {
		t.Errorf("HashKey() not deterministic")
	}

	if len(hash1) != 32 {
		t.Errorf("hash length = %v, want 32", len(hash1))
	}
}

func TestHashKeyHex(t *testing.T) {
	key := "sk_live_AbCdEfGh123456789012345678"

	hash := HashKeyHex(key)

	if len(hash) != 64 {
		t.Errorf("hash hex length = %v, want 64", len(hash))
	}
}

func TestFingerprint(t *testing.T) {
	key := "sk_live_AbCdEfGh123456789012345678"

	fp := Fingerprint(key)

	if len(fp) != 16 {
		t.Errorf("fingerprint length = %v, want 16", len(fp))
	}
}

func TestVerifyKey(t *testing.T) {
	key := "sk_live_AbCdEfGh123456789012345678"
	wrongKey := "sk_live_WrOnGeYvAlUe123456789"
	hash := HashKey(key)

	if !VerifyKey(key, hash) {
		t.Errorf("VerifyKey() failed for valid key")
	}

	if VerifyKey(wrongKey, hash) {
		t.Errorf("VerifyKey() succeeded for invalid key")
	}
}

func TestVerifyKeyHex(t *testing.T) {
	key := "sk_live_AbCdEfGh123456789012345678"
	wrongKey := "sk_live_WrOnGeYvAlUe123456789"
	hashHex := HashKeyHex(key)

	if !VerifyKeyHex(key, hashHex) {
		t.Errorf("VerifyKeyHex() failed for valid key")
	}

	if VerifyKeyHex(wrongKey, hashHex) {
		t.Errorf("VerifyKeyHex() succeeded for invalid key")
	}
}

func TestRecord(t *testing.T) {
	now := time.Now()
	record := &Record{
		ID:          "test-001",
		Prefix:      "sk_live",
		Hash:        []byte("testhash"),
		Fingerprint: "abc123",
		CreatedAt:   now,
	}

	if !record.IsActive() {
		t.Errorf("IsActive() = false, want true for non-revoked, non-disabled record")
	}

	record.Disabled = true
	if record.IsActive() {
		t.Errorf("IsActive() = true, want false for disabled record")
	}

	record.Disabled = false
	record.RevokedAt = &now
	if record.IsActive() {
		t.Errorf("IsActive() = true, want false for revoked record")
	}
}

func TestMemStore(t *testing.T) {
	store := NewMemStore()
	gen, _ := NewGenerator("sk_live", 32)

	key, _ := gen.Generate()
	record := &Record{
		ID:     "key-001",
		Prefix: "sk_live",
	}

	if err := store.Create(key, record); err != nil {
		t.Errorf("Create() error = %v", err)
	}

	retrieved, err := store.GetByID("key-001")
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if retrieved.ID != record.ID {
		t.Errorf("GetByID() = %v, want %v", retrieved.ID, record.ID)
	}

	retrievedByHash, err := store.GetByHash(record.Hash)
	if err != nil {
		t.Errorf("GetByHash() error = %v", err)
	}
	if retrievedByHash.ID != record.ID {
		t.Errorf("GetByHash() = %v, want %v", retrievedByHash.ID, record.ID)
	}

	if err := store.Revoke("key-001"); err != nil {
		t.Errorf("Revoke() error = %v", err)
	}

	retrieved, _ = store.GetByID("key-001")
	if retrieved.IsActive() {
		t.Errorf("IsActive() = true, want false after revocation")
	}
}
