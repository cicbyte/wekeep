package service

import (
	"context"
	api "github.com/cicbyte/wekeep/api/v1/articles"
	model "github.com/cicbyte/wekeep/internal/model"
)

type IArticles interface {
	// Add 新增文章
	Add(ctx context.Context, req *api.ArticlesAddReq) (err error)

	// Edit 编辑文章
	Edit(ctx context.Context, req *api.ArticlesEditReq) (err error)

	// Delete 删除文章
	Delete(ctx context.Context, id int) (err error)

	// BatchDelete 批量删除文章
	BatchDelete(ctx context.Context, ids []int) (success int, failed int, failedIds []int, err error)

	// List 获取文章列表（支持分页、筛选、搜索、排序）
	List(ctx context.Context, req *api.ArticlesListReq) (total interface{}, res []*model.ArticlesInfo, err error)

	// GetById 获取文章详情
	GetById(ctx context.Context, id int) (res *model.ArticlesInfo, err error)

	// GetTags 获取所有标签（去重，按使用频率排序）
	GetTags(ctx context.Context) (res []string, err error)

	// GetStats 获取综合统计数据
	GetStats(ctx context.Context) (res *model.ArticlesStatsInfo, err error)

	// ParseByURL 通过 URL 解析微信文章
	ParseByURL(ctx context.Context, req *api.ArticlesParseByURLReq) (res *api.ArticlesParseByURLRes, err error)

	// Parse 解析微信文章
	Parse(ctx context.Context, req *api.ArticlesParseReq) (res *api.ArticlesParseRes, err error)

	// GetAuthorStats 获取作者统计（内部使用）
	GetAuthorStats(ctx context.Context) (res []*model.AuthorStat, err error)

	// GetTagStats 获取标签统计（内部使用）
	GetTagStats(ctx context.Context) (res []*model.TagStat, err error)

	// GetDateTrend 获取时间趋势（内部使用）
	GetDateTrend(ctx context.Context) (res []*model.DateTrend, err error)

	// Reparse 重新解析文章（从原始URL重新获取内容）
	Reparse(ctx context.Context, req *api.ArticlesReparseReq) (res *api.ArticlesReparseRes, err error)
}

var localArticles IArticles

// Articles 返回文章管理服务的实例
func Articles() IArticles {
	if localArticles == nil {
		panic("implement not found for interface IArticles, forgot register?")
	}
	return localArticles
}

// RegisterArticles 注册文章管理服务实现
func RegisterArticles(i IArticles) {
	localArticles = i
}
