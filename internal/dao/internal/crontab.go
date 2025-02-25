// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CrontabDao is the data access object for the table crontab.
type CrontabDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of the current DAO.
	columns CrontabColumns // columns contains all the column names of Table for convenient usage.
}

// CrontabColumns defines and stores column names for the table crontab.
type CrontabColumns struct {
	Name       string //
	Expression string //
	CreatorId  string //
	BotId      string //
	Request    string //
	CreatedAt  string //
}

// crontabColumns holds the columns for the table crontab.
var crontabColumns = CrontabColumns{
	Name:       "name",
	Expression: "expression",
	CreatorId:  "creator_id",
	BotId:      "bot_id",
	Request:    "request",
	CreatedAt:  "created_at",
}

// NewCrontabDao creates and returns a new DAO object for table data access.
func NewCrontabDao() *CrontabDao {
	return &CrontabDao{
		group:   "default",
		table:   "crontab",
		columns: crontabColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CrontabDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CrontabDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CrontabDao) Columns() CrontabColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CrontabDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CrontabDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *CrontabDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
