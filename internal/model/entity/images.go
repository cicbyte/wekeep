// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Images is the golang structure for table images.
type Images struct {
	Id              uint        `json:"id"              orm:"id"               description:"主键ID"`         // 主键ID
	OriginalUrl     string      `json:"originalUrl"     orm:"original_url"     description:"原始图片URL"`      // 原始图片URL
	OriginalUrlHash string      `json:"originalUrlHash" orm:"original_url_hash" description:"URL哈希(SHA256)"` // URL哈希(SHA256)
	StoragePath     string      `json:"storagePath"     orm:"storage_path"     description:"RustFS存储路径"`   // RustFS存储路径
	StorageUrl      string      `json:"storageUrl"      orm:"storage_url"      description:"访问URL"`         // 访问URL
	FileSize        int         `json:"fileSize"        orm:"file_size"        description:"文件大小(字节)"`     // 文件大小(字节)
	MimeType        string      `json:"mimeType"        orm:"mime_type"        description:"MIME类型"`       // MIME类型
	RefCount        int         `json:"refCount"        orm:"ref_count"        description:"引用计数"`         // 引用计数
	DownloadStatus  int         `json:"downloadStatus"  orm:"download_status"  description:"下载状态"`         // 下载状态
	ErrorMessage    string      `json:"errorMessage"    orm:"error_message"    description:"错误信息"`         // 错误信息
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:"创建时间"`         // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       description:"更新时间"`         // 更新时间
}
