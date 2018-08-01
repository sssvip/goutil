/*
更加底层原生，提供原语级灵活调用
*/
package dbutil

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sssvip/goutil/dbutil/sqlutil"
	"github.com/sssvip/goutil/logutil"
	"github.com/sssvip/goutil/strutil"
)

//ErrorCode 默认的错误值，错误的count数等等
const ErrorCode = -1
const ErrorCount = -1

//NewDBByArg 通过参数获取db对象
func NewDBByArg(username, password, address, port, database string) *sql.DB {
	url := strutil.Format("%s:%s@tcp(%s:%s)/%s?collation=utf8mb4_unicode_ci&charset=utf8mb4", username, password, address, port, database)
	db, err := sql.Open("mysql", url)
	if err != nil {
		logutil.Error.Println(url, err)
		return nil
	}
	//db.SetMaxOpenConns(50)
	//db.SetMaxIdleConns(3)
	//db.Ping()
	return db
}

//GetRowBySQLStr 通过sqlStr获取数据库行记录
func GetRowBySQLStr(db *sql.DB, sqlStr string) (row []string, err error) {
	return justOneRow(GetRowsBySQLStr(db, sqlStr))
}

//QueryForObjectBySQLStr 通过sqlStr获取数据库行记录,直接写入传入的地址，不再copy数据
func QueryForObjectBySQLStr(db *sql.DB, sqlStr string, columns ...interface{}) (err error) {
	err = db.QueryRow(sqlStr).Scan(columns...)
	if err != nil {
		logutil.Error.Println(err, sqlStr)
		return err
	}
	return
}

//QueryForObject 通过sqlGen获取数据库行记录,直接写入传入的地址，不再copy数据
func QueryForObject(db *sql.DB, sqlGen *sqlutil.SQLGen, columns ...interface{}) (err error) {
	sqlGen.Limit(1)
	sqlStr, args, e := sqlGen.Query()
	if e != nil {
		logutil.Error.Println(e, sqlStr, args)
		return e
	}
	err = db.QueryRow(sqlStr, args...).Scan(columns...)
	if err != nil {
		logutil.Error.Println(err, sqlStr)
		return err
	}
	return
}

func GetRowsBySQLStr(db *sql.DB, sqlStr string, args ...interface{}) (rows [][]string, err error) {
	rs, e := db.Query(sqlStr, args...)
	if e != nil {
		logutil.Error.Println(e, sqlStr, args)
		return nil, err
	}
	defer func() {
		e := rs.Close()
		if e != nil {
			logutil.Error.Println(e)
		}
	}()
	for rs.Next() {
		cls, er := rs.Columns()
		if er != nil {
			logutil.Error.Println(er)
			return rows, er
		}
		//遍历取列
		buff := make([]interface{}, len(cls))
		var columns = make([]string, len(cls))
		for i, _ := range buff {
			buff[i] = &columns[i]
		}
		e := rs.Scan(buff...)
		if e != nil {
			logutil.Error.Println(e)
			return rows, e
		}
		rows = append(rows, columns)
	}
	return rows, err
}

func GetRowsBySQLGenPrintSql(db *sql.DB, sqlGen *sqlutil.SQLGen) (rows [][]string, err error) {
	sqlStr, args, e := sqlGen.Query()
	logutil.Console.Println(sqlStr)
	if e != nil {
		logutil.Error.Println(e, sqlStr, args)
	}
	return GetRowsBySQLStr(db, sqlStr, args...)
}

func GetRowsBySQLGen(db *sql.DB, sqlGen *sqlutil.SQLGen) (rows [][]string, err error) {
	sqlStr, args, e := sqlGen.Query()
	if e != nil {
		logutil.Error.Println(e, sqlStr, args)
		return nil, e
	}
	return GetRowsBySQLStr(db, sqlStr, args...)
}

func justOneRow(columns [][]string, e error) (row []string, err error) {
	if len(columns) > 0 {
		return columns[0], e
	}
	return []string{}, e
}

func GetRowBySQLGen(db *sql.DB, sqlGen *sqlutil.SQLGen) (row []string, err error) {
	sqlGen.Limit(1)
	return justOneRow(GetRowsBySQLGen(db, sqlGen))
}

func DeleteTableBySQLGen(db *sql.DB, sqlGen *sqlutil.SQLGen) (result int64, err error) {
	sqlStr, args, e := sqlGen.Delete()
	if e != nil {
		logutil.Error.Println(e, sqlStr, args)
		return ErrorCount, e
	}
	return Exec(db, sqlStr, args...)
}

func UpdateTableBySQLGen(db *sql.DB, sqlGen *sqlutil.SQLGen) (result int64, err error) {
	sqlStr, args, e := sqlGen.Update()
	if e != nil {
		logutil.Error.Println(e)
		return ErrorCount, e
	}
	return Exec(db, sqlStr, args...)
}

func InsertTableBySQLGen(db *sql.DB, sqlGen *sqlutil.SQLGen) (result int64, err error) {
	sqlStr, args, e := sqlGen.Insert()
	if e != nil {
		logutil.Error.Println(e)
		return ErrorCount, e
	}
	return Exec(db, sqlStr, args...)
}

func Exec(db *sql.DB, sql string, args ...interface{}) (result int64, err error) {
	rst, e := db.Exec(sql, args...)
	if e == nil {
		return rst.RowsAffected()
	} else {
		logutil.Error.Println(e, sql, args)
	}
	return ErrorCode, e
}

func CountBySQLStr(db *sql.DB, sqlStr string, args ...interface{}) (result int, err error) {
	var count int
	err = db.QueryRow(sqlStr, args...).Scan(&count)
	if err != nil {
		logutil.Error.Println(err, sqlStr, args)
		return ErrorCode, err
	}
	return count, nil
}

func CountBySQLGen(db *sql.DB, sqlGen *sqlutil.SQLGen) (result int, err error) {
	sqlStr, args, e := sqlGen.Count()
	if e != nil {
		logutil.Error.Println(e, sqlStr, args)
		return ErrorCode, e
	}
	return CountBySQLStr(db, sqlStr, args...)
}

func Count(db *sql.DB, tableName, condition string) (result int, err error) {
	return CountBySQLGen(db, sqlutil.NewSQLGen(tableName).CustomConditionAppend(condition))
}
