package jwt

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	"gin-web/internal/model"
)

func TestGenerateAndParseToken(t *testing.T) {
	user := &model.Auth{ID: 1, Username: "admin"}

	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken err: %v", err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken err: %v", err)
	}
	if claims.UserInfo == nil || claims.UserInfo.Username != "admin" {
		t.Fatalf("unexpected user info: %+v", claims.UserInfo)
	}
}

func TestParseTokenInvalid(t *testing.T) {
	if _, err := ParseToken("not-a-jwt"); !errors.Is(err, TokenInvalid) {
		t.Fatalf("expect TokenInvalid, got %v", err)
	}
}

// 签名被篡改的token应被拒绝
func TestParseTokenBadSignature(t *testing.T) {
	token, err := GenerateToken(&model.Auth{ID: 1, Username: "admin"})
	if err != nil {
		t.Fatalf("GenerateToken err: %v", err)
	}

	parts := strings.Split(token, ".")
	parts[2] = base64.RawURLEncoding.EncodeToString([]byte("bad-signature"))

	if _, err := ParseToken(strings.Join(parts, ".")); !errors.Is(err, TokenInvalid) {
		t.Fatalf("expect TokenInvalid, got %v", err)
	}
}

// 非HMAC算法签名的token应被拒绝，防止算法混淆攻击
func TestParseTokenRejectsNonHMACAlg(t *testing.T) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"user_info":null,"iss":"test","exp":4102444800}`))
	forged := header + "." + payload + "."

	if _, err := ParseToken(forged); !errors.Is(err, TokenInvalid) {
		t.Fatalf("expect TokenInvalid, got %v", err)
	}
}
