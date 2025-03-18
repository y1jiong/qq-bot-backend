// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TokenDao is the data access object for the table token.
type TokenDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TokenColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TokenColumns defines and stores column names for the table token.
type TokenColumns struct {
	Name        string //
	Token       string //
	OwnerId     string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
	LastLoginAt string //
	BotId       string //
}

// tokenColumns holds the columns for the table token.
var tokenColumns = TokenColumns{
	Name:        "name",
	Token:       "token",
	OwnerId:     "owner_id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	LastLoginAt: "last_login_at",
	BotId:       "bot_id",
}

// NewTokenDao creates and returns a new DAO object for table data access.
func NewTokenDao(handlers ...gdb.ModelHandler) *TokenDao {
	return &TokenDao{
		group:    "default",
		table:    "token",
		columns:  tokenColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TokenDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TokenDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TokenDao) Columns() TokenColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TokenDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TokenDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TokenDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
