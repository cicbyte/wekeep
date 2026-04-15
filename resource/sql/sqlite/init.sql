-- WeKeep SQLite 初始化脚本
-- 包含所有表的 CREATE 语句，不含示例数据
-- 执行: sqlite3 wekeep.db < init.sql

-- ============================================================
-- 1. categories 分类表
-- ============================================================
CREATE TABLE IF NOT EXISTS `categories` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` TEXT NOT NULL UNIQUE,
    `description` TEXT NOT NULL,
    `icon` TEXT,
    `sort` INTEGER NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS `idx_categories_sort` ON `categories` (`sort`);
CREATE INDEX IF NOT EXISTS `idx_categories_created_at` ON `categories` (`created_at`);

CREATE TRIGGER IF NOT EXISTS `tr_categories_updated_at`
AFTER UPDATE ON `categories`
FOR EACH ROW
BEGIN
    UPDATE `categories` SET `updated_at` = CURRENT_TIMESTAMP WHERE `id` = NEW.`id`;
END;

-- ============================================================
-- 2. authors 作者表
-- ============================================================
CREATE TABLE IF NOT EXISTS `authors` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` TEXT NOT NULL,
    `normalized_name` TEXT NOT NULL UNIQUE,
    `avatar` TEXT,
    `bio` TEXT,
    `website` TEXT,
    `article_count` INTEGER NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS `idx_authors_article_count` ON `authors` (`article_count` DESC);

CREATE TRIGGER IF NOT EXISTS `tr_authors_updated_at`
AFTER UPDATE ON `authors`
FOR EACH ROW
BEGIN
    UPDATE `authors` SET `updated_at` = CURRENT_TIMESTAMP WHERE `id` = NEW.`id`;
END;

-- ============================================================
-- 3. articles 文章表
-- ============================================================
CREATE TABLE IF NOT EXISTS `articles` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `title` TEXT NOT NULL,
    `url` TEXT,
    `summary` TEXT,
    `content` TEXT,
    `tags` TEXT,
    `date_added` INTEGER,
    `author_id` INTEGER,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS `idx_articles_date_added` ON `articles` (`date_added`);
CREATE INDEX IF NOT EXISTS `idx_articles_created_at` ON `articles` (`created_at`);
CREATE INDEX IF NOT EXISTS `idx_articles_author_id` ON `articles` (`author_id`);

CREATE TRIGGER IF NOT EXISTS `tr_articles_updated_at`
AFTER UPDATE ON `articles`
FOR EACH ROW
BEGIN
    UPDATE `articles` SET `updated_at` = CURRENT_TIMESTAMP WHERE `id` = NEW.`id`;
END;

-- ============================================================
-- 4. images 图片元信息表
-- ============================================================
CREATE TABLE IF NOT EXISTS `images` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `original_url` TEXT NOT NULL,
    `original_url_hash` TEXT NOT NULL UNIQUE,
    `storage_path` TEXT NOT NULL,
    `storage_url` TEXT,
    `file_size` INTEGER DEFAULT 0,
    `mime_type` TEXT,
    `ref_count` INTEGER NOT NULL DEFAULT 0,
    `download_status` INTEGER NOT NULL DEFAULT 0,
    `error_message` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS `idx_images_ref_count` ON `images` (`ref_count`);
CREATE INDEX IF NOT EXISTS `idx_images_download_status` ON `images` (`download_status`);

CREATE TRIGGER IF NOT EXISTS `tr_images_updated_at`
AFTER UPDATE ON `images`
FOR EACH ROW
BEGIN
    UPDATE `images` SET `updated_at` = CURRENT_TIMESTAMP WHERE `id` = NEW.`id`;
END;

-- ============================================================
-- 5. article_images 文章-图片关联表
-- ============================================================
CREATE TABLE IF NOT EXISTS `article_images` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `article_id` INTEGER NOT NULL,
    `image_id` INTEGER NOT NULL,
    `position` INTEGER DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (`article_id`, `image_id`),
    FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`image_id`) REFERENCES `images` (`id`) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS `idx_article_images_article_id` ON `article_images` (`article_id`);
CREATE INDEX IF NOT EXISTS `idx_article_images_image_id` ON `article_images` (`image_id`);
