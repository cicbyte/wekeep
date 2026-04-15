// Package images 图片管理 API
package images

import (
	commonApi "github.com/ciclebyte/wekeep/api/v1/common"
	"github.com/gogf/gf/v2/frame/g"
)

// ImagesListReq 图片列表请求
type ImagesListReq struct {
	g.Meta        `path:"/images/list" method:"get" tags:"图片" summary:"图片-列表"`
	commonApi.PageReq
	DownloadStatus *int `json:"downloadStatus" dc:"下载状态: 0-待下载, 1-下载中, 2-成功, 3-失败"`
}

// ImagesListRes 图片列表响应
type ImagesListRes struct {
	g.Meta `mime:"application/json" example:"string"`
	commonApi.ListRes
	List []ImagesListItem `json:"list"`
}

// ImagesListItem 图片列表项
type ImagesListItem struct {
	Id             uint   `json:"id" dc:"图片ID"`
	OriginalUrl    string `json:"originalUrl" dc:"原始URL"`
	StorageUrl     string `json:"storageUrl" dc:"存储URL"`
	FileSize       int    `json:"fileSize" dc:"文件大小(字节)"`
	MimeType       string `json:"mimeType" dc:"MIME类型"`
	RefCount       int    `json:"refCount" dc:"引用计数"`
	DownloadStatus int    `json:"downloadStatus" dc:"下载状态"`
	CreatedAt      string `json:"createdAt" dc:"创建时间"`
}

// ImagesDetailReq 图片详情请求
type ImagesDetailReq struct {
	g.Meta `path:"/images/detail" method:"get" tags:"图片" summary:"图片-详情"`
	Id     uint `json:"id" v:"required#图片ID不能为空" dc:"图片ID"`
}

// ImagesDetailRes 图片详情响应
type ImagesDetailRes struct {
	g.Meta `mime:"application/json"`
	*ImagesListItem
}

// ImagesDeleteReq 删除图片请求
type ImagesDeleteReq struct {
	g.Meta `path:"/images/delete" method:"delete" tags:"图片" summary:"图片-删除"`
	Id     uint `json:"id" v:"required#图片ID不能为空" dc:"图片ID"`
}

// ImagesDeleteRes 删除图片响应
type ImagesDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ImagesStatsReq 图片统计请求
type ImagesStatsReq struct {
	g.Meta `path:"/images/stats" method:"get" tags:"图片" summary:"图片-统计"`
}

// ImagesStatsRes 图片统计响应
type ImagesStatsRes struct {
	g.Meta       `mime:"application/json"`
	TotalImages  int `json:"totalImages" dc:"图片总数"`
	TotalSize    int `json:"totalSize" dc:"总大小(字节)"`
	PendingCount int `json:"pendingCount" dc:"待下载数"`
	SuccessCount int `json:"successCount" dc:"成功数"`
	FailedCount  int `json:"failedCount" dc:"失败数"`
}

// ImagesRetryReq 重试下载请求
type ImagesRetryReq struct {
	g.Meta `path:"/images/retry" method:"post" tags:"图片" summary:"图片-重试下载"`
	Id     uint `json:"id" v:"required#图片ID不能为空" dc:"图片ID"`
}

// ImagesRetryRes 重试下载响应
type ImagesRetryRes struct {
	g.Meta `mime:"application/json"`
}

// ImagesCleanupReq 清理无引用图片请求
type ImagesCleanupReq struct {
	g.Meta `path:"/images/cleanup" method:"post" tags:"图片" summary:"图片-清理无引用"`
}

// ImagesCleanupRes 清理无引用图片响应
type ImagesCleanupRes struct {
	g.Meta     `mime:"application/json"`
	DeletedNum int `json:"deletedNum" dc:"删除数量"`
}

// ImagesMigrateReq 图片迁移请求
type ImagesMigrateReq struct {
	g.Meta `path:"/images/migrate" method:"post" tags:"图片" summary:"图片-迁移到分层存储"`
}

// ImagesMigrateRes 图片迁移响应
type ImagesMigrateRes struct {
	g.Meta    `mime:"application/json"`
	Migrated  int `json:"migrated" dc:"迁移成功数"`
	Failed    int `json:"failed" dc:"迁移失败数"`
}

// ImagesFixContentReq 修复文章内容图片URL请求
type ImagesFixContentReq struct {
	g.Meta `path:"/images/fix-content" method:"post" tags:"图片" summary:"图片-修复文章内容URL"`
}

// ImagesFixContentRes 修复文章内容图片URL响应
type ImagesFixContentRes struct {
	g.Meta       `mime:"application/json"`
	UpdatedCount int `json:"updatedCount" dc:"更新文章数"`
}

// ImagesFileReq 图片文件访问请求
type ImagesFileReq struct {
	g.Meta `path:"/images/file/{path}" method:"get" tags:"图片" summary:"图片-文件访问"`
	Path   string `in:"path" json:"path"`
}
