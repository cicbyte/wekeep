-- WeKeep MySQL 初始化脚本
-- 包含所有表的 CREATE 语句，不含示例数据
-- 执行: mysql -u root -p wekeep < init.sql

-- ============================================================
-- 1. categories 分类表
-- ============================================================
CREATE TABLE IF NOT EXISTS `categories` (
    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` VARCHAR(255) NOT NULL COMMENT '分类名称，唯一',
    `description` TEXT NOT NULL COMMENT '分类描述',
    `icon` VARCHAR(512) DEFAULT NULL COMMENT '分类图标URL或标识',
    `sort` INT(11) NOT NULL DEFAULT 0 COMMENT '排序，数字越大越靠前',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`),
    KEY `idx_sort` (`sort`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- ============================================================
-- 2. authors 作者表
-- ============================================================
CREATE TABLE IF NOT EXISTS `authors` (
    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` VARCHAR(255) NOT NULL COMMENT '作者名称',
    `normalized_name` VARCHAR(255) NOT NULL COMMENT '标准化名称（去重用）',
    `avatar` VARCHAR(512) DEFAULT NULL COMMENT '头像URL',
    `bio` TEXT DEFAULT NULL COMMENT '作者简介',
    `website` VARCHAR(512) DEFAULT NULL COMMENT '个人网站',
    `article_count` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章数量',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_normalized_name` (`normalized_name`),
    KEY `idx_article_count` (`article_count` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='作者表';

-- ============================================================
-- 3. articles 文章表
-- ============================================================
CREATE TABLE IF NOT EXISTS `articles` (
    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `title` VARCHAR(512) NOT NULL COMMENT '文章标题',
    `url` VARCHAR(2048) DEFAULT NULL COMMENT '原文链接',
    `summary` TEXT DEFAULT NULL COMMENT '文章摘要',
    `content` MEDIUMTEXT DEFAULT NULL COMMENT 'Markdown内容',
    `tags` JSON DEFAULT NULL COMMENT '标签数组（JSON格式）',
    `date_added` BIGINT DEFAULT NULL COMMENT '添加时间戳(毫秒，保留前端原始时间)',
    `author_id` INT(11) UNSIGNED DEFAULT NULL COMMENT '作者ID（外键）',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_date_added` (`date_added`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_title_author` (`title`(100), `author_id`),
    KEY `idx_author_id` (`author_id`),
    CONSTRAINT `fk_articles_author_id` FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章表';

-- ============================================================
-- 4. images 图片元信息表
-- ============================================================
CREATE TABLE IF NOT EXISTS `images` (
    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `original_url` VARCHAR(2048) NOT NULL COMMENT '原始图片URL',
    `original_url_hash` VARCHAR(64) NOT NULL COMMENT 'URL哈希(SHA256)',
    `storage_path` VARCHAR(512) NOT NULL COMMENT '存储路径',
    `storage_url` VARCHAR(1024) DEFAULT NULL COMMENT '访问URL',
    `file_size` INT(11) UNSIGNED DEFAULT 0 COMMENT '文件大小(字节)',
    `mime_type` VARCHAR(64) DEFAULT NULL COMMENT 'MIME类型',
    `ref_count` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '引用计数',
    `download_status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '下载状态: 0-待下载, 1-下载中, 2-下载成功, 3-下载失败',
    `error_message` VARCHAR(512) DEFAULT NULL COMMENT '错误信息',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_url_hash` (`original_url_hash`),
    KEY `idx_ref_count` (`ref_count`),
    KEY `idx_download_status` (`download_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图片元信息表';

-- ============================================================
-- 5. article_images 文章-图片关联表
-- ============================================================
CREATE TABLE IF NOT EXISTS `article_images` (
    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `article_id` INT(11) UNSIGNED NOT NULL COMMENT '文章ID',
    `image_id` INT(11) UNSIGNED NOT NULL COMMENT '图片ID',
    `position` INT(11) UNSIGNED DEFAULT 0 COMMENT '图片位置(文章中的顺序)',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_article_image` (`article_id`, `image_id`),
    KEY `idx_article_id` (`article_id`),
    KEY `idx_image_id` (`image_id`),
    CONSTRAINT `fk_article_images_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_article_images_image` FOREIGN KEY (`image_id`) REFERENCES `images` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章-图片关联表';
