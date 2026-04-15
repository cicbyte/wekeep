package articles

import (
	commonApi "github.com/ciclebyte/wekeep/api/v1/common"
	model "github.com/ciclebyte/wekeep/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

// ArticlesAddReq 新增文章请求
type ArticlesAddReq struct {
	g.Meta    `path:"/articles/add" method:"post" tags:"文章管理" summary:"文章-新增"`
	Title     string   `json:"title" v:"required#文章标题不能为空|length:1,512#标题长度为1-512位"`
	AuthorId  int      `json:"authorId" v:"required#作者ID不能为空|min:1#作者ID必须大于0"`
	Url       string   `json:"url" v:"max-length:2048#链接最大长度为2048位"`
	Summary   string   `json:"summary" v:"max-length:2000#摘要最大长度为2000位"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	DateAdded int64    `json:"dateAdded" v:"min:0#时间戳不能为负数"`
}

type ArticlesAddRes struct {
	g.Meta `mime:"application/json"`
	Id     int `json:"id"` // 返回新增文章的ID
}

// ArticlesEditReq 编辑文章请求
type ArticlesEditReq struct {
	g.Meta    `path:"/articles/edit" method:"put" tags:"文章管理" summary:"文章-编辑"`
	Id        int      `json:"id" v:"required#文章ID不能为空|min:1#ID必须大于0"`
	Title     string   `json:"title" v:"required#文章标题不能为空|length:1,512#标题长度为1-512位"`
	AuthorId  int      `json:"authorId" v:"required#作者ID不能为空|min:1#作者ID必须大于0"`
	Url       string   `json:"url" v:"max-length:2048#链接最大长度为2048位"`
	Summary   string   `json:"summary" v:"max-length:2000#摘要最大长度为2000位"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	DateAdded int64    `json:"dateAdded" v:"min:0#时间戳不能为负数"`
}

type ArticlesEditRes struct {
	g.Meta `mime:"application/json"`
}

// ArticlesDelReq 删除文章请求
type ArticlesDelReq struct {
	g.Meta `path:"/articles/del" method:"delete" tags:"文章管理" summary:"文章-删除"`
	Id     int `json:"id" v:"required|min:1#ID必须大于0"`
}

type ArticlesDelRes struct {
	g.Meta `mime:"application/json"`
}

// ArticlesBatchDelReq 批量删除文章请求
type ArticlesBatchDelReq struct {
	g.Meta `path:"/articles/batchdel" method:"delete" tags:"文章管理" summary:"文章-批量删除"`
	Ids    []int `json:"ids" v:"required#请选择要删除的文章"`
}

type ArticlesBatchDelRes struct {
	g.Meta    `mime:"application/json"`
	Success   int    `json:"success"`   // 成功删除数量
	Failed    int    `json:"failed"`    // 失败数量
	FailedIds []int  `json:"failedIds"` // 失败的ID列表
}

// ArticlesListReq 文章列表请求
type ArticlesListReq struct {
	g.Meta     `path:"/articles/list" method:"get" tags:"文章管理" summary:"文章-列表"`
	commonApi.PageReq
	AuthorId  *int     `p:"authorId" dc:"按作者ID筛选"`
	Tags      []string `p:"tags[]" dc:"按标签筛选（支持多选）"`
	Keyword   string   `p:"keyword" dc:"搜索关键词（标题、摘要）"`
	TitleOnly bool     `p:"titleOnly" dc:"仅搜索标题（不使用全文检索）"`
	OrderBy   string   `p:"orderBy" dc:"排序方式，如：dateAdded desc"`
}

type ArticlesListRes struct {
	g.Meta `mime:"application/json"`
	commonApi.ListRes
	ArticlesList []*model.ArticlesInfo `json:"articlesList"`
}

// ArticlesDetailReq 文章详情请求
type ArticlesDetailReq struct {
	g.Meta `path:"/articles/detail" method:"get" tags:"文章管理" summary:"文章-详情"`
	Id     int `p:"id" v:"required|min:1#ID必须大于0"`
}

type ArticlesDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.ArticlesInfo
}

// ArticlesTagsReq 标签列表请求
type ArticlesTagsReq struct {
	g.Meta `path:"/articles/tags" method:"get" tags:"文章管理" summary:"获取标签列表"`
}

type ArticlesTagsRes struct {
	g.Meta `mime:"application/json"`
	Tags []string `json:"tags"` // 标签列表（去重，按使用频率排序）
}

// ArticlesParseByURLReq 通过 URL 解析微信文章请求
type ArticlesParseByURLReq struct {
	g.Meta `path:"/articles/parse-by-url" method:"post" tags:"文章管理" summary:"通过 URL 解析微信文章"`
	URL string `json:"url" v:"required|url#URL不能为空且格式不正确"`
}

type ArticlesParseByURLRes struct {
	g.Meta `mime:"application/json"`
	Title   string `json:"title"`   // 文章标题
	Author  string `json:"author"`  // 文章作者
	Content string `json:"content"` // 文章内容
	BaseURL string `json:"baseUrl"` // 基础 URL
}

// ArticlesParseReq 解析微信文章请求
type ArticlesParseReq struct {
	g.Meta `path:"/articles/parse" method:"post" tags:"文章管理" summary:"解析微信文章"`
	HtmlContent string `json:"htmlContent" v:"required#HTML内容不能为空"`
	BaseURL     string `json:"baseUrl"`
}

type ArticlesParseRes struct {
	g.Meta `mime:"application/json"`
	Title   string `json:"title"`   // 文章标题
	Author  string `json:"author"`  // 文章作者
	Content string `json:"content"` // 文章内容
	BaseURL string `json:"baseUrl"` // 基础 URL
}

// ArticlesStatsReq 统计数据请求
type ArticlesStatsReq struct {
	g.Meta `path:"/articles/stats" method:"get" tags:"文章管理" summary:"获取统计数据"`
}

type ArticlesStatsRes struct {
	g.Meta `mime:"application/json"`
	TotalCount      int                    `json:"totalCount"`      // 文章总数
	AuthorStats     []*model.AuthorStat    `json:"authorStats"`     // 作者统计
	TagStats        []*model.TagStat       `json:"tagStats"`        // 标签统计
	RecentArticles  []*model.ArticlesInfo  `json:"recentArticles"`  // 最近添加的文章（最多10条）
	DateTrend       []*model.DateTrend     `json:"dateTrend"`       // 时间趋势（最近30天）
}

// ArticlesReparseReq 重新解析文章请求
type ArticlesReparseReq struct {
	g.Meta `path:"/articles/reparse" method:"post" tags:"文章管理" summary:"文章-重新解析"`
	Id     int `json:"id" v:"required|min:1#文章ID不能为空"`
}

type ArticlesReparseRes struct {
	g.Meta      `mime:"application/json"`
	Title       string `json:"title"`       // 新标题
	Author      string `json:"author"`      // 新作者名
	Content     string `json:"content"`     // 新内容
	ImagesCount int    `json:"imagesCount"` // 处理的图片数量
}
