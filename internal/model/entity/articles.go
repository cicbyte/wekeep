// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Articles is the golang structure for table articles.
type Articles struct {
	Id        uint        `json:"id"        orm:"id"         description:"主键ID"`               // 主键ID
	Title     string      `json:"title"     orm:"title"      description:"文章标题"`               // 文章标题
	AuthorId  uint        `json:"authorId"  orm:"author_id"  description:"作者ID（外键）"`           // 作者ID（外键）
	Url       string      `json:"url"       orm:"url"        description:"原文链接"`               // 原文链接
	Summary   string      `json:"summary"   orm:"summary"    description:"文章摘要"`               // 文章摘要
	Content   string      `json:"content"   orm:"content"    description:"Markdown内容"`         // Markdown内容
	Tags      string      `json:"tags"      orm:"tags"       description:"标签数组（JSON格式）"`       // 标签数组（JSON格式）
	DateAdded int64       `json:"dateAdded" orm:"date_added" description:"添加时间戳(毫秒，保留前端原始时间)"` // 添加时间戳(毫秒，保留前端原始时间)
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`               // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`               // 更新时间
}
