package model

// AuthorsInfo 作者信息模型（用于返回给前端）
type AuthorsInfo struct {
	Id           int    `orm:"id" json:"id"`
	Name         string `orm:"name" json:"name"`
	Avatar       string `orm:"avatar" json:"avatar"`
	Bio          string `orm:"bio" json:"bio"`
	Website      string `orm:"website" json:"website"`
	ArticleCount int    `orm:"article_count" json:"articleCount"`
	CreatedAt    string `orm:"created_at" json:"createdAt"`
}

// AuthorSelectOption 作者选择选项（用于下拉框）
type AuthorSelectOption struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ArticleCount int    `json:"articleCount"`
}
