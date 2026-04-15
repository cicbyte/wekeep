// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Articles is the golang structure of table articles for DAO operations like Where/Data.
type Articles struct {
	g.Meta    `orm:"table:articles, do:true"`
	Id        any         // 主键ID
	Title     any         // 文章标题
	AuthorId  any         // 作者ID（外键）
	Url       any         // 原文链接
	Summary   any         // 文章摘要
	Content   any         // Markdown内容
	Tags      any         // 标签数组（JSON格式）
	DateAdded any         // 添加时间戳(毫秒，保留前端原始时间)
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
