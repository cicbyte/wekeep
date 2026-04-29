package images

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	api "github.com/cicbyte/wekeep/api/v1/images"
	dao "github.com/cicbyte/wekeep/internal/dao"
	model "github.com/cicbyte/wekeep/internal/model"
	do "github.com/cicbyte/wekeep/internal/model/do"
	entity "github.com/cicbyte/wekeep/internal/model/entity"
	"github.com/cicbyte/wekeep/internal/storage"
	service "github.com/cicbyte/wekeep/internal/service"
	liberr "github.com/cicbyte/wekeep/library/liberr"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	service.RegisterImages(New())
}

func New() *sImages {
	return &sImages{}
}

type sImages struct{}

// 微信图片域名列表
var wechatImageDomains = []string{
	"mmbiz.qpic.cn",
	"mmbiz.qlogo.cn",
	"wx.qlogo.cn",
	"thirdwx.qlogo.cn",
}

// isWechatImage 检查是否为微信图片
func isWechatImage(imageURL string) bool {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return false
	}
	for _, domain := range wechatImageDomains {
		if strings.Contains(parsedURL.Host, domain) {
			return true
		}
	}
	return false
}

// extractImageURLs 从 Markdown 内容中提取图片 URL
func extractImageURLs(content string) []string {
	// 匹配 Markdown 图片语法: ![alt](url)
	mdRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	matches := mdRegex.FindAllStringSubmatch(content, -1)

	urls := make([]string, 0)
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) >= 3 {
			imageURL := match[2]
			// 过滤空URL和已处理的URL
			if imageURL != "" && !seen[imageURL] && isWechatImage(imageURL) {
				urls = append(urls, imageURL)
				seen[imageURL] = true
			}
		}
	}

	return urls
}

// hashURL 计算URL的SHA256哈希
func hashURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	return hex.EncodeToString(hash[:])
}

// List 获取图片列表
func (s *sImages) List(ctx context.Context, req *api.ImagesListReq) (total interface{}, res *api.ImagesListRes, err error) {
	res = &api.ImagesListRes{}
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.Images.Ctx(ctx)

		// 按下载状态筛选
		if req.DownloadStatus != nil {
			m = m.Where("download_status = ?", *req.DownloadStatus)
		}

		m = m.Safe()
		total, err = m.Count()
		liberr.ErrIsNil(ctx, err, "获取图片列表失败")

		var entities []*entity.Images
		err = m.Page(req.PageNum, req.PageSize).Order("created_at DESC").Scan(&entities)
		liberr.ErrIsNil(ctx, err, "获取图片列表失败")

		res.Total = total
		res.CurrentPage = req.PageNum
		for _, e := range entities {
			res.List = append(res.List, api.ImagesListItem{
				Id:             e.Id,
				OriginalUrl:    e.OriginalUrl,
				StorageUrl:     e.StorageUrl,
				FileSize:       e.FileSize,
				MimeType:       e.MimeType,
				RefCount:       e.RefCount,
				DownloadStatus: e.DownloadStatus,
				CreatedAt:      e.CreatedAt.String(),
			})
		}
	})
	return
}

// GetById 获取图片详情
func (s *sImages) GetById(ctx context.Context, id uint) (res *model.ImagesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		var e *entity.Images
		err = dao.Images.Ctx(ctx).WherePri(id).Scan(&e)
		liberr.ErrIsNil(ctx, err, "图片不存在")

		res = &model.ImagesInfo{
			Id:              e.Id,
			OriginalUrl:     e.OriginalUrl,
			OriginalUrlHash: e.OriginalUrlHash,
			StoragePath:     e.StoragePath,
			StorageUrl:      e.StorageUrl,
			FileSize:        e.FileSize,
			MimeType:        e.MimeType,
			RefCount:        e.RefCount,
			DownloadStatus:  e.DownloadStatus,
			ErrorMessage:    e.ErrorMessage,
		}
	})
	return
}

// Delete 删除图片
func (s *sImages) Delete(ctx context.Context, id uint) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 获取图片信息
		imageInfo, err := s.GetById(ctx, id)
		liberr.ErrIsNil(ctx, err, "图片不存在")

		// 检查引用计数
		if imageInfo.RefCount > 0 {
			liberr.ErrIsNil(ctx, fmt.Errorf("图片仍有 %d 个引用，无法删除", imageInfo.RefCount), "删除失败")
		}

		// 从存储中删除文件
		if storage.IsInitialized() {
			provider := storage.GetProvider()
			deleteErr := provider.Delete(ctx, imageInfo.StoragePath)
			if deleteErr != nil {
				g.Log().Warningf(ctx, "Failed to delete file from storage: %v", deleteErr)
			}
		}

		// 删除数据库记录
		_, err = dao.Images.Ctx(ctx).WherePri(id).Delete()
		liberr.ErrIsNil(ctx, err, "删除图片失败")
	})
	return
}

// Stats 获取图片统计
func (s *sImages) Stats(ctx context.Context) (res *api.ImagesStatsRes, err error) {
	res = &api.ImagesStatsRes{}
	err = g.Try(ctx, func(ctx context.Context) {
		// 总数和总大小
		var stats struct {
			Total     int
			TotalSize int
		}
		err = dao.Images.Ctx(ctx).
			Fields("COUNT(*) as total, COALESCE(SUM(file_size), 0) as total_size").
			Scan(&stats)
		liberr.ErrIsNil(ctx, err, "获取统计失败")
		res.TotalImages = stats.Total
		res.TotalSize = stats.TotalSize

		// 按状态统计
		statusCounts := make(map[int]int)
		var statusStats []*struct {
			Status int `orm:"download_status"`
			Count  int `orm:"count"`
		}
		err = dao.Images.Ctx(ctx).
			Fields("download_status as status, COUNT(*) as count").
			Group("download_status").
			Scan(&statusStats)
		liberr.ErrIsNil(ctx, err, "获取状态统计失败")
		for _, s := range statusStats {
			statusCounts[s.Status] = s.Count
		}

		res.PendingCount = statusCounts[model.DownloadStatusPending]
		res.SuccessCount = statusCounts[model.DownloadStatusSuccess]
		res.FailedCount = statusCounts[model.DownloadStatusFailed]
	})
	return
}

// Retry 重试下载
func (s *sImages) Retry(ctx context.Context, id uint) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		imageInfo, err := s.GetById(ctx, id)
		liberr.ErrIsNil(ctx, err, "图片不存在")

		// 重置状态为待下载
		_, err = dao.Images.Ctx(ctx).WherePri(id).Update(do.Images{
			DownloadStatus: model.DownloadStatusPending,
			ErrorMessage:   "",
		})
		liberr.ErrIsNil(ctx, err, "重置状态失败")

		// 异步执行下载
		go func() {
			bgCtx := context.Background()
			_, downloadErr := s.DownloadAndUpload(bgCtx, imageInfo.OriginalUrl)
			if downloadErr != nil {
				g.Log().Errorf(bgCtx, "Retry download failed for image %d: %v", id, downloadErr)
			}
		}()
	})
	return
}

// Cleanup 清理无引用的图片
func (s *sImages) Cleanup(ctx context.Context) (deletedNum int, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 查找引用计数为0的图片
		var images []*entity.Images
		err = dao.Images.Ctx(ctx).Where("ref_count = 0").Scan(&images)
		liberr.ErrIsNil(ctx, err, "查询无引用图片失败")

		for _, img := range images {
			// 从存储中删除
			if storage.IsInitialized() && img.StoragePath != "" {
				deleteErr := storage.GetProvider().Delete(ctx, img.StoragePath)
				if deleteErr != nil {
					g.Log().Warningf(ctx, "Failed to delete file %s: %v", img.StoragePath, deleteErr)
				}
			}

			// 删除数据库记录
			_, err = dao.Images.Ctx(ctx).WherePri(img.Id).Delete()
			if err != nil {
				g.Log().Warningf(ctx, "Failed to delete image record %d: %v", img.Id, err)
				continue
			}
			deletedNum++
		}
	})
	return
}

// ProcessArticleImages 处理文章中的图片
func (s *sImages) ProcessArticleImages(ctx context.Context, articleId int, content string) (newContent string, imageIds []uint, err error) {
	newContent = content
	imageIds = make([]uint, 0)

	err = g.Try(ctx, func(ctx context.Context) {
		// 提取微信图片URL
		imageURLs := extractImageURLs(content)
		if len(imageURLs) == 0 {
			return
		}

		g.Log().Infof(ctx, "Found %d wechat images in article %d", len(imageURLs), articleId)

		// 处理每个图片URL
		type replaceInfo struct {
			oldURL string
			newURL string
			imageID uint
		}
		replacements := make([]replaceInfo, 0)

		for _, imageURL := range imageURLs {
			// 检查是否已存在
			urlHash := hashURL(imageURL)
			existingImage, err := s.GetByURLHash(ctx, urlHash)
			if err == nil && existingImage != nil {
				// 已存在，复用
				g.Log().Infof(ctx, "Image already exists: %d", existingImage.Id)
				if existingImage.DownloadStatus == model.DownloadStatusSuccess {
					replacements = append(replacements, replaceInfo{
						oldURL:  imageURL,
						newURL:  existingImage.StorageUrl,
						imageID: existingImage.Id,
					})
					imageIds = append(imageIds, existingImage.Id)
				}
				continue
			}

			// 下载并上传图片
			imageInfo, downloadErr := s.DownloadAndUpload(ctx, imageURL)
			if downloadErr != nil {
				g.Log().Warningf(ctx, "Failed to download image %s: %v", imageURL, downloadErr)
				continue
			}

			if imageInfo.DownloadStatus == model.DownloadStatusSuccess {
				replacements = append(replacements, replaceInfo{
					oldURL:  imageURL,
					newURL:  imageInfo.StorageUrl,
					imageID: imageInfo.Id,
				})
				imageIds = append(imageIds, imageInfo.Id)
			}
		}

		// 替换内容中的URL
		for _, r := range replacements {
			newContent = strings.ReplaceAll(newContent, r.oldURL, r.newURL)
		}

		// 创建文章-图片关联
		for i, imageId := range imageIds {
			_, err = dao.ArticleImages.Ctx(ctx).Insert(do.ArticleImages{
				ArticleId: uint(articleId),
				ImageId:   imageId,
				Position:  i,
			})
			if err != nil {
				g.Log().Warningf(ctx, "Failed to create article-image relation: %v", err)
			}
		}

		// 增加引用计数
		if len(imageIds) > 0 {
			err = s.IncrementRefCount(ctx, imageIds)
			if err != nil {
				g.Log().Warningf(ctx, "Failed to increment ref count: %v", err)
			}
		}
	})

	return
}

// RemoveArticleImages 移除文章关联的图片
func (s *sImages) RemoveArticleImages(ctx context.Context, articleId int) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 获取文章关联的所有图片ID
		var relations []*entity.ArticleImages
		err = dao.ArticleImages.Ctx(ctx).Where("article_id = ?", articleId).Scan(&relations)
		liberr.ErrIsNil(ctx, err, "查询文章图片关联失败")

		if len(relations) == 0 {
			return
		}

		// 收集图片ID
		imageIds := make([]uint, 0)
		for _, r := range relations {
			imageIds = append(imageIds, r.ImageId)
		}

		// 删除关联记录
		_, err = dao.ArticleImages.Ctx(ctx).Where("article_id = ?", articleId).Delete()
		liberr.ErrIsNil(ctx, err, "删除文章图片关联失败")

		// 减少引用计数
		err = s.DecrementRefCount(ctx, imageIds)
		if err != nil {
			g.Log().Warningf(ctx, "Failed to decrement ref count: %v", err)
		}
	})
	return
}

// DownloadAndUpload 下载图片并上传到存储
func (s *sImages) DownloadAndUpload(ctx context.Context, originalURL string) (imageInfo *model.ImagesInfo, err error) {
	urlHash := hashURL(originalURL)

	// 先检查是否已存在
	existingImage, err := s.GetByURLHash(ctx, urlHash)
	if err == nil && existingImage != nil {
		return existingImage, nil
	}

	// 创建图片记录（状态为下载中）
	result, err := dao.Images.Ctx(ctx).Insert(do.Images{
		OriginalUrl:     originalURL,
		OriginalUrlHash: urlHash,
		StoragePath:     "",
		StorageUrl:      "",
		FileSize:        0,
		MimeType:        "",
		RefCount:        0,
		DownloadStatus:  model.DownloadStatusDownloading,
	})
	liberr.ErrIsNil(ctx, err, "创建图片记录失败")

	imageId, _ := result.LastInsertId()

	// 更新状态为待下载（如果后续步骤失败，可以被重试）
	_, _ = dao.Images.Ctx(ctx).WherePri(imageId).Update(do.Images{
		DownloadStatus: model.DownloadStatusPending,
	})

	// 检查存储是否已初始化
	if !storage.IsInitialized() {
		// 存储未初始化，只记录URL不下载
		return s.GetById(ctx, uint(imageId))
	}

	// 下载图片
	g.Log().Infof(ctx, "Downloading image: %s", originalURL)

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	// 创建请求，模拟浏览器访问
	req, err := http.NewRequestWithContext(ctx, "GET", originalURL, nil)
	if err != nil {
		s.updateImageStatus(ctx, uint(imageId), model.DownloadStatusFailed, fmt.Sprintf("创建请求失败: %v", err))
		return nil, err
	}

	// 设置请求头模拟浏览器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Referer", "https://mp.weixin.qq.com/")

	resp, err := httpClient.Do(req)
	if err != nil {
		s.updateImageStatus(ctx, uint(imageId), model.DownloadStatusFailed, fmt.Sprintf("下载失败: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP status %d", resp.StatusCode)
		s.updateImageStatus(ctx, uint(imageId), model.DownloadStatusFailed, err.Error())
		return nil, err
	}

	// 读取图片内容
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		s.updateImageStatus(ctx, uint(imageId), model.DownloadStatusFailed, fmt.Sprintf("读取响应失败: %v", err))
		return nil, err
	}

	// 获取 MIME 类型
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg" // 默认类型
	}

	// 生成存储路径（按日期分层）
	storagePath := generateStoragePath(contentType)

	// 上传到存储
	provider := storage.GetProvider()
	storageURL, err := provider.Upload(ctx, bytes.NewReader(imageData), storagePath, contentType)
	if err != nil {
		s.updateImageStatus(ctx, uint(imageId), model.DownloadStatusFailed, fmt.Sprintf("上传失败: %v", err))
		return nil, err
	}

	// 更新图片记录（storage_url 存相对路径，前端渲染时动态拼接前缀）
	_, err = dao.Images.Ctx(ctx).WherePri(imageId).Update(do.Images{
		StoragePath:    storagePath,
		StorageUrl:     storageURL,
		FileSize:       len(imageData),
		MimeType:       contentType,
		DownloadStatus: model.DownloadStatusSuccess,
		ErrorMessage:   "",
	})
	liberr.ErrIsNil(ctx, err, "更新图片记录失败")

	g.Log().Infof(ctx, "Image downloaded and uploaded successfully: %s -> %s", originalURL, storageURL)

	return s.GetById(ctx, uint(imageId))
}

// updateImageStatus 更新图片状态
func (s *sImages) updateImageStatus(ctx context.Context, imageId uint, status int, errMsg string) {
	_, err := dao.Images.Ctx(ctx).WherePri(imageId).Update(do.Images{
		DownloadStatus: status,
		ErrorMessage:   errMsg,
	})
	if err != nil {
		g.Log().Warningf(ctx, "Failed to update image status: %v", err)
	}
}

// GetByURLHash 根据URL哈希获取图片
func (s *sImages) GetByURLHash(ctx context.Context, urlHash string) (res *model.ImagesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		var e *entity.Images
		err = dao.Images.Ctx(ctx).Where("original_url_hash = ?", urlHash).Scan(&e)
		liberr.ErrIsNil(ctx, err, "图片不存在")

		res = &model.ImagesInfo{
			Id:              e.Id,
			OriginalUrl:     e.OriginalUrl,
			OriginalUrlHash: e.OriginalUrlHash,
			StoragePath:     e.StoragePath,
			StorageUrl:      e.StorageUrl,
			FileSize:        e.FileSize,
			MimeType:        e.MimeType,
			RefCount:        e.RefCount,
			DownloadStatus:  e.DownloadStatus,
			ErrorMessage:    e.ErrorMessage,
		}
	})
	return
}

// IncrementRefCount 增加引用计数
func (s *sImages) IncrementRefCount(ctx context.Context, imageIds []uint) (err error) {
	if len(imageIds) == 0 {
		return nil
	}

	for _, id := range imageIds {
		_, err = dao.Images.Ctx(ctx).WherePri(id).Increment("ref_count", 1)
		if err != nil {
			g.Log().Warningf(ctx, "Failed to increment ref count for image %d: %v", id, err)
		}
	}
	return nil
}

// DecrementRefCount 减少引用计数
func (s *sImages) DecrementRefCount(ctx context.Context, imageIds []uint) (err error) {
	if len(imageIds) == 0 {
		return nil
	}

	for _, id := range imageIds {
		// 先减少引用计数
		_, err = dao.Images.Ctx(ctx).WherePri(id).Decrement("ref_count", 1)
		if err != nil {
			g.Log().Warningf(ctx, "Failed to decrement ref count for image %d: %v", id, err)
			continue
		}

		// 检查引用计数是否为0
		var image *entity.Images
		err = dao.Images.Ctx(ctx).WherePri(id).Scan(&image)
		if err != nil || image == nil {
			continue
		}

		if image.RefCount <= 0 {
		g.Log().Infof(ctx, "Image %d has no more references, deleting", id)
			// 从存储中删除文件
			if storage.IsInitialized() && image.StoragePath != "" {
				normalizedPath := strings.ReplaceAll(image.StoragePath, "\\", "/")
				if delErr := storage.GetProvider().Delete(ctx, normalizedPath); delErr != nil {
					g.Log().Warningf(ctx, "Failed to delete file %s: %v", image.StoragePath, delErr)
				}
			}
			// 删除数据库记录
			if _, delErr := dao.Images.Ctx(ctx).WherePri(id).Delete(); delErr != nil {
				g.Log().Warningf(ctx, "Failed to delete image record %d: %v", id, delErr)
			}
		}
	}
	return nil
}

// generateStoragePath 生成存储路径
func generateStoragePath(contentType string) string {
	ext := ".jpg"
	switch contentType {
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	}

	now := time.Now()
	dateDir := now.Format("2006/01/02")
	filename := fmt.Sprintf("%d%s", now.UnixNano(), ext)
	// 使用正斜杠拼接路径（S3/RustFS 只认 /，filepath.Join 在 Windows 下会产生 \）
	return dateDir + "/" + filename
}

// MigrateToHierarchical 迁移图片到分层存储结构
func (s *sImages) MigrateToHierarchical(ctx context.Context) (migrated int, failed int, err error) {
	g.Log().Info(ctx, "Starting image migration to hierarchical storage...")

	provider := storage.GetProvider()

	// 查询所有需要迁移的图片（storagePath 不包含路径分隔符，即平铺结构）
	var images []*entity.Images
	err = dao.Images.Ctx(ctx).
		Where("storage_path NOT LIKE ?", "%/%").
		Where("storage_path != ?", "").
		Where("download_status", model.DownloadStatusSuccess).
		Scan(&images)
	if err != nil {
		g.Log().Errorf(ctx, "Failed to query images for migration: %v", err)
		return 0, 0, err
	}

	if len(images) == 0 {
		g.Log().Info(ctx, "No images need migration")
		return 0, 0, nil
	}

	g.Log().Infof(ctx, "Found %d images to migrate", len(images))

	// 构建 URL 映射表
	urlMapping := make(map[string]string)

	for _, img := range images {
		oldPath := img.StoragePath

		// 生成新的分层路径
		var createdAt time.Time
		if img.CreatedAt != nil {
			createdAt = img.CreatedAt.Time
		} else {
			createdAt = time.Now()
		}

		// 使用图片创建时间构建目录结构
		dateDir := createdAt.Format("2006/01/02")
		filename := filepath.Base(oldPath)
		newPath := dateDir + "/" + filename

		// 如果新旧路径相同，跳过
		if oldPath == newPath {
			continue
		}

		g.Log().Infof(ctx, "Migrating image %d: %s -> %s", img.Id, oldPath, newPath)

		// 复制文件到新路径
		err := provider.Copy(ctx, oldPath, newPath)
		if err != nil {
			g.Log().Errorf(ctx, "Failed to copy image %d: %v", img.Id, err)
			failed++
			continue
		}

		// 生成新 URL 并记录映射
		newURL := provider.GetURL(newPath)
		urlMapping[img.StorageUrl] = newURL

		// 更新图片记录
		_, err = dao.Images.Ctx(ctx).WherePri(img.Id).Update(do.Images{
			StoragePath: newPath,
			StorageUrl:  newURL,
		})
		if err != nil {
			g.Log().Errorf(ctx, "Failed to update image record %d: %v", img.Id, err)
			// 尝试删除新文件，回滚
			provider.Delete(ctx, newPath)
			failed++
			continue
		}

		// 删除旧文件
		err = provider.Delete(ctx, oldPath)
		if err != nil {
			g.Log().Warningf(ctx, "Failed to delete old image %s: %v", oldPath, err)
			// 不计入失败，因为迁移已成功
		}

		migrated++
		g.Log().Infof(ctx, "Image %d migrated successfully", img.Id)
	}

	// 更新文章内容中的图片 URL
	if len(urlMapping) > 0 {
		g.Log().Info(ctx, "Updating article content with new image URLs...")
		updatedCount, err := s.updateArticleImageURLs(ctx, urlMapping)
		if err != nil {
			g.Log().Errorf(ctx, "Failed to update article content: %v", err)
		} else {
			g.Log().Infof(ctx, "Updated %d articles with new image URLs", updatedCount)
		}
	}

	g.Log().Infof(ctx, "Migration completed: %d migrated, %d failed", migrated, failed)
	return migrated, failed, nil
}

// updateArticleImageURLs 更新文章内容中的图片 URL
func (s *sImages) updateArticleImageURLs(ctx context.Context, urlMapping map[string]string) (int, error) {
	// 查询所有文章
	var articles []*entity.Articles
	err := dao.Articles.Ctx(ctx).Scan(&articles)
	if err != nil {
		return 0, err
	}

	updatedCount := 0
	for _, article := range articles {
		if article.Content == "" {
			continue
		}

		newContent := article.Content
		hasChanges := false

		// 替换所有旧 URL 为新 URL
		for oldURL, newURL := range urlMapping {
			if strings.Contains(newContent, oldURL) {
				newContent = strings.ReplaceAll(newContent, oldURL, newURL)
				hasChanges = true
			}
		}

		// 如果有变化，更新文章
		if hasChanges {
			_, err := dao.Articles.Ctx(ctx).WherePri(article.Id).Update(do.Articles{
				Content: newContent,
			})
			if err != nil {
				g.Log().Errorf(ctx, "Failed to update article %d: %v", article.Id, err)
				continue
			}
			updatedCount++
		}
	}

	return updatedCount, nil
}

// FixArticleImageURLs 修复文章内容中的图片URL
// 用于处理已迁移图片但文章内容未更新的情况
func (s *sImages) FixArticleImageURLs(ctx context.Context) (int, error) {
	g.Log().Info(ctx, "Fixing article content image URLs...")

	provider := storage.GetProvider()

	// 查询所有图片记录
	var images []*entity.Images
	err := dao.Images.Ctx(ctx).
		Where("download_status", model.DownloadStatusSuccess).
		Where("storage_path != ?", "").
		Scan(&images)
	if err != nil {
		return 0, err
	}

	// 构建文件名到新 URL 的映射
	filenameToURL := make(map[string]string)
	for _, img := range images {
		filename := filepath.Base(img.StoragePath)
		filenameToURL[filename] = img.StorageUrl
	}

	g.Log().Infof(ctx, "Found %d images to check", len(filenameToURL))

	// 查询所有文章
	var articles []*entity.Articles
	err = dao.Articles.Ctx(ctx).Scan(&articles)
	if err != nil {
		return 0, err
	}

	updatedCount := 0
	baseURL := strings.TrimSuffix(provider.GetURL(""), "/")

	for _, article := range articles {
		if article.Content == "" {
			continue
		}

		newContent := article.Content
		hasChanges := false

		// 查找内容中的图片 URL 并替换
		for filename, newURL := range filenameToURL {
			// 匹配旧格式的 URL: .../wekeep/文件名.扩展名
			oldPattern := baseURL + "/" + filename
			if strings.Contains(newContent, oldPattern) && oldPattern != newURL {
				newContent = strings.ReplaceAll(newContent, oldPattern, newURL)
				hasChanges = true
			}
		}

		if hasChanges {
			_, err := dao.Articles.Ctx(ctx).WherePri(article.Id).Update(do.Articles{
				Content: newContent,
			})
			if err != nil {
				g.Log().Errorf(ctx, "Failed to update article %d: %v", article.Id, err)
				continue
			}
			updatedCount++
		}
	}

	g.Log().Infof(ctx, "Fixed %d articles", updatedCount)
	return updatedCount, nil
}
