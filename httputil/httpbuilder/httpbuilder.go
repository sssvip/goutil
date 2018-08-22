package httpbuilder

import (
	"net/url"
)

//便于生成请求体参数

type BodyArgs struct {
	Values url.Values
}

func NewBodyArgs() *BodyArgs {
	return &BodyArgs{Values: url.Values{}}
}

func (arg *BodyArgs) Add(key, value string) *BodyArgs {
	arg.Values.Add(key, value)
	return arg
}

func (arg *BodyArgs) Encode() string {
	return arg.Values.Encode()
}

//便于生成请求头
type HttpHeader struct {
	header map[string]string
}

func NewHeader() *HttpHeader {
	return &HttpHeader{header: make(map[string]string)}
}

func (httpHeader *HttpHeader) Add(key, value string) *HttpHeader {
	httpHeader.header[key] = value
	return httpHeader
}
func (httpHeader *HttpHeader) AddCookie(value string) *HttpHeader {
	httpHeader.header["Cookie"] = value
	return httpHeader
}
func (httpHeader *HttpHeader) AddUserAgent(value string) *HttpHeader {
	httpHeader.header["User-Agent"] = value
	return httpHeader
}
func (httpHeader *HttpHeader) AddReferer(value string) *HttpHeader {
	httpHeader.header["Referer"] = value
	return httpHeader
}
func (httpHeader *HttpHeader) AddContentType(value string) *HttpHeader {
	httpHeader.header["Content-Type"] = value
	return httpHeader
}

func (httpHeader *HttpHeader) JSONContentType() string {
	return "application/json"
}

func (httpHeader *HttpHeader) URLEncodedContentType() string {
	return "application/x-www-form-urlencoded"
}
func (httpHeader *HttpHeader) Get() map[string]string {
	return httpHeader.header
}
