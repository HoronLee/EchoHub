package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/HoronLee/EchoHub/internal/model/user"
	jwtUtil "github.com/HoronLee/EchoHub/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)
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
			router := gin.New()
			router.Use(JWTAuthMiddleware())
			router.GET("/protected", func(c *gin.Context) {
				userID, exists := c.Get("user_id")
				assert.True(t, exists)
				assert.Equal(t, uint(1), userID)

				username, exists := c.Get("username")
				assert.True(t, exists)
				assert.Equal(t, "testuser", username)

				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// 创建测试请求
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// 执行请求
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// Feature: user-auth, Property 8: Authentication middleware token validation
// Validates: Requirements 3.1, 4.1, 4.5
func TestProperty_MiddlewareTokenValidation(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)
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
			router := gin.New()
			router.Use(JWTAuthMiddleware())

			contextUserID := uint(0)
			contextUsername := ""
			router.GET("/protected", func(c *gin.Context) {
				// Extract user ID from context
				if id, exists := c.Get("user_id"); exists {
					contextUserID = id.(uint)
				}
				if name, exists := c.Get("username"); exists {
					contextUsername = name.(string)
				}
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request with valid Bearer token
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify: Should return 200 and set correct user ID in context
			return w.Code == http.StatusOK && contextUserID == userID && contextUsername == username
		},
		gen.UIntRange(1, 1000000),
		gen.AlphaString(),
	))

	// Property 2: Requests without Bearer tokens should be rejected
	properties.Property("Requests without valid Bearer tokens are rejected", prop.ForAll(
		func(invalidHeader string) bool {
			// Create test router with middleware
			router := gin.New()
			router.Use(JWTAuthMiddleware())
			router.GET("/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request with invalid/missing authorization header
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if invalidHeader != "" {
				req.Header.Set("Authorization", invalidHeader)
			}

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify: Should return 401 Unauthorized
			return w.Code == http.StatusUnauthorized
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
			router := gin.New()
			router.Use(JWTAuthMiddleware())
			router.GET("/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request with invalid token
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+randomString)

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify: Should return 401 Unauthorized
			return w.Code == http.StatusUnauthorized
		},
		gen.Identifier(),
	))

	// Run all properties
	properties.TestingRun(t)
}
