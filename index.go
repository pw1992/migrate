package migrate

type Index struct {
	AutoIncrement bool
	//desc table的返回值
	Field string //字段名称
	//Type string		//字段类型
	//Null string		//是否为null
	Key string //索引类型 PRI  MUL  UNI
	//Default string	//默认值
	//Extra string

	//SELECT TABLE_NAME as TableName,INDEX_NAME as IndexName,COLUMN_NAME as ColumnName,INDEX_COMMENT as IndexComment from information_schema.`STATISTICS` where TABLE_NAME='user2'; 结构体
	TableName    string //表明
	IndexName    string //索引名称
	ColumnName   string //字段名称
	IndexComment string //索引说明
}
