package migrate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Table struct {
	Connection    string
	Db            *sql.DB
	TableName     string
	SqlAlterTable string
	OldColumn     map[string]*Column
	NewColumn     map[string]*Column
	DelColumn     map[string]*Column
}

func NewTable(tableName string) *Table {
	t := &Table{
		Connection:    "",
		Db:            nil,
		TableName:     tableName,
		SqlAlterTable: "",
		OldColumn:     map[string]*Column{},
		NewColumn:     map[string]*Column{},
		DelColumn:     map[string]*Column{},
	}
	return t
}

func (t *Table) Connect(dsn string) *Table {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	t.Db = db
	return t
}

// 判断表是否存在
func (t *Table) IsTableExists() bool {
	type showTable struct {
		TableName string
	}
	//sqlStr := fmt.Sprintf("SELECT COUNT(*) as cnt FROM information_schema.TABLES WHERE table_name ='%s'", t.TableName)
	sqlStr := fmt.Sprintf("show tables")
	query, _ := t.Db.Query(sqlStr)
	stables := make([]string, 0)
	for query.Next() {
		var stable showTable
		query.Scan(&stable.TableName)
		stables = append(stables, stable.TableName)
	}
	for _, v := range stables {
		if v == t.TableName {
			return true
		}
	}
	return false
}

// 新增表
func (t *Table) CreateTable() *Table {
	if t.IsTableExists() {
		return t
	}
	sql := fmt.Sprintf(`CREATE TABLE %s  (
							  id int NOT NULL AUTO_INCREMENT,
							  PRIMARY KEY (id)
							);`,
		t.TableName,
	)

	_, err := t.Db.Exec(sql)

	if err != nil {
		panic(err)
	}
	return t
}

// 创建字段
func (t *Table) CreateColumn(c *Column) *Table {
	t.NewColumn[c.Name] = c
	return t
}

// 删除字段
func (t *Table) DeleteColumn(c *Column) *Table {
	t.DelColumn[c.Name] = c
	return t
}

// 获取表结构
func (t *Table) DescTable() *Table {
	//显示表结构信息
	query, err := t.Db.Query("desc " + t.TableName)
	if err != nil {
		panic(err)
	}

	type DescTableStruct struct {
		Field   string //字段名称
		Type    string //字段类型
		Null    string //是否为null
		Key     string //索引类型 PRI  MUL  UNI
		Default string //默认值
		Extra   string
	}

	descs := make(map[string]DescTableStruct, 0)
	for query.Next() {
		var desc DescTableStruct
		query.Scan(&desc.Field, &desc.Type, &desc.Null, &desc.Key, &desc.Default, &desc.Extra)
		descs[desc.Field] = desc
	}

	//显示索引信息
	sql := fmt.Sprintf("SELECT TABLE_NAME as TableName,INDEX_NAME as IndexName,COLUMN_NAME as ColumnName,INDEX_COMMENT as IndexComment from information_schema.`STATISTICS` where TABLE_NAME='%s'", t.TableName)
	query, err = t.Db.Query(sql)
	type STATISTICSTableStruct struct {
		TableName    string //表明
		IndexName    string //索引名称
		ColumnName   string //字段名称
		IndexComment string //索引说明
	}
	if err != nil {
		panic(err)
	}
	indexMap := make(map[string]STATISTICSTableStruct, 0)
	for query.Next() {
		var index STATISTICSTableStruct
		query.Scan(&index.TableName, &index.IndexName, &index.ColumnName, &index.IndexComment)
		indexMap[index.ColumnName] = index
	}

	for _, desc := range descs {
		//查询长度
		reg := regexp.MustCompile(`\d+`)
		findString := reg.FindString(desc.Type)
		lenght := 0
		colType := ""
		if len(findString) > 0 {
			lenght, _ = strconv.Atoi(findString)
		}

		if strings.HasPrefix(desc.Type, "varchar") {
			colType = "varchar"
		} else if strings.HasPrefix(desc.Type, "char") {
			colType = "char"
		} else if strings.HasPrefix(desc.Type, "int") {
			colType = "integer"
		} else if strings.HasPrefix(desc.Type, "json") {
			colType = "json"
		} else if strings.HasPrefix(desc.Type, "timestamp") {
			colType = "timestamp"
		} else if strings.HasPrefix(desc.Type, "datetime") {
			colType = "datetime"
		}
		field := desc.Field

		index, _ := indexMap[field]

		t.OldColumn[field] = &Column{
			Name:    field,
			Type:    colType,
			Length:  lenght,
			Default: desc.Default,
			Comment: "",
			Null:    desc.Null,
			Extra:   desc.Extra,
			OptType: "update",
			Indexes: &Index{
				Field:        field,
				Key:          desc.Key,
				TableName:    index.TableName,
				IndexName:    index.IndexName,
				ColumnName:   index.ColumnName,
				IndexComment: index.IndexComment,
			},
		}
	}
	return t
}

// 升级表
func (t *Table) Upgrade() {
	t.DescTable() //获取表结构数据

	sql := ""
	addColumns := make([]*Column, 0)    //新增的字段
	updateColumns := make([]*Column, 0) //修改的字段
	for field, newCol := range t.NewColumn {
		_, ok := t.OldColumn[field]
		if !ok {
			addColumns = append(addColumns, newCol)
		} else {
			updateColumns = append(updateColumns, newCol)
		}
	}

	if len(addColumns) > 0 {
		sql += sqlStr(addColumns, "ADD")
	}
	if len(updateColumns) > 0 {
		sql += sqlStr(updateColumns, "MODIFY")
	}

	for column, _ := range t.DelColumn {
		sql += fmt.Sprintf(" DROP %s ,", column)
	}
	if len(sql) > 0 {
		sql = strings.TrimRight(sql, ",")
		sql = sql + ";"
		sql = "ALTER TABLE " + t.TableName + sql
		_, err := t.Db.Exec(sql)
		if err != nil {
			fmt.Println("执行错误:", err)
		}
	}

	t.OldColumn = map[string]*Column{}
	t.DescTable()
	//修改索引
	sql = ""

	keyMap := map[string]string{
		"UNI": "UNIQUE",
		"MUL": "INDEX",
		"PRI": "PRIMARY KEY",
	}

	for field, oldColumn := range t.OldColumn {
		newColumn, ok := t.NewColumn[field]
		if ok {
			if len(oldColumn.Indexes.IndexName) > 0 && len(newColumn.Indexes.IndexName) == 0 {
				//主键不动
				if oldColumn.Indexes.Key != "PRI" {
					sql += fmt.Sprintf("DROP INDEX %s ,", oldColumn.Indexes.IndexName)
				}
			} else if len(oldColumn.Indexes.IndexName) == 0 && len(newColumn.Indexes.IndexName) > 0 {
				if newColumn.Indexes.Key == "PRI" {
					sql += fmt.Sprintf("ADD PRIMARY KEY (%s),MODIFY COLUMN %s int(%d) NOT NULL AUTO_INCREMENT COMMENT '主键id' FIRST,",
						newColumn.Name,
						newColumn.Name,
						newColumn.Length,
					)
				} else {
					sql += fmt.Sprintf("ADD %s %s(%s) USING BTREE COMMENT '%s',",
						keyMap[newColumn.Indexes.Key],
						newColumn.Indexes.IndexName,
						newColumn.Indexes.ColumnName,
						newColumn.Indexes.IndexComment,
					)
				}
			} else if len(oldColumn.Indexes.ColumnName) > 0 && len(newColumn.Indexes.ColumnName) > 0 && oldColumn.Indexes.IndexName != newColumn.Indexes.IndexName {
				if newColumn.Indexes.Key == "PRI" {
					sql += fmt.Sprintf("DROP PRIMARY KEY, ADD PRIMARY KEY (%s)  , MODIFY COLUMN %s int(%d) NOT NULL AUTO_INCREMENT COMMENT '主键id' FIRST,",
						newColumn.Name,
						newColumn.Name,
						newColumn.Length,
					)
				} else {
					sql += fmt.Sprintf("DROP INDEX %s,ADD %s %s(%s) USING BTREE COMMENT '%s',",
						oldColumn.Indexes.IndexName,
						keyMap[newColumn.Indexes.Key],
						newColumn.Indexes.IndexName,
						newColumn.Indexes.ColumnName,
						newColumn.Indexes.IndexComment,
					)
				}
			}
		}
	}

	if len(sql) > 0 {
		sql = strings.TrimRight(sql, ",")
		sql = sql + ";"
		sql = "ALTER TABLE " + t.TableName + " " + sql

		_, err := t.Db.Exec(sql)
		if err != nil {
			fmt.Println("执行错误:", err)
		}
	}
}

func sqlStr(columns []*Column, addType string) string {
	sql := ``
	for _, column := range columns {
		switch column.Type {
		case "varchar", "char":
			sql = fmt.Sprintf("%s %s COLUMN %s %s(%d) COMMENT '%s' ",
				sql,
				addType,
				column.Name,
				column.Type,
				column.Length,
				column.Comment,
			)
			if column.Nullable {
				sql = fmt.Sprintf("%s NULL DEFAULT '%s' ,", sql, column.Default)
			} else {
				sql = fmt.Sprintf("%s NOT NULL DEFAULT '%s' ,", sql, column.Default)
			}

		case "text", "tinytext", "mediumtext", "longtext":
			sql = fmt.Sprintf("%s %s COLUMN %s %s COMMENT '%s' ",
				sql,
				addType,
				column.Name,
				column.Type,
				column.Comment,
			)
			if column.Nullable {
				sql += "NULL ,"
			} else {
				sql += "NOT NULL,"
			}
		case "int", "tinyint", "smallint", "mediumint", "bigint":
			sql = fmt.Sprintf("%s %s COLUMN %s %s(%d) COMMENT '%s' ",
				sql,
				addType,
				column.Name,
				column.Type,
				column.Length,
				column.Comment,
			)
			//如果不是int类型获取默认值不存在  设置0
			typeof := reflect.TypeOf(column.Default)
			if typeof.String() != "int" {
				column.Default = 0
			}
			nullable := "NOT NULL"
			if column.Nullable {
				nullable = "NULL"
			}
			sql = fmt.Sprintf("%s %s DEFAULT %v ,", sql, nullable, column.Default)
		case "timestamp", "datetime", "json":
			sql = fmt.Sprintf("%s %s COLUMN %s %s %s COMMENT '%s' ,",
				sql,
				addType,
				column.Name,
				column.Type,
				"NULL",
				column.Comment,
			)
		}

	}
	return sql

}

func DD(args ...interface{}) {
	Dump(args...)
	os.Exit(0)
}

func Dump(args ...interface{}) {
	format := ""
	for i := 0; i < len(args); i++ {
		format = format + "----\nType: %s\n%v\n----\n\n"
	}
	var inputs []interface{}
	for _, i := range args {
		if i == nil {
			inputs = append(inputs, "nil", "")
			continue
		}
		val := reflect.ValueOf(i)
		switch val.Kind() {
		case reflect.Map, reflect.Array, reflect.Interface, reflect.Slice:
			inputs = append(inputs, val.Type(), JsonEncode(i))
		default:
			inputs = append(inputs, val.Type(), i)
		}
	}
	fmt.Printf(format, inputs...)
}

func JsonEncode(i interface{}) (jsonstr string) {
	jsonbytes, err := json.Marshal(i)
	if err != nil {
		return "{}"
	}
	return string(jsonbytes)
}
