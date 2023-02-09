package migrate

type Column struct {
	Name     string
	Type     string
	Length   int
	Default  interface{}
	Comment  string
	Nullable bool
	Null     string
	Extra    string
	OptType  string //操作类型 add=新增  update=修改  del=删除
	Indexes  *Index //索引
}

func NewColumn() *Column {
	return &Column{
		Name:     "",
		Type:     "",
		Length:   0,
		Default:  nil,
		Comment:  "",
		Nullable: true,
		OptType:  "add",
		Indexes: &Index{
			Field:        "",
			Key:          "",
			TableName:    "",
			IndexName:    "",
			ColumnName:   "",
			IndexComment: "",
		},
	}
}

func (c *Column) GetColumn(columnName string) *Column {
	c.Name = columnName
	return c
}

func (c *Column) SetString(col string, defaultLen ...int) *Column {
	lenght := 255
	if len(defaultLen) > 0 {
		lenght = defaultLen[0]
	}
	c.Name = col
	c.Length = lenght
	c.Type = "varchar"
	c.OptType = "add"
	c.Default = ""
	return c
}

func (c *Column) SetText(col string) *Column {
	c.Name = col
	c.Length = 0
	c.Type = "text"
	c.OptType = "add"
	c.Default = ""
	return c
}

func (c *Column) SetTinyText(col string) *Column {
	c.Name = col
	c.Length = 0
	c.Type = "tinytext"
	c.OptType = "add"
	c.Default = ""
	return c
}

func (c *Column) SetMediumText(col string) *Column {
	c.Name = col
	c.Length = 0
	c.Type = "mediumtext"
	c.OptType = "add"
	c.Default = ""
	return c
}

func (c *Column) SetLongText(col string) *Column {
	c.Name = col
	c.Length = 0
	c.Type = "longtext"
	c.OptType = "add"
	c.Default = ""
	return c
}

func (c *Column) SetInteger(col string, defaultLen ...int) *Column {
	lenght := 20
	if len(defaultLen) > 0 {
		lenght = defaultLen[0]
	}
	c.Name = col
	c.Length = lenght
	c.Type = "int"
	c.OptType = "add"
	c.Default = 0
	return c
}

func (c *Column) SetTinyInteger(col string, defaultLen ...int) *Column {
	lenght := 4
	if len(defaultLen) > 0 {
		lenght = defaultLen[0]
	}
	if lenght >= 4 {
		lenght = 4
	}

	c.Name = col
	c.Length = lenght
	c.Type = "tinyint"
	c.OptType = "add"
	c.Default = 0
	return c
}

func (c *Column) SetSmallInteger(col string, defaultLen ...int) *Column {
	lenght := 6
	if len(defaultLen) > 0 {
		lenght = defaultLen[0]
	}
	if lenght >= 6 {
		lenght = 6
	}

	c.Name = col
	c.Length = lenght
	c.Type = "smallint"
	c.OptType = "add"
	c.Default = 0
	return c
}

func (c *Column) SetMediumInteger(col string, defaultLen ...int) *Column {
	lenght := 11
	if len(defaultLen) > 0 {
		lenght = defaultLen[0]
	}
	if lenght >= 11 {
		lenght = 11
	}

	c.Name = col
	c.Length = lenght
	c.Type = "mediumint"
	c.OptType = "add"
	c.Default = 0
	return c
}

func (c *Column) SetBigInteger(col string, defaultLen ...int) *Column {
	lenght := 20
	if len(defaultLen) > 0 {
		lenght = defaultLen[0]
	}
	if lenght >= 20 {
		lenght = 20
	}

	c.Name = col
	c.Length = lenght
	c.Type = "bigint"
	c.OptType = "add"
	c.Default = 0
	return c
}

func (c *Column) SetTimestamp(col string) *Column {
	c.Name = col
	c.Type = "timestamp"
	c.OptType = "add"
	return c
}

func (c *Column) SetDatetime(col string) *Column {
	c.Name = col
	c.Type = "datetime"
	c.OptType = "add"
	return c
}

func (c *Column) SetJson(col string) *Column {
	c.Name = col
	c.Type = "json"
	c.OptType = "add"
	return c
}

func (c *Column) SetDefault(d interface{}) *Column {
	c.Default = d
	return c
}

func (c *Column) SetComment(comment string) *Column {
	c.Comment = comment
	return c
}

func (c *Column) SetNullable(nullable bool) *Column {
	c.Nullable = nullable
	return c
}

//普通索引
func (c *Column) SetIndex() *Column {
	c.Indexes.Field = c.Name
	c.Indexes.Key = "MUL" //普通索引
	c.Indexes.IndexName = "index_" + c.Name
	c.Indexes.IndexComment = ""
	c.Indexes.ColumnName = c.Name
	return c
}

//唯一索引
func (c *Column) SetUnique() *Column {
	c.Indexes.Field = c.Name
	c.Indexes.Key = "UNI" //普通索引
	c.Indexes.IndexName = "unique_" + c.Name
	c.Indexes.IndexComment = ""
	c.Indexes.ColumnName = c.Name
	return c
}

//主键索引
func (c *Column) SetPrimary(auto_increments bool) *Column {
	c.Indexes.AutoIncrement = auto_increments
	c.Indexes.Field = c.Name
	c.Indexes.Key = "PRI" //普通索引
	c.Indexes.IndexName = "PRIMARY"
	c.Indexes.IndexComment = ""
	c.Indexes.ColumnName = c.Name
	return c
}

//索引说明
func (c *Column) SetIndexComment(comment string) *Column {
	c.Indexes.IndexComment = comment
	return c
}
