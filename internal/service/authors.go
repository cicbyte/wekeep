package service

import (
	"context"
	api "github.com/cicbyte/wekeep/api/v1/authors"
	model "github.com/cicbyte/wekeep/internal/model"
)

type IAuthors interface {
	// Add 新增作者
	Add(ctx context.Context, req *api.AuthorsAddReq) (id int, err error)

	// Edit 编辑作者
	Edit(ctx context.Context, req *api.AuthorsEditReq) (err error)

	// Delete 删除作者
	Delete(ctx context.Context, id int) (err error)

	// List 获取作者列表（支持分页、筛选、排序）
	List(ctx context.Context, req *api.AuthorsListReq) (total interface{}, res []*model.AuthorsInfo, err error)

	// GetById 获取作者详情
	GetById(ctx context.Context, id int) (res *model.AuthorsInfo, err error)

	// GetSelectOptions 获取作者选择选项（用于下拉框）
	GetSelectOptions(ctx context.Context) (res []*model.AuthorSelectOption, err error)

	// GetOrCreateByName 根据名称获取或创建作者（供 articles 使用）
	GetOrCreateByName(ctx context.Context, name string) (id int, err error)

	// UpdateArticleCount 更新作者的文章数量
	UpdateArticleCount(ctx context.Context, authorId int) (err error)
}

var localAuthors IAuthors

// Authors 返回作者管理服务的实例
func Authors() IAuthors {
	if localAuthors == nil {
		panic("implement not found for interface IAuthors, forgot register?")
	}
	return localAuthors
}

// RegisterAuthors 注册作者管理服务实现
func RegisterAuthors(i IAuthors) {
	localAuthors = i
}
