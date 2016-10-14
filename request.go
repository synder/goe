package goe

import (
	"net/http"
	"strings"
)

type Request struct {
	*http.Request

	params map[string]string
}


func (req *Request) Ip() string{
	return req.Get("x-forwarded-for")
}

func (req *Request) ClientIP() string {
	return req.RemoteAddr
}

func (req *Request) Xhr() bool{
	var val = req.Get("x-requested-with");
	return val != "" && strings.ToLower(val) == "xmlhttprequest";
}

func (req *Request) Secure() bool{
	return strings.Contains(strings.ToLower(req.Proto), "https")
}

func (req *Request) Set(key, value string) {
	req.Header.Set(key, value)
}

func (req *Request) Get(key string) string {
	key = strings.ToLower(key)
	return req.Headers()[key]
}

func (req *Request) Headers() map[string]string{
	headers := make(map[string]string)

	for key, value := range req.Header {
		headers[strings.ToLower(key)] = strings.Join(value, "")
	}

	return headers
}

func (req *Request) Query(key string) string  {
	return req.URL.Query().Get(key)
}

func (req *Request) Querys() map[string][]string{
	return req.URL.Query()
}

func (req *Request) Cookie(key string) *http.Cookie{
	cookies := req.Request.Cookies()

	for _, cookie := range cookies{
		if strings.ToLower(key) == strings.ToLower(cookie.Name){
			return cookie
		}
	}
	return nil
}

func (req *Request) Cookies(key string) []*http.Cookie {
	return req.Request.Cookies()
}

func (req *Request) Param(key string) string {
	return req.params[key]
}

func (req *Request) Params(key string) map[string]string {
	return req.params
}


func NewRequest(req *http.Request) *Request {
	return &Request{
		Request: req,
		params: nil,
	}
}