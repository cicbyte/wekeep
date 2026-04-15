// Package storage 提供对象存储抽象层
package storage

import (
	"context"
	"io"

	"github.com/gogf/gf/v2/frame/g"
)

// Provider 存储提供者接口
type Provider interface {
	// Upload 上传文件
	// reader: 文件内容读取器
	// objectKey: 对象存储路径/键名
	// contentType: 文件MIME类型
	Upload(ctx context.Context, reader io.Reader, objectKey string, contentType string) (storageURL string, err error)

	// Download 下载文件
	// 返回文件内容和错误
	Download(ctx context.Context, objectKey string) ([]byte, error)

	// Copy 复制文件
	// srcKey: 源文件路径
	// dstKey: 目标文件路径
	Copy(ctx context.Context, srcKey string, dstKey string) error

	// Delete 删除文件
	Delete(ctx context.Context, objectKey string) error

	// GetURL 获取文件访问URL
	GetURL(objectKey string) string

	// HealthCheck 健康检查
	HealthCheck(ctx context.Context) error
}

// Config 存储配置
type Config struct {
	// Type 存储类型 (local, rustfs, s3, minio 等)
	Type string `yaml:"type"`

	// Local 本地文件系统配置
	Local *LocalConfig `yaml:"local"`

	// RustFS RustFS 配置
	RustFS *RustFSConfig `yaml:"rustfs"`

	// Image 图片相关配置
	Image *ImageConfig `yaml:"image"`
}

// LocalConfig 本地文件系统配置
type LocalConfig struct {
	// basePath 本地存储根目录
	BasePath string `yaml:"basePath"`

	// BaseURL 文件访问基础URL
	BaseURL string `yaml:"baseURL"`
}

// RustFSConfig RustFS 配置
type RustFSConfig struct {
	// Endpoint 服务端点
	Endpoint string `yaml:"endpoint"`

	// Bucket 存储桶名称
	Bucket string `yaml:"bucket"`

	// Username 用户名
	Username string `yaml:"username"`

	// Password 密码
	Password string `yaml:"password"`

	// Timeout 超时时间(秒)
	Timeout int `yaml:"timeout"`
}

// ImageConfig 图片配置
type ImageConfig struct {
	// MaxFileSize 最大文件大小(字节)
	MaxFileSize int64 `yaml:"maxFileSize"`

	// PathPrefix 存储路径前缀
	PathPrefix string `yaml:"pathPrefix"`

	// AllowedMimeTypes 允许的MIME类型
	AllowedMimeTypes []string `yaml:"allowedMimeTypes"`
}

// 全局存储提供者实例
var globalProvider Provider

// 全局配置（用于显示）
var globalConfig *Config

// Init 初始化存储提供者
func Init(cfg *Config) error {
	// 保存配置
	globalConfig = cfg

	switch cfg.Type {
	case "local", "":
		// local 是默认存储类型
		p, err := NewLocalProvider(cfg.Local)
		if err != nil {
			return err
		}
		globalProvider = p
	case "rustfs":
		if cfg.RustFS == nil {
			return ErrConfigMissing
		}
		p, err := NewRustFSProvider(cfg.RustFS)
		if err != nil {
			return err
		}
		globalProvider = p
	default:
		return ErrUnsupportedProvider
	}
	return nil
}

// GetProvider 获取全局存储提供者
func GetProvider() Provider {
	if globalProvider == nil {
		panic("storage provider not initialized, forgot to call storage.Init()?")
	}
	return globalProvider
}

// IsInitialized 检查是否已初始化
func IsInitialized() bool {
	return globalProvider != nil
}

// GetCurrentConfig 获取当前配置
func GetCurrentConfig() *Config {
	return globalConfig
}

// IsMigrationEnabled 检查迁移功能是否启用
func IsMigrationEnabled() bool {
	ctx := context.Background()
	return g.Cfg().MustGet(ctx, "storage.migration.enabled", false).Bool()
}
