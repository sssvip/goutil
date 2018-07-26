package goutil

import (
	"goutil/timeutil/stopwatch"
	"goutil/strutil"
	"goutil/logutil"
	"goutil/dbutil"
	"goutil/dbutil/sqlutil"
)

func Example() {
	//sw
	stopwatch.NewStopWatch("test")
	//fmt
	strutil.Format("%s", "hello")
	//console
	logutil.Console.Println("heloo")
	//db
	dbutil.NewDB("", "", "", "", "")
	//sql
	sqlutil.Example()
}
