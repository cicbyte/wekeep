// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AuthorsDao is the data access object for the table authors.
type AuthorsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  AuthorsColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// AuthorsColumns defines and stores column names for the table authors.
type AuthorsColumns struct {
	Id             string // 主键ID
	Name           string // 作者名称
	NormalizedName string // 标准化名称（去重用）
	Avatar         string // 头像URL
	Bio            string // 作者简介
	Website        string // 个人网站
	ArticleCount   string // 文章数量
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
}

// authorsColumns holds the columns for the table authors.
var authorsColumns = AuthorsColumns{
	Id:             "id",
	Name:           "name",
	NormalizedName: "normalized_name",
	Avatar:         "avatar",
	Bio:            "bio",
	Website:        "website",
	ArticleCount:   "article_count",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewAuthorsDao creates and returns a new DAO object for table data access.
func NewAuthorsDao(handlers ...gdb.ModelHandler) *AuthorsDao {
	return &AuthorsDao{
		group:    "default",
		table:    "authors",
		columns:  authorsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AuthorsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AuthorsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AuthorsDao) Columns() AuthorsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AuthorsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AuthorsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *AuthorsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
