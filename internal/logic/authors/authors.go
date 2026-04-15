package authors

import (
	"context"
	"fmt"
	"strings"

	api "github.com/ciclebyte/wekeep/api/v1/authors"
	dao "github.com/ciclebyte/wekeep/internal/dao"
	model "github.com/ciclebyte/wekeep/internal/model"
	do "github.com/ciclebyte/wekeep/internal/model/do"
	service "github.com/ciclebyte/wekeep/internal/service"
	liberr "github.com/ciclebyte/wekeep/library/liberr"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	service.RegisterAuthors(New())
}

func New() *sAuthors {
	return &sAuthors{}
}

type sAuthors struct{}

// normalizeName 标准化名称（用于去重）
func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// Add 新增作者
func (s sAuthors) Add(ctx context.Context, req *api.AuthorsAddReq) (id int, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		normalized := normalizeName(req.Name)

		// 检查是否已存在
		count, err := dao.Authors.Ctx(ctx).Where("normalized_name", normalized).Count()
		liberr.ErrIsNil(ctx, err, "检查作者失败")
		if count > 0 {
			liberr.ErrIsNil(ctx, fmt.Errorf("作者名称已存在"), "新增作者失败")
		}

		result, err := dao.Authors.Ctx(ctx).Insert(do.Authors{
			Name:           req.Name,
			NormalizedName: normalized,
			Avatar:         req.Avatar,
			Bio:            req.Bio,
			Website:        req.Website,
			ArticleCount:   0,
		})
		liberr.ErrIsNil(ctx, err, "新增作者失败")

		lastId, err := result.LastInsertId()
		liberr.ErrIsNil(ctx, err, "获取作者ID失败")
		id = int(lastId)
	})
	return
}

// Edit 编辑作者
func (s sAuthors) Edit(ctx context.Context, req *api.AuthorsEditReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 验证作者是否存在
		_, err = s.GetById(ctx, req.Id)
		liberr.ErrIsNil(ctx, err, "作者不存在")

		normalized := normalizeName(req.Name)

		// 检查名称是否与其他作者重复
		count, err := dao.Authors.Ctx(ctx).
			Where("normalized_name", normalized).
			WhereNot("id", req.Id).
			Count()
		liberr.ErrIsNil(ctx, err, "检查作者失败")
		if count > 0 {
			liberr.ErrIsNil(ctx, fmt.Errorf("作者名称已存在"), "编辑作者失败")
		}

		_, err = dao.Authors.Ctx(ctx).WherePri(req.Id).Update(do.Authors{
			Name:           req.Name,
			NormalizedName: normalized,
			Avatar:         req.Avatar,
			Bio:            req.Bio,
			Website:        req.Website,
		})
		liberr.ErrIsNil(ctx, err, "修改作者失败")
	})
	return
}

// Delete 删除作者
func (s sAuthors) Delete(ctx context.Context, id int) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = s.GetById(ctx, id)
		liberr.ErrIsNil(ctx, err, "作者不存在")

		_, err = dao.Authors.Ctx(ctx).WherePri(id).Delete()
		liberr.ErrIsNil(ctx, err, "删除作者失败")
	})
	return
}

// List 获取作者列表
func (s sAuthors) List(ctx context.Context, req *api.AuthorsListReq) (total interface{}, authorsList []*model.AuthorsInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.Authors.Ctx(ctx)
		columns := dao.Authors.Columns()

		// 按名称模糊筛选
		if req.Name != "" {
			m = m.Where(fmt.Sprintf("%s LIKE ?", columns.Name), "%"+req.Name+"%")
		}

		// 排序
		orderBy := req.OrderBy
		if orderBy == "" {
			orderBy = fmt.Sprintf("%s DESC", columns.ArticleCount)
		}

		total, err = m.Count()
		liberr.ErrIsNil(ctx, err, "获取作者列表失败")

		err = m.Page(req.PageNum, req.PageSize).Order(orderBy).Scan(&authorsList)
		liberr.ErrIsNil(ctx, err, "获取作者列表失败")
	})
	return
}

// GetById 获取作者详情
func (s sAuthors) GetById(ctx context.Context, id int) (res *model.AuthorsInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.Authors.Ctx(ctx).WherePri(id).Scan(&res)
		liberr.ErrIsNil(ctx, err, "作者不存在")
		liberr.ValueIsNil(res, "作者不存在")
	})
	return
}

// GetSelectOptions 获取作者选择选项
func (s sAuthors) GetSelectOptions(ctx context.Context) (res []*model.AuthorSelectOption, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		columns := dao.Authors.Columns()
		err = dao.Authors.Ctx(ctx).
			Fields(fmt.Sprintf("id, %s as name, %s as article_count", columns.Name, columns.ArticleCount)).
			Order(fmt.Sprintf("%s DESC", columns.ArticleCount)).
			Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取作者选项失败")
	})
	return
}

// GetOrCreateByName 根据名称获取或创建作者
func (s sAuthors) GetOrCreateByName(ctx context.Context, name string) (id int, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		normalized := normalizeName(name)
		columns := dao.Authors.Columns()

		// 先查找是否存在
		var author *model.AuthorsInfo
		err = dao.Authors.Ctx(ctx).
			Where(fmt.Sprintf("%s = ?", columns.NormalizedName), normalized).
			Scan(&author)
		liberr.ErrIsNil(ctx, err, "查询作者失败")

		if author != nil {
			id = author.Id
			return
		}

		// 不存在则创建
		result, err := dao.Authors.Ctx(ctx).Insert(do.Authors{
			Name:           name,
			NormalizedName: normalized,
			ArticleCount:   0,
		})
		liberr.ErrIsNil(ctx, err, "创建作者失败")

		lastId, err := result.LastInsertId()
		liberr.ErrIsNil(ctx, err, "获取作者ID失败")
		id = int(lastId)
	})
	return
}

// UpdateArticleCount 更新作者的文章数量
func (s sAuthors) UpdateArticleCount(ctx context.Context, authorId int) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		articlesColumns := dao.Articles.Columns()

		// 统计该作者的文章数量
		count, err := dao.Articles.Ctx(ctx).
			Where(fmt.Sprintf("%s = ?", articlesColumns.AuthorId), authorId).
			Count()
		liberr.ErrIsNil(ctx, err, "统计文章数量失败")

		// 更新作者表
		authorsColumns := dao.Authors.Columns()
		_, err = dao.Authors.Ctx(ctx).
			WherePri(authorId).
			Update(g.Map{authorsColumns.ArticleCount: count})
		liberr.ErrIsNil(ctx, err, "更新文章数量失败")
	})
	return
}
