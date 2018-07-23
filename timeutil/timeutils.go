package timeutil

import (
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
