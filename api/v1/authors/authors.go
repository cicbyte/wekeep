package authors

import (
	commonApi "github.com/ciclebyte/wekeep/api/v1/common"
	model "github.com/ciclebyte/wekeep/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

// AuthorsAddReq 新增作者请求
type AuthorsAddReq struct {
	g.Meta   `path:"/authors/add" method:"post" tags:"作者管理" summary:"作者-新增"`
	Name     string `json:"name" v:"required#作者名称不能为空|length:1,255#作者名称长度为1-255位"`
	Avatar   string `json:"avatar" v:"max-length:512#头像URL最大长度为512位"`
	Bio      string `json:"bio"`
	Website  string `json:"website" v:"max-length:512#网站URL最大长度为512位"`
}

type AuthorsAddRes struct {
	g.Meta `mime:"application/json"`
	Id     int `json:"id"` // 返回新增作者的ID
}

// AuthorsEditReq 编辑作者请求
type AuthorsEditReq struct {
	g.Meta   `path:"/authors/edit" method:"put" tags:"作者管理" summary:"作者-编辑"`
	Id       int    `json:"id" v:"required#作者ID不能为空|min:1#ID必须大于0"`
	Name     string `json:"name" v:"required#作者名称不能为空|length:1,255#作者名称长度为1-255位"`
	Avatar   string `json:"avatar" v:"max-length:512#头像URL最大长度为512位"`
	Bio      string `json:"bio"`
	Website  string `json:"website" v:"max-length:512#网站URL最大长度为512位"`
}

type AuthorsEditRes struct {
	g.Meta `mime:"application/json"`
}

// AuthorsDelReq 删除作者请求
type AuthorsDelReq struct {
	g.Meta `path:"/authors/del" method:"delete" tags:"作者管理" summary:"作者-删除"`
	Id     int `json:"id" v:"required|min:1#ID必须大于0"`
}

type AuthorsDelRes struct {
	g.Meta `mime:"application/json"`
}

// AuthorsListReq 作者列表请求
type AuthorsListReq struct {
	g.Meta   `path:"/authors/list" method:"get" tags:"作者管理" summary:"作者-列表"`
	commonApi.PageReq
	Name    string `p:"name" dc:"按作者名称筛选（模糊匹配）"`
	OrderBy string `p:"orderBy" dc:"排序方式，如：article_count desc"`
}

type AuthorsListRes struct {
	g.Meta `mime:"application/json"`
	commonApi.ListRes
	AuthorsList []*model.AuthorsInfo `json:"authorsList"`
}

// AuthorsDetailReq 作者详情请求
type AuthorsDetailReq struct {
	g.Meta `path:"/authors/detail" method:"get" tags:"作者管理" summary:"作者-详情"`
	Id     int `p:"id" v:"required|min:1#ID必须大于0"`
}

type AuthorsDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.AuthorsInfo
}

// AuthorsSelectReq 作者选择选项请求（用于下拉框）
type AuthorsSelectReq struct {
	g.Meta `path:"/authors/select" method:"get" tags:"作者管理" summary:"作者-选择选项"`
}

type AuthorsSelectRes struct {
	g.Meta `mime:"application/json"`
	Options []*model.AuthorSelectOption `json:"options"`
}
