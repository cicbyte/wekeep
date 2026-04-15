// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ArticleImages is the golang structure of table article_images for DAO operations like Where/Data.
type ArticleImages struct {
	g.Meta    `orm:"table:article_images, do:true"`
	Id        any         // 主键ID
	ArticleId any         // 文章ID
	ImageId   any         // 图片ID
	Position  any         // 图片位置
	CreatedAt *gtime.Time // 创建时间
}
