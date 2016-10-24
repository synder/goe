package goe

import (
	"net/http"
	"strings"
)

type Request struct {
	*http.Request

	app *Application

	params map[string]string
	headers map[string]string
}


func (req *Request) Proxy() string{
	return req.Get("x-forwarded-for")
}

func (req *Request) IP() string {
	return req.RemoteAddr
}

func (req *Request) XHR() bool{
	var val = req.Get("x-requested-with");
	return val != "" && strings.ToLower(val) == "xmlhttprequest";
}

func (req *Request) Secure() bool{
	return strings.Contains(strings.ToLower(req.Proto), "https")
}

func (req *Request) Set(key, value string) {
	if req.headers == nil {
		req.headers = req.Headers()
	}

	req.headers[strings.ToLower(key)] = value
}

func (req *Request) Get(key string) string {
	if req.headers == nil {
		req.headers = req.Headers()
	}

	return req.headers[strings.ToLower(key)]
}

func (req *Request) Headers() map[string]string{
	if req.headers == nil {
		req.headers = make(map[string]string)
	}

	for key, value := range req.Request.Header {
		req.headers[strings.ToLower(key)] = strings.Join(value, "")
	}

	return req.headers
}

func (req *Request) Query(key string) string  {
	return req.URL.Query().Get(key)
}

func (req *Request) Querys() map[string][]string{
	return req.URL.Query()
}

func (req *Request) Param(key string) string {
	return req.params[key]
}

func (req *Request) Params(key string) map[string]string {
	return req.params
}

func NewRequest(req *http.Request, app *Application) *Request {
	return &Request{
		Request: req,
		app: app,
		params: nil,
		headers: nil,
	}
}