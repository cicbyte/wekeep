// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Categories is the golang structure of table categories for DAO operations like Where/Data.
type Categories struct {
	g.Meta      `orm:"table:categories, do:true"`
	Id          interface{} // 分类ID，自增主键
	Name        interface{} // 分类名称，唯一
	Description interface{} // 分类描述
	Icon        interface{} // 分类图标标识或URL
	Sort        interface{} // 数字越大越靠前
	CreatedAt   *gtime.Time // 记录创建时间
	UpdatedAt   *gtime.Time // 记录最后更新时间
}