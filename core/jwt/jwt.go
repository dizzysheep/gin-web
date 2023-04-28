package jwt

import (
	"errors"
	"gin-web/core/config"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte(config.GetString("jwt.secret"))

var (
	TokenExpired = errors.New("token已过期")
	TokenInvalid = errors.New("无效token")
)

type Claims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(userId int) (string, error) {
	validTime := config.GetInt64("jwt.expireTime")
	expireTime := time.Now().Unix() + validTime
	claims := Claims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    config.AppName,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if tokenClaims != nil {
		if Claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return Claims, nil
		}
		return nil, errors.New("无效token")
	}
	return nil, err
}
