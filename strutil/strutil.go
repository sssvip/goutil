package strutil

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/sssvip/goutil/logutil"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
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

//AtoIWithDefault 转换字符串道int 失败返回默认值
func AtoInt64WithDefault(text string, defaultValue int64) int64 {
	i, e := strconv.ParseInt(text, 10, 64)
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

//GetStrByRegexp 安全返回正则结果(正则分组返回结果，groups 0-> 包含正则条件中的串， 1 正则括号中的干净值)
func GetStrByRegexp(re *regexp.Regexp, text string, groups ...int) []string {
	arrRst := make([]string, len(groups))
	//正则实际结果
	arr := re.FindStringSubmatch(text)
	for index, data := range groups {
		if data > -1 && data < len(arr) {
			arrRst[index] = arr[data]
		}
	}
	return arrRst
}
func ArrayRandom(arr []string) (idx int, randStr string) {
	if len(arr) < 1 {
		return 0, ""
	}
	idx = rand.Intn(len(arr))
	return idx, arr[idx]
}
func ArrayRandomValue(arr []string) string {
	_, v := ArrayRandom(arr)
	return v
}
func RandNumStr(textLen int) string {
	var buffer bytes.Buffer
	for i := 0; i < textLen; i++ {
		buffer.WriteString(strconv.Itoa(rand.Intn(9) + 1))
	}
	return buffer.String()[:textLen]
}

func RandNumAlphabet(textLen int) string { //32位以内
	if textLen < 1 {
		return ""
	}
	cnt := textLen/32 + 1
	var buffer bytes.Buffer
	for i := 0; i < cnt; i++ {
		buffer.WriteString(strings.Replace(NewUUID(), "-", "", -1))
	}
	return buffer.String()[:textLen]
}
