package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	if token == "" {
		t.Fatal("Token is empty")
	}

	// Validate the token
	parsedID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if parsedID != userID {
		t.Fatalf("Parsed ID %v does not match %v", parsedID, userID)
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	parsedID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if parsedID != userID {
		t.Fatalf("Parsed ID %v does not match %v", parsedID, userID)
	}
}

func TestValidateJWT_InvalidSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Fatal("Expected error for invalid secret")
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := -time.Hour // expired

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatal("Expected error for expired token")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	_, err := ValidateJWT("invalid.token", "secret")
	if err == nil {
		t.Fatal("Expected error for invalid token")
	}
}
