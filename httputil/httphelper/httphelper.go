package httphelper

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const ContentType = "Content-Type"
const UserAgent = "User-Agent"

type HttpHelper struct {
	clt *http.Client
	err error
}

func NewHttpHelper() *HttpHelper {
	return &HttpHelper{clt: http.DefaultClient}
}

func (h *HttpHelper) NewProxyFunc(ip, port string) func(*http.Request) (*url.URL, error) {
	proxyAddr := func(_ *http.Request) (*url.URL, error) {
		uristr := fmt.Sprintf("http://%s:%s", ip, port)
		return url.Parse(uristr)
	}
	return proxyAddr
}

func (h *HttpHelper) NewProxyTransport(ip, port string) *http.Transport {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy:                 h.NewProxyFunc(ip, port)}
	return transport
}

func (h *HttpHelper) SetNewClient(ip string, port uint) *HttpHelper {
	h.clt = &http.Client{Transport: h.NewProxyTransport(ip, strconv.Itoa(int(port)))}
	return h
}

func (h *HttpHelper) GetClient() *http.Client {
	return h.clt
}

func (h *HttpHelper) SetNewClientWithTransport(trans *http.Transport) *HttpHelper {
	h.clt = &http.Client{Transport: trans}
	return h
}

func (h *HttpHelper) Get(urlText string) (resp string, httpCode int, respHeader http.Header) {
	return h.HttpRequestBase("GET", urlText, "", nil, true)
}

func (h *HttpHelper) GetWithHeader(urlText string, header map[string]string) (resp string, httpCode int, respHeader http.Header) {
	return h.HttpRequestBase("GET", urlText, "", header, true)
}

func (h *HttpHelper) GetNoRedirect(urlText string) (resp string, httpCode int, respHeader http.Header) {
	return h.HttpRequestBase("GET", urlText, "", nil, false)
}

func (h *HttpHelper) Post(urlText string, body string, header map[string]string) (resp string, httpCode int, respHeader http.Header) {
	return h.HttpRequestBase("POST", urlText, body, header, false)
}

func (h *HttpHelper) AutoTryGuessHeader(body string) (header map[string]string) {
	if body == "" {
		return
	}
	header = make(map[string]string)
	header[UserAgent] = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"
	header[ContentType] = "application/x-www-form-urlencoded"
	if string(body[0]) == "{" || string(body[0]) == "[" {
		header[ContentType] = "application/json"
	}
	return
}

func (h *HttpHelper) NewRequest(method, urlText, body string, header map[string]string) (*http.Request) {
	var request *http.Request
	var er error
	if len(body) > 0 {
		request, er = http.NewRequest(method, urlText, strings.NewReader(body))
	} else {
		request, er = http.NewRequest(method, urlText, nil)
	}
	if er != nil {
		return nil
	}
	for k, v := range h.AutoTryGuessHeader(body) { //尝试猜测需要加的header
		request.Header.Add(k, v)
	}
	for k, v := range header {
		request.Header.Add(k, v)
	}
	h.SetError(er)
	return request
}

func (h *HttpHelper) SetError(e error) *HttpHelper {
	h.err = e
	return h
}

func (h *HttpHelper) HasError() bool {
	return h.err != nil
}

func (h *HttpHelper) ClearError() *HttpHelper {
	h.err = nil
	return h
}

func (h *HttpHelper) HttpRequestBase(method, urlText, body string, header map[string]string, autoRedirect bool) (respText string, httpCode int, respHeader http.Header) {
	request := h.NewRequest(method, urlText, body, header)
	if h.HasError() {
		return
	}
	clt := h.clt
	var resp *http.Response
	if !autoRedirect {
		resp, h.err = clt.Transport.RoundTrip(request)
	} else {
		resp, h.err = clt.Do(request)
	}
	if h.HasError() || resp == nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var bodyByte []byte
	bodyByte, h.err = ioutil.ReadAll(resp.Body)
	return string(bodyByte), resp.StatusCode, resp.Header
}
