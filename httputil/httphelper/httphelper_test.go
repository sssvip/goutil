package httphelper

import (
	"net/http"
	"testing"
)

func TestHttpHelper_Timeout(t *testing.T) {
	httpHelper := NewHttpHelper()
	_, _, _, e := httpHelper.Get("http://192.168.2.200:5003/sleep")
	if e == nil {
		t.Error("timeout error")
	}
}

func TestHttpHelper_Get(t *testing.T) {
	httpHelper := NewHttpHelper()
	_, code, _, _ := httpHelper.Get("https://httpbin.org/ip")
	if code != http.StatusOK {
		t.Error("code error", code)
	}
}

func TestHttpHelper_GetWithProxy(t *testing.T) {
	// not test proxy now
	//httpHelper := NewHttpHelper()
	//httpHelper.SetNewClient("192.168.2.200", 1080)
	//_, code, _, _ := httpHelper.GetWithHeader("https://httpbin.org/ip", httpbuilder.NewHeader().Get())
	//if code != http.StatusOK {
	//	t.Error("code error", code)
	//
}

func BenchmarkHttpHelper_Get(b *testing.B) {
	for n := 0; n < b.N; n++ {
		httpHelper := NewHttpHelper()
		_, code, _, _ := httpHelper.Get("http://www.baidu.com")
		if code != http.StatusOK {
			b.Error("code error", code)
		}
	}

}

//go test -bench=. -benchmem
