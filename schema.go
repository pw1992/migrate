package migrate

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type SchemaBuild struct {
	Info   map[string]interface{}
	Fields []struct {
		Name     string
		Type     string
		Args     int
		Comment  string
		Default  string
		Example  interface{}
		Extra    interface{}
		Nullable bool
	}
	Indexes []map[string]interface{}
}

func NewSchemaBuild() *SchemaBuild {
	return new(SchemaBuild)
}

var files = make([]string, 0)

func (b *SchemaBuild) getFiles(dirPath string) []string {
	fileStat, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		panic(err)
	}
	if !fileStat.IsDir() {
		files = append(files, dirPath)
	}
	dirs, _ := ioutil.ReadDir(dirPath)
	for _, dir := range dirs {
		name := dir.Name()
		//只取json文件
		if !dir.IsDir() && strings.HasSuffix(name, ".json") {
			files = append(files, path.Join(dirPath, name))
		} else if dir.IsDir() {
			b.getFiles(path.Join(dirPath, name))
		} else {
			continue
		}
	}
	return files
}

// 数据库升级
func (b *SchemaBuild) Upgrade(dirPath string, dsn string) *SchemaBuild {
	modelFiles := b.getFiles(dirPath)
	for _, fileName := range modelFiles {
		var build SchemaBuild
		data, _ := ioutil.ReadFile(fileName)
		json.Unmarshal(data, &build)
		indexes := make(map[string]interface{}, 0)
		for _, index := range build.Indexes {
			indexes[index["field"].(string)] = index
		}

		tableName := build.Info["name"].(string)
		t := NewTable(tableName).Connect(dsn)
		if !t.IsTableExists() {
			t.CreateTable()
		}

		for _, field := range build.Fields {
			newColumn := NewColumn()
			switch field.Type {
			case "integer":
				newColumn.SetInteger(field.Name, field.Args).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "tinyInteger":
				newColumn.SetTinyInteger(field.Name, field.Args).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "smallInteger":
				newColumn.SetSmallInteger(field.Name, field.Args).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "mediumInteger":
				newColumn.SetMediumInteger(field.Name, field.Args).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "bigInteger":
				newColumn.SetBigInteger(field.Name, field.Args).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "string":
				newColumn.SetString(field.Name, field.Args).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "text":
				newColumn.SetText(field.Name).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "tinyText":
				newColumn.SetTinyText(field.Name).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "mediumText":
				newColumn.SetMediumText(field.Name).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "longText":
				newColumn.SetLongText(field.Name).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "timestamp":
				newColumn.SetTimestamp(field.Name).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "datetime":
				newColumn.SetDatetime(field.Name).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			case "json":
				newColumn.SetJson(field.Name).SetDefault(field.Default).SetComment(field.Comment).SetNullable(field.Nullable)
			}
			index, ok := indexes[field.Name]
			if ok {
				m := index.(map[string]interface{})
				switch m["type"].(string) {
				case "primary":
					comment, _ := m["comment"].(string)
					auto_increments, _ := m["auto_increments"].(bool)
					newColumn.SetPrimary(auto_increments).SetIndexComment(comment)
				case "unique":
					comment, _ := m["comment"].(string)
					newColumn.SetUnique().SetIndexComment(comment)
				case "index":
					comment, _ := m["comment"].(string)
					newColumn.SetIndex().SetIndexComment(comment)
				}
			}
			t.CreateColumn(newColumn)
		}

		//增加添加时间、修改时间、删除时间
		newColumn := NewColumn()
		newColumn.SetTimestamp("created_at").SetNullable(true).SetComment("添加时间").SetIndex()
		t.CreateColumn(newColumn)

		newColumn = NewColumn()
		newColumn.SetTimestamp("updated_at").SetNullable(true).SetComment("修改时间").SetIndex()
		t.CreateColumn(newColumn)

		newColumn = NewColumn()
		newColumn.SetTimestamp("deleted_at").SetNullable(true).SetComment("删除时间").SetIndex()
		t.CreateColumn(newColumn)

		t.Upgrade()
	}
	return b
}
