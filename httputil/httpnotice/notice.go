package httpnotice

import (
	"strings"
	"os"
	"github.com/sssvip/goutil/logutil"
	"github.com/sssvip/goutil/httputil/httpbuilder"
	"github.com/sssvip/goutil/httputil"
	"github.com/sssvip/goutil/timeutil"
	"github.com/sssvip/goutil/strutil"
	"sync"
)

var osName string

//noticeRootURL 发送通知的地址,自我实现
var noticeRootURL = ""

func init() {
	var err error
	osName, err = os.Hostname()
	if err != nil {
		osName = "Unknown Hostname"
		logutil.Error.Println(err)
	}
}

func SetNoticeUrl(noticeURL string) {
	noticeRootURL = noticeURL
}

//SendDavid 自用方法
func SendDavid(content, url string) {
	SendNotice("o44U9wt6O9P5j9M0L47RaGpxfe2o", content, url)
}

func SendNotice(openId, content, dstUrl string) (success bool) {
	if noticeRootURL == "" {
		logutil.Error.Println("please set noticeURL")
		return
	}
	data := httpbuilder.NewBodyArgs().
		Add("name", "golang notice").
		Add("openid", openId).
		Add("message", "默认内容").
		Add("rank", "1").
		Add("space", osName).
		Add("url", dstUrl).
		Add("remark", content).
		Add("content", "111")
	body, _, _ := httputils.Post(noticeRootURL, data.Encode(), nil)
	if !strings.Contains(body, `"errmsg":"ok"`) {
		logutil.Error.Println(body)
		return true
	} else {
		logutil.Info.Println(body)
	}
}

func keepAlive() {
	go func() {
		for {
			logutil.Info.Println(strutil.Format("%s 定时报平安", project))
			SendDavid(strutil.Format("%s I am live", project), "")
			timeutil.Sleep2Tomorrow(9, 0)
		}
	}()
}

var project string
var once sync.Once

func KeepAliveNotice(projectName string) {
	project = projectName
	once.Do(keepAlive)
}
