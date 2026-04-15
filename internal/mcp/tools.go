package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	api "github.com/ciclebyte/wekeep/api/v1/articles"
	"github.com/ciclebyte/wekeep/internal/consts"
	parser "github.com/ciclebyte/wekeep/internal/parser"
	service "github.com/ciclebyte/wekeep/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerTools 注册所有微信公众号相关工具到 MCP 服务器
func registerTools(s *server.MCPServer) {
	// 1. Echo 工具 - 测试连通性
	s.AddTool(mcp.NewTool("echo",
		mcp.WithDescription("回显用户输入的消息，用于测试 MCP 连接"),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("要回显的消息内容"),
		),
	), handleEcho)

	// 2. GetServerInfo 工具 - 获取服务器信息
	s.AddTool(mcp.NewTool("get_server_info",
		mcp.WithDescription("获取 MCP 服务器的运行信息，包括版本、运行时间、系统信息等"),
	), handleGetServerInfo)

	// 3. wechat_parse_url - 从URL解析微信公众号文章
	s.AddTool(mcp.NewTool("wechat_parse_url",
		mcp.WithDescription("从微信公众号文章URL解析文章内容（标题、作者、正文等），不保存到数据库"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("微信公众号文章URL，如：https://mp.weixin.qq.com/s/xxxxx"),
		),
	), handleWechatParseURL)

	// 4. wechat_save_article - 保存微信文章到数据库
	s.AddTool(mcp.NewTool("wechat_save_article",
		mcp.WithDescription("从微信公众号文章URL解析并保存文章到数据库。如果文章已存在（URL重复）则返回错误"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("微信公众号文章URL"),
		),
		mcp.WithArray("tags",
			mcp.Description("文章标签列表，如：[\"技术\", \"Go\"]"),
		),
	), handleWechatSaveArticle)

	// 5. wechat_list_articles - 查询文章列表
	s.AddTool(mcp.NewTool("wechat_list_articles",
		mcp.WithDescription("查询已收藏的微信公众号文章列表，支持分页、按作者筛选、按标签筛选"),
		mcp.WithNumber("page_num",
			mcp.Description("页码，从1开始，默认1"),
		),
		mcp.WithNumber("page_size",
			mcp.Description("每页数量，默认20，最大100"),
		),
		mcp.WithNumber("author_id",
			mcp.Description("作者ID，用于筛选特定作者的文章"),
		),
		mcp.WithArray("tags",
			mcp.Description("标签列表，用于筛选包含特定标签的文章"),
		),
		mcp.WithString("order_by",
			mcp.Description("排序方式，如：date_added desc（按收藏时间倒序）、date_added asc（正序）"),
		),
	), handleWechatListArticles)

	// 6. wechat_get_article - 获取文章详情
	s.AddTool(mcp.NewTool("wechat_get_article",
		mcp.WithDescription("根据文章ID获取已收藏的微信公众号文章完整内容"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("文章ID"),
		),
	), handleWechatGetArticle)

	// 7. wechat_search_articles - 搜索文章
	s.AddTool(mcp.NewTool("wechat_search_articles",
		mcp.WithDescription("搜索已收藏的微信公众号文章，支持标题和全文搜索"),
		mcp.WithString("keyword",
			mcp.Required(),
			mcp.Description("搜索关键词"),
		),
		mcp.WithBoolean("title_only",
			mcp.Description("是否只搜索标题，默认false（全文搜索）"),
		),
		mcp.WithNumber("page_num",
			mcp.Description("页码，从1开始，默认1"),
		),
		mcp.WithNumber("page_size",
			mcp.Description("每页数量，默认20，最大100"),
		),
	), handleWechatSearchArticles)

	// 8. wechat_get_tags - 获取所有标签
	s.AddTool(mcp.NewTool("wechat_get_tags",
		mcp.WithDescription("获取所有已使用的文章标签，按使用频率排序"),
	), handleWechatGetTags)

	// 9. wechat_get_stats - 获取统计数据
	s.AddTool(mcp.NewTool("wechat_get_stats",
		mcp.WithDescription("获取文章收藏统计数据，包括总数、作者统计、标签统计等"),
	), handleWechatGetStats)

	// 10. wechat_delete_article - 删除文章
	s.AddTool(mcp.NewTool("wechat_delete_article",
		mcp.WithDescription("根据文章ID删除已收藏的微信公众号文章"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("要删除的文章ID"),
		),
	), handleWechatDeleteArticle)
}

// handleEcho 处理 echo 工具调用
func handleEcho(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	message, ok := args["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message 参数必须是字符串")
	}

	return mcp.NewToolResultText(fmt.Sprintf("Echo: %s", message)), nil
}

// handleGetServerInfo 处理获取服务器信息
func handleGetServerInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	info := fmt.Sprintf(
		"WeKeep MCP Server Info:\n"+
			"  Name: WeKeep MCP Server\n"+
			"  Version: 1.0.0\n"+
			"  Description: 微信公众号文章收藏 MCP 服务\n"+
			"  Go Version: %s\n"+
			"  OS: %s\n"+
			"  Architecture: %s\n"+
			"  CPU Cores: %d\n"+
			"  Current Time: %s\n"+
			"\n可用工具:\n"+
			"  - wechat_parse_url: 解析微信文章URL\n"+
			"  - wechat_save_article: 保存微信文章\n"+
			"  - wechat_list_articles: 查询文章列表\n"+
			"  - wechat_get_article: 获取文章详情\n"+
			"  - wechat_search_articles: 搜索文章\n"+
			"  - wechat_get_tags: 获取所有标签\n"+
			"  - wechat_get_stats: 获取统计数据\n"+
			"  - wechat_delete_article: 删除文章",
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return mcp.NewToolResultText(info), nil
}

// handleWechatParseURL 处理从URL解析微信文章
func handleWechatParseURL(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	url, ok := args["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url 参数必须是字符串")
	}

	// 调用解析器
	article, err := parser.ParseFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("解析微信文章失败: %v", err)
	}

	result := fmt.Sprintf(
		"文章解析结果:\n"+
			"  标题: %s\n"+
			"  作者: %s\n"+
			"  原始链接: %s\n"+
			"  发布时间: %s\n"+
			"  内容长度: %d 字符\n\n"+
			"--- 正文预览 (前500字符) ---\n%s",
		article.Title,
		article.Author,
		article.BaseURL,
		article.PublishTime.Format("2006-01-02"),
		len(article.Content),
		truncateString(article.Content, 500),
	)

	return mcp.NewToolResultText(result), nil
}

// handleWechatSaveArticle 处理保存微信文章
func handleWechatSaveArticle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	url, ok := args["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url 参数必须是字符串")
	}

	// 解析文章
	article, err := parser.ParseFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("解析微信文章失败: %v", err)
	}

	// 获取或创建作者
	var authorId int
	if article.Author != "" {
		authorId, err = service.Authors().GetOrCreateByName(ctx, article.Author)
		if err != nil {
			g.Log().Warningf(ctx, "创建作者失败: %v, 使用默认作者", err)
			// 获取或创建默认作者
			authorId, _ = service.Authors().GetOrCreateByName(ctx, "未知作者")
		}
	} else {
		authorId, _ = service.Authors().GetOrCreateByName(ctx, "未知作者")
	}

	// 处理标签
	var tags []string
	if tagsInterface, ok := args["tags"].([]interface{}); ok {
		for _, t := range tagsInterface {
			if tag, ok := t.(string); ok {
				tags = append(tags, tag)
			}
		}
	}

	// 构建保存请求
	addReq := &api.ArticlesAddReq{
		Title:     article.Title,
		AuthorId:  authorId,
		Url:       url,
		Content:   article.Content,
		Tags:      tags,
		DateAdded: time.Now().UnixMilli(),
	}

	// 生成摘要（取前200字符）
	if len(article.Content) > 200 {
		addReq.Summary = article.Content[:200] + "..."
	} else {
		addReq.Summary = article.Content
	}

	// 保存文章
	err = service.Articles().Add(ctx, addReq)
	if err != nil {
		return nil, fmt.Errorf("保存文章失败: %v", err)
	}

	result := fmt.Sprintf(
		"✅ 文章保存成功!\n"+
			"  标题: %s\n"+
			"  作者: %s (ID: %d)\n"+
			"  标签: %v\n"+
			"  内容长度: %d 字符",
		article.Title,
		article.Author,
		authorId,
		tags,
		len(article.Content),
	)

	return mcp.NewToolResultText(result), nil
}

// handleWechatListArticles 处理查询文章列表
func handleWechatListArticles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	// 构建请求
	req := &api.ArticlesListReq{}

	// 设置默认值
	pageNum := 1
	pageSize := 20

	if pn, ok := args["page_num"].(float64); ok && pn > 0 {
		pageNum = int(pn)
	}
	if ps, ok := args["page_size"].(float64); ok && ps > 0 {
		pageSize = int(ps)
		if pageSize > 100 {
			pageSize = 100
		}
	}
	if pageSize == 0 {
		pageSize = consts.PageSize
	}
	if pageNum == 0 {
		pageNum = 1
	}

	// 设置分页参数
	req.PageNum = pageNum
	req.PageSize = pageSize

	if authorId, ok := args["author_id"].(float64); ok && authorId > 0 {
		aid := int(authorId)
		req.AuthorId = &aid
	}

	if orderBy, ok := args["order_by"].(string); ok {
		req.OrderBy = orderBy
	}

	// 处理标签
	if tagsInterface, ok := args["tags"].([]interface{}); ok {
		for _, t := range tagsInterface {
			if tag, ok := t.(string); ok {
				req.Tags = append(req.Tags, tag)
			}
		}
	}

	// 查询
	total, articles, err := service.Articles().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("查询文章列表失败: %v", err)
	}

	// 格式化输出
	var result string
	result += fmt.Sprintf("📚 文章列表 (共 %v 篇，第 %d 页，每页 %d 条)\n\n", total, req.PageNum, req.PageSize)

	if len(articles) == 0 {
		result += "暂无文章"
	} else {
		for i, article := range articles {
			result += fmt.Sprintf("%d. 【ID:%d】%s\n", i+1, article.Id, article.Title)
			if article.AuthorName != "" {
				result += fmt.Sprintf("   作者: %s\n", article.AuthorName)
			}
			if len(article.Tags) > 0 {
				result += fmt.Sprintf("   标签: %v\n", article.Tags)
			}
			if article.Summary != "" {
				result += fmt.Sprintf("   摘要: %s\n", truncateString(article.Summary, 100))
			}
			result += fmt.Sprintf("   收藏时间: %s\n", formatTimestamp(article.DateAdded))
			result += "\n"
		}
	}

	return mcp.NewToolResultText(result), nil
}

// handleWechatGetArticle 处理获取文章详情
func handleWechatGetArticle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	id, ok := args["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("id 参数必须是数字")
	}

	article, err := service.Articles().GetById(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("获取文章失败: %v", err)
	}

	result := fmt.Sprintf(
		"📖 文章详情\n\n"+
			"标题: %s\n"+
			"ID: %d\n",
		article.Title,
		article.Id,
	)

	if article.AuthorName != "" {
		result += fmt.Sprintf("作者: %s (ID: %d)\n", article.AuthorName, *article.AuthorId)
	}

	if article.Url != "" {
		result += fmt.Sprintf("原始链接: %s\n", article.Url)
	}

	if len(article.Tags) > 0 {
		result += fmt.Sprintf("标签: %v\n", article.Tags)
	}

	result += fmt.Sprintf("收藏时间: %s\n", formatTimestamp(article.DateAdded))

	if article.Summary != "" {
		result += fmt.Sprintf("\n📝 摘要:\n%s\n", article.Summary)
	}

	result += fmt.Sprintf("\n📄 正文内容:\n\n%s", article.Content)

	return mcp.NewToolResultText(result), nil
}

// handleWechatSearchArticles 处理搜索文章
func handleWechatSearchArticles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	keyword, ok := args["keyword"].(string)
	if !ok || keyword == "" {
		return nil, fmt.Errorf("keyword 参数不能为空")
	}

	// 构建请求
	req := &api.ArticlesListReq{
		Keyword: keyword,
	}

	// 设置分页参数
	pageNum := 1
	pageSize := 20

	if pn, ok := args["page_num"].(float64); ok && pn > 0 {
		pageNum = int(pn)
	}
	if ps, ok := args["page_size"].(float64); ok && ps > 0 {
		pageSize = int(ps)
		if pageSize > 100 {
			pageSize = 100
		}
	}
	if pageSize == 0 {
		pageSize = consts.PageSize
	}
	if pageNum == 0 {
		pageNum = 1
	}

	req.PageNum = pageNum
	req.PageSize = pageSize

	if titleOnly, ok := args["title_only"].(bool); ok {
		req.TitleOnly = titleOnly
	}

	// 搜索
	total, articles, err := service.Articles().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("搜索文章失败: %v", err)
	}

	// 格式化输出
	var result string
	searchType := "全文搜索"
	if req.TitleOnly {
		searchType = "标题搜索"
	}
	result += fmt.Sprintf("🔍 搜索结果: \"%s\" (%s)\n", keyword, searchType)
	result += fmt.Sprintf("找到 %v 篇相关文章\n\n", total)

	if len(articles) == 0 {
		result += "未找到匹配的文章"
	} else {
		for i, article := range articles {
			result += fmt.Sprintf("%d. 【ID:%d】%s\n", i+1, article.Id, article.Title)
			if article.AuthorName != "" {
				result += fmt.Sprintf("   作者: %s\n", article.AuthorName)
			}
			// 如果有高亮的摘要，显示高亮摘要
			if article.FormattedSummary != "" {
				result += fmt.Sprintf("   匹配摘要: %s\n", truncateString(article.FormattedSummary, 150))
			} else if article.ContextSnippet != "" {
				result += fmt.Sprintf("   匹配内容: %s\n", truncateString(article.ContextSnippet, 150))
			} else if article.Summary != "" {
				result += fmt.Sprintf("   摘要: %s\n", truncateString(article.Summary, 100))
			}
			result += "\n"
		}

		result += fmt.Sprintf("💡 提示: 使用 wechat_get_article 获取完整内容")
	}

	return mcp.NewToolResultText(result), nil
}

// handleWechatGetTags 处理获取所有标签
func handleWechatGetTags(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tags, err := service.Articles().GetTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取标签失败: %v", err)
	}

	result := "🏷️ 所有标签 (按使用频率排序)\n\n"

	if len(tags) == 0 {
		result += "暂无标签"
	} else {
		for i, tag := range tags {
			result += fmt.Sprintf("%d. %s\n", i+1, tag)
		}
	}

	return mcp.NewToolResultText(result), nil
}

// handleWechatGetStats 处理获取统计数据
func handleWechatGetStats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stats, err := service.Articles().GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取统计数据失败: %v", err)
	}

	result := fmt.Sprintf(
		"📊 收藏统计\n\n"+
			"📚 文章总数: %d\n",
		stats.TotalCount,
	)

	if len(stats.AuthorStats) > 0 {
		result += "\n✍️ 作者统计 (Top 10)\n"
		for i, author := range stats.AuthorStats {
			if i >= 10 {
				break
			}
			result += fmt.Sprintf("  %d. %s: %d 篇\n", i+1, author.AuthorName, author.Count)
		}
	}

	if len(stats.TagStats) > 0 {
		result += "\n🏷️ 标签统计 (Top 10)\n"
		for i, tag := range stats.TagStats {
			if i >= 10 {
				break
			}
			result += fmt.Sprintf("  %d. %s: %d 次\n", i+1, tag.Tag, tag.Count)
		}
	}

	return mcp.NewToolResultText(result), nil
}

// handleWechatDeleteArticle 处理删除文章
func handleWechatDeleteArticle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	id, ok := args["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("id 参数必须是数字")
	}

	// 先获取文章信息用于显示
	article, err := service.Articles().GetById(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("文章不存在: %v", err)
	}

	title := article.Title

	// 删除文章
	err = service.Articles().Delete(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("删除文章失败: %v", err)
	}

	result := fmt.Sprintf("✅ 文章已删除\n  ID: %d\n  标题: %s", int(id), title)

	return mcp.NewToolResultText(result), nil
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// formatTimestamp 格式化时间戳（毫秒）
func formatTimestamp(ts int64) string {
	if ts == 0 {
		return "未知"
	}
	return time.UnixMilli(ts).Format("2006-01-02 15:04:05")
}

// marshalJSON 将对象转换为 JSON 字符串（用于调试）
func marshalJSON(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("JSON 序列化失败: %v", err)
	}
	return string(data)
}
