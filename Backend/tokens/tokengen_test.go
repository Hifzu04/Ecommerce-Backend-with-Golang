package token

import (
	"os"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// setTestSecret sets a known secret key for deterministic tests.
func setTestSecret(t *testing.T) {
	t.Helper()
	os.Setenv("SECRET_LOVE", "test-secret-key-for-unit-tests")
	SECRET_KEY = os.Getenv("SECRET_LOVE")
}

// ===================================================================
// UNIT TESTS — Token generation
// ===================================================================

func TestTokenGenerator_Success(t *testing.T) {
	setTestSecret(t)

	token, refreshToken, err := TokenGenerator("test@example.com", "John", "Doe", "user123")

	if err != nil {
		t.Fatalf("TokenGenerator returned error: %v", err)
	}
	if token == "" {
		t.Fatal("TokenGenerator returned empty access token")
	}
	if refreshToken == "" {
		t.Fatal("TokenGenerator returned empty refresh token")
	}
}

func TestTokenGenerator_TokensAreDifferent(t *testing.T) {
	setTestSecret(t)

	token, refreshToken, err := TokenGenerator("a@b.com", "A", "B", "uid1")
	if err != nil {
		t.Fatalf("TokenGenerator error: %v", err)
	}

	if token == refreshToken {
		t.Fatal("Access token and refresh token should be different")
	}
}

func TestTokenGenerator_DifferentUsers_DifferentTokens(t *testing.T) {
	setTestSecret(t)

	token1, _, err1 := TokenGenerator("user1@test.com", "User", "One", "uid1")
	token2, _, err2 := TokenGenerator("user2@test.com", "User", "Two", "uid2")

	if err1 != nil || err2 != nil {
		t.Fatalf("TokenGenerator errors: %v, %v", err1, err2)
	}
	if token1 == token2 {
		t.Fatal("Different users should get different tokens")
	}
}

// ===================================================================
// UNIT TESTS — Token validation
// ===================================================================

func TestValidateToken_Success(t *testing.T) {
	setTestSecret(t)

	token, _, err := TokenGenerator("validate@test.com", "Val", "Test", "uid-val")
	if err != nil {
		t.Fatalf("TokenGenerator error: %v", err)
	}

	claims, msg := ValidateToken(token)

	if msg != "" {
		t.Fatalf("ValidateToken returned error message: %s", msg)
	}
	if claims == nil {
		t.Fatal("ValidateToken returned nil claims")
	}
	if claims.Email != "validate@test.com" {
		t.Fatalf("Expected email 'validate@test.com', got '%s'", claims.Email)
	}
	if claims.First_Name != "Val" {
		t.Fatalf("Expected first name 'Val', got '%s'", claims.First_Name)
	}
	if claims.Uid != "uid-val" {
		t.Fatalf("Expected uid 'uid-val', got '%s'", claims.Uid)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	setTestSecret(t)

	_, msg := ValidateToken("this.is.not.a.valid.token")

	if msg == "" {
		t.Fatal("ValidateToken should return error for invalid token")
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	setTestSecret(t)

	_, msg := ValidateToken("")

	if msg == "" {
		t.Fatal("ValidateToken should return error for empty token")
	}
}

func TestValidateToken_TamperedToken(t *testing.T) {
	setTestSecret(t)

	token, _, _ := TokenGenerator("tamper@test.com", "Tam", "Per", "uid-t")

	// Tamper with the token by modifying a character
	tampered := token[:len(token)-1] + "X"
	_, msg := ValidateToken(tampered)

	if msg == "" {
		t.Fatal("ValidateToken should reject tampered token")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	setTestSecret(t)

	token, _, _ := TokenGenerator("secret@test.com", "Sec", "Ret", "uid-s")

	// Change the secret key before validating
	SECRET_KEY = "completely-different-secret"
	_, msg := ValidateToken(token)

	if msg == "" {
		t.Fatal("ValidateToken should reject token signed with different secret")
	}

	// Restore for other tests
	setTestSecret(t)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	setTestSecret(t)

	// Manually create an expired token
	claims := &SignedDetails{
		Email:      "expired@test.com",
		First_Name: "Exp",
		Last_Name:  "Ired",
		Uid:        "uid-exp",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(), // expired 1 hour ago
		},
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenObj.SignedString([]byte(SECRET_KEY))
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	_, msg := ValidateToken(tokenString)

	if msg == "" {
		t.Fatal("ValidateToken should reject expired token")
	}
}

func TestValidateToken_ClaimsContainCorrectData(t *testing.T) {
	setTestSecret(t)

	token, _, _ := TokenGenerator("claims@test.com", "First", "Last", "uid-claims")
	claims, msg := ValidateToken(token)

	if msg != "" {
		t.Fatalf("Unexpected error: %s", msg)
	}
	if claims.Email != "claims@test.com" {
		t.Errorf("Email mismatch: got %s", claims.Email)
	}
	if claims.First_Name != "First" {
		t.Errorf("First_Name mismatch: got %s", claims.First_Name)
	}
	if claims.Last_Name != "Last" {
		t.Errorf("Last_Name mismatch: got %s", claims.Last_Name)
	}
	if claims.Uid != "uid-claims" {
		t.Errorf("Uid mismatch: got %s", claims.Uid)
	}
}
