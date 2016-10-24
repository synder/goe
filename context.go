package goe

import "fmt"

type Handler func(c *Context)


type Context struct {
	Application *Application
	Request     *Request
	Response    *Response
	Data        map[string]interface{}

	router      *Router
	render      *Render
	current     int
	params      map[string]string
	handlers    []Handler

}

func (context *Context) handle(err interface{}){

	if err != nil {

		if t, ok := err.(HttpError); ok {
			context.router.errorHandler(&HttpError{
				Code: t.Code,
				Msg: t.Error(),
			}, context)

		}else if t, ok := err.(error); ok {
			context.router.errorHandler(&HttpError{
				Code: 500,
				Msg: t.Error(),
			}, context)
		}else {
			context.router.errorHandler(&HttpError{
				Code: 500,
				Msg: fmt.Sprintf("%s", err),
			}, context)
		}
		return
	}

	length := len(context.handlers)

	if length == 0{
		context.router.errorHandler(&HttpError{
			Code: 404,
			Msg: "resource not found",
		}, context)
		return
	}

	if context.current < length{
		context.handlers[context.current](context)
		return
	}

	context.router.errorHandler(&HttpError{
		Code: 404,
		Msg: "resource not found",
	}, context)
}

func (context *Context) Next(err interface{}){
	if err != nil {
		context.handle(err)
	}else{
		context.current += 1
		context.handle(nil)
	}
}

func (context *Context) Abort(){
	context.current = len(context.handlers)
}

func (context *Context) Close()  {
	conn, _, err := context.Response.Hijack()

	if err != nil {
		panic(err)
	}

	conn.Close()
}


func (context *Context) Query(key string) string{
	return context.Request.Query(key)
}

func (context *Context) Param(key string) string{
	return context.Request.Param(key)
}

func (context *Context) Form(key string) string{
	context.Request.ParseForm()
	return context.Request.Form.Get(key)
}

func (context *Context) Status(code int) *Context {
	context.Response.WriteHeader(code)
	return context
}

func (context *Context) Json(obj interface{}) error{
	return context.Response.Json(obj)
}

func (context *Context) String(data string) error{
	_, err := context.Response.WriteString(data)
	return err
}

func (context *Context) Render(tpml string, obj interface{}) error{
	return context.Response.Render(tpml, obj, nil)
}


func NewContext(req *Request, res *Response, router *Router) *Context {

	handlers, params := router.match(req.Method, req.URL.Path)

	req.params = params

	var context = &Context{
		Application: router.App,
		Data: make(map[string]interface{}),
		Request: req,
		Response: res,
		router: router,
		current: 0,
		handlers: handlers,
		params: params,
	}

	return context
}