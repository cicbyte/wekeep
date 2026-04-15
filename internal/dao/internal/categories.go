// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CategoriesDao is the data access object for table categories.
type CategoriesDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns CategoriesColumns // columns contains all the column names of Table for convenient usage.
}

// CategoriesColumns defines and stores column names for table categories.
type CategoriesColumns struct {
	Id          string // 分类ID，自增主键
	Name        string // 分类名称，唯一
	Description string // 分类描述
	Icon        string // 分类图标标识或URL
	Sort        string // 数字越大越靠前
	CreatedAt   string // 记录创建时间
	UpdatedAt   string // 记录最后更新时间
}

// categoriesColumns holds the columns for table categories.
var categoriesColumns = CategoriesColumns{
	Id:          "id",
	Name:        "name",
	Description: "description",
	Icon:        "icon",
	Sort:        "sort",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewCategoriesDao creates and returns a new DAO object for table data access.
func NewCategoriesDao() *CategoriesDao {
	return &CategoriesDao{
		group:   "default",
		table:   "categories",
		columns: categoriesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CategoriesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CategoriesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CategoriesDao) Columns() CategoriesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CategoriesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CategoriesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CategoriesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}