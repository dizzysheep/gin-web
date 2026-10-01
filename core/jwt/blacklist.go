package jwt

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gin-web/core/redis"
)

const (
	blacklistKeyPrefix = "blog:jwt:blacklist:"
	verKeyPrefix       = "blog:user:ver:"
)

// Logout 将token加入黑名单（TTL为token剩余有效期）
func Logout(token string) error {
	claims, err := ParseToken(token)
	if err != nil {
		// 已过期/无效的token无需拉黑
		return nil
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}
	if redis.RedisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	return redis.RedisClient.Set(context.Background(), blacklistKeyPrefix+claims.ID, 1, ttl).Err()
}

// IsBlacklisted token是否已被拉黑（登出）。Redis异常时放行，保证可用性
func IsBlacklisted(jti string) bool {
	if jti == "" || redis.RedisClient == nil {
		return false
	}
	n, err := redis.RedisClient.Exists(context.Background(), blacklistKeyPrefix+jti).Result()
	if err != nil {
		return false
	}
	return n > 0
}

// CurrentVer 用户当前密码版本号，默认0
func CurrentVer(userID int64) int64 {
	if redis.RedisClient == nil {
		return 0
	}
	val, err := redis.RedisClient.Get(context.Background(), verKey(userID)).Result()
	if err != nil {
		return 0
	}
	ver, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0
	}
	return ver
}

// BumpVer 递增密码版本号，使该用户所有旧token失效
func BumpVer(userID int64) error {
	if redis.RedisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	ctx := context.Background()
	pipe := redis.RedisClient.TxPipeline()
	pipe.Incr(ctx, verKey(userID))
	pipe.Expire(ctx, verKey(userID), 7*24*time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

func verKey(userID int64) string {
	return fmt.Sprintf("%s%d", verKeyPrefix, userID)
}
