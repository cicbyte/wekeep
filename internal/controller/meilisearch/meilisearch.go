package meilisearch

import (
	"context"

	"github.com/cicbyte/wekeep/internal/logic/meilisearch"
	"github.com/gogf/gf/v2/frame/g"
)

var Search = meilisearchController{}

type meilisearchController struct{}

// Search 全文搜索接口
func (c *meilisearchController) Search(ctx context.Context, req *meilisearch.SearchReq) (res *meilisearch.SearchRes, err error) {
	res = new(meilisearch.SearchRes)
	res, err = meilisearch.Search(ctx, req)
	return
}

// StatusReq 状态请求
type StatusReq struct {
	g.Meta `path:"/search/status" method:"get" tags:"搜索" summary:"获取搜索引擎状态"`
}

type StatusRes struct {
	g.Meta       `mime:"application/json"`
	Enabled      bool  `json:"enabled"`      // 是否启用全文搜索
	IndexedCount int64 `json:"indexedCount"` // 已索引文章数
}

// Status 获取搜索引擎状态
func (c *meilisearchController) Status(ctx context.Context, req *StatusReq) (res *StatusRes, err error) {
	res = new(StatusRes)
	enabled, count, err := meilisearch.GetSearchStatus(ctx)
	if err != nil {
		return
	}
	res.Enabled = enabled
	res.IndexedCount = count
	return
}

// IndexAllReq 索引请求（空结构体）
type IndexAllReq struct {
	g.Meta `path:"/search/indexAll" method:"post" tags:"搜索" summary:"索引所有文章"`
}

type IndexAllRes struct {
	g.Meta      `mime:"application/json"`
	Msg         string `json:"msg"`
	IndexedCount int   `json:"indexedCount"` // 成功索引的文章数
}

// IndexAll 索引所有文章（管理接口）
func (c *meilisearchController) IndexAll(ctx context.Context, req *IndexAllReq) (res *IndexAllRes, err error) {
	res = new(IndexAllRes)
	count, err := meilisearch.IndexNewArticlesWithCount(ctx)
	if err != nil {
		return
	}
	res.Msg = "索引完成"
	res.IndexedCount = count
	return
}
