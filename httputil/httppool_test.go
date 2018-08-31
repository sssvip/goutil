package httputils

import (
	"fmt"
	"github.com/sssvip/goutil/jsonutil"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSimpleGet(t *testing.T) {
	//fmt.Println(SimpleGet("http://192.168.2.200/8wv0Q"))
}

func TestGet(t *testing.T) {
	body, code, _ := Get("http://httpbin.org/ip")
	assert.Equal(t, http.StatusOK, code)
	assert.Contains(t, body, "origin")
}

func TestPost(t *testing.T) {
	//暂时不测
}

type ReqIP struct {
	Origin string "json:origin"
}

func TestHttpBaseWithProxy(t *testing.T) {
	body, _, _ := Get("http://httpbin.org/ip")
	var realIP ReqIP
	assert.Nil(t, jsonutil.UnmarshalFromString(body, &realIP))
	fmt.Println(realIP.Origin)
	var proxyIP ReqIP
	body2, _, _ := ProxyGet("http://httpbin.org/ip")
	assert.Nil(t, jsonutil.UnmarshalFromString(body2, &proxyIP))
	fmt.Println(proxyIP.Origin)
	assert.NotEqual(t, realIP.Origin, proxyIP.Origin)
}
