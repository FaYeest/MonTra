package auth

import (
	"testing"
	"time"
)

func TestTokenManager(t *testing.T) {
	secret := "test-secret-at-least-32-chars-long-12345"
	tm, err := NewTokenManager(secret, 1*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error creating token manager: %v", err)
	}

	userID := "usr-uuid-1234"
	email := "user@example.com"

	token, err := tm.Generate(userID, email)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	claims, err := tm.Verify(token)
	if err != nil {
		t.Fatalf("unexpected error verifying token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %q, got %q", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected email %q, got %q", email, claims.Email)
	}

	// Test wrong secret verification
	otherTM, _ := NewTokenManager("different-secret-for-testing-purposes-9876", 1*time.Hour)
	_, err = otherTM.Verify(token)
	if err == nil {
		t.Errorf("expected error verifying with wrong secret, got nil")
	}

	// Test expired token
	expiredTM, _ := NewTokenManager(secret, -1*time.Hour)
	expiredToken, _ := expiredTM.Generate(userID, email)
	_, err = tm.Verify(expiredToken)
	if err == nil {
		t.Errorf("expected error verifying expired token, got nil")
	}
}
