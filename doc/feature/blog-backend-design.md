# 个人博客后端设计方案

> 版本：v1.0
> 日期：2026-10-01
> 状态：初稿（待评审）

---

## 1. 文档说明

### 1.1 背景

基于现有 gin-web 脚手架（Gin + GORM + MySQL + Redis）构建一个个人博客后端服务。现有代码已具备：登录认证（JWT）、标签管理、文章管理（列表/新增）的基础实现，本方案在其之上补全并扩展为完整的博客系统。

### 1.2 目标

- 为个人博客前端（Web 站点 + 管理后台）提供稳定的 RESTful API。
- 支持文章全生命周期管理（草稿 → 发布 → 归档）、分类/标签、评论、友情链接、站点配置、图片上传。
- 保持与现有分层架构、目录结构、编码规范一致，可渐进式落地。

### 1.3 范围

- 包含：API 设计、数据库设计、缓存设计、认证授权、日志与可观测性、部署方案、测试策略。
- 不包含：前端实现、搜索引擎（ES）部署（列为后续扩展）。

---

## 2. 需求分析

### 2.1 功能需求

| 模块 | 功能点 | 优先级 |
| --- | --- | --- |
| 认证 | 登录、登出、获取当前用户信息 | P0 |
| 文章 | 列表（分页/筛选/排序）、详情、新增、编辑、删除、发布/下架、置顶、浏览量统计 | P0 |
| 分类 | CRUD、分类下文章数统计 | P0 |
| 标签 | CRUD（已有，需改造为多对多）、标签下文章数统计 | P0 |
| 评论 | 发表评论、回复（二级）、评论列表（文章维度）、删除（管理员） | P1 |
| 文件 | 图片上传（本地/对象存储）、大小与格式校验 | P1 |
| 友情链接 | CRUD、前台展示 | P2 |
| 站点配置 | 站点标题、简介、备案号、社交链接等 KV 配置（已有雏形） | P2 |
| 归档 | 按年月归档文章 | P2 |
| 搜索 | 文章标题/摘要/内容模糊搜索 | P2 |
| 订阅 | RSS、Sitemap 输出 | P3 |

### 2.2 非功能需求

| 指标 | 目标 |
| --- | --- |
| 接口平均响应时间 | < 100ms（缓存命中时 < 10ms） |
| 可用性 | 单机部署 99.5%+，支持快速重启恢复 |
| 并发 | 个人博客量级，峰值 QPS 100 以内，读多写少 |
| 安全 | 密码 bcrypt 哈希、JWT 鉴权、参数校验、防 SQL 注入（GORM 参数化）、防 XSS |
| 可维护性 | 分层清晰、单测覆盖核心逻辑（service/dao 层） |

---

## 3. 技术选型

沿用现有脚手架技术栈，不引入新框架：

| 类别 | 选型 | 说明 |
| --- | --- | --- |
| 语言/框架 | Go 1.20+ / Gin | 现有 |
| ORM | GORM v1.25 | 现有，参数化查询防注入 |
| 数据库 | MySQL 8.0 | InnoDB，utf8mb4 |
| 缓存 | Redis 7（go-redis v9） | 缓存、计数器、分布式锁 |
| 认证 | JWT（golang-jwt/v5） | 现有 |
| 密码哈希 | bcrypt | 现有 |
| 依赖注入 | google/wire | 现有，编译期注入 |
| 日志 | logrus + 按天切分 + 敏感信息脱敏 | 现有 |
| 链路追踪 | SkyWalking（可选） | 现有 |
| 参数校验 | go-playground/validator | 现有 |
| API 文档 | swag（Swagger 2.0） | 现有 |
| 配置 | viper（app.toml） | 现有 |
| 测试 | go-sqlmock + testify | 现有 |
| 部署 | Docker + Nginx + GitHub Actions CI | 现有 ci.yml |

---

## 4. 系统架构

### 4.1 总体架构

```
            ┌────────────────────────────────────────────┐
            │                  Nginx                      │
            │   静态资源 / HTTPS / gzip / 反向代理 / 限流  │
            └───────────────────┬────────────────────────┘
                                │
┌───────────────────────────────▼───────────────────────────────┐
│                       gin-web (Go 进程)                        │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │ Middleware: RequestID → Logger → Recovery → CORS → JWT   │  │
│  └───────────────────────────┬─────────────────────────────┘  │
│  ┌───────────────────────────▼─────────────────────────────┐  │
│  │ Router  /api/v1/*                                        │  │
│  │  ├─ 公开路由（文章浏览、评论列表、登录…）                   │  │
│  │  └─ JWT 路由（管理：文章/分类/标签/评论/友链/配置）          │  │
│  └───────────────────────────┬─────────────────────────────┘  │
│  ┌───────────────────────────▼─────────────────────────────┐  │
│  │ Handler（app/handler）  参数绑定/校验 → 调 Service → 响应   │  │
│  └───────────────────────────┬─────────────────────────────┘  │
│  ┌───────────────────────────▼─────────────────────────────┐  │
│  │ Service（internal/service） 业务逻辑、事务编排、缓存读写     │  │
│  └───────────────────────────┬─────────────────────────────┘  │
│  ┌───────────────────────────▼─────────────────────────────┐  │
│  │ DAO（internal/dao）  GORM 数据访问，不含业务                │  │
│  └──────────┬────────────────────────────────┬─────────────┘  │
└─────────────┼────────────────────────────────┼────────────────┘
              ▼                                ▼
      ┌──────────────┐                 ┌──────────────┐
      │    MySQL     │                 │    Redis     │
      │  业务数据主存储 │                 │ 缓存/计数/锁  │
      └──────────────┘                 └──────────────┘
              ▲
      ┌──────────────┐
      │  本地磁盘/对象  │
      │  存储（图片）   │
      └──────────────┘
```

### 4.2 分层职责与依赖方向

依赖方向自上而下，单向依赖：

| 层 | 目录 | 职责 | 禁止 |
| --- | --- | --- | --- |
| Handler | `app/handler` | 绑定/校验参数（dto）、调用 service、组装响应 | 不写业务逻辑、不直接访问 DB |
| Service | `internal/service` | 业务逻辑、事务编排、缓存读写、错误码转换 | 不感知 HTTP（不 import gin） |
| DAO | `internal/dao` | GORM 数据访问、通用查询条件（gorm_cond）、分页 | 不含业务语义 |
| Model | `internal/model` | 数据库表结构体 | — |
| DTO | `dto` | 请求/响应传输对象，带 `binding` 校验标签 | — |

### 4.3 目录结构（新增部分）

```
gin-web/
├── app/
│   ├── handler/
│   │   ├── article.go      # 扩展：详情/编辑/删除/发布
│   │   ├── category.go     # 新增
│   │   ├── comment.go      # 新增
│   │   ├── upload.go       # 新增
│   │   ├── link.go         # 新增
│   │   ├── auth.go         # 扩展：登出/当前用户
│   │   └── ...
│   └── router/v1.go        # 注册新路由
├── dto/
│   ├── article.go          # 扩展
│   ├── category.go         # 新增
│   ├── comment.go          # 新增
│   └── ...
├── internal/
│   ├── dao/{article,category,comment,link,option}/
│   ├── model/{article,category,comment,link}.go
│   ├── service/{article,category,comment,link}/
│   └── errcode/errcode.go  # 新增错误码
├── core/
│   └── storage/            # 新增：文件存储抽象（本地/OSS）
├── cmd/api/main.go
└── doc/
    ├── feature/              # 设计文档
    ├── sql/db.sql            # 全量建表语句（全新建库直接执行）
    └── swagger/              # swag 生成的 API 文档（docs.go / swagger.json）
```

---

## 5. 数据库设计

### 5.1 设计约定

沿用现有约定：

- 表名前缀 `blog_`，全部小写下划线命名。
- 公共字段：`id`（自增主键）、`created_on`/`modified_on`（`int unsigned`，Unix 时间戳）、`created_by`/`modified_by`、`deleted_on`（软删除，0 表示未删除）、`state`（tinyint，0 禁用 / 1 启用）。
- 字符集 `utf8mb4`，引擎 InnoDB。
- GORM 模型通过现有 `gorm_hook.go` 统一填充时间戳与审计字段。

### 5.2 ER 关系

```
blog_user 1 ──── * blog_article * ──── * blog_tag        (blog_article_tag)
                  │
                  ├── * blog_category (多对一)
                  └── * blog_comment   (一篇文章多评论)
blog_option、blog_link 为独立表
```

### 5.3 表结构设计

> 存量本地数据不保留，以下均为最终 schema，直接按全新建库执行（同步重写 `doc/sql/db.sql`）。

#### 5.3.1 blog_user（账号，替代现有 blog_auth）

```sql
CREATE TABLE `blog_user`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT,
    `username`   varchar(50)  DEFAULT '' COMMENT '账号',
    `password`   varchar(100) DEFAULT '' COMMENT '密码(bcrypt哈希)',
    `nickname`   varchar(50)  DEFAULT '' COMMENT '昵称',
    `avatar`     varchar(255) DEFAULT '' COMMENT '头像URL',
    `email`      varchar(100) DEFAULT '' COMMENT '邮箱',
    `role`       tinyint unsigned DEFAULT '1' COMMENT '角色 1管理员(预留)',
    `last_login` int unsigned DEFAULT '0' COMMENT '最后登录时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='博客账号';
```

> 密码仍为 bcrypt 哈希（登录时 `bcrypt.CompareHashAndPassword` 校验），初始化脚本插入默认管理员账号。

#### 5.3.2 blog_category（分类，新增）

```sql
CREATE TABLE `blog_category`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(100) DEFAULT '' COMMENT '分类名称',
    `parent_id`   int unsigned DEFAULT '0' COMMENT '父分类ID，0为顶级',
    `sort`        int          DEFAULT '0' COMMENT '排序，越小越靠前',
    `created_on`  int unsigned DEFAULT '0' COMMENT '创建时间',
    `created_by`  varchar(100) DEFAULT '' COMMENT '创建人',
    `modified_on` int unsigned DEFAULT '0' COMMENT '修改时间',
    `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
    `deleted_on`  int unsigned DEFAULT '0',
    `state`       tinyint unsigned DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`, `deleted_on`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='文章分类';
```

#### 5.3.3 blog_tag（沿用现有结构，不变）

#### 5.3.4 blog_article（最终结构，含多对多标签改造）

```sql
CREATE TABLE `blog_article`
(
    `id`           int unsigned NOT NULL AUTO_INCREMENT,
    `category_id`  int unsigned DEFAULT '0' COMMENT '分类ID',
    `title`        varchar(100) DEFAULT '' COMMENT '文章标题',
    `slug`         varchar(150) DEFAULT '' COMMENT 'URL友好标识',
    `desc`         varchar(255) DEFAULT '' COMMENT '简述',
    `cover`        varchar(255) DEFAULT '' COMMENT '封面图URL',
    `content_md`   mediumtext COMMENT 'Markdown原文',
    `content_html` mediumtext COMMENT '渲染后HTML(预渲染)',
    `view_count`   int unsigned DEFAULT '0' COMMENT '浏览量',
    `is_top`       tinyint unsigned DEFAULT '0' COMMENT '是否置顶 0否1是',
    `is_draft`     tinyint unsigned DEFAULT '0' COMMENT '是否草稿 0否1是',
    `published_on` int unsigned DEFAULT '0' COMMENT '发布时间',
    `created_on`   int unsigned DEFAULT '0' COMMENT '创建时间',
    `created_by`   varchar(100) DEFAULT '' COMMENT '创建人',
    `modified_on`  int unsigned DEFAULT '0' COMMENT '修改时间',
    `modified_by`  varchar(255) DEFAULT '' COMMENT '修改人',
    `deleted_on`   int unsigned DEFAULT '0',
    `state`        tinyint unsigned DEFAULT '1' COMMENT '状态 0为禁用1为启用',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_slug` (`slug`, `deleted_on`),
    KEY `idx_category` (`category_id`),
    KEY `idx_published_on` (`published_on`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='文章管理';
```

> 标签关系由 `blog_article_tag` 关联表承载（见 5.3.5），不再使用单 `tag_id` 字段。

#### 5.3.5 blog_article_tag（文章-标签关联，新增）

```sql
CREATE TABLE `blog_article_tag`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT,
    `article_id` int unsigned NOT NULL COMMENT '文章ID',
    `tag_id`     int unsigned NOT NULL COMMENT '标签ID',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_article_tag` (`article_id`, `tag_id`),
    KEY `idx_tag` (`tag_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='文章标签关联';
```

#### 5.3.6 blog_comment（评论，新增）

```sql
CREATE TABLE `blog_comment`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `article_id`  int unsigned NOT NULL COMMENT '文章ID',
    `parent_id`   int unsigned DEFAULT '0' COMMENT '父评论ID，0为一级评论',
    `nickname`    varchar(50)  DEFAULT '' COMMENT '评论者昵称',
    `email`       varchar(100) DEFAULT '' COMMENT '评论者邮箱(不展示)',
    `content`     varchar(2000) DEFAULT '' COMMENT '评论内容',
    `is_admin`    tinyint unsigned DEFAULT '0' COMMENT '是否博主回复 0否1是',
    `ip`          varchar(64)  DEFAULT '' COMMENT '来源IP',
    `created_on`  int unsigned DEFAULT '0',
    `created_by`  varchar(100) DEFAULT '',
    `modified_on` int unsigned DEFAULT '0',
    `modified_by` varchar(100) DEFAULT '',
    `deleted_on`  int unsigned DEFAULT '0',
    `state`       tinyint unsigned DEFAULT '1' COMMENT '0待审核 1通过 2拒绝',
    PRIMARY KEY (`id`),
    KEY `idx_article` (`article_id`, `state`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='文章评论';
```

#### 5.3.7 blog_link（友情链接，新增）

```sql
CREATE TABLE `blog_link`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(100) DEFAULT '' COMMENT '站点名称',
    `url`         varchar(255) DEFAULT '' COMMENT '站点URL',
    `logo`        varchar(255) DEFAULT '' COMMENT '站点LOGO',
    `description` varchar(255) DEFAULT '' COMMENT '描述',
    `sort`        int          DEFAULT '0' COMMENT '排序',
    `created_on`  int unsigned DEFAULT '0',
    `created_by`  varchar(100) DEFAULT '',
    `modified_on` int unsigned DEFAULT '0',
    `modified_by` varchar(100) DEFAULT '',
    `deleted_on`  int unsigned DEFAULT '0',
    `state`       tinyint unsigned DEFAULT '1',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='友情链接';
```

#### 5.3.8 blog_option（站点配置 KV，新增）

```sql
CREATE TABLE `blog_option`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `key`         varchar(100) NOT NULL COMMENT '配置键',
    `value`       text COMMENT '配置值(JSON)',
    `remark`      varchar(255) DEFAULT '' COMMENT '说明',
    `modified_on` int unsigned DEFAULT '0',
    `modified_by` varchar(100) DEFAULT '',
    `deleted_on`  int unsigned DEFAULT '0',
    `state`       tinyint unsigned DEFAULT '1',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_key` (`key`, `deleted_on`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='站点配置';
```

> 现有 `GET /common/option` 改为读本表并做 Redis 缓存。

### 5.4 索引与查询匹配

| 高频查询 | 命中索引 |
| --- | --- |
| 文章列表（按发布时间倒序，过滤草稿/软删） | `idx_published_on` + `deleted_on`/`is_draft` 过滤 |
| 分类页文章列表 | `idx_category` |
| 文章详情（slug） | `uk_slug` |
| 标签页文章列表 | 关联表 `idx_tag` |
| 文章评论列表 | `idx_article` |

---

## 6. API 设计

### 6.1 通用约定

- 前缀：`/api/v1`，RESTful 风格，与管理端写操作走 JWT，公开读操作免认证。
- 统一响应结构（沿用现有 `response.Response`）：

```json
{
  "code": 0,
  "data": {},
  "msg": "ok"
}
```

- 分页请求：`page`（从 1 起）、`page_size`（默认 10，最大 100）；分页响应沿用现有 `dto.Pager`。
- 时间字段统一返回 Unix 时间戳（与现有 `created_on` 等保持一致）。
- 错误码通过 `internal/errcode` 管理，新增：文章不存在、分类已存在、评论不存在、文件类型不允许、文件超限 等。
- API 文档：swag 注解生成至 `doc/swagger/`，运行时访问 `http://{host}/swagger/index.html`（数据源 `/swagger/doc.json`）；重新生成命令：`swag init -g cmd/api/main.go -o doc/swagger -d ./ --parseInternal && rm -f doc/swagger/swagger.yaml`（不需要 yaml 格式，生成后删除）。

### 6.2 认证模块

| 方法 | 路径 | 认证 | 说明 |
| --- | --- | --- | --- |
| POST | /v1/user/login | 否 | 登录，返回 JWT（现有） |
| POST | /v1/user/logout | 是 | 登出（Redis 黑名单方式，见 8.3） |
| GET | /v1/user/info | 是 | 当前用户信息 |
| PATCH | /v1/user/password | 是 | 修改密码（旧密码校验 + 全端失效） |

### 6.3 文章模块

| 方法 | 路径 | 认证 | 说明 |
| --- | --- | --- | --- |
| GET | /v1/article | 否 | 文章列表（公开，仅已发布；支持 `category_id`/`tag_id`/`keyword`/`page`/`page_size`） |
| GET | /v1/article/:id | 否 | 文章详情（公开），Redis 计浏览量 |
| GET | /v1/article/slug/:slug | 否 | 按 slug 查详情 |
| GET | /v1/article/archive | 否 | 归档（按年月分组） |
| POST | /v1/article | 是 | 新建文章（支持草稿） |
| PATCH | /v1/article/:id | 是 | 编辑文章 |
| DELETE | /v1/article/:id | 是 | 删除文章（软删） |
| PATCH | /v1/article/:id/state | 是 | 发布/下架 |
| PATCH | /v1/article/:id/top | 是 | 置顶/取消置顶 |
| GET | /v1/admin/article | 是 | 管理端列表（含草稿，多条件筛选） |
| GET | /rss.xml | 否 | RSS 2.0 订阅（全部已发布文章） |
| GET | /sitemap.xml | 否 | Sitemap 站点地图（供爬虫抓取） |

文章创建/编辑 DTO 示例：

```go
type ArticleAddReq struct {
    Title      string   `json:"title" binding:"required,max=100"`
    Desc       string   `json:"desc" binding:"max=255"`
    ContentMd  string   `json:"content_md" binding:"required"`
    Cover      string   `json:"cover" binding:"omitempty,url,max=255"`
    CategoryID uint32   `json:"category_id" binding:"omitempty,gte=0"`
    TagIDs     []uint32 `json:"tag_ids" binding:"omitempty,max=10,dive,gte=1"`
    Slug       string   `json:"slug" binding:"omitempty,max=150"`
    IsDraft    bool     `json:"is_draft"`
}
```

### 6.4 分类 / 标签模块

| 方法 | 路径 | 认证 | 说明 |
| --- | --- | --- | --- |
| GET | /v1/category | 否 | 分类列表（含文章数） |
| POST / PATCH / DELETE | /v1/category(/:id) | 是 | 分类管理 |
| GET | /v1/tag | 否 | 标签列表（含文章数） |
| POST / PATCH / DELETE | /v1/tag(/:id) | 是 | 标签管理（现有，增加删除前置校验：存在关联文章则拒绝） |

### 6.5 评论模块

| 方法 | 路径 | 认证 | 说明 |
| --- | --- | --- | --- |
| GET | /v1/article/:id/comment | 否 | 评论列表（树形，两级） |
| POST | /v1/article/:id/comment | 否 | 发表评论（IP 限频，见 8.4） |
| GET | /v1/admin/comment | 是 | 管理端评论列表（按 state 筛选） |
| PATCH | /v1/admin/comment/:id/state | 是 | 审核（通过/拒绝） |
| DELETE | /v1/admin/comment/:id | 是 | 删除评论 |

### 6.6 文件 / 友链 / 配置模块

| 方法 | 路径 | 认证 | 说明 |
| --- | --- | --- | --- |
| POST | /v1/upload/image | 是 | 上传图片，返回 URL；`storage.Driver` 抽象，支持 local/OSS |
| GET | /v1/link | 否 | 友链列表（前台） |
| POST/PATCH/DELETE | /v1/link(/:id) | 是 | 友链管理 |
| GET | /v1/common/option | 否 | 站点配置（现有，改读 blog_option + 缓存） |
| PUT | /v1/admin/option | 是 | 批量保存站点配置 |

### 6.7 错误码规划（新增部分）

| 错误码 | 含义 |
| --- | --- |
| 1xxxx | 参数/校验类（沿用现有 ErrInvalidParams） |
| 20001 | 文章不存在 |
| 20002 | 分类不存在 / 分类名已存在 |
| 20003 | 标签不存在 / 标签已关联文章 |
| 20004 | 评论不存在 |
| 20005 | 文件类型不允许 / 超过大小限制 |
| 20006 | 旧密码错误 |

---

## 7. 认证与授权设计

- 登录成功签发 JWT（沿用现有 `core/jwt`），Claims 携带 `user_id`、`username`、`exp`。
- Token 有效期：Access Token 2h；配置于 `app.toml` 的 `jwt` 段。
- 登出：以 `jti` 为键写入 Redis 黑名单（TTL = 剩余有效期），JWT 中间件先查黑名单再验签，实现服务端可控的失效。
- 单用户模型：个人博客只有博主一个管理员，暂不做 RBAC；`blog_user` 预留 `role` 字段（tinyint，默认 1）以便后续扩展。
- 密码：bcrypt（cost 10），沿用现有 `core/crypto`。
- 修改密码后：清除该用户所有 Token（Redis 以 `user:{id}:token_version` 版本号控制，JWT 携带版本号比对）。

---

## 8. 缓存设计（Redis）

### 8.1 Key 规划

| Key | 类型 | TTL | 说明 |
| --- | --- | --- | --- |
| `blog:article:{id}` | string(JSON) | 30min + 随机抖动 | 文章详情缓存 |
| `blog:article:list:{query_hash}` | string(JSON) | 5min | 列表页缓存 |
| `blog:view:count:{date}` | hash(article_id→cnt) | 7d | 当日浏览量缓冲 |
| `blog:option` | string(JSON) | 10min | 站点配置 |
| `blog:archive` | string(JSON) | 1h | 归档数据 |
| `blog:jwt:blacklist:{jti}` | string | token 剩余期 | 登出黑名单 |
| `blog:user:ver:{id}` | string | 7d | 密码版本号 |

### 8.2 缓存策略（Cache-Aside）

- 读：先查 Redis，未命中查 DB 后回填；回填使用 `SETNX` 防击穿，空值缓存 60s 防穿透。
- 写：文章更新/删除/发布时，删除对应详情缓存与列表缓存（按 `blog:article:list:*` 模式或版本号失效），由 DB 为准延迟双删兜底。
- 统计类（站点总文章数/总浏览量）：定时任务每 5min 重算。

### 8.3 浏览量计数

- 详情接口内 `HINCRBY blog:view:count:{date} {article_id} 1`，纯内存自增，不打 DB。
- 定时任务（每 10min）将当日 hash 汇总 `UPDATE blog_article SET view_count = view_count + ? WHERE id = ?`，批量落库后删除已同步字段。

### 8.4 评论防刷

- `INCR blog:comment:limit:{ip}`，TTL 60s，超过阈值（如 5 次/分钟）返回错误码。
- 评论默认 `state=0` 待审核，管理端审核通过后前台可见（可配置免审）。

---

## 9. 日志与可观测性

- 沿用现有：logrus + 按天切分 + RequestID 中间件串联请求日志 + 敏感字段脱敏。
- 慢 SQL：现有 `gorm_log.go` 输出慢查询日志，阈值配置化（默认 200ms）。
- 健康检查：现有 `/health` 接口扩展为同时探测 DB（`SELECT 1`）与 Redis（`PING`），供 Nginx/K8s 探活。
- 指标（后续扩展）：`expvar` 或 prometheus client 输出 QPS、时延、错误率。
- 链路追踪：SkyWalking 已集成，按需开启。

---

## 10. 部署方案

### 10.1 单机部署（推荐，个人博客够用）

```
Nginx (443/80, gzip, 静态资源, 反代 :8080)
  └─ gin-web (Docker 容器, 单副本)
       ├─ MySQL 8.0 (容器/宿主机)
       └─ Redis 7 (容器/宿主机)
```

- 配置通过 `app.toml` 挂载进容器，密钥不入库不入镜像。
- 图片存储：首选对象存储（防宿主机迁移丢数据），本地存储则挂载 volume 并由 Nginx 直接托管静态目录。
- 备份：MySQL 每日 mysqldump 定时备份 + 保留 7 天；对象存储开启版本控制。

### 10.2 CI（复用现有 .github/workflows/ci.yml）

1. `go vet && go test ./...`（含 sqlmock 单测）。
2. `docker build` 打镜像，推送镜像仓库。
3. 服务器 `docker compose pull && up -d`，健康检查通过后切流量。

---

## 11. 测试策略

| 层级 | 工具 | 覆盖点 |
| --- | --- | --- |
| DAO 单测 | go-sqlmock（现有模式） | 条件构造、分页、软删过滤 |
| Service 单测 | mock DAO 接口 | 业务分支：发布/草稿、标签多对多事务、缓存读写与失效 |
| Handler 集成 | httptest | 参数校验（binding 标签）、JWT 拦截、响应结构 |
| 压测 | wrk/vegeta | 详情页缓存命中场景 QPS 与时延 |

---

## 12. 实施计划（建议顺序）

1. **阶段一（基础扩展）**：按最终 schema 全新建库（重写 `doc/sql/db.sql`，含默认管理员）；`blog_user` 接入；文章 CRUD 补全与草稿/发布流。
2. **阶段二（组织与展示）**：分类管理、标签多对多改造、归档、搜索（LIKE 版）。
3. **阶段三（互动与运营）**：评论（含审核、限频）、图片上传、友链、站点配置。
4. **阶段四（优化与扩展）**：Redis 缓存全面接入、浏览量落库定时任务、RSS/Sitemap、Prometheus 指标。

每阶段完成后更新本设计文档状态，并在 `doc/sql/db.sql` 同步维护全量建表语句。

---

## 13. 风险与后续扩展

| 项 | 风险/说明 | 应对 |
| --- | --- | --- |
| Markdown XSS | 存原文 + 前端渲染存在 XSS 面 | 预渲染 HTML 时用白名单过滤（bluemonday）；评论纯文本转义 |
| 搜索 LIKE 性能 | 数据量大后全表扫描 | 文章量 > 5w 再引入 Meilisearch/ES |
| Redis 故障 | 缓存层不可用拖垮接口 | 缓存读写全部超时降级直查 DB，不抛错 |
| 单点 | 单机部署不可用 | 已有每日备份 + 镜像可快速重建；量级不足暂不做高可用 |
