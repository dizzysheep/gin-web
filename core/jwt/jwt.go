package jwt

import (
	"errors"
	"gin-web/core/config"
	"gin-web/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte(config.GetString("jwt.secret"))

var (
	TokenExpired = errors.New("token已过期")
	TokenInvalid = errors.New("无效token")
)

type Claims struct {
	UserInfo *model.User `json:"user_info"`
	// Ver 密码版本号，修改密码后递增使旧token失效
	Ver int64 `json:"ver"`
	jwt.RegisteredClaims
}

func GenerateToken(userInfo *model.User, ver int64) (string, error) {
	validTime := config.GetInt64("jwt.expireTime")
	claims := Claims{
		UserInfo: userInfo,
		Ver:      ver,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(validTime) * time.Second)),
			Issuer:    config.AppName,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(jwtSecret)
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// 限制签名算法，防止算法混淆攻击
		if t.Method != jwt.SigningMethodHS256 {
			return nil, TokenInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, TokenExpired
		}
		return nil, TokenInvalid
	}

	claims, ok := tokenClaims.Claims.(*Claims)
	if !ok || !tokenClaims.Valid || claims.UserInfo == nil || claims.UserInfo.ID <= 0 || claims.ID == "" {
		return nil, TokenInvalid
	}
	return claims, nil
}
