package service

import (
	"context"
	api "github.com/ciclebyte/wekeep/api/v1/images"
	model "github.com/ciclebyte/wekeep/internal/model"
)

type IImages interface {
	// List 获取图片列表
	List(ctx context.Context, req *api.ImagesListReq) (total interface{}, res *api.ImagesListRes, err error)

	// GetById 获取图片详情
	GetById(ctx context.Context, id uint) (res *model.ImagesInfo, err error)

	// Delete 删除图片
	Delete(ctx context.Context, id uint) (err error)

	// Stats 获取图片统计
	Stats(ctx context.Context) (res *api.ImagesStatsRes, err error)

	// Retry 重试下载
	Retry(ctx context.Context, id uint) (err error)

	// Cleanup 清理无引用的图片
	Cleanup(ctx context.Context) (deletedNum int, err error)

	// ProcessArticleImages 处理文章中的图片
	// 返回处理后的内容（URL已替换）和图片ID列表
	ProcessArticleImages(ctx context.Context, articleId int, content string) (newContent string, imageIds []uint, err error)

	// RemoveArticleImages 移除文章关联的图片（删除文章时调用）
	RemoveArticleImages(ctx context.Context, articleId int) (err error)

	// DownloadAndUpload 下载图片并上传到存储
	DownloadAndUpload(ctx context.Context, originalURL string) (imageInfo *model.ImagesInfo, err error)

	// GetByURLHash 根据URL哈希获取图片
	GetByURLHash(ctx context.Context, urlHash string) (res *model.ImagesInfo, err error)

	// IncrementRefCount 增加引用计数
	IncrementRefCount(ctx context.Context, imageIds []uint) (err error)

	// DecrementRefCount 减少引用计数
	DecrementRefCount(ctx context.Context, imageIds []uint) (err error)

	// MigrateToHierarchical 迁移图片到分层存储结构
	// 返回迁移数量和错误
	MigrateToHierarchical(ctx context.Context) (migrated int, failed int, err error)

	// FixArticleImageURLs 修复文章内容中的图片URL
	// 返回更新的文章数量
	FixArticleImageURLs(ctx context.Context) (updatedCount int, err error)
}

var localImages IImages

// Images 返回图片管理服务的实例
func Images() IImages {
	if localImages == nil {
		panic("implement not found for interface IImages, forgot register?")
	}
	return localImages
}

// RegisterImages 注册图片管理服务实现
func RegisterImages(i IImages) {
	localImages = i
}
