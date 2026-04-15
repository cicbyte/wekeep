package libMeilisearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/meilisearch/meilisearch-go"
)

var (
	client      meilisearch.ServiceManager
	initialized bool
	enabled     bool // 是否启用全文搜索
)

// IsSearchEnabled 检查是否启用全文搜索
func IsSearchEnabled(ctx context.Context) bool {
	return g.Cfg().MustGet(ctx, "search.enabled").Bool()
}

// GetMeilisearchAddress 获取 Meilisearch 地址
func GetMeilisearchAddress(ctx context.Context) string {
	address := g.Cfg().MustGet(ctx, "search.meilisearch.address").String()
	if address == "" {
		address = "http://localhost:7700"
	}
	return address
}

// Init 初始化 Meilisearch 客户端
func Init(ctx context.Context) error {
	if initialized {
		return nil
	}

	// 检查是否启用全文搜索
	if !IsSearchEnabled(ctx) {
		g.Log().Info(ctx, "全文搜索未启用，使用数据库查询")
		initialized = true
		enabled = false
		return nil
	}

	// 从配置读取 Meilisearch 地址
	address := GetMeilisearchAddress(ctx)

	// 创建 Meilisearch 客户端
	client = meilisearch.New(address)

	// 验证连接
	_, err := client.ListIndexes(&meilisearch.IndexesQuery{})
	if err != nil {
		return fmt.Errorf("Meilisearch 连接失败: %w", err)
	}

	// 创建索引
	indexName := "articles"
	primaryKey := "id"

	// 检查索引是否存在，不存在则创建
	_, err = client.GetIndex(indexName)
	if err != nil {
		// 索引不存在，创建新索引
		_, err = client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        indexName,
			PrimaryKey: primaryKey,
		})
		if err != nil {
			return fmt.Errorf("创建 Meilisearch 索引失败: %w", err)
		}

		g.Log().Infof(ctx, "创建索引: %s", indexName)
	}

	// 总是更新索引设置（确保配置正确）
	index := client.Index(indexName)
	_, err = index.UpdateSettings(&meilisearch.Settings{
		SearchableAttributes: []string{"title", "author", "summary", "content"},
		SortableAttributes:   []string{"date_added", "created_at"},
		DisplayedAttributes:  []string{"id", "title", "author", "summary", "content", "tags", "date_added"},
		FilterableAttributes: []string{"author", "date_added"},
	})
	if err != nil {
		return fmt.Errorf("配置索引字段失败: %w", err)
	}
	g.Log().Infof(ctx, "配置索引字段成功")

	initialized = true
	enabled = true
	g.Log().Info(ctx, "Meilisearch 客户端初始化成功，全文搜索已启用")
	return nil
}

// GetClient 获取 Meilisearch 客户端
func GetClient() (meilisearch.ServiceManager, error) {
	if !initialized {
		return nil, fmt.Errorf("Meilisearch 客户端未初始化")
	}
	if !enabled || client == nil {
		return nil, fmt.Errorf("全文搜索未启用")
	}
	return client, nil
}

// IndexDocument 索引文档
func IndexDocument(ctx context.Context, document interface{}) error {
	cli, err := GetClient()
	if err != nil {
		return err
	}

	index := cli.Index("articles")
	_, err = index.AddDocuments(document, nil)
	if err != nil {
		return fmt.Errorf("索引文档失败: %w", err)
	}

	g.Log().Infof(ctx, "文档索引成功")
	return nil
}

// DeleteDocument 删除文档
func DeleteDocument(ctx context.Context, documentId int) error {
	cli, err := GetClient()
	if err != nil {
		return err
	}

	index := cli.Index("articles")
	_, err = index.DeleteDocument(strconv.Itoa(documentId), nil)
	if err != nil {
		return fmt.Errorf("删除文档失败: %w", err)
	}

	g.Log().Infof(ctx, "文档删除成功: ID=%d", documentId)
	return nil
}

// SearchResult 搜索结果结构
type SearchResult struct {
	Id             int                    `json:"id"`
	Title          string                 `json:"title"`
	Author         string                 `json:"author"`
	Summary        string                 `json:"summary"`
	FormattedTitle string                 `json:"formattedTitle"`    // 高亮后的标题
	FormattedSummary string               `json:"formattedSummary"`  // 高亮后的摘要
	ContextSnippet string                 `json:"contextSnippet"`    // 正文匹配上下文
	MatchFields    []string               `json:"matchFields"`       // 匹配的字段
	Score          float64                `json:"score"`
}

// SearchWithHighlight 搜索文档（带高亮和上下文）
func SearchWithHighlight(ctx context.Context, query string, limit int64, filter string) ([]*SearchResult, error) {
	cli, err := GetClient()
	if err != nil {
		return nil, err
	}

	index := cli.Index("articles")

	// 构建搜索请求，启用高亮和裁剪
	searchReq := &meilisearch.SearchRequest{
		Query:                query,
		Limit:                limit,
		AttributesToHighlight: []string{"title", "summary", "content"},
		AttributesToCrop:     []string{"content"},
		CropLength:           30,  // 上下文长度
		CropMarker:           "...",
		HighlightPreTag:      "<mark>",
		HighlightPostTag:     "</mark>",
	}

	if filter != "" {
		searchReq.Filter = filter
	}

	// 执行搜索
	searchRes, err := index.Search(query, searchReq)
	if err != nil {
		return nil, fmt.Errorf("搜索失败: %w", err)
	}

	g.Log().Infof(ctx, "搜索成功: query=%s, hits=%d", query, searchRes.EstimatedTotalHits)

	// 解析结果
	results := make([]*SearchResult, 0, len(searchRes.Hits))
	for i, hit := range searchRes.Hits {
		result := &SearchResult{}

		// 调试：打印第一个 hit 的所有字段
		if i == 0 {
			g.Log().Infof(ctx, "First hit keys: %v", func() []string {
				keys := make([]string, 0, len(hit))
				for k := range hit {
					keys = append(keys, k)
				}
				return keys
			}())
		}

		// 解析原始字段
		if rawVal, ok := hit["id"]; ok {
			var id float64
			if err := json.Unmarshal(rawVal, &id); err == nil {
				result.Id = int(id)
			}
		}
		if rawVal, ok := hit["title"]; ok {
			json.Unmarshal(rawVal, &result.Title)
		}
		if rawVal, ok := hit["author"]; ok {
			json.Unmarshal(rawVal, &result.Author)
		}
		if rawVal, ok := hit["summary"]; ok {
			json.Unmarshal(rawVal, &result.Summary)
		}
		if rawVal, ok := hit["_score"]; ok {
			json.Unmarshal(rawVal, &result.Score)
		}

		// 解析格式化后的字段（来自 _formatted）
		if formatted, ok := hit["_formatted"]; ok {
			var formattedMap map[string]interface{}
			if err := json.Unmarshal(formatted, &formattedMap); err == nil {
				g.Log().Debugf(ctx, "Formatted data: %+v", formattedMap)
				if title, ok := formattedMap["title"].(string); ok {
					result.FormattedTitle = title
				}
				if summary, ok := formattedMap["summary"].(string); ok {
					result.FormattedSummary = summary
				}
				// 正文裁剪后的上下文
				if content, ok := formattedMap["content"].(string); ok {
					result.ContextSnippet = content
				}
			}
		}

		// 解析匹配的字段信息（来自 _matchesPosition）
		if matches, ok := hit["_matchesPosition"]; ok {
			var matchesMap map[string]interface{}
			if err := json.Unmarshal(matches, &matchesMap); err == nil {
				for field := range matchesMap {
					result.MatchFields = append(result.MatchFields, field)
				}
			}
		}

		results = append(results, result)

		// 调试：打印第一个结果的完整信息
		if i == 0 {
			g.Log().Infof(ctx, "SearchResult: id=%d, formattedTitle=%s, formattedSummary=%s, contextSnippet=%s",
				result.Id, result.FormattedTitle, result.FormattedSummary, result.ContextSnippet)
		}
	}

	return results, nil
}

// Search 搜索文档（兼容旧接口）
func Search(ctx context.Context, query string, limit int64, filter string) ([]interface{}, error) {
	results, err := SearchWithHighlight(ctx, query, limit, filter)
	if err != nil {
		return nil, err
	}

	// 转换为 []interface{} 兼容旧接口
	hits := make([]interface{}, len(results))
	for i, r := range results {
		hit := map[string]interface{}{
			"id":               r.Id,
			"title":            r.Title,
			"author":           r.Author,
			"summary":          r.Summary,
			"formattedTitle":   r.FormattedTitle,
			"formattedSummary": r.FormattedSummary,
			"contextSnippet":   r.ContextSnippet,
			"matchFields":      r.MatchFields,
			"_score":           r.Score,
		}
		hits[i] = hit
		// 调试：打印第一个结果的 map
		if i == 0 {
			g.Log().Infof(ctx, "Search hit[0]: formattedTitle=%s", r.FormattedTitle)
		}
	}
	return hits, nil
}

// UpdateDocument 更新文档
func UpdateDocument(ctx context.Context, documentId int, document interface{}) error {
	cli, err := GetClient()
	if err != nil {
		return err
	}

	index := cli.Index("articles")
	_, err = index.UpdateDocuments(document, nil)
	if err != nil {
		return fmt.Errorf("更新文档失败: %w", err)
	}

	g.Log().Infof(ctx, "文档更新成功: ID=%d", documentId)
	return nil
}

// IsEnabled 检查 Meilisearch 是否启用（已初始化且配置开启）
func IsEnabled() bool {
	return initialized && enabled && client != nil
}

// IsConfigEnabled 检查配置是否启用全文搜索（不需要初始化）
func IsConfigEnabled(ctx context.Context) bool {
	return IsSearchEnabled(ctx)
}

// GetIndexedCount 获取索引文档数量
func GetIndexedCount(ctx context.Context) (int64, error) {
	cli, err := GetClient()
	if err != nil {
		return 0, err
	}

	index := cli.Index("articles")
	stats, err := index.GetStats()
	if err != nil {
		return 0, fmt.Errorf("获取索引统计失败: %w", err)
	}

	return int64(stats.NumberOfDocuments), nil
}
