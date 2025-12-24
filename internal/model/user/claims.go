package user

import "github.com/golang-jwt/jwt/v5"

// Claims JWT Claims 结构体
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
