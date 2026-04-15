package model

// ArticlesInfo 文章信息模型（用于返回给前端）
type ArticlesInfo struct {
	Id               int      `orm:"id" json:"id"`
	Title            string   `orm:"title" json:"title"`
	AuthorId         *uint    `orm:"author_id" json:"authorId"`       // 作者ID（外键）
	AuthorName       string   `orm:"-" json:"authorName"`             // 作者名称（关联查询）
	Url              string   `orm:"url" json:"url"`
	Summary          string   `orm:"summary" json:"summary"`
	Content          string   `orm:"content" json:"content"`
	Tags             []string `orm:"tags" json:"tags"`                // JSON 数组，从 string 转换
	DateAdded        int64    `orm:"date_added" json:"dateAdded"`
	// 搜索高亮相关字段
	FormattedTitle   string   `orm:"-" json:"formattedTitle,omitempty"`   // 高亮后的标题
	FormattedSummary string   `orm:"-" json:"formattedSummary,omitempty"` // 高亮后的摘要
	ContextSnippet   string   `orm:"-" json:"contextSnippet,omitempty"`   // 正文匹配上下文
	MatchFields      []string `orm:"-" json:"matchFields,omitempty"`      // 匹配的字段列表
}

// AuthorStat 作者统计模型
type AuthorStat struct {
	AuthorId    *uint  `json:"authorId"`    // 作者ID
	AuthorName  string `json:"authorName"`  // 作者名称
	Count       int    `json:"count"`       // 文章数量
	LatestAdded int64  `json:"latestAdded"` // 最新文章添加时间戳
}

// TagStat 标签统计模型
type TagStat struct {
	Tag   string `json:"tag"`   // 标签名称
	Count int    `json:"count"` // 使用次数
}

// DateTrend 时间趋势模型
type DateTrend struct {
	Date  string `json:"date"`  // 日期 (YYYY-MM-DD)
	Count int    `json:"count"` // 当天文章数量
}

// ArticlesStatsInfo 综合统计模型（可选，用于复杂统计）
type ArticlesStatsInfo struct {
	TotalCount     int           `json:"totalCount"`
	AuthorStats    []*AuthorStat `json:"authorStats"`
	TagStats       []*TagStat    `json:"tagStats"`
	RecentArticles []*ArticlesInfo `json:"recentArticles"`
	DateTrend      []*DateTrend  `json:"dateTrend"`
}
