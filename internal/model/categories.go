package model

type CategoriesInfo struct {
	Id          int    `orm:"id"  json:"id"`                   // 分类ID，自增主键
	Name        string `orm:"name"  json:"name"`               // 分类名称，唯一
	Description string `orm:"description"  json:"description"` // 分类描述
	Icon        string `orm:"icon"  json:"icon"`               // 分类图标标识或URL
	Sort        int    `orm:"sort"  json:"sort"`               // 数字越大越靠前
}