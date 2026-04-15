// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Authors is the golang structure for table authors.
type Authors struct {
	Id             uint        `json:"id"             orm:"id"              description:"主键ID"`       // 主键ID
	Name           string      `json:"name"           orm:"name"            description:"作者名称"`       // 作者名称
	NormalizedName string      `json:"normalizedName" orm:"normalized_name" description:"标准化名称（去重用）"` // 标准化名称（去重用）
	Avatar         string      `json:"avatar"         orm:"avatar"          description:"头像URL"`      // 头像URL
	Bio            string      `json:"bio"            orm:"bio"             description:"作者简介"`       // 作者简介
	Website        string      `json:"website"        orm:"website"         description:"个人网站"`       // 个人网站
	ArticleCount   uint        `json:"articleCount"   orm:"article_count"   description:"文章数量"`       // 文章数量
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:"创建时间"`       // 创建时间
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:"更新时间"`       // 更新时间
}
