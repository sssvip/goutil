package goutil

import (
	"github.com/sssvip/goutil/timeutil/stopwatch"
	"github.com/sssvip/goutil/strutil"
	"github.com/sssvip/goutil/logutil"
	"github.com/sssvip/goutil/dbutil"
	"github.com/sssvip/goutil/dbutil/sqlutil"
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
