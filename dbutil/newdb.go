package dbutil

import (
	"database/sql"
	"github.com/sssvip/goutil/dbutil/sqlutil"
)

/*

基于dbutil 包装

*/

type DBWrapper struct {
	OriginDB *sql.DB
}

func NewDB(username, password, address, port, database string) *DBWrapper {
	return &DBWrapper{OriginDB: NewDBByArg(username, password, address, port, database)}
}

func (db *DBWrapper) QueryForObjectBySQLStr(sqlStr string, columnsAddress ...interface{}) error {
	return QueryForObjectBySQLStr(db.OriginDB, sqlStr, columnsAddress...)
}

func (db *DBWrapper) QueryForObject(sqlGen *sqlutil.SQLGen, columnsAddress ...interface{}) error {
	return QueryForObject(db.OriginDB, sqlGen, columnsAddress...)
}

func (db *DBWrapper) CountBySQLGen(sqlGen *sqlutil.SQLGen) (result int, err error) {
	return CountBySQLGen(db.OriginDB, sqlGen)
}

func (db *DBWrapper) CountBySQLStr(sqlStr string, args ...interface{}) (result int, err error) {
	return CountBySQLStr(db.OriginDB, sqlStr, args...)
}

func (db *DBWrapper) InsertTableBySQLGen(sqlGen *sqlutil.SQLGen) (result int64, err error) {
	return InsertTableBySQLGen(db.OriginDB, sqlGen)
}

func (db *DBWrapper) DeleteTableBySQLGen(sqlGen *sqlutil.SQLGen) (result int64, err error) {
	return DeleteTableBySQLGen(db.OriginDB, sqlGen)
}

func (db *DBWrapper) UpdateTableBySQLGen(sqlGen *sqlutil.SQLGen) (result int64, err error) {
	return UpdateTableBySQLGen(db.OriginDB, sqlGen)
}

func (db *DBWrapper) Exec(sqlStr string, args ...interface{}) (result int64, err error) {
	return Exec(db.OriginDB, sqlStr, args...)
}

func (db *DBWrapper) GetRowBySQLGen(sqlGen *sqlutil.SQLGen) ([]string, error) {
	return GetRowBySQLGen(db.OriginDB, sqlGen)
}

func (db *DBWrapper) GetRowBySQLStr(sqlStr string, args ...interface{}) (row []string, err error) {
	return GetRowBySQLStr(db.OriginDB, sqlStr, args...)
}

func (db *DBWrapper) GetRowsBySQLGenPrintSql(sqlGen *sqlutil.SQLGen) ([][]string, error) {
	return GetRowsBySQLGenPrintSql(db.OriginDB, sqlGen)
}

func (db *DBWrapper) GetRowsBySQLGen(sqlGen *sqlutil.SQLGen) ([][]string, error) {
	return GetRowsBySQLGen(db.OriginDB, sqlGen)
}

func (db *DBWrapper) GetRowsBySQLStr(sqlStr string, args ...interface{}) ([][]string, error) {
	return GetRowsBySQLStr(db.OriginDB, sqlStr, args...)
}
func (db *DBWrapper) GetRowsBySQLStrWithQueryColumnsCount(sqlStr string, queryColumnsCount int) ([][]string, error) {
	return GetRowsBySQLStr(db.OriginDB, sqlStr)
}
