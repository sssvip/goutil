package testutil

import (
	"fmt"
	"github.com/sssvip/goutil/strutil"
	"github.com/sssvip/goutil/timeutil"
	"github.com/sssvip/goutil/timeutil/stopwatch"
	"github.com/xcltapestry/xclpkg/clcolor"
	"os"
	"runtime"
	"time"
	"unicode/utf8"
)

var successChar = clcolor.Green("✔️")
var failedChar = clcolor.Red("✖️")
var swAll = stopwatch.NewStopWatch("swAll")
var stepCount = 0
var allShowTryTextLen = 0
var backupTryTimeNum = 0

func init() {
	if IsWindows() {
		successChar = "✔️"
		failedChar = "✖️"
	}
}
func RedErrorStr(text string) string {
	if !IsWindows() {
		return clcolor.Red(text)
	}
	return text
}
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func BackTimes(tryTimes int) {
	for {
		if tryTimes <= 0 {
			break
		}
		tryTimes--
		fmt.Print("\b")
	}
}

func ShowTryText(time int) {
	ClearShowTry()
	str := strutil.Format("try %d times ", time)
	backupTryTimeNum = time
	allShowTryTextLen += utf8.RuneCountInString(str)
	fmt.Print(clcolor.Yellow(str))
}

func PrintLine(text string) {
	Print(text + "\n")
}

func Print(text string) {
	ClearShowTry()
	fmt.Print(text)
	ReShowTryTimes()
}

func ReShowTryTimes() {
	if backupTryTimeNum == 0 {
		return
	}
	ShowTryText(backupTryTimeNum)
}

func ClearShowTry() {
	BackTimes(allShowTryTextLen)
	allShowTryTextLen = 0
}

func TryMoreTime(f func() error, times int, name string, periodPerExecMill ...int) {
	periodPerExecReal := 1000
	if len(periodPerExecMill) > 0 {
		periodPerExecReal = periodPerExecMill[0]
	}
	stepCount++
	sw := stopwatch.NewStopWatch("t")
	fmt.Print(strutil.Format("%d.[%s]\n", stepCount, name))
	err := f()
	tryTimes := 1
	for {
		times--
		if times <= 0 {
			break
		}
		if err == nil {
			break
		}
		tryTimes++
		ShowTryText(tryTimes)
		err = f()
		timeutil.Sleep(periodPerExecReal)
	}
	//回退
	ClearShowTry()
	timeStr := "time"
	if tryTimes > 1 {
		timeStr += "s"
	}
	useTimeStr := strutil.Format("current step use %dms(%d %s),total use %ds", sw.ElapsedMilliSeconds(), tryTimes, timeStr, swAll.ElapsedSeconds())
	if err != nil {
		fmt.Println(strutil.Format("current step [%s] test result: %s", name, failedChar))
		fmt.Println(strutil.Format("Failed reason:[%s], please check...", RedErrorStr(err.Error())))
		fmt.Println(useTimeStr)
		os.Exit(-1)
	} else {
		fmt.Println(strutil.Format("current step [%s] test result: %s", name, successChar))
		fmt.Println(useTimeStr)
	}
}

func StartTest() {
	fmt.Println(strutil.Format("start to test all,now:%s", timeutil.FormatDateTimeString(time.Now())))
}
func EndTest() {
	endTimeStr := timeutil.FormatDateTimeString(time.Now())
	fmt.Println(strutil.Format("pass all tests... total use %ds, end time:%s", swAll.ElapsedSeconds(), endTimeStr))
}
