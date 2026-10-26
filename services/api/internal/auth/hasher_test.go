package auth

import (
	"strings"
	"testing"
)

func TestHasher(t *testing.T) {
	h := NewHasher()

	// Short password
	_, err := h.Hash("short")
	if err == nil {
		t.Errorf("expected error for short password, got nil")
	}

	// Long password (> 72 bytes)
	_, err = h.Hash(strings.Repeat("a", 73))
	if err == nil {
		t.Errorf("expected error for oversized password, got nil")
	}

	// Valid password
	password := "rahasia123"
	hash, err := h.Hash(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if len(hash) != 60 {
		t.Errorf("expected bcrypt hash length 60, got %d", len(hash))
	}

	if !h.Compare(hash, password) {
		t.Errorf("expected hash comparison to succeed for correct password")
	}

	if h.Compare(hash, "wrong-password") {
		t.Errorf("expected hash comparison to fail for wrong password")
	}
}
