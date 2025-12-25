package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/HoronLee/EchoHub/internal/model/user"
	jwtUtil "github.com/HoronLee/EchoHub/internal/util/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware(t *testing.T) {
	// 设置测试环境
	config.JWT_SECRET = []byte("test-secret-key")

	// 创建 JWT 服务
	jwtService := jwtUtil.NewJWT[user.Claims](&jwtUtil.Config{
		SecretKey: string(config.JWT_SECRET),
	})

	// 生成有效的测试 Token
	claims := &user.Claims{
		UserID:   1,
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
			Issuer:    "echohub",
			Subject:   "testuser",
			Audience:  []string{"echohub-api"},
		},
	}
	validToken, err := jwtService.GenerateToken(claims)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedMsg:    "",
		},
		{
			name:           "Missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Token not found",
		},
		{
			name:           "Invalid token format - no Bearer prefix",
			authHeader:     validToken,
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Token format invalid",
		},
		{
			name:           "Invalid token format - wrong prefix",
			authHeader:     "Basic " + validToken,
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Token format invalid",
		},
		{
			name:           "Empty token after Bearer",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Token not found",
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Token invalid or expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			e := echo.New()
			e.Use(JwtAuth())
			e.GET("/protected", func(c echo.Context) error {
				userID := c.Get("user_id")
				assert.NotNil(t, userID)
				assert.Equal(t, uint(1), userID)

				username := c.Get("username")
				assert.NotNil(t, username)
				assert.Equal(t, "testuser", username)

				return c.JSON(http.StatusOK, map[string]string{"message": "success"})
			})

			// 创建测试请求
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// 执行请求
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

// Feature: user-auth, Property 8: Authentication middleware token validation
// Validates: Requirements 3.1, 4.1, 4.5
func TestProperty_MiddlewareTokenValidation(t *testing.T) {
	// 设置测试环境
	config.JWT_SECRET = []byte("test-secret-key-for-property-testing")

	// 创建 JWT 服务
	jwtService := jwtUtil.NewJWT[user.Claims](&jwtUtil.Config{
		SecretKey: string(config.JWT_SECRET),
	})

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 1: Valid tokens should be accepted and set user context
	properties.Property("Valid tokens are accepted and set user ID in context", prop.ForAll(
		func(userID uint, username string) bool {
			// Skip empty usernames
			if username == "" {
				return true
			}

			// Create valid claims
			claims := &user.Claims{
				UserID:   userID,
				Username: username,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
					Issuer:    "echohub",
					Subject:   username,
					Audience:  []string{"echohub-api"},
				},
			}

			// Generate valid token
			token, err := jwtService.GenerateToken(claims)
			if err != nil {
				return false
			}

			// Create test router with middleware
			e := echo.New()
			e.Use(JwtAuth())

			contextUserID := uint(0)
			contextUsername := ""
			e.GET("/protected", func(c echo.Context) error {
				// Extract user ID from context
				if id := c.Get("user_id"); id != nil {
					contextUserID = id.(uint)
				}
				if name := c.Get("username"); name != nil {
					contextUsername = name.(string)
				}
				return c.JSON(http.StatusOK, map[string]string{"message": "success"})
			})

			// Create request with valid Bearer token
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			// Execute request
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			// Verify: Should return 200 and set correct user ID in context
			return rec.Code == http.StatusOK && contextUserID == userID && contextUsername == username
		},
		gen.UIntRange(1, 1000000),
		gen.AlphaString(),
	))

	// Property 2: Requests without Bearer tokens should be rejected
	properties.Property("Requests without valid Bearer tokens are rejected", prop.ForAll(
		func(invalidHeader string) bool {
			// Create test router with middleware
			e := echo.New()
			e.Use(JwtAuth())
			e.GET("/protected", func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]string{"message": "success"})
			})

			// Create request with invalid/missing authorization header
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if invalidHeader != "" {
				req.Header.Set("Authorization", invalidHeader)
			}

			// Execute request
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			// Verify: Should return 401 Unauthorized
			return rec.Code == http.StatusUnauthorized
		},
		gen.OneGenOf(gen.Const(""), gen.AlphaString(), gen.Const("Basic token"), gen.Const("token")),
	))

	// Property 3: Invalid tokens should be rejected
	properties.Property("Invalid tokens are rejected with 401", prop.ForAll(
		func(randomString string) bool {
			// Skip empty strings and very short strings
			if randomString == "" || len(randomString) < 10 {
				return true
			}

			// Create test router with middleware
			e := echo.New()
			e.Use(JwtAuth())
			e.GET("/protected", func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]string{"message": "success"})
			})

			// Create request with invalid token
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+randomString)

			// Execute request
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			// Verify: Should return 401 Unauthorized
			return rec.Code == http.StatusUnauthorized
		},
		gen.Identifier(),
	))

	// Run all properties
	properties.TestingRun(t)
}
