package auth

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mukailasam/igo"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

type Claims struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

type contextKey string

func contextWithUID(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, contextKey("uid"), uid)
}

func Auth(c *igo.Context, next igo.HandlerFunc) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Abort(403, "missing token")
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := ParseToken(tokenStr)
	if err != nil {
		c.Abort(403, "invalid token")
		return
	}
	ctx := c.Request.Context()
	ctx = contextWithUID(ctx, claims.UID)
	c.Request = c.Request.WithContext(ctx)
	next(c)
}

func GenerateToken(uid string) (string, error) {
	claims := &Claims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
