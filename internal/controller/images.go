package controller

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"

	api "github.com/ciclebyte/wekeep/api/v1/images"
	"github.com/ciclebyte/wekeep/internal/consts"
	storagePkg "github.com/ciclebyte/wekeep/internal/storage"
	service "github.com/ciclebyte/wekeep/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Images = imagesController{}

type imagesController struct {
	BaseController
}

// File 图片文件访问代理（直接返回图片二进制，不走 JSON 序列化）
func File(r *ghttp.Request) {
	// 从 URL 路径中提取文件相对路径
	filePath := strings.TrimPrefix(r.URL.Path, "/api/v1/images/file/")
	filePath = strings.TrimPrefix(filePath, "/api/v1/images/file")
	if filePath == "" {
		r.Response.WriteStatus(http.StatusBadRequest)
		return
	}

	// 安全检查：防止路径遍历
	if strings.Contains(filePath, "..") {
		r.Response.WriteStatus(http.StatusForbidden)
		return
	}

	// 统一路径分隔符
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	if !storagePkg.IsInitialized() {
		r.Response.WriteStatus(http.StatusServiceUnavailable)
		return
	}

	data, err := storagePkg.GetProvider().Download(r.Context(), filePath)
	if err != nil {
		r.Response.WriteStatus(http.StatusNotFound)
		return
	}

	// 设置 Content-Type
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		r.Response.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		r.Response.Header().Set("Content-Type", "image/png")
	case ".gif":
		r.Response.Header().Set("Content-Type", "image/gif")
	case ".webp":
		r.Response.Header().Set("Content-Type", "image/webp")
	case ".svg":
		r.Response.Header().Set("Content-Type", "image/svg+xml")
	default:
		r.Response.Header().Set("Content-Type", "application/octet-stream")
	}

	// 缓存控制（图片不常变）
	r.Response.Header().Set("Cache-Control", "public, max-age=86400")
	r.Response.Write(data)
}

// List 图片列表
func (c *imagesController) List(ctx context.Context, req *api.ImagesListReq) (res *api.ImagesListRes, err error) {
	res = new(api.ImagesListRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}

	_, res, err = service.Images().List(ctx, req)
	return
}

// Detail 图片详情
func (c *imagesController) Detail(ctx context.Context, req *api.ImagesDetailReq) (res *api.ImagesDetailRes, err error) {
	res = new(api.ImagesDetailRes)
	info, err := service.Images().GetById(ctx, req.Id)
	if err != nil {
		return
	}

	res.ImagesListItem = &api.ImagesListItem{
		Id:             info.Id,
		OriginalUrl:    info.OriginalUrl,
		StorageUrl:     info.StorageUrl,
		FileSize:       info.FileSize,
		MimeType:       info.MimeType,
		RefCount:       info.RefCount,
		DownloadStatus: info.DownloadStatus,
		CreatedAt:      "",
	}
	return
}

// Delete 删除图片
func (c *imagesController) Delete(ctx context.Context, req *api.ImagesDeleteReq) (res *api.ImagesDeleteRes, err error) {
	res = new(api.ImagesDeleteRes)
	err = service.Images().Delete(ctx, req.Id)
	return
}

// Stats 图片统计
func (c *imagesController) Stats(ctx context.Context, req *api.ImagesStatsReq) (res *api.ImagesStatsRes, err error) {
	res = new(api.ImagesStatsRes)
	res, err = service.Images().Stats(ctx)
	return
}

// Retry 重试下载
func (c *imagesController) Retry(ctx context.Context, req *api.ImagesRetryReq) (res *api.ImagesRetryRes, err error) {
	res = new(api.ImagesRetryRes)
	err = service.Images().Retry(ctx, req.Id)
	return
}

// Cleanup 清理无引用图片
func (c *imagesController) Cleanup(ctx context.Context, req *api.ImagesCleanupReq) (res *api.ImagesCleanupRes, err error) {
	res = new(api.ImagesCleanupRes)
	deletedNum, err := service.Images().Cleanup(ctx)
	res.DeletedNum = deletedNum
	return
}

// Migrate 迁移图片到分层存储
func (c *imagesController) Migrate(ctx context.Context, req *api.ImagesMigrateReq) (res *api.ImagesMigrateRes, err error) {
	res = new(api.ImagesMigrateRes)
	res.Migrated, res.Failed, err = service.Images().MigrateToHierarchical(ctx)
	return
}

// FixContent 修复文章内容中的图片URL
func (c *imagesController) FixContent(ctx context.Context, req *api.ImagesFixContentReq) (res *api.ImagesFixContentRes, err error) {
	res = new(api.ImagesFixContentRes)
	res.UpdatedCount, err = service.Images().FixArticleImageURLs(ctx)
	return
}
