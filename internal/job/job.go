// Package job 后台定时任务
package job

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gin-web/core/config"
	"gin-web/core/log"
	"gin-web/core/redis"
	"gin-web/internal/dao/article"
	"gin-web/pkg/boostrap"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

const viewCountKeyPattern = "blog:view:count:*"

// RunViewCountFlusher 周期性将Redis中的浏览量缓冲刷盘到MySQL
func RunViewCountFlusher(ctx context.Context) {
	interval := config.GetInt("job.viewFlushSec")
	if interval <= 0 {
		interval = 600
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			flushCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			flushViewCount(flushCtx)
			cancel()
			return
		case <-ticker.C:
			flushViewCount(ctx)
		}
	}
}

// flushViewCount 扫描所有日期的浏览量key，逐文章累加落库后删除key
func flushViewCount(ctx context.Context) {
	client := redis.RedisClient
	if client == nil {
		return
	}

	iterator := client.Scan(ctx, 0, viewCountKeyPattern, 100).Iterator()
	for iterator.Next(ctx) {
		key := iterator.Val()
		processingKey := key
		if !strings.Contains(key, ":processing:") {
			processingKey = fmt.Sprintf("%s:processing:%s", key, uuid.NewString())
			if err := client.Rename(ctx, key, processingKey).Err(); err != nil {
				if err != goredis.Nil {
					log.Get(nil).Errorf("view flush rename %s err: %s", key, err.Error())
				}
				continue
			}
		}

		entries, err := client.HGetAll(ctx, processingKey).Result()
		if err != nil {
			log.Get(nil).Errorf("view flush hgetall %s err: %s", processingKey, err.Error())
			continue
		}

		dao := article.NewArticleDao(boostrap.InitDBEngine())
		for field, deltaStr := range entries {
			id, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				continue
			}
			delta, err := strconv.ParseInt(deltaStr, 10, 64)
			if err != nil || delta <= 0 {
				continue
			}
			if err := dao.IncrViewCount(ctx, id, delta); err != nil {
				log.Get(nil).Errorf("view flush incr article %d err: %s", id, err.Error())
				continue
			}
			if err := client.HDel(ctx, processingKey, field).Err(); err != nil {
				log.Get(nil).Errorf("view flush ack article %d err: %s", id, err.Error())
			}
		}

		if n, err := client.HLen(ctx, processingKey).Result(); err == nil && n == 0 {
			if err := client.Del(ctx, processingKey).Err(); err != nil {
				log.Get(nil).Errorf("view flush del %s err: %s", processingKey, err.Error())
			}
		}
	}
	if err := iterator.Err(); err != nil {
		log.Get(nil).Errorf("view flush scan err: %s", err.Error())
	}
}
