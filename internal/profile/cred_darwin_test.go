//go:build darwin

package profile

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestParseAccessToken(t *testing.T) {
	// 1. Plain JSON with nested token
	plainJSON := `{"token": {"access_token": "token-12345"}}`
	if got := parseAccessToken(plainJSON); got != "token-12345" {
		t.Errorf("parseAccessToken(plainJSON) = %q, want %q", got, "token-12345")
	}

	// 2. Base64 prefixed token (go-keyring-base64:)
	b64Payload := base64.StdEncoding.EncodeToString([]byte(plainJSON))
	keyringToken := "go-keyring-base64:" + b64Payload
	if got := parseAccessToken(keyringToken); got != "token-12345" {
		t.Errorf("parseAccessToken(keyringToken) = %q, want %q", got, "token-12345")
	}

	// 3. Flat token
	flatJSON := `{"access_token": "flat-token-67890"}`
	if got := parseAccessToken(flatJSON); got != "flat-token-67890" {
		t.Errorf("parseAccessToken(flatJSON) = %q, want %q", got, "flat-token-67890")
	}
}

func TestParseAccountEmail(t *testing.T) {
	// 1. Direct email field
	directJSON := `{"email": "user@example.com"}`
	if got := parseAccountEmail(directJSON); got != "user@example.com" {
		t.Errorf("parseAccountEmail(directJSON) = %q, want %q", got, "user@example.com")
	}

	// 2. Nested in token
	nestedJSON := `{"token": {"email": "nested@example.com"}}`
	if got := parseAccountEmail(nestedJSON); got != "nested@example.com" {
		t.Errorf("parseAccountEmail(nestedJSON) = %q, want %q", got, "nested@example.com")
	}

	// 3. Extracted from unpadded JWT id_token (RawURLEncoding)
	claims := `{"email": "jwt@example.com", "sub": "12345"}`
	claimsB64 := base64.RawURLEncoding.EncodeToString([]byte(claims))
	fakeJWT := fmt.Sprintf("header.%s.signature", claimsB64)
	jwtJSON := fmt.Sprintf(`{"token": {"access_token": "abc"}, "id_token": "%s"}`, fakeJWT)

	if got := parseAccountEmail(jwtJSON); got != "jwt@example.com" {
		t.Errorf("parseAccountEmail(jwtJSON) = %q, want %q", got, "jwt@example.com")
	}

	// 4. Wrapped in go-keyring-base64:
	b64Keyring := "go-keyring-base64:" + base64.StdEncoding.EncodeToString([]byte(jwtJSON))
	if got := parseAccountEmail(b64Keyring); got != "jwt@example.com" {
		t.Errorf("parseAccountEmail(b64Keyring) = %q, want %q", got, "jwt@example.com")
	}
}
