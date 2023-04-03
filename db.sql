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
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='文章标签管理';

CREATE TABLE `blog_auth`
(
    `id`       int unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(50) DEFAULT '' COMMENT '账号',
    `password` varchar(50) DEFAULT '' COMMENT '密码',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;

CREATE TABLE `blog_article`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `tag_id`      int unsigned DEFAULT '0' COMMENT '标签ID',
    `title`       varchar(100) DEFAULT '' COMMENT '文章标题',
    `desc`        varchar(255) DEFAULT '' COMMENT '简述',
    `content`     text,
    `created_on`  int          DEFAULT NULL,
    `created_by`  varchar(100) DEFAULT '' COMMENT '创建人',
    `modified_on` int unsigned DEFAULT '0' COMMENT '修改时间',
    `modified_by` varchar(255) DEFAULT '' COMMENT '修改人',
    `deleted_on`  int unsigned DEFAULT '0',
    `state`       tinyint unsigned DEFAULT '1' COMMENT '状态 0为禁用1为启用',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;

INSERT INTO `blog2`.`blog_auth`(`username`, `password`) VALUES ('admin', '123456');