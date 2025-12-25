package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestJwtAuth_NotFoundRoute 测试访问不存在的路由时返回404而非401
func TestJwtAuth_NotFoundRoute(t *testing.T) {
	e := echo.New()

	// 模拟路由组结构
	v1 := e.Group("/api/v1")
	private := v1.Group("")
	private.Use(JwtAuth())

	// 注册一个私有路由
	private.GET("/protected", func(c echo.Context) error {
		return c.String(http.StatusOK, "protected")
	})

	tests := []struct {
		name           string
		path           string
		method         string
		token          string
		expectedStatus int
	}{
		{
			name:           "不存在的路由应返回404",
			path:           "/api/v1/notexist",
			method:         http.MethodGet,
			token:          "",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "存在的路由无token应返回401",
			path:           "/api/v1/protected",
			method:         http.MethodGet,
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "不在v1下的路由应返回404",
			path:           "/api/notexist",
			method:         http.MethodGet,
			token:          "",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

// TestJwtAuth_WildcardRoute 测试真正的通配符路由仍需认证
func TestJwtAuth_WildcardRoute(t *testing.T) {
	e := echo.New()

	v1 := e.Group("/api/v1")
	private := v1.Group("")
	private.Use(JwtAuth())

	// 注册一个真正的通配符路由
	private.GET("/files/*", func(c echo.Context) error {
		return c.String(http.StatusOK, "files: "+c.Param("*"))
	})

	tests := []struct {
		name           string
		path           string
		token          string
		expectedStatus int
	}{
		{
			name:           "通配符路由无token应返回401",
			path:           "/api/v1/files/test.txt",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "通配符路由根路径无token应返回401",
			path:           "/api/v1/files/",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
