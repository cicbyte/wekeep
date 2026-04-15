// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Categories is the golang structure for table categories.
type Categories struct {
	Id          int         `json:"id"          description:"分类ID，自增主键"`
	Name        string      `json:"name"        description:"分类名称，唯一"`
	Description string      `json:"description" description:"分类描述"`
	Icon        string      `json:"icon"        description:"分类图标标识或URL"`
	Sort        int         `json:"sort"        description:"数字越大越靠前"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"记录创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"记录最后更新时间"`
}