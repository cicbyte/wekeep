package categories

import (
	"context"
	"fmt"

	api "github.com/ciclebyte/wekeep/api/v1/categories"
	dao "github.com/ciclebyte/wekeep/internal/dao"
	model "github.com/ciclebyte/wekeep/internal/model"
	do "github.com/ciclebyte/wekeep/internal/model/do"
	service "github.com/ciclebyte/wekeep/internal/service"
	liberr "github.com/ciclebyte/wekeep/library/liberr"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	service.RegisterCategories(New())
}

func New() *sCategories {
	return &sCategories{}
}

type sCategories struct {
}

func (s sCategories) List(ctx context.Context, req *api.CategoriesListReq) (total interface{}, categoriesList []*model.CategoriesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.Categories.Ctx(ctx)
		columns := dao.Categories.Columns()
		if req.Name != "" {
			m = m.Where(fmt.Sprintf("%s like ?", columns.Name), "%"+req.Name+"%")
		}
		total, err = m.Count()
		liberr.ErrIsNil(ctx, err, "获取分类列表失败")
		orderBy := req.OrderBy
		if orderBy == "" {
			orderBy = "created_at desc"
		}
		err = m.Page(req.PageNum, req.PageSize).Order(orderBy).Scan(&categoriesList)
		liberr.ErrIsNil(ctx, err, "获取分类列表失败")
	})
	return
}

func (s sCategories) Add(ctx context.Context, req *api.CategoriesAddReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// TODO 查询是否已经存在

		// add
		_, err = dao.Categories.Ctx(ctx).Insert(do.Categories{
			Name:        req.Name,        // 分类名称，唯一
			Description: req.Description, // 分类描述
			Icon:        req.Icon,        // 分类图标标识或URL
			Sort:        req.Sort,        // 数字越大越靠前
		})
		liberr.ErrIsNil(ctx, err, "新增分类失败")
	})
	return
}

func (s sCategories) Edit(ctx context.Context, req *api.CategoriesEditReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = s.GetById(ctx, req.Id)
		liberr.ErrIsNil(ctx, err, "获取分类失败")
		//TODO 根据名称等查询是否存在

		//编辑
		_, err = dao.Categories.Ctx(ctx).WherePri(req.Id).Update(do.Categories{
			Id:          req.Id,          // 分类ID，自增主键
			Name:        req.Name,        // 分类名称，唯一
			Description: req.Description, // 分类描述
			Icon:        req.Icon,        // 分类图标标识或URL
			Sort:        req.Sort,        // 数字越大越靠前
		})
		liberr.ErrIsNil(ctx, err, "修改分类失败")
	})
	return
}

func (s sCategories) Delete(ctx context.Context, id int) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.Categories.Ctx(ctx).WherePri(id).Delete()
		liberr.ErrIsNil(ctx, err, "删除分类失败")
	})
	return
}

func (s sCategories) BatchDelete(ctx context.Context, ids []int) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.Categories.Ctx(ctx).Where(dao.Categories.Columns().Id+" in(?)", ids).Delete()
		liberr.ErrIsNil(ctx, err, "批量删除分类失败")
	})
	return
}

func (s sCategories) GetById(ctx context.Context, id int) (res *model.CategoriesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.Categories.Ctx(ctx).Where(fmt.Sprintf("%s=?", dao.Categories.Columns().Id), id).Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取分类失败")
	})
	return
}