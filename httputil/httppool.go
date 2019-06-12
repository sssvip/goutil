package httputils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
	//"github.com/Workiva/go-datastructures/queue"
	"github.com/sssvip/goutil/logutil"
	"net"
	"strings"
	"sync/atomic"
)

var AUTO_TRY = true
var PRINT_HTTP_NOT_OK_LOG = true

const HttpErrorCode = 500

const ProxyError = "ProxyError"

const HttpReqError = "HttpReqError"

var transChan = make(chan *TransWithExpireTime, 1000)

var clientChan = make(chan *http.Client, 1000)

//var clientQueue = queue.Queue{}

var totalClient int32 = 0

//var UseCacheClient = true

var defaultRetryTimes = 3

//var defaultProxyRetryTimes = 15

func PoolStatistic() string {
	return fmt.Sprintf("transChan:%d\nclientChan:%d\ntotalClient:%d\n", len(transChan), len(clientChan), totalClient)
}

type TransWithExpireTime struct {
	trans      *http.Transport
	expireTime time.Time
}

func (t TransWithExpireTime) Produce(trans *http.Transport, expireTime time.Time) *TransWithExpireTime {
	return &TransWithExpireTime{trans: trans, expireTime: expireTime}
}
func (t TransWithExpireTime) ProduceBasic(ip, port string, expireTime time.Time) *TransWithExpireTime {
	return &TransWithExpireTime{trans: ProxyTransport(ip, port), expireTime: expireTime}
}

type ProxyDTO struct {
	ERRORCODE string `json:"ERRORCODE"`
	RESULT    []struct {
		IP   string `json:"ip"`
		Port string `json:"port"`
	} `json:"RESULT"`
}

var makeProxyRootURL = ""

func SetProxyRootURL(proxyURL string) {
	makeProxyRootURL = proxyURL
}

func ProxyTransport(ip, port string) *http.Transport {
	proxyAddr := func(_ *http.Request) (*url.URL, error) {
		uristr := fmt.Sprintf("http://%s:%s", ip, port)
		return url.Parse(uristr)
	}
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy:                 proxyAddr}
	return transport
}

func DefaultProduceProxyTrans() {
	if makeProxyRootURL == "" {
		logutil.Error.Println("please set makeProxyRootUrl before use proxy")
		return
	}
	resp, err := http.Get(makeProxyRootURL)
	if err != nil {
		return
	}
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		logutil.Error.Println(err)
		return
	}
	bodyStr := string(body)
	err = resp.Body.Close()
	if err != nil {
		logutil.Error.Println(err)
		return
	}
	if bodyStr == "" {
		return
	}
	var proxy = ProxyDTO{}
	err = json.Unmarshal([]byte(bodyStr), &proxy)
	if err != nil {
		logutil.Error.Println(err)
		return
	}
	if proxy.RESULT != nil {
		for _, data := range proxy.RESULT {
			transport := ProxyTransport(data.IP, data.Port)
			transChan <- &TransWithExpireTime{transport, time.Now()}
		}
	}
}

type CustomProxyProduceFunc = func() *TransWithExpireTime

var customProxyFuncs = []CustomProxyProduceFunc{}

func RegisterCustomProxyProduceFunc(f CustomProxyProduceFunc) {
	customProxyFuncs = append(customProxyFuncs, f)
}

func getTransFromQueue() *http.Transport {
	if len(transChan) == 0 {
		time.Sleep(200 * time.Millisecond)
		if len(transChan) == 0 {
			if len(customProxyFuncs) > 0 {
				go func() {
					for _, pf := range customProxyFuncs {
						transChan <- pf()
					}
				}()
			} else {
				go DefaultProduceProxyTrans()
			}
		}
	}
	twt := <-transChan
	if twt.expireTime.After(time.Now()) {
		return twt.trans
	} else {
		return getTransFromQueue()
	}
}

func MakeNewClient() {
	atomic.AddInt32(&totalClient, 1)
	clt := &http.Client{Transport: getTransFromQueue()}
	clientChan <- clt
}

func Release(clt *http.Client, valid bool, desc string) {
	//无效后更换新的代理
	if !valid {
		logutil.Error.Println("client invalid....descpition:{}", desc)
		clt.Transport = getTransFromQueue()
	}
	clientChan <- clt
}

func Require() *http.Client {
	if len(clientChan) == 0 {
		//log.Println("client is empty.......")
		time.Sleep(100 * time.Millisecond)
		if len(clientChan) == 0 {
			go MakeNewClient()
		}
	}
	return <-clientChan
}

type ReqHelper struct {
	resp *http.Response
	err  error
	clt  *http.Client
}

func Get(url string) (string, int, http.Header) {
	return HttpBase("GET", url, "", false, defaultRetryTimes, nil, nil, false)
}

func GetWithWarning(url string) (string, int, http.Header) {
	return HttpBaseWithWarning("GET", url, "", false, defaultRetryTimes, nil, nil, false, true)
}

//不自动跳转
func SimpleGet(url string) (body string, statusCode int, respHeader http.Header) {
	return HttpBase("GET", url, "", false, defaultRetryTimes, nil, nil, true)
}

func ProxyGet(url string) (string, int, http.Header) {
	return HttpBase("GET", url, "", true, defaultRetryTimes, nil, nil, false)
}

func GetWithHeader(url string, headers map[string]string) (string, int, http.Header) {
	return HttpBase("GET", url, "", false, defaultRetryTimes, nil, headers, false)
}

//仅url必要
func Post(url string, bodyArgs string, headers map[string]string) (string, int, http.Header) {
	return HttpBase("POST", url, bodyArgs, false, defaultRetryTimes, nil, headers, false)
}

const ContentType = "Content-Type"
const UserAgent = "User-Agent"

func HttpBase(method string, url string, body string, useProxy bool, retryTime int, expectTexts []string, headers map[string]string, noAutoRedirect bool) (string, int, http.Header) {
	return HttpBaseWithWarning(method, url, body, useProxy, retryTime, expectTexts, headers, noAutoRedirect, false)
}

//body 可用NewBodyArgs生成
//expectTexts:满足一个text则正常请求
//返回respBodyStr,StatusCode,respHeader
func HttpBaseWithWarning(method string, url string, body string, useProxy bool, retryTime int, expectTexts []string, headers map[string]string, noAutoRedirect, printWarning bool) (string, int, http.Header) {
	reqHelper := ReqHelper{}
	var request *http.Request
	var er error
	if len(body) > 0 {
		request, er = http.NewRequest(method, url, strings.NewReader(body))
	} else {
		request, er = http.NewRequest(method, url, nil)
	}
	if er != nil {
		logutil.Error.Println(er)
		return "Create Request Error", 500, nil
	}
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}
	if len(body) > 0 && request.Header.Get(ContentType) == "" {
		if string(body[0]) == "{" {
			request.Header.Add(ContentType, "application/json")
		} else {
			request.Header.Add(ContentType, "application/x-www-form-urlencoded")
		}
	}
	//agent
	if request.Header.Get(UserAgent) == "" {
		//设置默认
		request.Header.Add(UserAgent, "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36")
	}
	if useProxy {
		clt := Require()
		reqHelper.clt = clt
		reqHelper.resp, reqHelper.err = reqHelper.clt.Do(request)
	} else if noAutoRedirect {
		reqHelper.resp, reqHelper.err = http.DefaultTransport.RoundTrip(request)
	} else {
		reqHelper.resp, reqHelper.err = http.DefaultClient.Do(request)
	}
	if reqHelper.resp != nil {
		defer func() {
			e := reqHelper.resp.Body.Close()
			if e != nil {
				logutil.Error.Println(e)
			}
		}()
	}
	if reqHelper.resp == nil || reqHelper.err != nil || reqHelper.resp.StatusCode == 403 {
		if useProxy {
			Release(reqHelper.clt, false, "403 or no resp body")
		}
		if reqHelper.err != nil {
			logutil.Error.Println(reqHelper.err, body)
		}
		if AUTO_TRY && retryTime > 0 {
			retryTime--
			return HttpBase(method, url, body, useProxy, retryTime, expectTexts, headers, noAutoRedirect)
		} else {
			if useProxy && reqHelper.err != nil && (strings.Contains(reqHelper.err.Error(), "timeout") || strings.Contains(reqHelper.err.Error(), "proxyconnect")) {
				return ProxyError, HttpErrorCode, nil
			}
			if reqHelper.resp != nil {
				bodyByte, _ := ioutil.ReadAll(reqHelper.resp.Body)
				return string(bodyByte), reqHelper.resp.StatusCode, reqHelper.resp.Header
			}
			return HttpReqError, HttpErrorCode, nil
		}
	}
	statusCode := reqHelper.resp.StatusCode
	bodyByte, e := ioutil.ReadAll(reqHelper.resp.Body)
	bodyStr := string(bodyByte)
	if PRINT_HTTP_NOT_OK_LOG && printWarning && statusCode != 200 {
		logutil.Warning.Printf("http status warn:"+strconv.Itoa(reqHelper.resp.StatusCode)+" %s,%s ,data:%s", method, url, bodyStr)
	}
	if useProxy {
		if !containsTexts(bodyStr, expectTexts) {
			Release(reqHelper.clt, false, fmt.Sprintf("not contain expect string: %s", expectTexts))
		} else {
			Release(reqHelper.clt, true, "")
		}
	}
	if e != nil {
		logutil.Error.Println(e)
	}
	return bodyStr, statusCode, reqHelper.resp.Header
}

func containsTexts(resp string, checkTexts []string) bool {
	if checkTexts == nil {
		return true
	} else {
		for _, d := range checkTexts {
			if strings.Contains(resp, d) {
				return true
			}
		}
		return false
	}
}

var SogouIndexHeaders = make(map[string]string)

var WxIndexHeaders = make(map[string]string)

func init() {
	SogouIndexHeaders["Referer"] = "http://index.sogou.com/index/media/wechat?kwdNamesStr=%E7%BE%8E%E4%B8%BD&timePeriodType=MONTH&dataType=MEDIA_WECHAT&queryType=INPUT"
	SogouIndexHeaders["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	SogouIndexHeaders["Accept-Encoding"] = "gzip, deflate"
	SogouIndexHeaders["Accept-Language"] = "zh-CN,zh;q=0.9"
	SogouIndexHeaders["Connection"] = "keep-alive"
	SogouIndexHeaders["DNT"] = "1"
	SogouIndexHeaders["Host"] = "index.sogou.com"
	SogouIndexHeaders["Upgrade-Insecure-Requests"] = "1"
	WxIndexHeaders["Referer"] = "https://servicewechat.com/wxc026e7662ec26a3a/4/page-frame.html"
}
