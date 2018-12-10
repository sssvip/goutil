package dbutil

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sssvip/goutil/dbutil/sqlutil"
	"github.com/sssvip/goutil/strutil"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const dbFileName = "temp.db"

const tableName = "tb_user_info"

var createTableSql = strutil.Format(`CREATE TABLE %s (
    "uid" INTEGER PRIMARY KEY AUTOINCREMENT,
    "username" VARCHAR(64) NULL,
    "depart_name" VARCHAR(64) NULL,
    "created_on" DATE NULL
);`, tableName)

var tempTestDB *sql.DB

func createTestDB(t *testing.T) {
	db, e := sql.Open("sqlite3", dbFileName)
	if e != nil {
		t.Error("创建临时数据库失败", e)
	}
	db.Exec(createTableSql)
	tempTestDB = db
}
func removeTestDB(t *testing.T) {
	e := os.Remove(dbFileName)
	if e != nil {
		t.Error(e)
	}
}

func TestCRUD(t *testing.T) {
	createTestDB(t)
	var usernameVale = "test"
	var departNameValue = "test_dep"
	//新增
	sqlGen := sqlutil.NewSQLGen(tableName).InsertColumn("username", usernameVale).InsertColumn("depart_name", departNameValue)
	//插入两条
	cnt, e := InsertTableBySQLGen(tempTestDB, sqlGen)
	_, e = InsertTableBySQLGen(tempTestDB, sqlGen)
	assert.Nil(t, e)
	assert.Equal(t, int64(1), cnt)
	//查询
	var name string
	var departName string
	e = QueryForObject(tempTestDB, sqlutil.NewSQLGen(tableName).QueryColumns("username", "depart_name"), &name, &departName)
	assert.Nil(t, e)
	assert.Equal(t, usernameVale, name)
	assert.Equal(t, departNameValue, departName)
	//	查询方式2
	row, e := GetRowBySQLGen(tempTestDB, sqlutil.NewSQLGen(tableName).QueryColumns("username", "depart_name"))
	assert.Nil(t, e)
	assert.Equal(t, 2, len(row))
	assert.Equal(t, usernameVale, row[0])
	assert.Equal(t, departNameValue, row[1])
	//	查询 count
	count, e := CountBySQLGen(tempTestDB, sqlutil.NewSQLGen(tableName))
	assert.Nil(t, e)
	assert.Equal(t, 2, count)
	//修改
	var newUserameValue = "new_test"
	var newdepartNameValue = "new_test_dep"
	updateCount, e := UpdateTableBySQLGen(tempTestDB, sqlutil.NewSQLGen(tableName).UpdateColumn("username", newUserameValue).UpdateColumn("depart_name", newdepartNameValue).ForceExecOnNoCondition())
	assert.Nil(t, e)
	assert.Equal(t, int64(2), updateCount)
	//		查询修改后的结果
	count, e = CountBySQLGen(tempTestDB, sqlutil.NewSQLGen(tableName).And("username", newUserameValue).And("depart_name", newdepartNameValue))
	assert.Nil(t, e)
	assert.Equal(t, 2, count)
	//删除
	deleteCount, e := DeleteTableBySQLGen(tempTestDB, sqlutil.NewSQLGen(tableName).And("username", newUserameValue))
	assert.Nil(t, e)
	assert.Equal(t, int64(2), deleteCount)
	//      查询删除后的结果
	count, e = CountBySQLGen(tempTestDB, sqlutil.NewSQLGen(tableName))
	assert.Nil(t, e)
	assert.Equal(t, 0, count)
	removeTestDB(t)
}
