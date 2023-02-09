package migrate

import (
	"testing"
)

func TestTable_CreateColumn(t *testing.T) {
	dsn := "username:passwordk@tcp(host:port)/dbname"
	table := NewTable("user_333").Connect(dsn)
	table.CreateTable()

	//新增字段
	column := NewColumn()
	column.SetString("c1").SetDefault("aaa").SetComment("c111111").SetNullable(false).SetIndex().SetIndexComment("aaaaaaIndex")
	table.CreateColumn(column)

	//修改字段
	column = NewColumn()
	column.SetString("name").SetDefault("name").SetComment("name").SetNullable(false).SetUnique().SetIndexComment("name")
	table.CreateColumn(column)

	//新增字段
	column = NewColumn()
	column.SetTimestamp("c2").SetDefault("aaa").SetComment("c111111").SetNullable(false).SetIndex().SetIndexComment("aaaaaaIndex")
	table.CreateColumn(column)

	//table.DeleteColumn( NewColumn().GetColumn("address") )

	table.Upgrade()
}
