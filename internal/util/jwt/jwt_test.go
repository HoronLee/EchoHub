package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TestClaims 测试用的 Claims 结构
type TestClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func TestGenerateToken(t *testing.T) {
	jwtService := NewJWT[TestClaims](&Config{
		SecretKey: "test-secret-key",
	})

	claims := &TestClaims{
		UserID:   1,
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwtService.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("Generated token should not be empty")
	}

	// Token should have 3 parts separated by dots
	parts := 0
	for _, c := range token {
		if c == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Errorf("JWT token should have 3 parts (2 dots), got %d dots", parts)
	}
}

func TestParseToken(t *testing.T) {
	jwtService := NewJWT[TestClaims](&Config{
		SecretKey: "test-secret-key",
	})

	// Generate a token
	originalClaims := &TestClaims{
		UserID:   1,
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwtService.GenerateToken(originalClaims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Parse the token
	parsedClaims, err := jwtService.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if parsedClaims.UserID != originalClaims.UserID {
		t.Errorf("Expected UserID %d, got %d", originalClaims.UserID, parsedClaims.UserID)
	}

	if parsedClaims.Username != originalClaims.Username {
		t.Errorf("Expected Username %s, got %s", originalClaims.Username, parsedClaims.Username)
	}
}

func TestParseTokenInvalid(t *testing.T) {
	jwtService := NewJWT[TestClaims](&Config{
		SecretKey: "test-secret-key",
	})

	// Test with invalid token
	_, err := jwtService.ParseToken("invalid.token.string")
	if err == nil {
		t.Error("ParseToken should fail with invalid token")
	}

	// Test with empty token
	_, err = jwtService.ParseToken("")
	if err == nil {
		t.Error("ParseToken should fail with empty token")
	}
}

func TestParseTokenExpired(t *testing.T) {
	jwtService := NewJWT[TestClaims](&Config{
		SecretKey: "test-secret-key",
	})

	// Generate a token that expires immediately
	claims := &TestClaims{
		UserID:   1,
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Already expired
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwtService.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Try to parse the expired token
	_, err = jwtService.ParseToken(token)
	if err == nil {
		t.Error("ParseToken should fail with expired token")
	}
}

func TestRoundTrip(t *testing.T) {
	jwtService := NewJWT[TestClaims](&Config{
		SecretKey: "test-secret-key",
	})

	// Create claims
	originalClaims := &TestClaims{
		UserID:   42,
		Username: "roundtripuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate token
	token, err := jwtService.GenerateToken(originalClaims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Parse token
	parsedClaims, err := jwtService.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	// Verify round-trip consistency
	if parsedClaims.UserID != originalClaims.UserID {
		t.Errorf("Round-trip failed: UserID mismatch. Expected %d, got %d",
			originalClaims.UserID, parsedClaims.UserID)
	}

	if parsedClaims.Username != originalClaims.Username {
		t.Errorf("Round-trip failed: Username mismatch. Expected %s, got %s",
			originalClaims.Username, parsedClaims.Username)
	}
}

func TestDifferentClaimsTypes(t *testing.T) {
	// Test that we can use different claims types
	type CustomClaims struct {
		Role string `json:"role"`
		jwt.RegisteredClaims
	}

	jwtService := NewJWT[CustomClaims](&Config{
		SecretKey: "test-secret-key",
	})

	claims := &CustomClaims{
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token, err := jwtService.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	parsedClaims, err := jwtService.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if parsedClaims.Role != claims.Role {
		t.Errorf("Expected Role %s, got %s", claims.Role, parsedClaims.Role)
	}
}
