// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ArticleImages is the golang structure for table article_images.
type ArticleImages struct {
	Id        uint        `json:"id"        orm:"id"         description:"主键ID"`    // 主键ID
	ArticleId uint        `json:"articleId" orm:"article_id" description:"文章ID"`    // 文章ID
	ImageId   uint        `json:"imageId"   orm:"image_id"   description:"图片ID"`    // 图片ID
	Position  int         `json:"position"  orm:"position"   description:"图片位置"`    // 图片位置
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`    // 创建时间
}
