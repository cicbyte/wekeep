// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Authors is the golang structure of table authors for DAO operations like Where/Data.
type Authors struct {
	g.Meta         `orm:"table:authors, do:true"`
	Id             any         // 主键ID
	Name           any         // 作者名称
	NormalizedName any         // 标准化名称（去重用）
	Avatar         any         // 头像URL
	Bio            any         // 作者简介
	Website        any         // 个人网站
	ArticleCount   any         // 文章数量
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
}
