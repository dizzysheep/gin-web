package jwt

import (
	"errors"
	"gin-web/core/config"
	"gin-web/internal/model"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.GetString("jwt.secret"))

var (
	TokenExpired = errors.New("token已过期")
	TokenInvalid = errors.New("无效token")
)

type Claims struct {
	UserInfo *model.Auth `json:"user_info"`
	jwtlib.RegisteredClaims
}

func GenerateToken(userInfo *model.Auth) (string, error) {
	validTime := config.GetInt64("jwt.expireTime")
	claims := Claims{
		UserInfo: userInfo,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(time.Duration(validTime) * time.Second)),
			Issuer:    config.AppName,
		},
	}
	tokenClaims := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return tokenClaims.SignedString(jwtSecret)
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwtlib.ParseWithClaims(token, &Claims{}, func(t *jwtlib.Token) (interface{}, error) {
		// 限制签名算法，防止算法混淆攻击
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, TokenInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwtlib.ErrTokenExpired) {
			return nil, TokenExpired
		}
		return nil, TokenInvalid
	}

	claims, ok := tokenClaims.Claims.(*Claims)
	if !ok || !tokenClaims.Valid {
		return nil, TokenInvalid
	}
	return claims, nil
}
