package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	api "github.com/cicbyte/wekeep/api/v1/articles"
	dao "github.com/cicbyte/wekeep/internal/dao"
	model "github.com/cicbyte/wekeep/internal/model"
	do "github.com/cicbyte/wekeep/internal/model/do"
	entity "github.com/cicbyte/wekeep/internal/model/entity"
	parser "github.com/cicbyte/wekeep/internal/parser"
	service "github.com/cicbyte/wekeep/internal/service"
	liberr "github.com/cicbyte/wekeep/library/liberr"
	libMeilisearch "github.com/cicbyte/wekeep/library/libMeilisearch"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	service.RegisterArticles(New())
}

func New() *sArticles {
	return &sArticles{}
}

type sArticles struct{}

// intToUint 将 int 转换为 uint
func intToUint(v int) uint {
	return uint(v)
}

// uintToUintPtr 将 uint 转换为 *uint
func uintToUintPtr(v uint) *uint {
	return &v
}

// Add 新增文章
func (s sArticles) Add(ctx context.Context, req *api.ArticlesAddReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		g.Log().Info(ctx, "开始添加文章", "title", req.Title, "authorId", req.AuthorId)

		// 验证作者是否存在
		_, err := service.Authors().GetById(ctx, req.AuthorId)
		liberr.ErrIsNil(ctx, err, "作者不存在")
		g.Log().Info(ctx, "作者验证通过")

		// 将 Tags 数组转换为 JSON 字符串
		// 注意：数据库字段是 JSON 类型，空数组用 "[]"，不能用空字符串或 null
		tagsJson := "[]"
		if len(req.Tags) > 0 {
			bytes, _ := json.Marshal(req.Tags)
			tagsJson = string(bytes)
		}

		authorId := intToUint(req.AuthorId)
		g.Log().Info(ctx, "准备插入文章", "authorId", authorId, "title", req.Title)

		// 检查 URL 是否已存在（重复收藏检查）- 在插入前先检查
		if req.Url != "" {
			existingArticle, _ := dao.Articles.Ctx(ctx).Where("url = ?", req.Url).One()
			if existingArticle != nil {
				g.Log().Warningf(ctx, "URL already exists: %s", req.Url)
				liberr.ErrIsNil(ctx, fmt.Errorf("该文章已经收藏过了"), "")
				return // 不应该到达这里，但为了安全
			}
		}

		// 先插入文章
		result, err := dao.Articles.Ctx(ctx).Insert(do.Articles{
			Title:     req.Title,
			AuthorId:  authorId,
			Url:       req.Url,
			Summary:   req.Summary,
			Content:   req.Content,
			Tags:      tagsJson,
			DateAdded: req.DateAdded,
		})
		liberr.ErrIsNil(ctx, err, "新增文章失败")

		// 获取插入的文章ID
		articleId, _ := result.LastInsertId()
		g.Log().Info(ctx, "文章插入成功", "articleId", articleId)

		// 更新作者的文章数量
		g.Log().Info(ctx, "准备更新作者文章数量", "authorId", req.AuthorId)
		err = service.Authors().UpdateArticleCount(ctx, req.AuthorId)
		liberr.ErrIsNil(ctx, err, "更新作者文章数量失败")
		g.Log().Info(ctx, "作者文章数量更新成功")

		// 异步处理文章中的微信图片（本地化存储）
		if req.Content != "" {
			originalContent := req.Content
			go func() {
				bgCtx := context.Background()
				newContent, imageIds, processErr := service.Images().ProcessArticleImages(bgCtx, int(articleId), originalContent)
				if processErr != nil {
					g.Log().Warningf(bgCtx, "Failed to process images for article %d: %v", articleId, processErr)
					return
				}

				// 如果图片URL有变化，更新文章内容
				if newContent != originalContent {
					_, updateErr := dao.Articles.Ctx(bgCtx).WherePri(articleId).Update(do.Articles{
						Content: newContent,
					})
					if updateErr != nil {
						g.Log().Warningf(bgCtx, "Failed to update article %d content: %v", articleId, updateErr)
					} else {
						g.Log().Infof(bgCtx, "Article %d content updated with %d local images", articleId, len(imageIds))
					}
				}
			}()
		}
	})
	return
}

// Edit 编辑文章
func (s sArticles) Edit(ctx context.Context, req *api.ArticlesEditReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 验证文章是否存在
		oldArticle, err := s.GetById(ctx, req.Id)
		liberr.ErrIsNil(ctx, err, "文章不存在")

		// 验证新作者是否存在
		_, err = service.Authors().GetById(ctx, req.AuthorId)
		liberr.ErrIsNil(ctx, err, "作者不存在")

		// 将 Tags 数组转换为 JSON 字符串
		// 注意：数据库字段是 JSON 类型，空数组用 "[]"，不能用空字符串
		tagsJson := "[]"
		if len(req.Tags) > 0 {
			bytes, _ := json.Marshal(req.Tags)
			tagsJson = string(bytes)
		}

		authorId := intToUint(req.AuthorId)

		// 处理文章中的微信图片（异步处理）
		processedContent := req.Content
		if req.Content != "" && req.Content != oldArticle.Content {
			go func() {
				bgCtx := context.Background()
				newContent, _, processErr := service.Images().ProcessArticleImages(bgCtx, req.Id, req.Content)
				if processErr != nil {
					g.Log().Warningf(bgCtx, "Failed to process images for article %d: %v", req.Id, processErr)
				} else if newContent != req.Content {
					// 更新文章内容（图片URL已替换）
					_, updateErr := dao.Articles.Ctx(bgCtx).WherePri(req.Id).Update(do.Articles{
						Content: newContent,
					})
					if updateErr != nil {
						g.Log().Warningf(bgCtx, "Failed to update article content with new image URLs: %v", updateErr)
					}
				}
			}()
		}

		_, err = dao.Articles.Ctx(ctx).WherePri(req.Id).Update(do.Articles{
			Title:     req.Title,
			AuthorId:  authorId,
			Url:       req.Url,
			Summary:   req.Summary,
			Content:   processedContent,
			Tags:      tagsJson,
			DateAdded: req.DateAdded,
		})
		liberr.ErrIsNil(ctx, err, "修改文章失败")

		// 更新旧作者的文章数量
		if oldArticle.AuthorId != nil {
			err = service.Authors().UpdateArticleCount(ctx, int(*oldArticle.AuthorId))
			liberr.ErrIsNil(ctx, err, "更新旧作者文章数量失败")
		}

		// 更新新作者的文章数量（如果作者变更）
		if oldArticle.AuthorId == nil || *oldArticle.AuthorId != authorId {
			err = service.Authors().UpdateArticleCount(ctx, req.AuthorId)
			liberr.ErrIsNil(ctx, err, "更新新作者文章数量失败")
		}
	})
	return
}

// Delete 删除文章
func (s sArticles) Delete(ctx context.Context, id int) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 获取文章信息以更新作者文章数量
		article, err := s.GetById(ctx, id)
		liberr.ErrIsNil(ctx, err, "文章不存在")

		// 移除文章关联的图片（引用计数减1）
		removeErr := service.Images().RemoveArticleImages(ctx, id)
		if removeErr != nil {
			g.Log().Warningf(ctx, "Failed to remove article images: %v", removeErr)
			// 不阻断删除流程
		}

		_, err = dao.Articles.Ctx(ctx).WherePri(id).Delete()
		liberr.ErrIsNil(ctx, err, "删除文章失败")

		// 更新作者的文章数量
		if article.AuthorId != nil {
			err = service.Authors().UpdateArticleCount(ctx, int(*article.AuthorId))
			liberr.ErrIsNil(ctx, err, "更新作者文章数量失败")
		}
	})
	return
}

// BatchDelete 批量删除文章
func (s sArticles) BatchDelete(ctx context.Context, ids []int) (success int, failed int, failedIds []int, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		if len(ids) == 0 {
			liberr.ErrIsNil(ctx, fmt.Errorf("请选择要删除的文章"), "批量删除失败")
		}

		// 收集需要更新的作者ID
		authorIdMap := make(map[uint]bool)
		for _, id := range ids {
			article, err := s.GetById(ctx, id)
			if err == nil && article.AuthorId != nil {
				authorIdMap[*article.AuthorId] = true
			}

			// 移除文章关联的图片（引用计数减1）
			removeErr := service.Images().RemoveArticleImages(ctx, id)
			if removeErr != nil {
				g.Log().Warningf(ctx, "Failed to remove images for article %d: %v", id, removeErr)
			}
		}

		// 分批删除（每批100个）
		batchSize := 100
		for i := 0; i < len(ids); i += batchSize {
			end := i + batchSize
			if end > len(ids) {
				end = len(ids)
			}
			batch := ids[i:end]

			_, err = dao.Articles.Ctx(ctx).WherePri(batch).Delete()
			if err != nil {
				failed += len(batch)
				failedIds = append(failedIds, batch...)
			} else {
				success += len(batch)
			}
		}

		// 更新所有受影响作者的文章数量
		for authorId := range authorIdMap {
			_ = service.Authors().UpdateArticleCount(ctx, int(authorId))
		}
	})
	return
}

// List 获取文章列表
func (s sArticles) List(ctx context.Context, req *api.ArticlesListReq) (total interface{}, articlesList []*model.ArticlesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.Articles.Ctx(ctx)
		columns := dao.Articles.Columns()

		// 按作者ID筛选
		if req.AuthorId != nil {
			m = m.Where(fmt.Sprintf("%s = ?", columns.AuthorId), *req.AuthorId)
		}

		// 按标签筛选（JSON_CONTAINS）
		if len(req.Tags) > 0 {
			for _, tag := range req.Tags {
				m = m.Where(fmt.Sprintf("JSON_CONTAINS(%s, '\"%s\"')", columns.Tags, tag))
			}
		}

		// 关键词搜索
		var searchHighlightMap map[int]map[string]interface{} // 搜索高亮信息
		if req.Keyword != "" {
			// 如果指定了 TitleOnly，或者未启用全文搜索，则只搜索标题
			if req.TitleOnly || !libMeilisearch.IsConfigEnabled(ctx) {
				// 只搜索标题
				keyword := "%" + req.Keyword + "%"
				m = m.Where(fmt.Sprintf("%s LIKE ?", columns.Title), keyword)
			} else {
				// 使用 Meilisearch 全文搜索
				articleIds, highlightMap := s.searchWithMeilisearch(ctx, req.Keyword)
				searchHighlightMap = highlightMap
				if len(articleIds) > 0 {
					m = m.WherePri(articleIds)
				} else {
					// 没有搜索结果，返回空
					total = 0
					articlesList = []*model.ArticlesInfo{}
					return
				}
			}
		}

		// 排序
		orderBy := req.OrderBy
		if orderBy == "" {
			orderBy = "date_added desc"
		}

		// 使用 Safe() 确保链式调用正确
		m = m.Safe()

		total, err = m.Count()
		liberr.ErrIsNil(ctx, err, "获取文章列表失败")

		// 使用 entity 扫描，然后转换
		var entities []*entity.Articles
		err = m.Page(req.PageNum, req.PageSize).Order(orderBy).Scan(&entities)
		liberr.ErrIsNil(ctx, err, "获取文章列表失败")

		// 收集所有作者ID
		authorIds := make([]uint, 0)
		for _, e := range entities {
			if e.AuthorId > 0 {
				authorIds = append(authorIds, e.AuthorId)
			}
		}

		// 批量获取作者信息
		authorMap := make(map[uint]string)
		if len(authorIds) > 0 {
			var authorEntities []*entity.Authors
			err = dao.Authors.Ctx(ctx).WherePri(authorIds).Scan(&authorEntities)
			if err == nil {
				for _, a := range authorEntities {
					authorMap[a.Id] = a.Name
				}
			}
		}

		// 转换为 model
		for _, e := range entities {
			info := &model.ArticlesInfo{
				Id:        int(e.Id),
				Title:     e.Title,
				Url:       e.Url,
				Summary:   e.Summary,
				Content:   e.Content,
				DateAdded: e.DateAdded,
			}

			// 处理 AuthorId 和 AuthorName
			if e.AuthorId > 0 {
				aid := e.AuthorId
				info.AuthorId = &aid
				if name, ok := authorMap[e.AuthorId]; ok {
					info.AuthorName = name
				}
			}

			// 解析 Tags JSON 字段
			if e.Tags != "" {
				var tags []string
				err = json.Unmarshal([]byte(e.Tags), &tags)
				if err == nil {
					info.Tags = tags
				} else {
					info.Tags = []string{}
				}
			} else {
				info.Tags = []string{}
			}

			// 添加搜索高亮信息
			if searchHighlightMap != nil {
				if highlight, ok := searchHighlightMap[int(e.Id)]; ok {
					g.Log().Debugf(ctx, "Applying highlight to article %d: %+v", e.Id, highlight)
					if ft, ok := highlight["formattedTitle"].(string); ok && ft != "" {
						info.FormattedTitle = ft
					}
					if fs, ok := highlight["formattedSummary"].(string); ok && fs != "" {
						info.FormattedSummary = fs
					}
					if cs, ok := highlight["contextSnippet"].(string); ok && cs != "" {
						info.ContextSnippet = cs
					}
					if mf, ok := highlight["matchFields"].([]string); ok {
						info.MatchFields = mf
					}
				}
			}

			articlesList = append(articlesList, info)
		}
	})
	return
}

// GetById 获取文章详情
func (s sArticles) GetById(ctx context.Context, id int) (res *model.ArticlesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 先从数据库查询 entity
		var articleEntity *entity.Articles
		err = dao.Articles.Ctx(ctx).WherePri(id).Scan(&articleEntity)
		liberr.ErrIsNil(ctx, err, "文章不存在")

		// 转换为 ArticlesInfo
		res = &model.ArticlesInfo{
			Id:        int(articleEntity.Id),
			Title:     articleEntity.Title,
			Url:       articleEntity.Url,
			Summary:   articleEntity.Summary,
			Content:   articleEntity.Content,
			DateAdded: articleEntity.DateAdded,
		}

		// 处理 AuthorId（uint -> *uint，0 表示 NULL）
		if articleEntity.AuthorId > 0 {
			res.AuthorId = uintToUintPtr(articleEntity.AuthorId)
		}

		// 解析 Tags JSON 字段
		if articleEntity.Tags != "" {
			var tags []string
			err = json.Unmarshal([]byte(articleEntity.Tags), &tags)
			if err == nil {
				res.Tags = tags
			} else {
				res.Tags = []string{}
			}
		} else {
			res.Tags = []string{}
		}
	})
	return
}

// GetTags 获取所有标签（去重，按使用频率排序）
func (s sArticles) GetTags(ctx context.Context) (res []string, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 查询所有文章
		var articles []*entity.Articles
		err = dao.Articles.Ctx(ctx).Scan(&articles)
		liberr.ErrIsNil(ctx, err, "获取标签失败")

		// 统计标签频率
		tagMap := make(map[string]int)
		for _, article := range articles {
			if article.Tags != "" {
				var tags []string
				err = json.Unmarshal([]byte(article.Tags), &tags)
				if err == nil {
					for _, tag := range tags {
						tagMap[tag]++
					}
				}
			}
		}

		// 转换为结果数组并排序
		type tagCount struct {
			Tag   string
			Count int
		}
		var tagCounts []tagCount
		for tag, count := range tagMap {
			tagCounts = append(tagCounts, tagCount{Tag: tag, Count: count})
		}

		// 按频率降序排序
		sort.Slice(tagCounts, func(i, j int) bool {
			return tagCounts[i].Count > tagCounts[j].Count
		})

		// 提取标签名
		for _, tc := range tagCounts {
			res = append(res, tc.Tag)
		}
	})
	return
}

// GetStats 获取综合统计数据
func (s sArticles) GetStats(ctx context.Context) (res *model.ArticlesStatsInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		res = &model.ArticlesStatsInfo{}

		// 文章总数
		count, err := dao.Articles.Ctx(ctx).Count()
		liberr.ErrIsNil(ctx, err, "获取文章总数失败")
		res.TotalCount = count

		// 作者统计（从 authors 表获取）
		authorStats, err := s.GetAuthorStats(ctx)
		liberr.ErrIsNil(ctx, err, "获取作者统计失败")
		res.AuthorStats = authorStats

		// 标签统计
		tagStats, err := s.GetTagStats(ctx)
		liberr.ErrIsNil(ctx, err, "获取标签统计失败")
		res.TagStats = tagStats

		// 时间趋势
		dateTrend, err := s.GetDateTrend(ctx)
		liberr.ErrIsNil(ctx, err, "获取时间趋势失败")
		res.DateTrend = dateTrend

		// 最近添加的文章
		var recentArticles []*model.ArticlesInfo
		columns := dao.Articles.Columns()
		err = dao.Articles.Ctx(ctx).
			Order(fmt.Sprintf("%s DESC", columns.DateAdded)).
			Limit(10).
			Scan(&recentArticles)
		liberr.ErrIsNil(ctx, err, "获取最近文章失败")
		res.RecentArticles = recentArticles

		// 解析最近文章的 Tags
		for _, article := range recentArticles {
			if article.Tags == nil {
				article.Tags = []string{}
			}
		}
	})
	return
}

// GetAuthorStats 获取作者统计（从 authors 表关联查询）
func (s sArticles) GetAuthorStats(ctx context.Context) (res []*model.AuthorStat, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 从 authors 表获取统计数据
		authorsColumns := dao.Authors.Columns()
		articlesColumns := dao.Articles.Columns()

		// 关联查询：作者信息 + 最新文章时间
		err = dao.Authors.Ctx(ctx).
			Fields(fmt.Sprintf("authors.id as author_id, authors.name as author_name, authors.article_count as count, MAX(articles.%s) as latest_added",
				articlesColumns.DateAdded)).
			LeftJoin(fmt.Sprintf("articles ON authors.id = articles.%s", articlesColumns.AuthorId)).
			Group("authors.id").
			Order(fmt.Sprintf("authors.%s DESC", authorsColumns.ArticleCount)).
			Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取作者统计失败")
	})
	return
}

// GetTagStats 获取标签统计
func (s sArticles) GetTagStats(ctx context.Context) (res []*model.TagStat, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 查询所有文章
		var articles []*entity.Articles
		err = dao.Articles.Ctx(ctx).Scan(&articles)
		liberr.ErrIsNil(ctx, err, "获取标签统计失败")

		// 统计标签频率
		tagMap := make(map[string]int)
		for _, article := range articles {
			if article.Tags != "" {
				var tags []string
				err = json.Unmarshal([]byte(article.Tags), &tags)
				if err == nil {
					for _, tag := range tags {
						tagMap[tag]++
					}
				}
			}
		}

		// 转换为结果数组
		for tag, count := range tagMap {
			res = append(res, &model.TagStat{Tag: tag, Count: count})
		}

		// 按频率降序排序
		sort.Slice(res, func(i, j int) bool {
			return res[i].Count > res[j].Count
		})
	})
	return
}

// ParseByURL 通过 URL 解析微信文章
func (s sArticles) ParseByURL(ctx context.Context, req *api.ArticlesParseByURLReq) (res *api.ArticlesParseByURLRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 调用解析器
		article, err := parser.ParseFromURL(req.URL)
		liberr.ErrIsNil(ctx, err, "解析微信文章失败")

		res = &api.ArticlesParseByURLRes{
			Title:   article.Title,
			Author:  article.Author,
			Content: article.Content,
			BaseURL: article.BaseURL,
		}
	})
	return
}

// Parse 解析微信文章
func (s sArticles) Parse(ctx context.Context, req *api.ArticlesParseReq) (res *api.ArticlesParseRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 调用解析器
		article, err := parser.ParseWechatArticle(req.HtmlContent, req.BaseURL)
		liberr.ErrIsNil(ctx, err, "解析微信文章失败")

		res = &api.ArticlesParseRes{
			Title:   article.Title,
			Author:  article.Author,
			Content: article.Content,
			BaseURL: article.BaseURL,
		}
	})
	return
}

// GetDateTrend 获取时间趋势（最近30天）
func (s sArticles) GetDateTrend(ctx context.Context) (res []*model.DateTrend, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		columns := dao.Articles.Columns()
		// 将毫秒时间戳转换为日期并统计
		err = dao.Articles.Ctx(ctx).
			Fields(fmt.Sprintf("FROM_UNIXTIME(%s/1000, '%%Y-%%m-%%d') as date, COUNT(*) as count",
				columns.DateAdded)).
			Where(fmt.Sprintf("%s IS NOT NULL", columns.DateAdded)).
			Group("date").
			Order("date DESC").
			Limit(30).
			Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取时间趋势失败")
	})
	return
}

// Reparse 重新解析文章（从原始URL重新获取内容）
func (s sArticles) Reparse(ctx context.Context, req *api.ArticlesReparseReq) (res *api.ArticlesReparseRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 获取原文章信息
		oldArticle, err := s.GetById(ctx, req.Id)
		liberr.ErrIsNil(ctx, err, "文章不存在")

		// 检查是否有原始URL
		if oldArticle.Url == "" {
			liberr.ErrIsNil(ctx, fmt.Errorf("该文章没有原始URL，无法重新解析"), "")
			return
		}

		g.Log().Infof(ctx, "Reparse article %d from URL: %s", req.Id, oldArticle.Url)

		// 从原始URL重新解析
		article, err := parser.ParseFromURL(oldArticle.Url)
		liberr.ErrIsNil(ctx, err, "重新解析文章失败")

		// 先移除旧的图片关联（减少引用计数）
		err = service.Images().RemoveArticleImages(ctx, req.Id)
		if err != nil {
			g.Log().Warningf(ctx, "Failed to remove old article images: %v", err)
		}

		// 处理新内容中的图片（不会重复添加，通过URL hash判断）
		newContent := article.Content
		var imageIds []uint
		newContent, imageIds, err = service.Images().ProcessArticleImages(ctx, req.Id, article.Content)
		if err != nil {
			g.Log().Warningf(ctx, "Failed to process images: %v", err)
		}

		// 获取或创建作者
		var authorId int
		if oldArticle.AuthorId != nil {
			authorId = int(*oldArticle.AuthorId)
		}

		if article.Author != "" {
			// 使用 GetOrCreateByName 获取或创建作者
			newAuthorId, createErr := service.Authors().GetOrCreateByName(ctx, article.Author)
			if createErr == nil {
				authorId = newAuthorId
			}
		}

		// 更新文章
		_, err = dao.Articles.Ctx(ctx).WherePri(req.Id).Update(do.Articles{
			Title:    article.Title,
			AuthorId: uintToUintPtr(uint(authorId)),
			Content:  newContent,
		})
		liberr.ErrIsNil(ctx, err, "更新文章失败")

		// 更新作者的文章数量
		service.Authors().UpdateArticleCount(ctx, authorId)

		g.Log().Infof(ctx, "Article %d reparsed successfully, %d images processed", req.Id, len(imageIds))

		res = &api.ArticlesReparseRes{
			Title:       article.Title,
			Author:      article.Author,
			Content:     newContent,
			ImagesCount: len(imageIds),
		}
	})
	return
}

// searchWithMeilisearch 使用 Meilisearch 进行全文搜索
// 返回匹配的文章ID列表和高亮信息映射
func (s sArticles) searchWithMeilisearch(ctx context.Context, keyword string) ([]int, map[int]map[string]interface{}) {
	// 初始化 Meilisearch（如果未初始化）
	if err := libMeilisearch.Init(ctx); err != nil {
		g.Log().Warningf(ctx, "Meilisearch 初始化失败，回退到数据库查询: %v", err)
		return nil, nil
	}

	// 执行搜索
	hits, err := libMeilisearch.Search(ctx, keyword, 1000, "") // 限制1000条
	if err != nil {
		g.Log().Warningf(ctx, "Meilisearch 搜索失败: %v", err)
		return nil, nil
	}

	// 提取文章ID和高亮信息
	articleIds := make([]int, 0)
	highlightMap := make(map[int]map[string]interface{}) // articleId -> highlight info

	for i, hit := range hits {
		if hitMap, ok := hit.(map[string]interface{}); ok {
			var articleId int
			if idVal, ok := hitMap["id"]; ok {
				switch v := idVal.(type) {
				case float64:
					articleId = int(v)
				case int:
					articleId = v
				case int64:
					articleId = int(v)
				case string:
					if intVal, parseErr := strconv.Atoi(v); parseErr == nil {
						articleId = intVal
					}
				}
				if articleId > 0 {
					articleIds = append(articleIds, articleId)
					// 保存高亮信息
					formattedTitle, _ := hitMap["formattedTitle"].(string)
					formattedSummary, _ := hitMap["formattedSummary"].(string)
					contextSnippet, _ := hitMap["contextSnippet"].(string)
					matchFields, _ := hitMap["matchFields"].([]string)

					highlightMap[articleId] = map[string]interface{}{
						"formattedTitle":   formattedTitle,
						"formattedSummary": formattedSummary,
						"contextSnippet":   contextSnippet,
						"matchFields":      matchFields,
					}

					// 调试：打印第一个结果的高亮信息
					if i == 0 {
						g.Log().Infof(ctx, "First hit highlight: id=%d, formattedTitle=%s, contextSnippet=%s",
							articleId, formattedTitle, contextSnippet)
					}
				}
			}
		}
	}

	g.Log().Infof(ctx, "Meilisearch 搜索: keyword=%s, found=%d", keyword, len(articleIds))
	return articleIds, highlightMap
}
