// Package storage 存储迁移 API
package storage

import (
	commonApi "github.com/ciclebyte/wekeep/api/v1/common"
	"github.com/gogf/gf/v2/frame/g"
)

// StorageMigrateReq 存储迁移请求
type StorageMigrateReq struct {
	g.Meta          `path:"/storage/migrate" method:"post" tags:"存储管理" summary:"迁移存储"`
	SourceStorage   string       `json:"sourceStorage"`                                                         // 源存储类型，留空则使用当前配置
	TargetStorage   string       `json:"targetStorage" v:"required|in:local,rustfs#目标存储类型不能为空|目标存储类型必须是 local 或 rustfs"`
	LocalConfig     *LocalConfig `json:"localConfig"`
	RustFSConfig    *RustFSConfig `json:"rustfsConfig"`
	UpdateMarkdown  bool         `json:"updateMarkdown" d:"true"` // 是否更新 Markdown 内容
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	BasePath string `json:"basePath"` // 存储根目录
	BaseURL  string `json:"baseURL"`  // 访问URL前缀
}

// RustFSConfig RustFS 配置
type RustFSConfig struct {
	Endpoint string `json:"endpoint" v:"required#Endpoint 不能为空"`
	Bucket   string `json:"bucket" v:"required#Bucket 不能为空"`
	Username string `json:"username" v:"required#用户名不能为空"`
	Password string `json:"password" v:"required#密码不能为空"`
	Timeout  int    `json:"timeout" d:"60"`
}

// StorageMigrateRes 迁移响应
type StorageMigrateRes struct {
	g.Meta `mime:"application/json"`
	Message string `json:"message"`
}

// StorageStatusReq 存储状态请求
type StorageStatusReq struct {
	g.Meta `path:"/storage/status" method:"get" tags:"存储管理" summary:"获取存储状态"`
}

// StorageStatusRes 存储状态响应
type StorageStatusRes struct {
	g.Meta `mime:"application/json"`
	// 当前存储类型
	CurrentStorage string `json:"currentStorage"`
	// 是否已初始化
	Initialized bool `json:"initialized"`
	// 迁移状态
	Migration *MigrationStatus `json:"migration"`
	// 当前存储配置（隐藏敏感信息）
	CurrentConfig *StorageConfig `json:"currentConfig,omitempty"`
	// 是否启用迁移功能
	MigrationEnabled bool `json:"migrationEnabled"`
}

// StorageConfig 存储配置（用于显示）
type StorageConfig struct {
	// RustFS 配置
	RustFS *RustFSConfigDisplay `json:"rustfs,omitempty"`
	// 本地配置
	Local *LocalConfigDisplay `json:"local,omitempty"`
}

// RustFSConfigDisplay RustFS 配置（隐藏密码）
type RustFSConfigDisplay struct {
	Endpoint string `json:"endpoint"`
	Bucket   string `json:"bucket"`
	Username string `json:"username"`
	Timeout  int    `json:"timeout"`
}

// LocalConfigDisplay 本地存储配置
type LocalConfigDisplay struct {
	BasePath string `json:"basePath"`
	BaseURL  string `json:"baseURL"`
}

// MigrationStatus 迁移状态
type MigrationStatus struct {
	Running      bool    `json:"running"`
	Source       string  `json:"source"`
	Target       string  `json:"target"`
	Total        int     `json:"total"`
	Completed    int     `json:"completed"`
	Failed       int     `json:"failed"`
	CurrentFile  string  `json:"currentFile"`
	StartTime    string  `json:"startTime"`
	EndTime      string  `json:"endTime"`
	Error        string  `json:"error"`
}

// StorageValidateReq 验证存储配置请求
type StorageValidateReq struct {
	g.Meta          `path:"/storage/validate" method:"post" tags:"存储管理" summary:"验证存储配置"`
	TargetStorage   string       `json:"targetStorage" v:"required|in:local,rustfs#目标存储类型不能为空|目标存储类型必须是 local 或 rustfs"`
	LocalConfig     *LocalConfig `json:"localConfig"`
	RustFSConfig    *RustFSConfig `json:"rustfsConfig"`
}

// StorageValidateRes 验证响应
type StorageValidateRes struct {
	g.Meta `mime:"application/json"`
	Valid  bool   `json:"valid"`
	Error  string `json:"error,omitempty"`
}

// StorageStatsReq 存储统计请求
type StorageStatsReq struct {
	g.Meta `path:"/storage/stats" method:"get" tags:"存储管理" summary:"获取存储统计"`
}

// StorageStatsRes 存储统计响应
type StorageStatsRes struct {
	g.Meta `mime:"application/json"`
	// 图片总数
	TotalImages int `json:"totalImages"`
	// 总大小（字节）
	TotalSize int64 `json:"totalSize"`
	// 本地存储大小
	LocalSize int64 `json:"localSize"`
	// 按存储类型统计
	ByStorageType map[string]int64 `json:"byStorageType"`
}

// StorageSwitchReq 切换存储请求
type StorageSwitchReq struct {
	g.Meta         `path:"/storage/switch" method:"post" tags:"存储管理" summary:"切换存储（更新配置文件）"`
	TargetStorage  string       `json:"targetStorage" v:"required|in:local,rustfs#目标存储类型不能为空|目标存储类型必须是 local 或 rustfs"`
	LocalConfig    *LocalConfig `json:"localConfig"`
	RustFSConfig   *RustFSConfig `json:"rustfsConfig"`
}

// StorageSwitchRes 切换响应
type StorageSwitchRes struct {
	g.Meta   `mime:"application/json"`
	Message  string `json:"message"`
	Requires string `json:"requires,omitempty"` // 提示需要重启
}

// StorageUpdateRefsReq 更新图片引用请求
type StorageUpdateRefsReq struct {
	g.Meta `path:"/storage/update-refs" method:"post" tags:"存储管理" summary:"更新图片引用"`
}

// StorageUpdateRefsRes 更新引用响应
type StorageUpdateRefsRes struct {
	g.Meta   `mime:"application/json"`
	Message  string `json:"message"`
	Updated  int    `json:"updated"`  // 更新的文章数
	Total    int    `json:"total"`    // 总图片映射数
}

// ListReq 通用列表请求
type ListReq struct {
	g.Meta `path:"/storage" method:"get" tags:"存储管理" summary:"存储信息"`
	commonApi.PageReq
}

// ListRes 通用列表响应
type ListRes struct {
	g.Meta `mime:"application/json"`
	commonApi.ListRes
}
