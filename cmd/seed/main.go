// Command seed inserts representative blog data into the configured database.
// It is safe to run repeatedly: rows are identified by their existing unique keys
// or by stable test-data markers.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"gin-web/core/config"
	_ "github.com/go-sql-driver/mysql"
)

type articleSeed struct {
	Category string
	TagNames []string
	Slug     string
	Title    string
	Desc     string
	Markdown string
	HTML     string
	Draft    int
	Top      int
	Views    int
}

func main() {
	dbConfig := config.NewDbConfig("blog")
	db, err := sql.Open("mysql", dbConfig.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("connect database: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	now := time.Now().Unix()
	categoryIDs := map[string]int64{}
	for _, name := range []string{"技术实践", "生活随笔", "项目记录"} {
		if _, err := tx.Exec(`INSERT IGNORE INTO blog_category
			(name, parent_id, sort, created_on, created_by, modified_on, modified_by, state)
			VALUES (?, 0, 0, ?, 'seed', ?, 'seed', 1)`, name, now, now); err != nil {
			log.Fatal(err)
		}
		categoryIDs[name], err = rowID(tx, "SELECT id FROM blog_category WHERE name = ? AND deleted_on = 0 LIMIT 1", name)
		if err != nil {
			log.Fatal(err)
		}
	}

	tagIDs := map[string]int64{}
	for _, name := range []string{"Go", "Gin", "MySQL", "Redis", "测试数据"} {
		if _, err := tx.Exec(`INSERT IGNORE INTO blog_tag
			(name, created_on, created_by, modified_on, modified_by, state)
			VALUES (?, ?, 'seed', ?, 'seed', 1)`, name, now, now); err != nil {
			log.Fatal(err)
		}
		tagIDs[name], err = rowID(tx, "SELECT id FROM blog_tag WHERE name = ? AND deleted_on = 0 LIMIT 1", name)
		if err != nil {
			log.Fatal(err)
		}
	}

	articles := []articleSeed{
		{
			Category: "技术实践", TagNames: []string{"Go", "Gin", "测试数据"}, Slug: "test-gin-basics",
			Title: "Gin 博客接口入门", Desc: "用于验证公开文章列表和详情接口的测试文章。",
			Markdown: "# Gin 博客接口入门\n\n这是一篇用于本地联调的测试文章。", HTML: "<h1>Gin 博客接口入门</h1><p>这是一篇用于本地联调的测试文章。</p>", Top: 1, Views: 128,
		},
		{
			Category: "技术实践", TagNames: []string{"Go", "MySQL", "测试数据"}, Slug: "test-mysql-migrations",
			Title: "MySQL 数据库迁移实践", Desc: "用于验证分类、标签和文章关联查询的测试文章。",
			Markdown: "# MySQL 数据库迁移实践\n\n记录一次数据库结构调整。", HTML: "<h1>MySQL 数据库迁移实践</h1><p>记录一次数据库结构调整。</p>", Views: 64,
		},
		{
			Category: "项目记录", TagNames: []string{"Redis", "Go", "测试数据"}, Slug: "test-redis-view-count",
			Title: "Redis 浏览量统计", Desc: "用于验证浏览量和文章归档数据的测试文章。",
			Markdown: "# Redis 浏览量统计\n\n浏览量先写入 Redis，再异步刷入 MySQL。", HTML: "<h1>Redis 浏览量统计</h1><p>浏览量先写入 Redis，再异步刷入 MySQL。</p>", Views: 32,
		},
		{
			Category: "生活随笔", TagNames: []string{"测试数据"}, Slug: "test-draft-article",
			Title: "待发布的草稿文章", Desc: "用于验证草稿不会出现在公开列表中的测试文章。",
			Markdown: "# 待发布的草稿文章\n\n这篇文章仍处于草稿状态。", HTML: "<h1>待发布的草稿文章</h1><p>这篇文章仍处于草稿状态。</p>", Draft: 1,
		},
	}

	articleIDs := map[string]int64{}
	for index, article := range articles {
		publishedOn := uint32(0)
		if article.Draft == 0 {
			publishedOn = uint32(now - int64((index+1)*86400))
		}
		if _, err := tx.Exec(`INSERT IGNORE INTO blog_article
			(category_id, title, slug, `+"`desc`"+`, cover, content_md, content_html, view_count,
			 is_top, is_draft, published_on, created_on, created_by, modified_on, modified_by, state)
			VALUES (?, ?, ?, ?, '', ?, ?, ?, ?, ?, ?, ?, 'seed', ?, 'seed', 1)`,
			categoryIDs[article.Category], article.Title, article.Slug, article.Desc, article.Markdown, article.HTML,
			article.Views, article.Top, article.Draft, publishedOn, now, now); err != nil {
			log.Fatal(err)
		}
		articleIDs[article.Slug], err = rowID(tx, "SELECT id FROM blog_article WHERE slug = ? AND deleted_on = 0 LIMIT 1", article.Slug)
		if err != nil {
			log.Fatal(err)
		}
		for _, tagName := range article.TagNames {
			if _, err := tx.Exec(`INSERT IGNORE INTO blog_article_tag (article_id, tag_id) VALUES (?, ?)`, articleIDs[article.Slug], tagIDs[tagName]); err != nil {
				log.Fatal(err)
			}
		}
	}

	comments := []struct {
		article  string
		nickname string
		email    string
		content  string
		admin    int
		parent   string
	}{
		{articles[0].Slug, "访客小明", "xiaoming@example.com", "这篇文章对接口联调很有帮助。", 0, ""},
		{articles[0].Slug, "博主", "admin@example.com", "感谢反馈，后续会继续补充示例。", 1, "这篇文章对接口联调很有帮助。"},
		{articles[1].Slug, "测试读者", "reader@example.com", "请问数据库初始化脚本在哪里？", 0, ""},
	}
	for _, comment := range comments {
		parentID := int64(0)
		if comment.parent != "" {
			parentID, err = rowID(tx, "SELECT id FROM blog_comment WHERE article_id = ? AND content = ? AND deleted_on = 0 LIMIT 1", articleIDs[comment.article], comment.parent)
			if err != nil {
				log.Fatal(err)
			}
		}
		if _, err := tx.Exec(`INSERT INTO blog_comment
			(article_id, parent_id, nickname, email, content, is_admin, ip, created_on, created_by, modified_on, modified_by, state)
			SELECT ?, ?, ?, ?, ?, ?, '127.0.0.1', ?, 'seed', ?, 'seed', 1
			WHERE NOT EXISTS (SELECT 1 FROM blog_comment WHERE article_id = ? AND content = ? AND deleted_on = 0)`,
			articleIDs[comment.article], parentID, comment.nickname, comment.email, comment.content, comment.admin,
			now, now, articleIDs[comment.article], comment.content); err != nil {
			log.Fatal(err)
		}
	}

	links := []struct{ name, url, description string }{
		{"Gin", "https://gin-gonic.com/", "Gin Web Framework 官方网站"},
		{"Go", "https://go.dev/", "Go 官方网站"},
	}
	for _, link := range links {
		if _, err := tx.Exec(`INSERT INTO blog_link
			(name, url, logo, description, sort, created_on, created_by, modified_on, modified_by, state)
			SELECT ?, ?, '', ?, 0, ?, 'seed', ?, 'seed', 1
			WHERE NOT EXISTS (SELECT 1 FROM blog_link WHERE name = ? AND url = ? AND deleted_on = 0)`,
			link.name, link.url, link.description, now, now, link.name, link.url); err != nil {
			log.Fatal(err)
		}
	}

	options := []struct{ key, value, remark string }{
		{"site_title", `"gin-web 测试博客"`, "测试站点标题"},
		{"site_desc", `"用于本地联调的博客站点"`, "测试站点描述"},
		{"comment_audit", `"0"`, "测试环境关闭评论审核"},
	}
	for _, option := range options {
		if _, err := tx.Exec(`INSERT IGNORE INTO blog_option (`+"`key`"+`, value, remark, modified_on, modified_by, state)
			VALUES (?, ?, ?, ?, 'seed', 1)`, option.key, option.value, option.remark, now); err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	for _, table := range []string{"blog_user", "blog_category", "blog_tag", "blog_article", "blog_article_tag", "blog_comment", "blog_link", "blog_option"} {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-18s %d\n", table, count)
	}
}

func rowID(tx *sql.Tx, query string, args ...any) (int64, error) {
	var id int64
	err := tx.QueryRow(query, args...).Scan(&id)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		return 0, fmt.Errorf("seed lookup returned no rows: %w", err)
	}
	return id, err
}
