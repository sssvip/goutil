package timeutil

import (
	"github.com/sssvip/goutil/strutil"
	"time"
)

//毫秒
func Sleep(millSeconds int) {
	time.Sleep(time.Duration(millSeconds) * time.Millisecond)
}
func SleepSeconds(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func Sleep2Tomorrow(hour int, min int) {
	now := time.Now()
	// 计算下一个零点
	tomorrow := now.Add(time.Hour * 24)
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), hour, min, 0, 0, tomorrow.Location())
	t := time.NewTimer(tomorrow.Sub(now))
	<-t.C
}

func UTCNow() time.Time {
	t, _ := time.Parse(TimeFormatStr, time.Now().Format(TimeFormatStr))
	return t
}

func Add(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

func Sleep2NextHour(min, second int) {
	now := time.Now()
	// 计算下一个零点
	nextHour := now.Add(time.Hour)
	nextHour = time.Date(nextHour.Year(), nextHour.Month(), nextHour.Day(), nextHour.Hour(), min, second, 0, nextHour.Location())
	sleepSeconds := int(nextHour.Unix() - now.Unix())
	SleepSeconds(sleepSeconds)
}
func UnixNanoStrWithLen(numLen int) string {
	s := strutil.Format("%v", time.Now().UnixNano())
	return s[:numLen]
}
