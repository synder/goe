package goe

import (
	"io"
	"bytes"
	"bufio"
	"strconv"
	"net/url"
	"net/http"
	"net"
	"encoding/json"
	"path"
)

const (
	DEFAULT_SIZE = 0
	DEFAULT_STATUS = -1
)

type Response struct {
	http.ResponseWriter

	app *Application

	size int  //content type
	status int //response status
}


//Set header
func (res *Response) Set(key, value string) {
	res.Header().Set(key, value)
}

func (res *Response) Get(key string) string {
	return res.Header().Get(key)
}

func (res *Response) Status() int {
	return res.status
}

func (res *Response) Size() int {
	return res.size
}

func (res *Response) WriteHeader(code int) {
	if res.status > 0 {
		return
	}

	res.status = code

	res.ResponseWriter.WriteHeader(res.status)
}

func (res *Response) Write(data []byte) (int, error) {
	n, err := res.ResponseWriter.Write(data)
	res.size += n
	return n, err
}

func (res *Response) WriteString(data string) (int, error){
	n, err := io.WriteString(res.ResponseWriter, data)
	res.size += n
	return n, err
}

func (res *Response) Json(obj interface{}) error {
	res.Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(res).Encode(obj)
}

func (res *Response) Jsonp(obj interface{}) error {

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("X-Content-Type-Options", "nosniff")
	results, err := json.Marshal(obj)

	if err != nil {
		return err
	}

	temp := [][]byte{
		[]byte("/**/ typeof"),
		[]byte("callback"),
		[]byte(" === 'function' &&"),
		[]byte("callback("),
		results,
		[]byte(");"),
	}

	body := bytes.Join(temp, []byte(" "))

	res.Write(body)

	return nil
}

func (res *Response) Render(tmpl string, obj interface{}, owr io.Writer) error {

	tmpl = path.Join(res.app.ViewPath, tmpl)

	//var wr io.Writer
	//
	//if owr != nil {
	//	wr = io.MultiWriter(res, owr)
	//}else{
	//	wr = res
	//}

	err := res.app.render.Execute(res, res.app.ViewCache, tmpl, obj)

	if err != nil {
		return err
	}

	return nil
}

func (res *Response) Redirect(status int, url string) {
	if status > 300 && status < 400 {
		res.Set("Location", url)
		res.WriteHeader(status)
	}else{
		res.Set("Location", url)
		res.WriteHeader(302)
	}
}

func (res *Response) Refresh(url string, time int) {
	value := strconv.Itoa(time) + ";url='" + url +"'"
	res.Set("Refresh", value)
}

func (res *Response) Download(data []byte, filename string) {
	res.Set("Content-Disposition", "attachment;filename=" + url.QueryEscape(filename));
	res.Write(data)
}

func (res *Response) Cookie(cookie *http.Cookie) {
	res.Header().Add("Set-Cookie", cookie.String())
}

func (res *Response) ClearCookie(name, path string) {
	cookie := &http.Cookie{
		Name: name,
		Value: "",
		Path: path,
		MaxAge: -1,
	}
	res.Cookie(cookie)
}

func (res *Response) Flush(){
	if t, ok := res.ResponseWriter.(http.Flusher); ok {
		t.Flush()
	}
}

func (res *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if res.size < 0 {
		res.size = 0
	}

	return res.ResponseWriter.(http.Hijacker).Hijack()
}

func (res *Response) CloseNotify() <-chan bool {
	return res.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func NewResponse(res http.ResponseWriter, app *Application) *Response {
	return &Response{
		ResponseWriter: res,
		app: app,
		size: DEFAULT_SIZE,
		status: DEFAULT_STATUS,
	}
}
