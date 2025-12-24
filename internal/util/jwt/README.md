# JWT 工具包

这是一个通用的 JWT 工具包，使用 Go 泛型实现，可以支持任意类型的 Claims。

## 设计理念

参照 krathub 项目的设计模式：
- **JWT 工具包保持纯粹的泛型实现**，不依赖任何具体的业务模型
- **通过依赖注入在 service 层使用**，而不是在工具包中硬编码
- **支持多种 Claims 类型**，提高代码复用性

## 使用方式

### 1. 定义你的 Claims 结构

```go
import "github.com/golang-jwt/jwt/v5"

type MyClaims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims  // 必须嵌入这个结构体
}
```

### 2. 在 Service 中注入 JWT 实例

```go
import jwtpkg "github.com/HoronLee/EchoHub/internal/util/jwt"

type MyService struct {
    jwt *jwtpkg.JWT[MyClaims]
}

func NewMyService() *MyService {
    jwtService := jwtpkg.NewJWT[MyClaims](&jwtpkg.Config{
        SecretKey: "your-secret-key",
    })
    
    return &MyService{
        jwt: jwtService,
    }
}
```

### 3. 生成和解析 Token

```go
// 生成 Token
claims := &MyClaims{
    UserID:   1,
    Username: "testuser",
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    },
}

token, err := s.jwt.GenerateToken(claims)

// 解析 Token
parsedClaims, err := s.jwt.ParseToken(token)
```

## 示例

参考 `internal/service/user.go` 中的实现，了解如何在实际业务中使用这个工具包。

## 优势

1. **类型安全**：使用泛型确保编译时类型检查
2. **灵活性**：可以为不同的业务场景定义不同的 Claims
3. **解耦**：工具包不依赖任何业务模型
4. **可测试**：易于编写单元测试
