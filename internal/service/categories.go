package service

import (
	"context"

	api "github.com/cicbyte/wekeep/api/v1/categories"
	model "github.com/cicbyte/wekeep/internal/model"
)

type ICategories interface {
	List(ctx context.Context, req *api.CategoriesListReq) (total interface{}, res []*model.CategoriesInfo, err error)
	Add(ctx context.Context, req *api.CategoriesAddReq) (err error)
	Edit(ctx context.Context, req *api.CategoriesEditReq) (err error)
	Delete(ctx context.Context, id int) (err error)
	BatchDelete(ctx context.Context, ids []int) (err error)
	GetById(ctx context.Context, id int) (res *model.CategoriesInfo, err error)
}

var localCategories ICategories

func Categories() ICategories {
	if localCategories == nil {
		panic("implement not found for interface ICategories, forgot register?")
	}
	return localCategories
}

func RegisterCategories(i ICategories) {
	localCategories = i
}