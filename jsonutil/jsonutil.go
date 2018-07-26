package jsonutil

import (
	"github.com/json-iterator/go"
	"goutil/logutil"
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

func UnmarshalFromString(text string, a interface{}) error {
	e := json.UnmarshalFromString(text, a)
	if e != nil {
		logutil.Error.Println(text, e)
	}
	return e
}
