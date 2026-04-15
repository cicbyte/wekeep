package controller

import (
	"context"
	api "github.com/ciclebyte/wekeep/api/v1/authors"
	consts "github.com/ciclebyte/wekeep/internal/consts"
	service "github.com/ciclebyte/wekeep/internal/service"
)

var Authors = authorsController{}

type authorsController struct {
	BaseController
}

// Add 新增作者
func (c *authorsController) Add(ctx context.Context, req *api.AuthorsAddReq) (res *api.AuthorsAddRes, err error) {
	res = new(api.AuthorsAddRes)
	id, err := service.Authors().Add(ctx, req)
	if err != nil {
		return
	}
	res.Id = id
	return
}

// Edit 编辑作者
func (c *authorsController) Edit(ctx context.Context, req *api.AuthorsEditReq) (res *api.AuthorsEditRes, err error) {
	res = new(api.AuthorsEditRes)
	err = service.Authors().Edit(ctx, req)
	return
}

// Delete 删除作者
func (c *authorsController) Delete(ctx context.Context, req *api.AuthorsDelReq) (res *api.AuthorsDelRes, err error) {
	res = new(api.AuthorsDelRes)
	err = service.Authors().Delete(ctx, req.Id)
	return
}

// List 作者列表
func (c *authorsController) List(ctx context.Context, req *api.AuthorsListReq) (res *api.AuthorsListRes, err error) {
	res = new(api.AuthorsListRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	total, list, err := service.Authors().List(ctx, req)
	res.Total = total
	res.CurrentPage = req.PageNum
	res.AuthorsList = list
	return
}

// Detail 作者详情
func (c *authorsController) Detail(ctx context.Context, req *api.AuthorsDetailReq) (res *api.AuthorsDetailRes, err error) {
	res = new(api.AuthorsDetailRes)
	info, err := service.Authors().GetById(ctx, req.Id)
	if err != nil {
		return
	}
	res.AuthorsInfo = info
	return
}

// Select 作者选择选项
func (c *authorsController) Select(ctx context.Context, req *api.AuthorsSelectReq) (res *api.AuthorsSelectRes, err error) {
	res = new(api.AuthorsSelectRes)
	options, err := service.Authors().GetSelectOptions(ctx)
	if err != nil {
		return
	}
	res.Options = options
	return
}
