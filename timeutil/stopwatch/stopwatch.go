package stopwatch

import (
	"github.com/sssvip/goutil/logutil"
	"github.com/sssvip/goutil/strutil"
	"github.com/sssvip/goutil/timeutil"
	"time"
)

func CalcTime(t *time.Time, explain ...interface{}) {
	d := (time.Now().UnixNano() - t.UnixNano()) / 1000 / 1000
	if len(explain) > 0 {
		logutil.Info.Println(strutil.Format("[%v] use time %d ms", explain[0], d))
	} else {
		logutil.Info.Println(strutil.Format("use time %d ms", d))
	}
}

func CalcFuncUseTime(f func(), explain ...interface{}) {
	now := time.Now()
	f()
	if len(explain) < 1 {
		explain = append(explain, "func()")
	}
	CalcTime(&now, explain...)
}

type StopWatch struct {
	name      string
	startTime time.Time
}

func NewStopWatch(name ...string) *StopWatch {
	sName := ""
	if len(name) < 1 || name[0] == "" {
		sName = strutil.Format("started at %s", timeutil.FormatDateTimeString(time.Now()))
	} else {
		sName = name[0]
	}
	return &StopWatch{name: sName, startTime: time.Now()}
}

func (s *StopWatch) GetName() string {
	return s.name
}

func (s *StopWatch) ResetNameAndRestoreStartTime(name string) {
	s.name = name
	s.RestoreStartTime()
}

func (s *StopWatch) ResetName(name string) {
	s.name = name
}

func (s *StopWatch) RestoreStartTime() {
	s.startTime = time.Now()
}

func (s *StopWatch) ElapsedMilliSeconds() int64 {
	return (time.Now().UnixNano() - s.startTime.UnixNano()) / 1000 / 1000
}

func (s *StopWatch) ElapsedSeconds() int64 {
	return s.ElapsedMilliSeconds() / 1000
}

//ConsoleElapsedMilliSeconds 仅仅输出到std out
func (s *StopWatch) ConsoleElapsedMilliSeconds(alterArg ...string) {
	arg := ""
	if len(alterArg) > 0 {
		arg = strutil.Format(", alter arg: %s", alterArg[0])
	}
	logutil.Console.Println(strutil.Format("StopWatch [%s] elapsed %d ms%s", s.name, s.ElapsedMilliSeconds(), arg))
}

//ConsoleElapsedMilliSeconds 可能会输出到文件（在设置非debug模式下）
func (s *StopWatch) LogElapsedMilliSeconds(alterArg ...string) {
	arg := ""
	if len(alterArg) > 0 {
		arg = strutil.Format(", alter arg: %s", alterArg[0])
	}
	logutil.Info.Println(strutil.Format("StopWatch [%s] elapsed %d ms%s", s.name, s.ElapsedMilliSeconds(), arg))
}
