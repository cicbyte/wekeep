package meilisearch

import (
	"github.com/gogf/gf/v2/frame/g"
)

type SearchReq struct {
	g.Meta  `path:"/search" method:"get" tags:"搜索" summary:"全文搜索"`
	Query  string `p:"query" v:"required#搜索关键词不能为空"`
	Limit  int    `p:"limit" d:"10"`
	Author string `p:"author"`
}

type SearchRes struct {
	g.Meta `mime:"application/json"`
	Total int           `json:"total"`
	Hits  []*SearchHit `json:"hits"`
}

type SearchHit struct {
	Id        int      `json:"id"`
	Title     string   `json:"title"`
	Author    string   `json:"author"`
	Summary   string   `json:"summary"`
	Url       string   `json:"url"`
	Tags      []string `json:"tags"`
	DateAdded int64     `json:"dateAdded"`
	Score     float64  `json:"score"`
}
