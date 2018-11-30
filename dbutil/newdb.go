package dbutil

import (
	"database/sql"
	"github.com/sssvip/goutil/dbutil/sqlutil"
	"github.com/sssvip/goutil/timeutil/stopwatch"
	"github.com/sssvip/goutil/jsonutil"
	"github.com/sssvip/goutil/logutil"
)

/*

基于dbutil 包装

*/

type DBWrapper struct {
	slowSQLSeconds   int
	OriginDB         *sql.DB
	openCheckSlowSQL bool
}

var slowSQLCheck = "slowSQLCheck"

func (db *DBWrapper) SetSlowSQLSeconds(seconds int) {
	db.slowSQLSeconds = seconds
}

func (db *DBWrapper) OpenCheckSlowSQL(open bool) {
	db.openCheckSlowSQL = open
}

func (db *DBWrapper) check(sw *stopwatch.StopWatch, o interface{}) {
	if sw == nil {
		return
	}
	if db.openCheckSlowSQL && int(sw.ElapsedSeconds()) > db.slowSQLSeconds {
		logutil.Warning.Println("found slow sql:", jsonutil.MarshalToString(o))
	}
}

func NewDB(username, password, address, port, database string) *DBWrapper {
	return &DBWrapper{OriginDB: NewDBByArg(username, password, address, port, database),
		slowSQLSeconds: 10, openCheckSlowSQL: false}
}

func (db *DBWrapper) QueryForObjectBySQLStr(sqlStr string, columnsAddress ...interface{}) error {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	e := QueryForObjectBySQLStr(db.OriginDB, sqlStr, columnsAddress...)
	db.check(t, sqlStr)
	return e
}

func (db *DBWrapper) QueryForObject(sqlGen *sqlutil.SQLGen, columnsAddress ...interface{}) error {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	e := QueryForObject(db.OriginDB, sqlGen, columnsAddress...)
	db.check(t, sqlGen)
	return e
}

func (db *DBWrapper) CountBySQLGen(sqlGen *sqlutil.SQLGen) (result int, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = CountBySQLGen(db.OriginDB, sqlGen)
	db.check(t, sqlGen)
	return result, err
}

func (db *DBWrapper) CountBySQLStr(sqlStr string, args ...interface{}) (result int, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = CountBySQLStr(db.OriginDB, sqlStr, args...)
	db.check(t, sqlStr)
	return result, err
}

func (db *DBWrapper) InsertTableBySQLGen(sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = InsertTableBySQLGen(db.OriginDB, sqlGen)
	db.check(t, sqlGen)
	return
}

func (db *DBWrapper) InsertTableBySQLGenTx(tx *sql.Tx, sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = InsertTableBySQLGenTx(tx, sqlGen)
	db.check(t, sqlGen)
	return
}

func (db *DBWrapper) DeleteTableBySQLGen(sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = DeleteTableBySQLGen(db.OriginDB, sqlGen)
	db.check(t, sqlGen)
	return
}

func (db *DBWrapper) DeleteTableBySQLGenTx(tx *sql.Tx, sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = DeleteTableBySQLGenTx(tx, sqlGen)
	db.check(t, sqlGen)
	return
}

func (db *DBWrapper) UpdateTableBySQLGen(sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = UpdateTableBySQLGen(db.OriginDB, sqlGen)
	db.check(t, sqlGen)
	return
}

func (db *DBWrapper) UpdateTableBySQLGenTx(tx *sql.Tx, sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = UpdateTableBySQLGenTx(tx, sqlGen)
	db.check(t, sqlGen)
	return
}

func (db *DBWrapper) Exec(sqlStr string, args ...interface{}) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = Exec(db.OriginDB, sqlStr, args...)
	db.check(t, sqlStr)
	return
}

func (db *DBWrapper) ExecTx(tx *sql.Tx, sqlStr string, args ...interface{}) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = ExecTx(tx, sqlStr, args...)
	db.check(t, sqlStr)
	return
}

func (db *DBWrapper) GetRowBySQLGen(sqlGen *sqlutil.SQLGen) (result []string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowBySQLGen(db.OriginDB, sqlGen)
	db.check(t, sqlGen)
	return
}

func (db *DBWrapper) GetRowBySQLGenTx(tx *sql.Tx, sqlGen *sqlutil.SQLGen) (result []string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowBySQLGenTx(tx, sqlGen)
	db.check(t, sqlGen)
	return
}

func (db *DBWrapper) GetRowBySQLStr(sqlStr string, args ...interface{}) (result []string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowBySQLStr(db.OriginDB, sqlStr, args...)
	db.check(t, sqlStr)
	return
}

func (db *DBWrapper) GetRowsBySQLGenPrintSql(sqlGen *sqlutil.SQLGen) (result [][]string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowsBySQLGenPrintSql(db.OriginDB, sqlGen)
	db.check(t, sqlGen)
	return
}

func (db *DBWrapper) GetRowsBySQLGen(sqlGen *sqlutil.SQLGen) (result [][]string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowsBySQLGen(db.OriginDB, sqlGen)
	db.check(t, sqlGen)
	return
}

func (db *DBWrapper) GetRowsBySQLGenTx(tx *sql.Tx, sqlGen *sqlutil.SQLGen) (result [][]string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowsBySQLGenTx(tx, sqlGen)
	db.check(t, sqlGen)
	return
}

func (db *DBWrapper) GetRowsBySQLStr(sqlStr string, args ...interface{}) (result [][]string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowsBySQLStr(db.OriginDB, sqlStr, args...)
	db.check(t, sqlStr)
	return
}

func (db *DBWrapper) GetRowsBySQLStrWithQueryColumnsCount(sqlStr string, queryColumnsCount int) (result [][]string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowsBySQLStr(db.OriginDB, sqlStr)
	db.check(t, sqlStr)
	return
}
