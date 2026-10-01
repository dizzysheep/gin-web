-- =================================================---------------
-- 个人博客最终 schema（详见 doc/blog-backend-design.md）
-- 全新建库脚本，直接按顺序执行即可
-- ================================================================

-- 博客账号
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

-- 文章分类
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

-- 文章标签
CREATE TABLE `blog_tag`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(100) DEFAULT '' COMMENT '标签名称',
    `created_on`  int unsigned DEFAULT '0' COMMENT '创建时间',
    `created_by`  varchar(100) DEFAULT '' COMMENT '创建人',
    `modified_on` int unsigned DEFAULT '0' COMMENT '修改时间',
    `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
    `deleted_on`  int unsigned DEFAULT '0',
    `state`       tinyint unsigned DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`, `deleted_on`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='文章标签管理';

-- 文章
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

-- 文章-标签关联（多对多）
CREATE TABLE `blog_article_tag`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT,
    `article_id` int unsigned NOT NULL COMMENT '文章ID',
    `tag_id`     int unsigned NOT NULL COMMENT '标签ID',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_article_tag` (`article_id`, `tag_id`),
    KEY `idx_tag` (`tag_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='文章标签关联';

-- 文章评论
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
    `created_on`  int unsigned DEFAULT '0' COMMENT '创建时间',
    `created_by`  varchar(100) DEFAULT '' COMMENT '创建人',
    `modified_on` int unsigned DEFAULT '0' COMMENT '修改时间',
    `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
    `deleted_on`  int unsigned DEFAULT '0',
    `state`       tinyint unsigned DEFAULT '1' COMMENT '0待审核 1通过 2拒绝',
    PRIMARY KEY (`id`),
    KEY `idx_article` (`article_id`, `state`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='文章评论';

-- 友情链接
CREATE TABLE `blog_link`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(100) DEFAULT '' COMMENT '站点名称',
    `url`         varchar(255) DEFAULT '' COMMENT '站点URL',
    `logo`        varchar(255) DEFAULT '' COMMENT '站点LOGO',
    `description` varchar(255) DEFAULT '' COMMENT '描述',
    `sort`        int          DEFAULT '0' COMMENT '排序',
    `created_on`  int unsigned DEFAULT '0' COMMENT '创建时间',
    `created_by`  varchar(100) DEFAULT '' COMMENT '创建人',
    `modified_on` int unsigned DEFAULT '0' COMMENT '修改时间',
    `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
    `deleted_on`  int unsigned DEFAULT '0',
    `state`       tinyint unsigned DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='友情链接';

-- 站点配置KV
CREATE TABLE `blog_option`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `key`         varchar(100) NOT NULL COMMENT '配置键',
    `value`       text COMMENT '配置值(JSON)',
    `remark`      varchar(255) DEFAULT '' COMMENT '说明',
    `modified_on` int unsigned DEFAULT '0' COMMENT '修改时间',
    `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
    `deleted_on`  int unsigned DEFAULT '0',
    `state`       tinyint unsigned DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_key` (`key`, `deleted_on`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='站点配置';

-- ================================================================
-- 初始数据
-- ================================================================

-- 密码为 123456 的 bcrypt 哈希，登录时由 bcrypt.CompareHashAndPassword 校验
INSERT INTO `blog_user`(`username`, `password`, `nickname`)
VALUES ('admin', '$2a$10$p2C7meW8dLv833ASNs3gveTArVo5GK/isvAMc2hHwUCojYie9RxIm', 'admin');

INSERT INTO `blog_option`(`key`, `value`, `remark`)
VALUES ('site_title', '"gin-web 个人博客"', '站点标题'),
       ('site_desc', '"基于 gin-web 的个人博客"', '站点描述'),
       ('site_icp', '""', '备案号'),
       ('comment_audit', '"1"', '评论是否需要审核 0否1是'),
       ('page_size', '"10"', '默认分页大小');
