package goe

type Handler func(c *Context)


type Context struct {
	Application *Application

	Request     *Request
	Response    *Response

	router      *Router
	current     int
	params      map[string]string
	handlers    []Handler
}

func (context *Context) handle(err *HttpError){

	if err != nil {
		context.router.errorHandler(&HttpError{
			Code: 500,
			Msg: err.Error(),
		}, context)
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

func (context *Context) Next(err *HttpError){
	if err != nil {
		context.handle(err)
	}else{
		context.current += 1
		context.handle(nil)
	}
}

func NewContext(req *Request, res *Response, router *Router) *Context {

	handlders, params := router.match(req.Method, req.URL.Path)

	req.params = params

	return &Context{
		Application: router.App,
		Request: req,
		Response: res,
		router: router,
		current: 0,
		handlers: handlders,
		params: params,
	}
}