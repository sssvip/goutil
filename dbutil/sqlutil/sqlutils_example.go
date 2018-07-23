package sqlutil

import (
	"goutil/strutil"
	"strings"
)

func count() string {
	sql, tArgs, _ := NewSQLGen("tb_test").
		QueryColumns("tiny_url", "created_on", "origin_url", "source").
		And("tiny_url", "test value").
		Count()
	return strutil.Format("%s,%v", sql, tArgs)
}
func insert() string {
	sql, values, _ := NewSQLGen("tb_test").
		InsertColumn("source", "test").
		InsertColumn("tiny_url", "test").
		InsertColumn("created_on", "test").
		Insert()
	return strutil.Format("%v,%v", sql, values)
}
func update() string {
	sql, tArgs, _ := NewSQLGen("tb_test").
		UpdateColumn("source", "test value").
		And("tiny_url", "test value").
		Update()
	return strutil.Format("%s,%v", sql, tArgs)
}
func query() string {
	sql, tArgs, _ := NewSQLGen("tb_test").
		QueryColumns("tiny_url", "created_on", "origin_url", "source").
		And("tiny_url", "test value").
		Or("created_on", "or value").
		Or("source", "2").
		Limit(10).
		OrderByAsc("created_on").
		OrderByDesc("tiny_url").
		Query()
	return strutil.Format("%s,%v", sql, tArgs)
}

func delete() string {
	sql, tArgs, _ := NewSQLGen("tb_test").
		QueryColumns("tiny_url", "created_on", "origin_url", "source").
		And("tiny_url", "test value").Delete()
	return strutil.Format("%s,%v", sql, tArgs)
}

func Example() string {
	return exampleBase(true)
}

func exampleBase(newLineSplit bool) string {
	f := "%s;%s;%s;%s;%s;"
	if newLineSplit {
		f = strings.Replace(f, ";", ";\n", -1)
	}
	return strutil.Format(f, insert(), update(), query(), count(), delete())
}
