package goe

import (
	"net/http"
	"strconv"
	"os"
	"bufio"
	"io"
	"net/url"
)

type Response struct {
	http.ResponseWriter
	status int
}


//Set header
func (res *Response) Set(key, value string) {
	res.Header().Set(key, value)
}

func (res *Response) Get(key string) string {
	return res.Header().Get(key)
}

func (res *Response) Status(code int) int{
	if(code != 0){
		res.status = code
	}
	return res.status
}

func (res *Response) Send(data []byte) {
	res.WriteHeader(res.status)
	res.Write(data)
}

func (res *Response) Json(json string) {
	res.Set("Content-Type", "application/json")
	res.Send([]byte(json))
}

func (res *Response) Jsonp(json string) {

}

func (res *Response) Location(url string, code int) {
	res.Set("Location", url)
	res.Status(code)
}

func (res *Response) Redirect(url string) {
	res.Location(url, 302)
}

func (res *Response) Refresh(url string, time int) {
	value := strconv.Itoa(time) + ";url='" + url +"'"
	res.Set("Refresh", value)
}

func (res *Response) Attachment(path, filename string) {
	file, err := os.Open(path)

	if err != nil {
		//抛出错误
	}

	res.Set("Content-Disposition", "attachment;filename=" + url.QueryEscape(filename));

	br := bufio.NewReader(file)
	var temp = make([]byte, 512)

	for{
		n, err := br.Read(temp)
		if err != nil {
			if err == io.EOF {
				break;
			}else {
				//抛出错误
			}
		}
		res.Write(temp[:n])
	}
}

func (res *Response) Download(data []byte, filename string) {
	res.Write(data)
	res.Set("Content-Disposition", "attachment;filename=" + url.QueryEscape(filename));
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

func (res *Response) Render(path string, send bool) []byte {

	data := []byte{}

	if send == true {
		res.Set("Content-Type", "text/html")
		res.Send(data)
	}

	return data
}


func NewResponse(res http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: res,
	}
}
