package dbutil

import (
	"database/sql"
	"github.com/sssvip/goutil/dbutil/sqlutil"
	"github.com/sssvip/goutil/logutil"
	"github.com/sssvip/goutil/strutil"
	"github.com/sssvip/goutil/timeutil"
	"github.com/sssvip/goutil/timeutil/stopwatch"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*

基于dbutil 包装

*/

type SQLStatisticSummary struct {
	//统计起止时间
	StatisticStartTime string
	StatisticEndTime   string

	//SQL执行次数
	TotalCmdCount int64

	QueryCmdCount          int
	QueryLongestAveTimeSQL string //执行最长平均时间SQL
	QueryLongestAveTime    int64  //执行最长平均时间
	QueryMaxTimesSQL       string //执行次数最多
	QueryMaxTimesCount     int64

	//Insert
	InsertCmdCount          int
	InsertLongestAveTimeSQL string
	InsertLongestAveTime    int64
	InsertMaxTimesSQL       string
	InsertMaxTimesCount     int64

	//Update
	UpdateCmdCount          int
	UpdateLongestAveTimeSQL string
	UpdateLongestAveTime    int64
	UpdateMaxTimesSQL       string
	UpdateMaxTimesCount     int64

	//Delete
	DeleteCmdCount          int
	DeleteLongestAveTimeSQL string
	DeleteLongestAveTime    int64
	DeleteMaxTimesSQL       string
	DeleteMaxTimesCount     int64
}

type SQLStatistic struct {
	SQLStr                  string
	MaxExecTimeMilliSeconds int64
	MinExecMilliSeconds     int64
	ExecCount               int64
	ExecTotalTime           int64
	ExecAveTime             int64
}

func (sqlStatistic *SQLStatistic) ExecOnce(sqlStr string, useMilliSeconds int64) {
	sqlStatistic.SQLStr = sqlStr
	if sqlStatistic.MaxExecTimeMilliSeconds < useMilliSeconds {
		sqlStatistic.MaxExecTimeMilliSeconds = useMilliSeconds
	}
	if sqlStatistic.MinExecMilliSeconds > useMilliSeconds {
		sqlStatistic.MinExecMilliSeconds = useMilliSeconds
	}
	atomic.AddInt64(&sqlStatistic.ExecCount, 1)
	atomic.AddInt64(&sqlStatistic.ExecTotalTime, useMilliSeconds)
}

type DBWrapper struct {
	slowSQLSeconds     int
	OriginDB           *sql.DB
	openCheckSlowSQL   bool
	statisticSQL       bool
	statistics         sync.Map
	statisticStartTime string
}

var slowSQLCheck = "slowSQLCheck"

func (db *DBWrapper) SetSlowSQLSeconds(seconds int) {
	db.slowSQLSeconds = seconds
}

func (db *DBWrapper) OpenCheckSlowSQL(open bool) {
	db.openCheckSlowSQL = open
}

func (db *DBWrapper) StatisticSQL(open bool) {
	db.statisticSQL = open
	db.statisticStartTime = timeutil.FormatDateTimeString(time.Now())
}

func (db *DBWrapper) ClearStatistic() {
	db.statistics = sync.Map{}
	db.statisticStartTime = timeutil.FormatDateTimeString(time.Now())
}

func (db *DBWrapper) Statistic() (stats []SQLStatistic) {
	db.statistics.Range(func(key, value interface{}) bool {
		s := *value.(*SQLStatistic)
		s.ExecAveTime = s.ExecTotalTime / s.ExecCount
		stats = append(stats, s)
		return true
	})
	return
}
func (db *DBWrapper) StatisticSummary() (summary SQLStatisticSummary) {
	summary = SQLStatisticSummary{}
	summary.StatisticStartTime = db.statisticStartTime
	stats := db.Statistic()
	for _, s := range stats {
		summary.TotalCmdCount += s.ExecCount
		if strings.HasPrefix(s.SQLStr, "select") {
			summary.QueryCmdCount++
			if s.ExecCount > summary.QueryMaxTimesCount {
				summary.QueryMaxTimesSQL = s.SQLStr
				summary.QueryMaxTimesCount = s.ExecCount
			}
			if s.ExecAveTime > summary.QueryLongestAveTime {
				summary.QueryLongestAveTimeSQL = s.SQLStr
				summary.QueryLongestAveTime = s.ExecAveTime

			}
			continue
		}
		if strings.HasPrefix(s.SQLStr, "update") {
			summary.UpdateCmdCount++
			if s.ExecCount > summary.UpdateMaxTimesCount {
				summary.UpdateMaxTimesSQL = s.SQLStr
				summary.UpdateMaxTimesCount = s.ExecCount
			}
			if s.ExecAveTime > summary.UpdateLongestAveTime {
				summary.UpdateLongestAveTimeSQL = s.SQLStr
				summary.UpdateLongestAveTime = s.ExecAveTime

			}
			continue
		}
		if strings.HasPrefix(s.SQLStr, "insert") {
			summary.InsertCmdCount++
			if s.ExecCount > summary.InsertMaxTimesCount {
				summary.InsertMaxTimesSQL = s.SQLStr
				summary.InsertMaxTimesCount = s.ExecCount
			}
			if s.ExecAveTime > summary.InsertLongestAveTime {
				summary.InsertLongestAveTimeSQL = s.SQLStr
				summary.InsertLongestAveTime = s.ExecAveTime

			}
			continue
		}
		if strings.HasPrefix(s.SQLStr, "delete") {
			summary.DeleteCmdCount++
			if s.ExecCount > summary.DeleteMaxTimesCount {
				summary.DeleteMaxTimesSQL = s.SQLStr
				summary.DeleteMaxTimesCount = s.ExecCount
			}
			if s.ExecAveTime > summary.DeleteLongestAveTime {
				summary.DeleteLongestAveTimeSQL = s.SQLStr
				summary.DeleteLongestAveTime = s.ExecAveTime

			}
			continue
		}
	}
	summary.StatisticEndTime = timeutil.FormatDateTimeString(time.Now())
	return
}

func (db *DBWrapper) checkSQLStr(sw *stopwatch.StopWatch, sql string) {
	if sw == nil {
		return
	}
	if db.openCheckSlowSQL {
		execTime := sw.ElapsedMilliSeconds()
		//统计sql
		if db.statisticSQL {
			key := strutil.Md5(sql)
			v, ok := db.statistics.Load(key)
			if !ok {
				var stat SQLStatistic
				stat.MinExecMilliSeconds = execTime
				stat.MaxExecTimeMilliSeconds = execTime
				stat.ExecOnce(sql, execTime)
				db.statistics.Store(key, &stat)
			} else {
				v.(*SQLStatistic).ExecOnce(sql, execTime)
			}
		}
		if int(execTime)/1000 > db.slowSQLSeconds {
			logutil.Warning.Println(strutil.Format("found slow sql exec time [%d ms] sqlStr:%s", execTime, sql))
		}
	}
}

func NewDB(username, password, address, port, database string) *DBWrapper {
	return &DBWrapper{OriginDB: NewDBByArg(username, password, address, port, database),
		slowSQLSeconds: 10, openCheckSlowSQL: false, statisticSQL: false, statistics: sync.Map{}}
}

func (db *DBWrapper) QueryForObjectBySQLStr(sqlStr string, columnsAddress ...interface{}) error {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	e := QueryForObjectBySQLStr(db.OriginDB, sqlStr, columnsAddress...)
	if db.openCheckSlowSQL {
		db.checkSQLStr(t, sqlStr)
	}
	return e
}

func (db *DBWrapper) QueryForObject(sqlGen *sqlutil.SQLGen, columnsAddress ...interface{}) error {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	e := QueryForObject(db.OriginDB, sqlGen, columnsAddress...)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Query()
		db.checkSQLStr(t, sqlStr)
	}
	return e
}

func (db *DBWrapper) CountBySQLGen(sqlGen *sqlutil.SQLGen) (result int, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = CountBySQLGen(db.OriginDB, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Count()
		db.checkSQLStr(t, sqlStr)
	}
	return result, err
}

func (db *DBWrapper) CountBySQLStr(sqlStr string, args ...interface{}) (result int, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = CountBySQLStr(db.OriginDB, sqlStr, args...)
	if db.openCheckSlowSQL {
		db.checkSQLStr(t, sqlStr)
	}
	return result, err
}

func (db *DBWrapper) InsertTableBySQLGen(sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = InsertTableBySQLGen(db.OriginDB, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Insert()
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) InsertTableBySQLGenTx(tx *sql.Tx, sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = InsertTableBySQLGenTx(tx, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Insert()
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) DeleteTableBySQLGen(sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = DeleteTableBySQLGen(db.OriginDB, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Delete()
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) DeleteTableBySQLGenTx(tx *sql.Tx, sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = DeleteTableBySQLGenTx(tx, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Delete()
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) UpdateTableBySQLGen(sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = UpdateTableBySQLGen(db.OriginDB, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Update()
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) UpdateTableBySQLGenTx(tx *sql.Tx, sqlGen *sqlutil.SQLGen) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = UpdateTableBySQLGenTx(tx, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Update()
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) Exec(sqlStr string, args ...interface{}) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = Exec(db.OriginDB, sqlStr, args...)
	db.checkSQLStr(t, sqlStr)
	return
}

func (db *DBWrapper) ExecTx(tx *sql.Tx, sqlStr string, args ...interface{}) (result int64, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = ExecTx(tx, sqlStr, args...)
	if db.openCheckSlowSQL {
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) GetRowBySQLGen(sqlGen *sqlutil.SQLGen) (result []string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowBySQLGen(db.OriginDB, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Query()
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) GetRowBySQLGenTx(tx *sql.Tx, sqlGen *sqlutil.SQLGen) (result []string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowBySQLGenTx(tx, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Query()
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) GetRowBySQLStr(sqlStr string, args ...interface{}) (result []string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowBySQLStr(db.OriginDB, sqlStr, args...)
	if db.openCheckSlowSQL {
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) GetRowsBySQLGenPrintSql(sqlGen *sqlutil.SQLGen) (result [][]string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowsBySQLGenPrintSql(db.OriginDB, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Query()
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) GetRowsBySQLGen(sqlGen *sqlutil.SQLGen) (result [][]string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowsBySQLGen(db.OriginDB, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Query()
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) GetRowsBySQLGenTx(tx *sql.Tx, sqlGen *sqlutil.SQLGen) (result [][]string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowsBySQLGenTx(tx, sqlGen)
	if db.openCheckSlowSQL {
		sqlStr, _, _ := sqlGen.Query()
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) GetRowsBySQLStr(sqlStr string, args ...interface{}) (result [][]string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowsBySQLStr(db.OriginDB, sqlStr, args...)
	if db.openCheckSlowSQL {
		db.checkSQLStr(t, sqlStr)
	}
	return
}

func (db *DBWrapper) GetRowsBySQLStrWithQueryColumnsCount(sqlStr string, queryColumnsCount int) (result [][]string, err error) {
	var t *stopwatch.StopWatch
	if db.openCheckSlowSQL {
		t = stopwatch.NewStopWatch(slowSQLCheck)
	}
	result, err = GetRowsBySQLStr(db.OriginDB, sqlStr)
	if db.openCheckSlowSQL {
		db.checkSQLStr(t, sqlStr)
	}
	return
}
