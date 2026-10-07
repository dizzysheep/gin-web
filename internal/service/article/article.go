package article

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gin-web/core/markdown"
	"gin-web/core/redis"
	"gin-web/core/xtime"
	"gin-web/dto"
	"gin-web/internal/dao"
	"gin-web/internal/dao/common"
	"gin-web/internal/errcode"
	"gin-web/internal/model"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

const (
	// viewCountKeyPrefix 浏览量Redis缓冲key，按天分hash
	viewCountKeyPrefix = "blog:view:count:"
	archiveCacheKey    = "blog:archive"
	archiveCacheTTL    = time.Hour
	articleCachePrefix = "blog:article:"
	slugCachePrefix    = "blog:article:slug:"
	listCachePrefix    = "blog:article:list:"
	listVersionKey     = "blog:article:list:version"
	detailVersionKey   = "blog:article:detail:version"
	articleCacheTTL    = 30 * time.Minute
	listCacheTTL       = 5 * time.Minute
)

type articleService struct {
	daos *dao.Daos
}

func NewArticleService(daos *dao.Daos) ArticleService {
	return &articleService{daos}
}

// List 公开文章列表（仅已发布且启用）
func (s *articleService) List(ctx context.Context, reqDTO *dto.ListArticleReqDTO) (*dto.ListArticleRespDTO, error) {
	cacheKey := s.listCacheKey(ctx, reqDTO)
	var cached dto.ListArticleRespDTO
	if err := redis.GetStruct(cacheKey, &cached); err == nil && cached.Pager != nil && cached.Items != nil {
		return &cached, nil
	}

	conditions := common.GormConditions{
		&common.EqCond{Field: "state", Value: model.StateEnabled},
		&common.EqCond{Field: "is_draft", Value: 0},
	}
	if reqDTO.ID > 0 {
		// 精确 ID 查询时其他筛选条件无意义，直接跳过关键词模糊匹配
		conditions = append(conditions, &common.EqCond{Field: "id", Value: reqDTO.ID})
	} else {
		conditions = s.appendFilters(conditions, reqDTO.CategoryID, reqDTO.Keyword)
	}

	result, err := s.list(ctx, conditions, reqDTO.TagID, reqDTO.Pager)
	if err != nil {
		return nil, err
	}
	_, _ = redis.SaveStructNX(cacheKey, result, listCacheTTL)
	return result, nil
}

// AdminList 管理端文章列表（含草稿）
func (s *articleService) AdminList(ctx context.Context, reqDTO *dto.AdminListArticleReqDTO) (*dto.ListArticleRespDTO, error) {
	conditions := common.GormConditions{
		&common.EqCond{Field: "state", Value: reqDTO.State},
		&common.EqCond{Field: "is_draft", Value: reqDTO.IsDraft},
		&common.GteCond{Field: "published_on", Value: reqDTO.PublishedFrom},
		&common.LtCond{Field: "published_on", Value: reqDTO.PublishedTo},
	}
	conditions = s.appendFilters(conditions, reqDTO.CategoryID, reqDTO.Keyword)

	return s.list(ctx, conditions, reqDTO.TagID, reqDTO.Pager)
}

// AllPublished returns all published articles for feeds and sitemaps.
func (s *articleService) AllPublished(ctx context.Context) (*dto.ListArticleRespDTO, error) {
	conditions := common.GormConditions{
		&common.EqCond{Field: "state", Value: model.StateEnabled},
		&common.EqCond{Field: "is_draft", Value: 0},
	}
	return s.list(ctx, conditions, 0, nil)
}

func (s *articleService) appendFilters(conditions common.GormConditions, categoryID int64, keyword string) common.GormConditions {
	if categoryID > 0 {
		conditions = append(conditions, &common.EqCond{Field: "category_id", Value: categoryID})
	}
	return append(conditions, &common.OrConditions{GormCond: []common.GormCond{
		&common.LikeCond{Field: "title", Value: keyword},
		&common.LikeCond{Field: "desc", Value: keyword},
		&common.LikeCond{Field: "content_md", Value: keyword},
	}})
}

func (s *articleService) list(ctx context.Context, conditions common.GormConditions, tagID int64, pager *dto.Pager) (*dto.ListArticleRespDTO, error) {
	var (
		total      int64
		articlePOs []*model.Article
		eg         errgroup.Group
	)

	// 排序仅用于列表查询，不能带入COUNT
	selectConditions := append(append(common.GormConditions{}, conditions...), &common.OrderCond{Columns: []string{"is_top DESC", "published_on DESC", "id DESC"}})
	var p *common.Pagination
	if pager != nil {
		p = &common.Pagination{Offset: pager.Offset, PageSize: pager.PageSize}
	}

	eg.Go(func() error {
		var (
			count int64
			err   error
		)
		if tagID > 0 {
			count, err = s.daos.Article.CountByTagID(ctx, conditions, tagID)
		} else {
			count, err = s.daos.Article.Count(ctx, conditions)
		}
		if err != nil {
			return errors.Wrap(err, "select article count")
		}
		total = count
		return nil
	})

	eg.Go(func() error {
		var (
			pos []*model.Article
			err error
		)
		if tagID > 0 {
			pos, err = s.daos.Article.SelectManyByTagID(ctx, selectConditions, p, tagID)
		} else {
			pos, err = s.daos.Article.SelectMany(ctx, selectConditions, p)
		}
		if err != nil {
			return errors.Wrap(err, "select article list")
		}
		articlePOs = pos
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	items, err := s.enrich(ctx, articlePOs)
	if err != nil {
		return nil, err
	}

	if pager != nil {
		pager.Total = total
	}
	return &dto.ListArticleRespDTO{
		Items: items,
		Pager: pager,
	}, nil
}

// enrich 批量聚合文章的分类与标签
func (s *articleService) enrich(ctx context.Context, articles []*model.Article) ([]*dto.ArticleItem, error) {
	items := make([]*dto.ArticleItem, 0, len(articles))
	if len(articles) == 0 {
		return items, nil
	}

	articleIDs := make([]int64, 0, len(articles))
	categoryIDSet := make(map[int64]struct{})
	for _, a := range articles {
		articleIDs = append(articleIDs, a.ID)
		if a.CategoryID > 0 {
			categoryIDSet[a.CategoryID] = struct{}{}
		}
	}

	// 标签：关联关系 + 标签详情
	relations, err := s.daos.Tag.SelectRelationsByArticleIDs(ctx, articleIDs)
	if err != nil {
		return nil, errors.Wrap(err, "select article tag relations")
	}
	tagIDSet := make(map[int64]struct{})
	for _, rel := range relations {
		tagIDSet[rel.TagID] = struct{}{}
	}
	tagIDs := make([]int64, 0, len(tagIDSet))
	for id := range tagIDSet {
		tagIDs = append(tagIDs, id)
	}
	tags, err := s.daos.Tag.SelectManyByIDs(ctx, tagIDs)
	if err != nil {
		return nil, errors.Wrap(err, "select tags")
	}
	tagMap := make(map[int64]*model.Tag, len(tags))
	for _, t := range tags {
		tagMap[t.ID] = t
	}
	articleTagsMap := make(map[int64][]*model.Tag)
	for _, rel := range relations {
		if tag, ok := tagMap[rel.TagID]; ok {
			articleTagsMap[rel.ArticleID] = append(articleTagsMap[rel.ArticleID], tag)
		}
	}

	// 分类
	categoryIDs := make([]int64, 0, len(categoryIDSet))
	for id := range categoryIDSet {
		categoryIDs = append(categoryIDs, id)
	}
	categories, err := s.daos.Category.SelectManyByIDs(ctx, categoryIDs)
	if err != nil {
		return nil, errors.Wrap(err, "select categories")
	}
	categoryMap := make(map[int64]*model.Category, len(categories))
	for _, c := range categories {
		categoryMap[c.ID] = c
	}

	for _, a := range articles {
		tags := articleTagsMap[a.ID]
		if tags == nil {
			tags = []*model.Tag{}
		}
		items = append(items, &dto.ArticleItem{
			PO:       a,
			Category: categoryMap[a.CategoryID],
			Tags:     tags,
		})
	}
	return items, nil
}

// Detail 公开文章详情（含草稿校验与浏览量计数）
func (s *articleService) Detail(ctx context.Context, reqDTO *dto.IDReqDTO) (*dto.ArticleDetailRespDTO, error) {
	cacheVersion := articleCacheVersion(ctx)
	cacheKey := articleIDCacheKey(cacheVersion, reqDTO.ID)
	if cached := s.getCachedDetail(cacheKey); cached != nil {
		cached.ViewDelta = s.incrView(cached.PO.ID)
		return cached, nil
	}

	po, err := s.daos.Article.SelectOne(ctx, reqDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.NewCustomError(errcode.ErrArticleNotFound)
		}
		return nil, errors.Wrap(err, "select article fail")
	}
	if po.IsDraft == 1 || po.State != model.StateEnabled {
		return nil, errcode.NewCustomError(errcode.ErrArticleNotFound)
	}

	items, err := s.enrich(ctx, []*model.Article{po})
	if err != nil {
		return nil, err
	}

	result := &dto.ArticleDetailRespDTO{
		PO:       po,
		Category: items[0].Category,
		Tags:     items[0].Tags,
	}
	s.cacheDetail(result, cacheVersion)
	result.ViewDelta = s.incrView(po.ID)
	return result, nil
}

// DetailBySlug 按slug查询公开文章详情
func (s *articleService) DetailBySlug(ctx context.Context, reqDTO *dto.ArticleSlugReqDTO) (*dto.ArticleDetailRespDTO, error) {
	cacheVersion := articleCacheVersion(ctx)
	cacheKey := articleSlugCacheKey(cacheVersion, reqDTO.Slug)
	if cached := s.getCachedDetail(cacheKey); cached != nil {
		cached.ViewDelta = s.incrView(cached.PO.ID)
		return cached, nil
	}

	po, err := s.daos.Article.SelectOneBySlug(ctx, reqDTO.Slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.NewCustomError(errcode.ErrArticleNotFound)
		}
		return nil, errors.Wrap(err, "select article by slug fail")
	}
	if po.IsDraft == 1 || po.State != model.StateEnabled {
		return nil, errcode.NewCustomErrorWithMessage(errcode.ErrArticleNotFound, "slug: "+reqDTO.Slug)
	}

	items, err := s.enrich(ctx, []*model.Article{po})
	if err != nil {
		return nil, err
	}

	result := &dto.ArticleDetailRespDTO{
		PO:       po,
		Category: items[0].Category,
		Tags:     items[0].Tags,
	}
	s.cacheDetail(result, cacheVersion)
	result.ViewDelta = s.incrView(po.ID)
	return result, nil
}

// AdminDetail 管理端文章详情（含草稿，不计浏览量）
func (s *articleService) AdminDetail(ctx context.Context, reqDTO *dto.IDReqDTO) (*dto.ArticleDetailRespDTO, error) {
	po, err := s.daos.Article.SelectOne(ctx, reqDTO.ID)
	if err != nil {
		return nil, errors.Wrap(err, "select article fail")
	}

	items, err := s.enrich(ctx, []*model.Article{po})
	if err != nil {
		return nil, err
	}

	return &dto.ArticleDetailRespDTO{
		PO:       po,
		Category: items[0].Category,
		Tags:     items[0].Tags,
	}, nil
}

// Add 新建文章
func (s *articleService) Add(ctx context.Context, reqDTO *dto.AddArticleReqDTO) error {
	if strings.TrimSpace(reqDTO.Title) == "" || strings.TrimSpace(reqDTO.ContentMd) == "" {
		return errcode.NewCustomErrorWithMessage(errcode.ErrInvalidParams, "标题和正文不能为空")
	}
	if err := s.validateRef(ctx, reqDTO.CategoryID, reqDTO.TagIDs); err != nil {
		return err
	}
	reqDTO.TagIDs = uniqueInt64(reqDTO.TagIDs)
	slug := buildSlug(reqDTO.Slug, reqDTO.Title)
	exists, err := s.daos.Article.ExistSlug(ctx, slug, 0)
	if err != nil {
		return errors.Wrap(err, "check article slug fail")
	}
	if exists {
		return errcode.NewCustomErrorWithMessage(errcode.ErrInvalidParams, "slug 已存在")
	}
	contentHTML, err := markdown.Render(reqDTO.ContentMd)
	if err != nil {
		return errors.Wrap(err, "render article markdown")
	}

	now := uint32(xtime.GetTimestamp())
	po := &model.Article{
		CategoryID:  reqDTO.CategoryID,
		Title:       reqDTO.Title,
		Slug:        slug,
		Desc:        reqDTO.Desc,
		Cover:       reqDTO.Cover,
		ContentMd:   reqDTO.ContentMd,
		ContentHtml: contentHTML,
		IsTop:       0,
		IsDraft:     boolToInt8(reqDTO.IsDraft),
		PublishedOn: 0,
	}
	po.State = model.StateEnabled
	po.CreatedBy = reqDTO.Username
	po.ModifiedBy = reqDTO.Username
	po.CreatedOn = now
	po.ModifiedOn = now
	if !reqDTO.IsDraft {
		po.PublishedOn = now
	}

	if err := s.daos.Article.Create(ctx, po, reqDTO.TagIDs); err != nil {
		return errors.Wrap(err, "create article fail")
	}
	s.invalidateArchiveCache()
	s.invalidateListCache()
	return nil
}

// Edit 编辑文章
func (s *articleService) Edit(ctx context.Context, reqDTO *dto.EditArticleReqDTO) error {
	po, err := s.daos.Article.SelectOne(ctx, reqDTO.ID)
	if err != nil {
		return errcode.NewCustomError(errcode.ErrArticleNotFound)
	}

	if err := s.validateRef(ctx, reqDTO.CategoryID, reqDTO.TagIDs); err != nil {
		return err
	}
	reqDTO.TagIDs = uniqueInt64(reqDTO.TagIDs)
	if strings.TrimSpace(reqDTO.Title) == "" || strings.TrimSpace(reqDTO.ContentMd) == "" {
		return errcode.NewCustomErrorWithMessage(errcode.ErrInvalidParams, "标题和正文不能为空")
	}
	slug := po.Slug
	if strings.TrimSpace(reqDTO.Slug) != "" || slug == "" {
		slug = buildSlug(reqDTO.Slug, reqDTO.Title)
	}
	exists, err := s.daos.Article.ExistSlug(ctx, slug, reqDTO.ID)
	if err != nil {
		return errors.Wrap(err, "check article slug fail")
	}
	if exists {
		return errcode.NewCustomErrorWithMessage(errcode.ErrInvalidParams, "slug 已存在")
	}
	contentHTML, err := markdown.Render(reqDTO.ContentMd)
	if err != nil {
		return errors.Wrap(err, "render article markdown")
	}

	now := uint32(xtime.GetTimestamp())
	oldSlug := po.Slug
	po.CategoryID = reqDTO.CategoryID
	po.Title = reqDTO.Title
	po.Slug = slug
	po.Desc = reqDTO.Desc
	po.Cover = reqDTO.Cover
	po.ContentMd = reqDTO.ContentMd
	po.ContentHtml = contentHTML
	po.ModifiedBy = reqDTO.Username
	po.ModifiedOn = now

	// 草稿 -> 发布时补充发布时间
	if po.IsDraft == 1 && !reqDTO.IsDraft && po.PublishedOn == 0 {
		po.PublishedOn = now
	}
	po.IsDraft = boolToInt8(reqDTO.IsDraft)

	if err := s.daos.Article.Update(ctx, po, reqDTO.TagIDs, true); err != nil {
		return errors.Wrap(err, "update article fail")
	}
	s.invalidateArchiveCache()
	s.invalidateContentCaches(po.ID, oldSlug, po.Slug)
	return nil
}

// Del 删除文章（软删）
func (s *articleService) Del(ctx context.Context, reqDTO *dto.IDReqDTO) error {
	po, err := s.daos.Article.SelectOne(ctx, reqDTO.ID)
	if err != nil {
		return errcode.NewCustomError(errcode.ErrArticleNotFound)
	}

	if err := s.daos.Article.Delete(ctx, reqDTO.ID); err != nil {
		return errors.Wrap(err, "delete article fail")
	}
	s.invalidateArchiveCache()
	s.invalidateContentCaches(po.ID, po.Slug)
	return nil
}

// Publish 发布/下架
func (s *articleService) Publish(ctx context.Context, reqDTO *dto.PublishArticleReqDTO) error {
	po, err := s.daos.Article.SelectOne(ctx, reqDTO.ID)
	if err != nil {
		return errcode.NewCustomErrorWithMessage(errcode.ErrArticleNotFound, "id: "+strconv.FormatInt(reqDTO.ID, 10))
	}

	cols := map[string]interface{}{
		"is_draft":    boolToInt8(!reqDTO.Publish),
		"modified_on": xtime.GetTimestamp(),
		"modified_by": reqDTO.Username,
	}
	// 首次发布补充发布时间
	if reqDTO.Publish && po.PublishedOn == 0 {
		cols["published_on"] = xtime.GetTimestamp()
	}

	if err := s.daos.Article.UpdateColumns(ctx, reqDTO.ID, cols); err != nil {
		return errors.Wrap(err, "publish article fail")
	}
	s.invalidateArchiveCache()
	s.invalidateContentCaches(po.ID, po.Slug)
	return nil
}

// State 启用/禁用
func (s *articleService) State(ctx context.Context, reqDTO *dto.StateReqDTO) error {
	po, err := s.daos.Article.SelectOne(ctx, reqDTO.ID)
	if err != nil {
		return errcode.NewCustomError(errcode.ErrArticleNotFound)
	}

	cols := map[string]interface{}{
		"state":       reqDTO.State,
		"modified_on": xtime.GetTimestamp(),
		"modified_by": reqDTO.Username,
	}
	if err := s.daos.Article.UpdateColumns(ctx, reqDTO.ID, cols); err != nil {
		return errors.Wrap(err, "update article state fail")
	}
	s.invalidateArchiveCache()
	s.invalidateContentCaches(reqDTO.ID, po.Slug)
	return nil
}

// Top 置顶/取消置顶
func (s *articleService) Top(ctx context.Context, reqDTO *dto.TopArticleReqDTO) error {
	po, err := s.daos.Article.SelectOne(ctx, reqDTO.ID)
	if err != nil {
		return errcode.NewCustomError(errcode.ErrArticleNotFound)
	}

	cols := map[string]interface{}{
		"is_top":      reqDTO.IsTop,
		"modified_on": xtime.GetTimestamp(),
		"modified_by": reqDTO.Username,
	}
	if err := s.daos.Article.UpdateColumns(ctx, reqDTO.ID, cols); err != nil {
		return errors.Wrap(err, "update article top fail")
	}
	s.invalidateContentCaches(reqDTO.ID, po.Slug)
	return nil
}

// Archive 按月归档（Redis缓存1小时）
func (s *articleService) Archive(ctx context.Context) (*dto.ArchiveRespDTO, error) {
	var rows []*model.ArchiveRow
	if err := redis.GetStruct(archiveCacheKey, &rows); err == nil && rows != nil {
		return &dto.ArchiveRespDTO{Rows: rows}, nil
	}

	rows, err := s.daos.Article.Archive(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "select archive fail")
	}

	// 缓存失败不影响响应
	_ = redis.SaveStruct(archiveCacheKey, rows, archiveCacheTTL)
	return &dto.ArchiveRespDTO{Rows: rows}, nil
}

// validateRef 校验文章引用的分类与标签存在
func (s *articleService) validateRef(ctx context.Context, categoryID int64, tagIDs []int64) error {
	if categoryID > 0 {
		if _, err := s.daos.Category.SelectOne(ctx, categoryID); err != nil {
			return errcode.NewCustomError(errcode.ErrCategoryNotFound)
		}
	}

	if len(tagIDs) > 0 {
		tags, err := s.daos.Tag.SelectManyByIDs(ctx, tagIDs)
		if err != nil {
			return errors.Wrap(err, "select tags fail")
		}
		if len(tags) != len(uniqueInt64(tagIDs)) {
			return errcode.NewCustomError(errcode.ErrTagNotFound)
		}
	}
	return nil
}

// incrView 浏览量Redis计数，返回当日累计增量（Redis异常时返回0，不影响详情）
func (s *articleService) incrView(articleID int64) int64 {
	if redis.RedisClient == nil {
		return 0
	}
	key := viewCountKeyPrefix + time.Now().Format("20060102")
	field := strconv.FormatInt(articleID, 10)
	n, err := redis.RedisClient.HIncrBy(context.Background(), key, field, 1).Result()
	if err != nil {
		return 0
	}
	if n == 1 {
		_ = redis.RedisClient.Expire(context.Background(), key, 7*24*time.Hour).Err()
	}
	return n
}

// invalidateArchiveCache 文章变动后清除归档缓存
func (s *articleService) invalidateArchiveCache() {
	if redis.RedisClient != nil {
		_ = redis.RedisClient.Del(context.Background(), archiveCacheKey).Err()
	}
	s.invalidateListCache()
}

func (s *articleService) invalidateListCache() {
	if redis.RedisClient != nil {
		_ = redis.RedisClient.Incr(context.Background(), listVersionKey).Err()
	}
}

func (s *articleService) invalidateContentCaches(articleID int64, slugs ...string) {
	InvalidateArticleCaches()
}

func (s *articleService) cacheDetail(detail *dto.ArticleDetailRespDTO, version string) {
	if detail == nil || detail.PO == nil {
		return
	}
	copy := *detail
	copy.ViewDelta = 0
	_, _ = redis.SaveStructNX(articleIDCacheKey(version, detail.PO.ID), &copy, articleCacheTTL)
	if detail.PO.Slug != "" {
		_, _ = redis.SaveStructNX(articleSlugCacheKey(version, detail.PO.Slug), &copy, articleCacheTTL)
	}
}

func articleIDCacheKey(version string, id int64) string {
	return articleCachePrefix + strconv.FormatInt(id, 10) + ":v" + version
}

func articleSlugCacheKey(version, slug string) string {
	return slugCachePrefix + slug + ":v" + version
}

func articleCacheVersion(ctx context.Context) string {
	if redis.RedisClient != nil {
		if value, err := redis.RedisClient.Get(ctx, detailVersionKey).Result(); err == nil && value != "" {
			return value
		}
	}
	return "0"
}

// InvalidateArticleCaches invalidates public article lists and detail entries.
// Versioned keys avoid blocking wildcard deletion and expire naturally.
func InvalidateArticleCaches() {
	if redis.RedisClient == nil {
		return
	}
	ctx := context.Background()
	pipe := redis.RedisClient.TxPipeline()
	pipe.Incr(ctx, listVersionKey)
	pipe.Incr(ctx, detailVersionKey)
	_, _ = pipe.Exec(ctx)
}

func (s *articleService) getCachedDetail(key string) *dto.ArticleDetailRespDTO {
	var cached dto.ArticleDetailRespDTO
	if err := redis.GetStruct(key, &cached); err != nil || cached.PO == nil {
		return nil
	}
	return &cached
}

func (s *articleService) listCacheKey(ctx context.Context, req *dto.ListArticleReqDTO) string {
	version := "0"
	if redis.RedisClient != nil {
		if value, err := redis.RedisClient.Get(ctx, listVersionKey).Result(); err == nil && value != "" {
			version = value
		}
	}
	query := fmt.Sprintf("%s|%d|%d|%d|%d|%d|%s", version, req.Pager.Page, req.Pager.PageSize, req.ID, req.CategoryID, req.TagID, req.Keyword)
	hash := sha256.Sum256([]byte(query))
	return fmt.Sprintf("%s%x", listCachePrefix, hash[:])
}

func buildSlug(slug, title string) string {
	value := strings.ToLower(strings.TrimSpace(slug))
	if value == "" {
		value = strings.ToLower(strings.TrimSpace(title))
	}
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case b.Len() > 0 && !lastDash:
			b.WriteByte('-')
			lastDash = true
		}
	}
	result := strings.Trim(b.String(), "-")
	if result == "" {
		hash := sha256.Sum256([]byte(title))
		result = "article-" + fmt.Sprintf("%x", hash[:6])
	}
	if len(result) > 150 {
		result = result[:150]
	}
	return result
}

func boolToInt8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func uniqueInt64(ids []int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}
