package controller

import (
	"context"
	api "github.com/cicbyte/wekeep/api/v1/articles"
	consts "github.com/cicbyte/wekeep/internal/consts"
	service "github.com/cicbyte/wekeep/internal/service"
)

var Articles = articlesController{}

type articlesController struct {
	BaseController
}

// Add 新增文章
func (c *articlesController) Add(ctx context.Context, req *api.ArticlesAddReq) (res *api.ArticlesAddRes, err error) {
	res = new(api.ArticlesAddRes)
	err = service.Articles().Add(ctx, req)
	// TODO: 返回新增文章的ID
	return
}

// Edit 编辑文章
func (c *articlesController) Edit(ctx context.Context, req *api.ArticlesEditReq) (res *api.ArticlesEditRes, err error) {
	res = new(api.ArticlesEditRes)
	err = service.Articles().Edit(ctx, req)
	return
}

// Delete 删除文章
func (c *articlesController) Delete(ctx context.Context, req *api.ArticlesDelReq) (res *api.ArticlesDelRes, err error) {
	res = new(api.ArticlesDelRes)
	err = service.Articles().Delete(ctx, req.Id)
	return
}

// BatchDelete 批量删除文章
func (c *articlesController) BatchDelete(ctx context.Context, req *api.ArticlesBatchDelReq) (res *api.ArticlesBatchDelRes, err error) {
	res = new(api.ArticlesBatchDelRes)
	success, failed, failedIds, err := service.Articles().BatchDelete(ctx, req.Ids)
	if err != nil {
		return
	}
	res.Success = success
	res.Failed = failed
	res.FailedIds = failedIds
	return
}

// List 文章列表
func (c *articlesController) List(ctx context.Context, req *api.ArticlesListReq) (res *api.ArticlesListRes, err error) {
	res = new(api.ArticlesListRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	total, list, err := service.Articles().List(ctx, req)
	res.Total = total
	res.CurrentPage = req.PageNum
	res.ArticlesList = list
	return
}

// Detail 文章详情
func (c *articlesController) Detail(ctx context.Context, req *api.ArticlesDetailReq) (res *api.ArticlesDetailRes, err error) {
	res = new(api.ArticlesDetailRes)
	info, err := service.Articles().GetById(ctx, req.Id)
	if err != nil {
		return
	}
	res.ArticlesInfo = info
	return
}

// Tags 标签列表
func (c *articlesController) Tags(ctx context.Context, req *api.ArticlesTagsReq) (res *api.ArticlesTagsRes, err error) {
	res = new(api.ArticlesTagsRes)
	tags, err := service.Articles().GetTags(ctx)
	if err != nil {
		return
	}
	res.Tags = tags
	return
}

// ParseByURL 通过 URL 解析微信文章
func (c *articlesController) ParseByURL(ctx context.Context, req *api.ArticlesParseByURLReq) (res *api.ArticlesParseByURLRes, err error) {
	res = new(api.ArticlesParseByURLRes)
	result, err := service.Articles().ParseByURL(ctx, req)
	if err != nil {
		return
	}
	res.Title = result.Title
	res.Author = result.Author
	res.Content = result.Content
	res.BaseURL = result.BaseURL
	return
}

// Parse 解析微信文章
func (c *articlesController) Parse(ctx context.Context, req *api.ArticlesParseReq) (res *api.ArticlesParseRes, err error) {
	res = new(api.ArticlesParseRes)
	result, err := service.Articles().Parse(ctx, req)
	if err != nil {
		return
	}
	res.Title = result.Title
	res.Author = result.Author
	res.Content = result.Content
	res.BaseURL = result.BaseURL
	return
}

// Stats 统计数据
func (c *articlesController) Stats(ctx context.Context, req *api.ArticlesStatsReq) (res *api.ArticlesStatsRes, err error) {
	res = new(api.ArticlesStatsRes)
	stats, err := service.Articles().GetStats(ctx)
	if err != nil {
		return
	}
	res.TotalCount = stats.TotalCount
	res.AuthorStats = stats.AuthorStats
	res.TagStats = stats.TagStats
	res.RecentArticles = stats.RecentArticles
	res.DateTrend = stats.DateTrend
	return
}

// Reparse 重新解析文章
func (c *articlesController) Reparse(ctx context.Context, req *api.ArticlesReparseReq) (res *api.ArticlesReparseRes, err error) {
	res = new(api.ArticlesReparseRes)
	result, err := service.Articles().Reparse(ctx, req)
	if err != nil {
		return
	}
	res.Title = result.Title
	res.Author = result.Author
	res.Content = result.Content
	res.ImagesCount = result.ImagesCount
	return
}
