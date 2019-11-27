package httphelper

import (
	"github.com/sssvip/goutil/httputil/httpbuilder"
	"github.com/sssvip/goutil/logutil"
	"net/http"
	"testing"
)

func TestHttpHelper_Get(t *testing.T) {
	httpHelper := NewHttpHelper()
	resp, code, _ := httpHelper.Get("https://httpbin.org/ip")
	if code != http.StatusOK {
		t.Error("code error", code)
	}
	logutil.Console.Println(resp)
}

func TestHttpHelper_GetWithProxy(t *testing.T) {
	httpHelper := NewHttpHelper()
	httpHelper.SetNewClient("192.168.2.200", 1080)
	resp, code, _ := httpHelper.GetWithHeader("https://httpbin.org/ip", httpbuilder.NewHeader().Get())
	if code != http.StatusOK {
		t.Error("code error", code)
	}
	logutil.Console.Println(resp)
}

func BenchmarkHttpHelper_Get(b *testing.B) {
	for n := 0; n < b.N; n++ {
		httpHelper := NewHttpHelper()
		_, code, _ := httpHelper.Get("http://www.baidu.com")
		if code != http.StatusOK {
			b.Error("code error", code)
		}
	}

}

//go test -bench=. -benchmem
