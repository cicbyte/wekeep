package meilisearch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cicbyte/wekeep/internal/model"
	service "github.com/cicbyte/wekeep/internal/service"
	"github.com/cicbyte/wekeep/library/libMeilisearch"
	liberr "github.com/cicbyte/wekeep/library/liberr"
	"github.com/gogf/gf/v2/frame/g"
)

type SearchReq struct {
	g.Meta   `path:"/search" method:"get" tags:"搜索" summary:"全文搜索"`
	Query    string `p:"query" v:"required#搜索关键词不能为空"`
	Limit    int    `p:"limit" d:"10"`
	AuthorId *int   `p:"authorId" dc:"按作者ID筛选"`
}

type SearchRes struct {
	g.Meta `mime:"application/json"`
	Total  int          `json:"total"` // 总结果数
	Hits   []*SearchHit `json:"hits"`  // 搜索结果
}

type SearchHit struct {
	Id        int      `json:"id"`
	Title     string   `json:"title"`
	AuthorId  *int     `json:"authorId"`
	Author    string   `json:"author"` // 作者名称（冗余存储便于搜索）
	Summary   string   `json:"summary"`
	Url       string   `json:"url"`
	Tags      []string `json:"tags"`
	DateAdded int64    `json:"dateAdded"`
	Score     float64  `json:"score"` // 相关性得分
}

// init 注册到 Service 层
func init() {
	// 注册到现有 service
}

// Search 全文搜索
func Search(ctx context.Context, req *SearchReq) (res *SearchRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 初始化 Meilisearch
		err = libMeilisearch.Init(ctx)
		liberr.ErrIsNil(ctx, err, "Meilisearch 初始化失败")

		// 构建过滤条件
		filter := ""
		if req.AuthorId != nil {
			filter = fmt.Sprintf("authorId = %d", *req.AuthorId)
		}

		// 执行搜索
		hits, err := libMeilisearch.Search(ctx, req.Query, int64(req.Limit), filter)
		liberr.ErrIsNil(ctx, err, "搜索失败")

		// 转换结果
		var searchHits []*SearchHit
		for _, hit := range hits {
			hitMap, ok := hit.(map[string]interface{})
			if !ok {
				continue
			}

			// 解析 tags
			var tags []string
			if tagsVal, ok := hitMap["tags"]; ok {
				if tagsStr, ok := tagsVal.(string); ok {
					json.Unmarshal([]byte(tagsStr), &tags)
				}
			}

			searchHit := &SearchHit{
				Score: hitMap["_score"].(float64),
			}

			// 解析字段
			if idVal, ok := hitMap["id"]; ok {
				switch v := idVal.(type) {
				case float64:
					searchHit.Id = int(v)
				case int:
					searchHit.Id = v
				}
			}
			if titleVal, ok := hitMap["title"]; ok {
				searchHit.Title = titleVal.(string)
			}
			if authorIdVal, ok := hitMap["authorId"]; ok {
				switch v := authorIdVal.(type) {
				case float64:
					aid := int(v)
					searchHit.AuthorId = &aid
				case int:
					aid := v
					searchHit.AuthorId = &aid
				}
			}
			if authorVal, ok := hitMap["author"]; ok {
				searchHit.Author = authorVal.(string)
			}
			if summaryVal, ok := hitMap["summary"]; ok {
				searchHit.Summary = summaryVal.(string)
			}
			if urlVal, ok := hitMap["url"]; ok {
				searchHit.Url = urlVal.(string)
			}
			if dateAddedVal, ok := hitMap["dateAdded"]; ok {
				searchHit.DateAdded = int64(dateAddedVal.(float64))
			}
			searchHit.Tags = tags

			searchHits = append(searchHits, searchHit)
		}

		res = new(SearchRes)
		res.Total = len(searchHits)
		res.Hits = searchHits

		g.Log().Infof(ctx, "搜索完成: query=%s, results=%d", req.Query, res.Total)
	})
	return
}

// IndexArticle 索引文章到 Meilisearch
func IndexArticle(ctx context.Context, article *model.ArticlesInfo) error {
	// 获取作者名称
	authorName := ""
	if article.AuthorId != nil {
		author, err := service.Authors().GetById(ctx, int(*article.AuthorId))
		if err == nil {
			authorName = author.Name
		}
	}

	document := map[string]interface{}{
		"id":        article.Id,
		"title":     article.Title,
		"authorId":  article.AuthorId,
		"author":    authorName, // 冗余存储作者名称，便于搜索
		"summary":   article.Summary,
		"content":   article.Content,
		"url":       article.Url,
		"tags":      g.NewVar(article.Tags).String(),
		"dateAdded": article.DateAdded,
	}

	return libMeilisearch.IndexDocument(ctx, document)
}
