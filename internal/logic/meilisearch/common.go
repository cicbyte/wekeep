package meilisearch

import (
	"context"

	articlesApi "github.com/ciclebyte/wekeep/api/v1/articles"
	service "github.com/ciclebyte/wekeep/internal/service"
	"github.com/ciclebyte/wekeep/library/libMeilisearch"
	liberr "github.com/ciclebyte/wekeep/library/liberr"
	"github.com/gogf/gf/v2/frame/g"
)

type sMeilisearch struct{}

func New() *sMeilisearch {
	return &sMeilisearch{}
}

func init() {
	// 暂时不注册独立service，逻辑整合到 articles 中
}

// GetSearchStatus 获取搜索引擎状态
func GetSearchStatus(ctx context.Context) (enabled bool, indexedCount int64, err error) {
	// 检查配置是否启用
	enabled = libMeilisearch.IsConfigEnabled(ctx)
	if !enabled {
		return false, 0, nil
	}

	// 尝试初始化并获取索引统计
	err = libMeilisearch.Init(ctx)
	if err != nil {
		g.Log().Warningf(ctx, "Meilisearch 初始化失败: %v", err)
		return true, 0, nil // 配置启用但初始化失败
	}

	// 获取索引文档数量
	indexedCount, err = libMeilisearch.GetIndexedCount(ctx)
	if err != nil {
		g.Log().Warningf(ctx, "获取索引数量失败: %v", err)
		return true, 0, nil
	}

	return enabled, indexedCount, nil
}

// IndexNewArticles 索引所有现有文章到 Meilisearch
func IndexNewArticles(ctx context.Context) error {
	_, err := IndexNewArticlesWithCount(ctx)
	return err
}

// IndexNewArticlesWithCount 索引所有现有文章并返回成功数量
func IndexNewArticlesWithCount(ctx context.Context) (int, error) {
	var successCount int
	err := g.Try(ctx, func(ctx context.Context) {
		// 初始化 Meilisearch
		err := libMeilisearch.Init(ctx)
		liberr.ErrIsNil(ctx, err, "Meilisearch 初始化失败")

		// 获取所有文章 - List 返回
		total, articlesList, err := service.Articles().List(ctx, &articlesApi.ArticlesListReq{})
		liberr.ErrIsNil(ctx, err, "获取文章列表失败")

		for _, article := range articlesList {
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
				"tags":      article.Tags,
				"dateAdded": article.DateAdded,
			}
			// 调试：检查 content 是否为空
			if article.Content == "" {
				g.Log().Warningf(ctx, "文章 ID=%d content 为空", article.Id)
			}
			err = libMeilisearch.IndexDocument(ctx, document)
			if err != nil {
				g.Log().Warningf(ctx, "索引文章失败: ID=%d, err=%v", article.Id, err)
			} else {
				successCount++
			}
		}

		g.Log().Infof(ctx, "批量索引完成，共索引 %d/%v 篇文章", successCount, total)
	})
	return successCount, err
}
