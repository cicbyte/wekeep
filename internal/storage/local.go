// Package storage 提供对象存储抽象层
package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// LocalProvider 本地文件系统存储提供者
type LocalProvider struct {
	// basePath 本地存储根目录
	basePath string

	// baseURL 文件访问基础URL
	baseURL string
}

// NewLocalProvider 创建本地文件系统存储提供者
func NewLocalProvider(cfg *LocalConfig) (*LocalProvider, error) {
	if cfg == nil {
		return nil, ErrConfigMissing
	}

	// 设置默认值
	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "./uploads"
	}

	// 确保路径是绝对路径
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	// 设置默认的 base URL
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "/uploads"
	}

	// 移除末尾的斜杠
	baseURL = strings.TrimSuffix(baseURL, "/")

	p := &LocalProvider{
		basePath: absPath,
		baseURL:  baseURL,
	}

	g.Log().Infof(context.Background(), "Local storage initialized: path=%s, url=%s", absPath, baseURL)

	return p, nil
}

// Upload 上传文件到本地文件系统
func (p *LocalProvider) Upload(ctx context.Context, reader io.Reader, objectKey string, contentType string) (storageURL string, err error) {
	// 构建完整文件路径
	fullPath := filepath.Join(p.basePath, objectKey)

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		g.Log().Errorf(ctx, "Failed to create directory: %v", err)
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		g.Log().Errorf(ctx, "Failed to create file: %v", err)
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 写入文件内容
	_, err = io.Copy(file, reader)
	if err != nil {
		g.Log().Errorf(ctx, "Failed to write file: %v", err)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// 返回相对路径（objectKey），由图片代理接口统一访问
	g.Log().Infof(ctx, "File uploaded successfully: %s", objectKey)
	return objectKey, nil
}

// Download 从本地文件系统下载文件
func (p *LocalProvider) Download(ctx context.Context, objectKey string) ([]byte, error) {
	fullPath := filepath.Join(p.basePath, objectKey)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		g.Log().Errorf(ctx, "Failed to read file: %v", err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// Copy 复制文件
func (p *LocalProvider) Copy(ctx context.Context, srcKey string, dstKey string) error {
	srcPath := filepath.Join(p.basePath, srcKey)
	dstPath := filepath.Join(p.basePath, dstKey)

	// 确保目标目录存在
	dir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 读取源文件
	srcData, err := os.ReadFile(srcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// 写入目标文件
	if err := os.WriteFile(dstPath, srcData, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}

	g.Log().Infof(ctx, "File copied: %s -> %s", srcKey, dstKey)
	return nil
}

// Delete 从本地文件系统删除文件
func (p *LocalProvider) Delete(ctx context.Context, objectKey string) error {
	fullPath := filepath.Join(p.basePath, objectKey)

	err := os.Remove(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，视为成功
			g.Log().Warningf(ctx, "File not found, already deleted: %s", objectKey)
			return nil
		}
		g.Log().Errorf(ctx, "Failed to delete file: %v", err)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	g.Log().Infof(ctx, "File deleted: %s", objectKey)
	return nil
}

// GetURL 获取文件访问URL
func (p *LocalProvider) GetURL(objectKey string) string {
	return p.baseURL + "/" + objectKey
}

// HealthCheck 健康检查
func (p *LocalProvider) HealthCheck(ctx context.Context) error {
	// 检查目录是否存在且可写
	info, err := os.Stat(p.basePath)
	if err != nil {
		return fmt.Errorf("storage directory not accessible: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("storage path is not a directory")
	}

	// 尝试创建临时文件测试写入权限
	testFile := filepath.Join(p.basePath, ".healthcheck")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("storage directory is not writable: %w", err)
	}
	f.Close()
	os.Remove(testFile)

	return nil
}

// GetBasePath 获取存储根目录
func (p *LocalProvider) GetBasePath() string {
	return p.basePath
}

// ServeHTTP 实现 http.Handler 接口，用于提供静态文件服务
func (p *LocalProvider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取请求路径，移除前缀
	relPath := strings.TrimPrefix(r.URL.Path, p.baseURL)
	relPath = strings.TrimPrefix(relPath, "/")

	// 构建完整文件路径
	fullPath := filepath.Join(p.basePath, relPath)

	// 安全检查：确保路径在存储目录内
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(absPath, p.basePath) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// 检查文件是否存在
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 如果是目录，返回 404
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	// 设置 Content-Type
	ext := strings.ToLower(filepath.Ext(absPath))
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".webp":
		w.Header().Set("Content-Type", "image/webp")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// 发送文件
	http.ServeFile(w, r, absPath)
}
