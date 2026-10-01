package option

import (
	"context"
	"encoding/json"
	"time"

	"gin-web/core/redis"
	"gin-web/dto"
	"gin-web/internal/dao"
	"github.com/pkg/errors"
)

const (
	optionCacheKey = "blog:option"
	optionCacheTTL = 10 * time.Minute
)

type optionService struct {
	daos *dao.Daos
}

func NewOptionService(daos *dao.Daos) OptionService {
	return &optionService{daos}
}

// Get 站点配置（Redis缓存10分钟，写时失效）
func (s *optionService) Get(ctx context.Context) (*dto.GetOptionsRespDTO, error) {
	options, err := s.loadOptions(ctx)
	if err != nil {
		return nil, err
	}
	return &dto.GetOptionsRespDTO{Options: options}, nil
}

// Value 读取单个配置值（供其他服务使用）
func (s *optionService) Value(ctx context.Context, key string) string {
	options, err := s.loadOptions(ctx)
	if err != nil {
		return ""
	}
	return options[key]
}

// Save 批量保存站点配置并清除缓存
func (s *optionService) Save(ctx context.Context, reqDTO *dto.SaveOptionsReqDTO) error {
	for _, opt := range reqDTO.Options {
		if err := s.daos.Option.Upsert(ctx, opt); err != nil {
			return errors.Wrap(err, "upsert option fail")
		}
	}

	if redis.RedisClient != nil {
		_ = redis.RedisClient.Del(ctx, optionCacheKey).Err()
	}
	return nil
}

func (s *optionService) loadOptions(ctx context.Context) (map[string]string, error) {
	var options map[string]string
	if err := redis.GetStruct(optionCacheKey, &options); err == nil && options != nil {
		return options, nil
	}

	pos, err := s.daos.Option.SelectAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "select options fail")
	}

	options = make(map[string]string, len(pos))
	for _, po := range pos {
		options[po.Key] = normalizeValue(po.Value)
	}

	// 缓存失败不影响业务
	_ = redis.SaveStruct(optionCacheKey, options, optionCacheTTL)
	return options, nil
}

func normalizeValue(value string) string {
	var decoded string
	if err := json.Unmarshal([]byte(value), &decoded); err == nil {
		return decoded
	}
	return value
}
