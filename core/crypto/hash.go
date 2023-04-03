package crypto

import (
    "crypto/hmac"
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
)

// 加盐 MD5 值
func Md5WithSalt(s, salt string) string {
    h := md5.Sum([]byte(s + salt))
    return hex.EncodeToString(h[:])
}

// MD5 值
func Md5(s string) string {
    return Md5WithSalt(s, "")
}

// sha1
func Sha1(s string) string {
    r := sha1.Sum([]byte(s))
    return hex.EncodeToString(r[:])
}

// sha256
func Sha256(s string) string {
    r := sha256.Sum256([]byte(s))
    return hex.EncodeToString(r[:])
}

// 老接口
func HmacSha1(input, key string) string {
    h := hmac.New(sha1.New, []byte(key))
    h.Write([]byte(input))
    s := base64.StdEncoding.EncodeToString(h.Sum(nil))
    return s
}
