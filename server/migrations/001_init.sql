-- Leslie Blog 初始化数据库结构
-- 含：users / categories / posts / post_translations / post_tags /
--     projects / tech_stack / media / comments / settings
-- i18n：文章多语言通过 post_translations 实现（zh-CN / en-US）

SET NAMES utf8mb4;

-- 后台用户
CREATE TABLE IF NOT EXISTS `users` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username`   VARCHAR(64)  NOT NULL,
  `password`   VARCHAR(128) NOT NULL,                -- bcrypt hash
  `nickname`   VARCHAR(64)  NOT NULL DEFAULT '',
  `email`      VARCHAR(128) NOT NULL DEFAULT '',
  `avatar`     VARCHAR(512) NOT NULL DEFAULT '',
  `role`       VARCHAR(32)  NOT NULL DEFAULT 'admin',
  `created_at` DATETIME     NOT NULL,
  `updated_at` DATETIME     NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分类
CREATE TABLE IF NOT EXISTS `categories` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(64)  NOT NULL,
  `slug`        VARCHAR(64)  NOT NULL,
  `description` VARCHAR(255) NOT NULL DEFAULT '',
  `parent_id`   BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `sort_order`  INT          NOT NULL DEFAULT 0,
  `created_at`  DATETIME     NOT NULL,
  `updated_at`  DATETIME     NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_categories_slug` (`slug`),
  KEY `idx_categories_parent` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文章（语言无关字段）
CREATE TABLE IF NOT EXISTS `posts` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `slug`            VARCHAR(128) NOT NULL,                       -- i18n 共享 slug
  `cover`           VARCHAR(512) NOT NULL DEFAULT '',
  `category_id`     BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `status`          TINYINT      NOT NULL DEFAULT 0,             -- 0 草稿 1 已发布 2 已下线
  `is_top`          TINYINT(1)   NOT NULL DEFAULT 0,
  `view_count`      INT UNSIGNED NOT NULL DEFAULT 0,
  `comment_count`   INT UNSIGNED NOT NULL DEFAULT 0,
  `seo_title`       VARCHAR(255) NOT NULL DEFAULT '',
  `seo_description` VARCHAR(500) NOT NULL DEFAULT '',
  `seo_keywords`    VARCHAR(255) NOT NULL DEFAULT '',
  `published_at`    DATETIME     NULL,
  `created_at`      DATETIME     NOT NULL,
  `updated_at`      DATETIME     NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_posts_slug` (`slug`),
  KEY `idx_posts_category` (`category_id`),
  KEY `idx_posts_status` (`status`),
  KEY `idx_posts_top` (`is_top`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文章多语言翻译（标题 / 摘要 / 正文）
CREATE TABLE IF NOT EXISTS `post_translations` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `post_id`    BIGINT UNSIGNED NOT NULL,
  `locale`     VARCHAR(16) NOT NULL,                            -- zh-CN / en-US
  `title`      VARCHAR(255) NOT NULL,
  `summary`    VARCHAR(500) NOT NULL DEFAULT '',
  `content`    LONGTEXT     NOT NULL,
  `created_at` DATETIME    NOT NULL,
  `updated_at` DATETIME    NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_post_translations` (`post_id`, `locale`),
  KEY `idx_pt_locale` (`locale`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文章标签
CREATE TABLE IF NOT EXISTS `post_tags` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `post_id`    BIGINT UNSIGNED NOT NULL,
  `name`       VARCHAR(64) NOT NULL,
  `created_at` DATETIME    NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_post_tags` (`post_id`, `name`),
  KEY `idx_pt_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 项目
CREATE TABLE IF NOT EXISTS `projects` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(128) NOT NULL,
  `description` TEXT         NOT NULL,
  `cover`       VARCHAR(512) NOT NULL DEFAULT '',
  `demo_url`    VARCHAR(512) NOT NULL DEFAULT '',
  `repo_url`    VARCHAR(512) NOT NULL DEFAULT '',
  `tech_stack`  JSON         NULL,                            -- ["Go","Vue",...]
  `status`      TINYINT      NOT NULL DEFAULT 0,               -- 0 进行中 1 已完成 2 已归档
  `sort_order`  INT          NOT NULL DEFAULT 0,
  `created_at`  DATETIME     NOT NULL,
  `updated_at`  DATETIME     NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 技术栈
CREATE TABLE IF NOT EXISTS `tech_stack` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(64)  NOT NULL,
  `icon`        VARCHAR(512) NOT NULL DEFAULT '',
  `category`    VARCHAR(64)  NOT NULL DEFAULT '',
  `level`       INT          NOT NULL DEFAULT 0,             -- 0-100
  `description` VARCHAR(255) NOT NULL DEFAULT '',
  `sort_order`  INT          NOT NULL DEFAULT 0,
  `created_at`  DATETIME     NOT NULL,
  `updated_at`  DATETIME     NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 媒体
CREATE TABLE IF NOT EXISTS `media` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`       VARCHAR(255) NOT NULL,
  `url`        VARCHAR(512) NOT NULL,
  `type`       VARCHAR(32)  NOT NULL,                          -- image / video / file
  `size`       BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `mime_type`  VARCHAR(128) NOT NULL DEFAULT '',
  `created_at` DATETIME     NOT NULL,
  `updated_at` DATETIME     NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_media_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 评论
CREATE TABLE IF NOT EXISTS `comments` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `post_id`    BIGINT UNSIGNED NOT NULL,
  `author`     VARCHAR(64)  NOT NULL,
  `email`      VARCHAR(128) NOT NULL DEFAULT '',
  `avatar`     VARCHAR(512) NOT NULL DEFAULT '',
  `content`    TEXT         NOT NULL,
  `parent_id`  BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `status`     TINYINT      NOT NULL DEFAULT 0,                 -- 0 待审 1 通过 2 拒绝
  `ip`         VARCHAR(64)  NOT NULL DEFAULT '',
  `created_at` DATETIME     NOT NULL,
  `updated_at` DATETIME     NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_comments_post` (`post_id`),
  KEY `idx_comments_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 站点设置（单行）
CREATE TABLE IF NOT EXISTS `settings` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `site_name`      VARCHAR(128) NOT NULL DEFAULT 'Leslie Blog',
  `logo`           VARCHAR(512) NOT NULL DEFAULT '',
  `description`    VARCHAR(500) NOT NULL DEFAULT '',
  `keywords`       VARCHAR(255) NOT NULL DEFAULT '',
  `author`         VARCHAR(64)  NOT NULL DEFAULT '',
  `icp`            VARCHAR(128) NOT NULL DEFAULT '',
  `social_github`  VARCHAR(255) NOT NULL DEFAULT '',
  `social_email`   VARCHAR(255) NOT NULL DEFAULT '',
  `social_twitter` VARCHAR(255) NOT NULL DEFAULT '',
  `updated_at`     DATETIME     NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 默认设置行
INSERT INTO `settings` (`id`, `site_name`, `updated_at`)
VALUES (1, 'Leslie Blog', NOW())
ON DUPLICATE KEY UPDATE `id` = `id`;
