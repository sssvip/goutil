package strutil

import (
	"crypto/md5"
	"io"
	"strings"
	"encoding/hex"
	"fmt"
	"strconv"
	"github.com/google/uuid"
	"goutil/logutil"
)

//NewUUID 返回UUID
func NewUUID() string {
	return uuid.New().String()
}

//Md5 尝试md5加密,失败返回空字符串
func Md5(data string) string {
	h := md5.New()
	_, e := io.Copy(h, strings.NewReader(data))
	if e != nil {
		logutil.Error.Println(e)
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

//Format 简单格式化包装
func Format(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

//AtoIWithDefault 转换字符串道int 失败返回默认值
func AtoIWithDefault(text string, defaultValue int) int {
	i, e := strconv.Atoi(text)
	if e != nil {
		return defaultValue
	}
	return i
}

//AtoFloat64WithDefault 转换字符串道float64 失败返回默认值
func AtoFloat64WithDefault(text string, defaultValue float64) float64 {
	f, e := strconv.ParseFloat(text, 64)
	if e != nil {
		return defaultValue
	}
	return f
}

//SafeCutString 按需截断字符串,如果text超过小于size,返回text,否则按需截取
func SafeCutString(text string, size int) string {
	if text == "" || size < 1 {
		return text
	}
	tempRune := []rune(text)
	if len(tempRune) > size {
		return string(tempRune[:size])
	}
	return text
}
