// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ArticlesDao is the data access object for the table articles.
type ArticlesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ArticlesColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ArticlesColumns defines and stores column names for the table articles.
type ArticlesColumns struct {
	Id        string // 主键ID
	Title     string // 文章标题
	AuthorId  string // 作者ID（外键）
	Url       string // 原文链接
	Summary   string // 文章摘要
	Content   string // Markdown内容
	Tags      string // 标签数组（JSON格式）
	DateAdded string // 添加时间戳(毫秒，保留前端原始时间)
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// articlesColumns holds the columns for the table articles.
var articlesColumns = ArticlesColumns{
	Id:        "id",
	Title:     "title",
	AuthorId:  "author_id",
	Url:       "url",
	Summary:   "summary",
	Content:   "content",
	Tags:      "tags",
	DateAdded: "date_added",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewArticlesDao creates and returns a new DAO object for table data access.
func NewArticlesDao(handlers ...gdb.ModelHandler) *ArticlesDao {
	return &ArticlesDao{
		group:    "default",
		table:    "articles",
		columns:  articlesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ArticlesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ArticlesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ArticlesDao) Columns() ArticlesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ArticlesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ArticlesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ArticlesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
