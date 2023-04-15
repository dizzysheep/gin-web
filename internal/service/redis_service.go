package service

import (
	"github.com/yangxx0612/plugins/config"
	"strconv"
)

type RedisService struct {
	Delimiter  string
	ModuleName string
}

func NewRedisService(moduleName string) *RedisService {
	return &RedisService{
		Delimiter:  ".",
		ModuleName: moduleName,
	}
}

func (s *RedisService) GetPrefix() string {
	return config.AppName + s.Delimiter + s.ModuleName
}

func (s *RedisService) GetDetailKey(id int) string {
	cacheKey := s.GetPrefix() + s.Delimiter + "id" + s.Delimiter + strconv.Itoa(id)
	return cacheKey
}
