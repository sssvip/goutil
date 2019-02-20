package jsonutil

import (
	"github.com/json-iterator/go"
	"github.com/sssvip/goutil/logutil"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

/*func init() {
	extra.RegisterFuzzyDecoders()
}*/

func Marshal(a interface{}) []byte {
	b, e := json.Marshal(a)
	if e != nil {
		logutil.Error.Println(e)
		return nil
	}
	return b
}

func MarshalToString(a interface{}) string {
	b, e := json.MarshalToString(a)
	if e != nil {
		logutil.Error.Println(e)
		return ""
	}
	return b
}

func MarshalToIndentString(a interface{}, indent ...string) string {
	i := "  "
	if len(indent) > 0 {
		i = indent[0]
	}
	b, e := json.MarshalIndent(a, "", i)
	if e != nil {
		logutil.Error.Println(e)
		return ""
	}
	return string(b)
}

func UnmarshalFromString(text string, a interface{}) error {
	e := json.UnmarshalFromString(text, a)
	if e != nil {
		logutil.Error.Println(text, e)
	}
	return e
}
