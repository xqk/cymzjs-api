package pkg

import (
	"context"
	"git.zc0901.com/go/god/lib/gconv"
	"git.zc0901.com/go/god/lib/stringx"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func JwtFromUid(secret string, expire, uid int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Unix() + expire
	claims["userId"] = uid
	tokenizer := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenizer.SignedString([]byte(secret))
}

func UidFromJwt(ctx context.Context) int64 {
	return gconv.Int64(ctx.Value("userId"))
}

func JwtDecode(tokenString, secret string) (jwt.MapClaims, error) {
	tokenString = stringx.ReplaceByMap(tokenString, map[string]string{
		"Bearer ": "",
		"bearer ": "",
	})
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
