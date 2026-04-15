// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Images is the golang structure of table images for DAO operations like Where/Data.
type Images struct {
	g.Meta          `orm:"table:images, do:true"`
	Id              any         // 主键ID
	OriginalUrl     any         // 原始图片URL
	OriginalUrlHash any         // URL哈希(SHA256)
	StoragePath     any         // RustFS存储路径
	StorageUrl      any         // 访问URL
	FileSize        any         // 文件大小(字节)
	MimeType        any         // MIME类型
	RefCount        any         // 引用计数
	DownloadStatus  any         // 下载状态
	ErrorMessage    any         // 错误信息
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}
