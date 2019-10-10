package executil

import (
	"github.com/sssvip/goutil/strutil"
	"os"
	"os/exec"
)

func GetShellPath() string {
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh, err := exec.LookPath("sh")
		if err == nil {
			return sh
		}
		sh = "/system/bin/sh"
	}
	return sh
}

func ExecBase(name string, args ...string) (output string, err error) {
	cmd := exec.Command(name, args...)
	o, err := cmd.CombinedOutput()
	return string(o), err
}

func ExecWithError(command string, args ...interface{}) (result string, err error) {
	realCmd := ""
	if len(args) < 1 {
		realCmd = command
	} else {
		realCmd = strutil.Format(command, args...)
	}
	var newArr []string
	newArr = append(newArr, "-c")
	newArr = append(newArr, realCmd)
	return ExecBase(GetShellPath(), newArr...)
}

type ExecResp struct {
	Error    bool   //是否错误
	ErrorStr string //错误信息
	Content  string // 正常信息
}

func Exec(command string, args ...interface{}) ExecResp {
	var m ExecResp
	r, e := ExecWithError(command, args...)
	if e != nil {
		m.ErrorStr = e.Error()
	}
	m.Error = e != nil
	m.Content = r
	return m
}
