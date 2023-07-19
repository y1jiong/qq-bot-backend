// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NamespaceDao is the data access object for table namespace.
type NamespaceDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns NamespaceColumns // columns contains all the column names of Table for convenient usage.
}

// NamespaceColumns defines and stores column names for table namespace.
type NamespaceColumns struct {
	Namespace   string //
	OwnerId     string //
	SettingJson string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// namespaceColumns holds the columns for table namespace.
var namespaceColumns = NamespaceColumns{
	Namespace:   "namespace",
	OwnerId:     "owner_id",
	SettingJson: "setting_json",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewNamespaceDao creates and returns a new DAO object for table data access.
func NewNamespaceDao() *NamespaceDao {
	return &NamespaceDao{
		group:   "default",
		table:   "namespace",
		columns: namespaceColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *NamespaceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *NamespaceDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *NamespaceDao) Columns() NamespaceColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *NamespaceDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *NamespaceDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *NamespaceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
