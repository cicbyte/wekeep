// Package storage 提供对象存储抽象层
package storage

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MigrationStatus 迁移状态
type MigrationStatus struct {
	Running      bool         `json:"running"`       // 是否正在运行
	Source       string       `json:"source"`        // 源存储类型
	Target       string       `json:"target"`        // 目标存储类型
	Total        int          `json:"total"`         // 总文件数
	Completed    int          `json:"completed"`     // 已完成数
	Failed       int          `json:"failed"`        // 失败数
	CurrentFile  string       `json:"currentFile"`   // 当前处理的文件
	StartTime    *gtime.Time  `json:"startTime"`     // 开始时间
	EndTime      *gtime.Time  `json:"endTime"`       // 结束时间
	Error        string       `json:"error"`         // 错误信息
}

// MigrationConfig 迁移配置
type MigrationConfig struct {
	// SourceStorage 源存储类型 (local, rustfs)，留空则使用当前配置
	SourceStorage string `json:"sourceStorage"`

	// TargetStorage 目标存储类型 (local, rustfs)
	TargetStorage string `json:"targetStorage"`

	// LocalConfig 本地存储配置（目标为 local 时使用，或源为 local 时使用）
	LocalConfig *LocalConfig `json:"localConfig"`

	// RustFSConfig RustFS 配置（目标为 rustfs 时使用，或源为 rustfs 时使用）
	RustFSConfig *RustFSConfig `json:"rustfsConfig"`

	// UpdateMarkdown 是否更新 Markdown 内容中的图片引用
	UpdateMarkdown bool `json:"updateMarkdown"`
}

// 全局迁移状态
var (
	migrationStatus = &MigrationStatus{Running: false}
	migrationMutex  sync.Mutex
)

// GetMigrationStatus 获取迁移状态
func GetMigrationStatus() *MigrationStatus {
	migrationMutex.Lock()
	defer migrationMutex.Unlock()
	return migrationStatus
}

// MigrateStorage 执行存储迁移
func MigrateStorage(ctx context.Context, cfg *MigrationConfig) error {
	migrationMutex.Lock()

	// 检查是否已有迁移任务在运行
	if migrationStatus.Running {
		migrationMutex.Unlock()
		return fmt.Errorf("migration already in progress")
	}

	// 确定源存储类型
	sourceStorage := cfg.SourceStorage
	if sourceStorage == "" {
		sourceStorage = GetCurrentStorageType()
	}

	// 初始化状态
	migrationStatus = &MigrationStatus{
		Running:   true,
		Source:    sourceStorage,
		Target:    cfg.TargetStorage,
		StartTime: gtime.Now(),
	}
	migrationMutex.Unlock()

	// 异步执行迁移（使用独立的 context，不受 HTTP 请求生命周期影响）
	go func() {
		defer func() {
			migrationMutex.Lock()
			migrationStatus.Running = false
			migrationStatus.EndTime = gtime.Now()
			migrationMutex.Unlock()
		}()

		// 使用 background context，不受 HTTP 请求取消影响
		migrationCtx := context.Background()
		err := doMigration(migrationCtx, cfg)
		if err != nil {
			migrationMutex.Lock()
			migrationStatus.Error = err.Error()
			migrationMutex.Unlock()
			g.Log().Errorf(migrationCtx, "Migration failed: %v", err)
		}
	}()

	return nil
}

// doMigration 执行实际的迁移逻辑
func doMigration(ctx context.Context, cfg *MigrationConfig) error {
	g.Log().Infof(ctx, "Starting storage migration: %s -> %s", migrationStatus.Source, migrationStatus.Target)

	// 1. 初始化源存储提供者
	sourceProvider, err := createSourceProvider(cfg)
	if err != nil {
		return fmt.Errorf("failed to create source provider: %w", err)
	}

	// 2. 初始化目标存储提供者
	targetProvider, err := createTargetProvider(cfg)
	if err != nil {
		return fmt.Errorf("failed to create target provider: %w", err)
	}

	// 2. 健康检查
	if err := targetProvider.HealthCheck(ctx); err != nil {
		return fmt.Errorf("target storage health check failed: %w", err)
	}

	// 3. 源存储健康检查
	if err := sourceProvider.HealthCheck(ctx); err != nil {
		return fmt.Errorf("source storage health check failed: %w", err)
	}

	// 4. 查询所有图片记录
	var images []struct {
		Id          int    `orm:"id"`
		ArticleId   int    `orm:"article_id"`
		StoragePath string `orm:"storage_path"`
		StorageUrl  string `orm:"storage_url"`
	}

	err = g.DB().Model("images").
		Where("storage_path != ?", "").
		Where("storage_url != ?", "").
		Scan(&images)
	if err != nil {
		return fmt.Errorf("failed to query images: %w", err)
	}

	migrationMutex.Lock()
	migrationStatus.Total = len(images)
	migrationMutex.Unlock()

	g.Log().Infof(ctx, "Found %d images to migrate", len(images))

	if len(images) == 0 {
		g.Log().Info(ctx, "No images to migrate")
		return nil
	}

	// 5. 迁移图片
	urlMapping := make(map[string]string) // oldURL -> newURL

	for _, img := range images {
			// 统一路径分隔符为正斜杠（Windows 反斜杠会导致 S3 路径异常）
		normalizedPath := strings.ReplaceAll(img.StoragePath, "\\", "/")

		migrationMutex.Lock()
		migrationStatus.CurrentFile = normalizedPath
		migrationMutex.Unlock()

		// 从源存储下载
		data, err := sourceProvider.Download(ctx, normalizedPath)
		if err != nil {
			g.Log().Warningf(ctx, "Failed to download %s: %v", normalizedPath, err)
			migrationMutex.Lock()
			migrationStatus.Failed++
			migrationMutex.Unlock()
			continue
		}

		// 上传到目标存储
		contentType := getContentType(normalizedPath)
		newURL, err := targetProvider.Upload(ctx, strings.NewReader(string(data)), normalizedPath, contentType)
		if err != nil {
			g.Log().Warningf(ctx, "Failed to upload %s: %v", normalizedPath, err)
			migrationMutex.Lock()
			migrationStatus.Failed++
			migrationMutex.Unlock()
			continue
		}

		// 记录 URL 映射
		urlMapping[img.StorageUrl] = newURL

		// 更新数据库
		_, err = g.DB().Model("images").
			Where("id", img.Id).
			Update(g.Map{
				"storage_path": normalizedPath,
				"storage_url":  newURL,
			})
		if err != nil {
			g.Log().Warningf(ctx, "Failed to update image record %d: %v", img.Id, err)
		}

		migrationMutex.Lock()
		migrationStatus.Completed++
		g.Log().Infof(ctx, "Migrated %d/%d: %s", migrationStatus.Completed, migrationStatus.Total, normalizedPath)
		migrationMutex.Unlock()

		// 避免过快请求
		time.Sleep(10 * time.Millisecond)
	}

	// 6. 更新 Markdown 内容中的图片引用
	if cfg.UpdateMarkdown && len(urlMapping) > 0 {
		g.Log().Info(ctx, "Updating markdown content...")
		if err := updateMarkdownContent(ctx, urlMapping); err != nil {
			g.Log().Warningf(ctx, "Failed to update markdown content: %v", err)
		}
	}

	g.Log().Infof(ctx, "Migration completed: %d success, %d failed", migrationStatus.Completed, migrationStatus.Failed)
	return nil
}

// createSourceProvider 创建源存储提供者
func createSourceProvider(cfg *MigrationConfig) (Provider, error) {
	sourceStorage := cfg.SourceStorage
	if sourceStorage == "" {
		sourceStorage = GetCurrentStorageType()
	}

	switch sourceStorage {
	case "local", "":
		localCfg := cfg.LocalConfig
		if localCfg == nil {
			// 使用默认配置或当前配置
			localCfg = &LocalConfig{
				BasePath: "./uploads",
				BaseURL:  "/uploads",
			}
		}
		return NewLocalProvider(localCfg)
	case "rustfs":
		if cfg.RustFSConfig == nil {
			return nil, fmt.Errorf("rustfs config is required for source storage")
		}
		// 记录配置信息（隐藏密码）
		g.Log().Infof(context.Background(), "Creating source RustFS provider: endpoint=%s, bucket=%s, username=%s",
			cfg.RustFSConfig.Endpoint, cfg.RustFSConfig.Bucket, cfg.RustFSConfig.Username)
		return NewRustFSProvider(cfg.RustFSConfig)
	default:
		return nil, fmt.Errorf("unsupported source storage: %s", sourceStorage)
	}
}

// createTargetProvider 创建目标存储提供者
func createTargetProvider(cfg *MigrationConfig) (Provider, error) {
	switch cfg.TargetStorage {
	case "local":
		localCfg := cfg.LocalConfig
		if localCfg == nil {
			localCfg = &LocalConfig{
				BasePath: "./uploads",
				BaseURL:  "/uploads",
			}
		}
		return NewLocalProvider(localCfg)
	case "rustfs":
		if cfg.RustFSConfig == nil {
			return nil, fmt.Errorf("rustfs config is required")
		}
		return NewRustFSProvider(cfg.RustFSConfig)
	default:
		return nil, fmt.Errorf("unsupported target storage: %s", cfg.TargetStorage)
	}
}

// BuildImageURLMapping 从数据库构建图片 URL 映射
// 返回 storage_url -> 当前实际 URL 的映射
func BuildImageURLMapping(ctx context.Context) (map[string]string, error) {
	// 查询所有图片记录
	var images []struct {
		Id          int    `orm:"id"`
		StoragePath string `orm:"storage_path"`
		StorageUrl  string `orm:"storage_url"`
	}

	err := g.DB().Model("images").
		Where("storage_path != ?", "").
		Where("storage_url != ?", "").
		Scan(&images)
	if err != nil {
		return nil, fmt.Errorf("failed to query images: %w", err)
	}

	urlMapping := make(map[string]string)

	// 获取当前存储提供者
	if !IsInitialized() {
		return nil, fmt.Errorf("storage provider not initialized")
	}

	provider := GetProvider()

	for _, img := range images {
		// 计算当前应该的 URL
		path := strings.ReplaceAll(img.StoragePath, "\\", "/")
		currentURL := provider.GetURL(path)
		// 映射：旧 URL -> 新 URL
		if img.StorageUrl != currentURL {
			urlMapping[img.StorageUrl] = currentURL
		}
	}

	return urlMapping, nil
}

// UpdateMarkdownContent 更新文章 Markdown 内容中的图片引用（导出函数）
func UpdateMarkdownContent(ctx context.Context, urlMapping map[string]string) (int, error) {
	if len(urlMapping) == 0 {
		return 0, nil
	}

	// 查询所有文章
	var articles []struct {
		Id      int    `orm:"id"`
		Content string `orm:"content"`
	}

	err := g.DB().Model("articles").
		Where("content IS NOT NULL").
		Where("content != ?", "").
		Scan(&articles)
	if err != nil {
		return 0, fmt.Errorf("failed to query articles: %w", err)
	}

	g.Log().Infof(ctx, "Checking %d articles for image URL updates", len(articles))

	updated := 0
	for _, article := range articles {
		newContent := article.Content
		changed := false

		// 替换所有旧的 URL
		for oldURL, newURL := range urlMapping {
			if strings.Contains(newContent, oldURL) {
				newContent = strings.ReplaceAll(newContent, oldURL, newURL)
				changed = true
			}
		}

		if changed {
			_, err := g.DB().Model("articles").
				Where("id", article.Id).
				Update(g.Map{
					"content": newContent,
				})
			if err != nil {
				g.Log().Warningf(ctx, "Failed to update article %d: %v", article.Id, err)
				continue
			}
			updated++
		}
	}

	g.Log().Infof(ctx, "Updated %d articles with new image URLs", updated)
	return updated, nil
}

// updateMarkdownContent 内部函数，调用导出版本
func updateMarkdownContent(ctx context.Context, urlMapping map[string]string) error {
	_, err := UpdateMarkdownContent(ctx, urlMapping)
	return err
}

// getContentType 根据文件扩展名获取 Content-Type
func getContentType(path string) string {
	ext := strings.ToLower(path)
	if idx := strings.LastIndex(ext, "."); idx >= 0 {
		ext = ext[idx:]
	}

	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

// GetCurrentStorageType 获取当前存储类型
func GetCurrentStorageType() string {
	// 从配置读取
	ctx := context.Background()
	storageType := g.Cfg().MustGet(ctx, "storage.type", "local").String()
	return storageType
}

// ValidateStorageConfig 验证存储配置
func ValidateStorageConfig(ctx context.Context, cfg *MigrationConfig) error {
	// 创建临时提供者进行验证
	provider, err := createTargetProvider(cfg)
	if err != nil {
		return err
	}

	// 健康检查
	return provider.HealthCheck(ctx)
}

// ExtractImageURLs 从 Markdown 内容中提取图片 URL
func ExtractImageURLs(content string) []string {
	// 匹配 ![alt](url) 格式
	re := regexp.MustCompile(`!\[.*?\]\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(content, -1)

	urls := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	// 匹配 <img src="url"> 格式
	re2 := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["']`)
	matches2 := re2.FindAllStringSubmatch(content, -1)
	for _, match := range matches2 {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	return urls
}

// ReadFileAsReader 辅助函数：将 []byte 转换为 io.Reader
func ReadFileAsReader(data []byte) io.Reader {
	return strings.NewReader(string(data))
}
