package httphelper

import (
	"net/http"
	"testing"
)

func TestHttpHelper_Get(t *testing.T) {
	httpHelper := NewHttpHelper()
	_, code, _ := httpHelper.Get("http://www.baidu.com")
	if code != http.StatusOK {
		t.Error("code error", code)
	}
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
