package timeutil

import (
	"time"
	"fmt"
	"github.com/sssvip/goutil/logutil"
)

const TimeFormatStr = "2006-01-02 15:04:05"
const DateFormatStr = "2006-01-02"

//second
func FormatTimeStringByTimestamp(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format(TimeFormatStr)
}

func FormatDateTimeString(t time.Time) string {
	return t.Format(TimeFormatStr)
}

//检验日期格式 如 2018-02-12
func ParseTimeByDateStr(dateStr string) *time.Time {
	t, e := time.Parse(DateFormatStr, dateStr)
	if e != nil {
		logutil.Error.Println(e)
		return nil
	}
	return &t
}

func ParseTimeByDateTimeStr(dateTimeStr string) *time.Time {
	t, e := time.Parse(TimeFormatStr, dateTimeStr)
	if e != nil {
		logutil.Error.Println(e)
		return nil
	}
	return &t
}

func FormatDateString(t time.Time) string {
	d, y, m := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", d, y, m)
}
