package redis

import (
	"gin-web/core/config"
	"strconv"
)

type RedisKey struct {
	Delimiter  string
	ModuleName string
}

func NewRedisService(moduleName string) *RedisKey {
	return &RedisKey{
		Delimiter:  ".",
		ModuleName: moduleName,
	}
}

func (s *RedisKey) GetPrefix() string {
	return config.AppName + s.Delimiter + s.ModuleName
}

func (s *RedisKey) GetDetailKey(id int) string {
	cacheKey := s.GetPrefix() + s.Delimiter + "id" + s.Delimiter + strconv.Itoa(id)
	return cacheKey
}
