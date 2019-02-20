package logutil

import (
	"io"
	"log"
	"os"
)

var (
	Debug   *log.Logger
	Console *log.Logger
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 非常严重的问题

)
var formatMask = log.Ldate | log.Ltime | log.Lshortfile

//默认初始化不写入文件
func init() {
	Debug = log.New(os.Stdout, "DEBUG", formatMask)
	Console = log.New(os.Stdout, "CONSOLE: ", formatMask)
	Info = log.New(os.Stdout, "INFO: ", formatMask)
	Warning = log.New(os.Stderr, "WARNING: ", formatMask)
	Error = log.New(os.Stderr, "ERROR: ", formatMask)
}

//日志持久化到文件
func LogToFile() {
	LogToFileSpec("log.txt")
}

func LogToFileSpec(fileName string) {
	var logFile *os.File
	var err error
	logFile, err = os.OpenFile(fileName,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	//前期都先写入文件
	Info = log.New(io.MultiWriter(logFile, os.Stdout), "INFO: ", formatMask)
	Warning = log.New(io.MultiWriter(logFile, os.Stderr), "WARNING: ", formatMask)
	Error = log.New(io.MultiWriter(logFile, os.Stderr), "ERROR: ", formatMask)
}
