package apikey

import (
	"fmt"
	"testing"
	"time"
)

func TestExamples(t *testing.T) {
	t.Run("Generator", func(t *testing.T) {
		gen, err := NewGenerator("sk_live", 32)
		if err != nil {
			t.Fatalf("failed to create generator: %v", err)
		}

		key, err := gen.Generate()
		if err != nil {
			t.Fatalf("failed to generate key: %v", err)
		}

		fmt.Printf("Generated key: %s\n", key)
		fmt.Printf("Prefix: %s\n", gen.Prefix())
	})

	t.Run("GenerateHashStore", func(t *testing.T) {
		gen, _ := NewGenerator("sk_live", 32)
		key, _ := gen.Generate()

		hash := HashKey(key)
		hashHex := HashKeyHex(key)
		fingerprint := Fingerprint(key)

		fmt.Printf("Key: %s\n", key)
		fmt.Printf("Hash (hex): %s\n", hashHex)
		fmt.Printf("Fingerprint: %s\n", fingerprint)

		record := &Record{
			ID:          "key-001",
			Prefix:      gen.Prefix(),
			Hash:        hash,
			Fingerprint: fingerprint,
			CreatedAt:   time.Now(),
		}

		store := NewMemStore()
		if err := store.Create(key, record); err != nil {
			t.Fatalf("failed to store: %v", err)
		}

		fmt.Printf("Stored record ID: %s\n", record.ID)
		fmt.Printf("Fingerprint: %s\n", record.Fingerprint)
	})

	t.Run("VerifyKey", func(t *testing.T) {
		gen, _ := NewGenerator("sk_live", 32)
		key, _ := gen.Generate()

		hash := HashKey(key)

		valid := VerifyKey(key, hash)
		fmt.Printf("Valid (same key): %v\n", valid)

		invalidKey := "sk_live_wrongkey123456789"
		valid = VerifyKey(invalidKey, hash)
		fmt.Printf("Valid (wrong key): %v\n", valid)
	})
}
