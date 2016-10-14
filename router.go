package goe

import (
	"net/http"
	"strings"
	"synder.me/goe/lib"
)

//Params define types
type Params map[string]string

type Next func(err error)

type HttpError struct {
	Code int
	Msg  string
}

func (e HttpError) Error() string {
	return e.Msg
}

type ErrorHandler func(err *HttpError, context *Context)

func DefaultErrorHandler(err *HttpError, context *Context) {
	context.Response.Status(err.Code)
	context.Response.Send([]byte(err.Error()))
	return
}



//router route tree----------------------------------------------------------
type node struct {
	path         string
	param        string
	handlers     []Handler
	interceptors []Handler
	childs       map[string]*node
}

func newNode(path, param string) *node {
	return &node{
		path: path,
		param: param,
		handlers: nil,
		interceptors: nil,
		childs: nil,
	}
}

type tree struct {
	root *node
}

func (t *tree) match(pth string) ([]Handler, Params) {

	pths := lib.Split(pth)
	length := len(pths)

	//root
	if length == 0 {
		return t.root.handlers, nil
	}

	var params map[string]string
	var child *node
	var root *node = t.root

	for i := 0; i < length; i++ {

		temp := pths[i]

		child = root.childs[temp];

		if child == nil {
			child = root.childs["*"]
			if child != nil {
				if child.param != "" {
					if params == nil {
						params = make(map[string]string)
					}
					params[child.param] = temp
				}
			} else {
				return nil, nil
			}
		}

		if i == (length - 1) {
			return child.handlers, params
		}

		root = child
	}

	return nil, nil
}

func (t *tree) add(pth string, handler Handler) {

	pths := lib.Split(pth)
	var length int = len(pths)

	//add intercetors or routers in root path "/"
	if length == 0 {
		if t.root.handlers == nil {
			t.root.handlers = []Handler{handler}
		} else {
			t.root.handlers = append(t.root.handlers, handler)
		}
		return
	}

	//add router or interceptor to node
	var child *node
	var root *node = t.root

	for i := 0; i < length; i++ {

		temp := pths[i]
		param := ""

		if strings.HasPrefix(temp, ":") {
			param = temp[1:]
			temp = "*"
		}

		if root.childs == nil {
			root.childs = make(map[string]*node)
		}

		if _, ok := root.childs[temp]; ok == false {
			root.childs[temp] = newNode(temp, param)
		}

		child = root.childs[temp]

		if i == (length - 1) {
			if child.handlers == nil {
				child.handlers = []Handler{handler}
			} else {
				child.handlers = append(child.handlers, handler)
			}
		}

		root = child
	}
}

func newTree() *tree {
	return &tree{
		root: newNode("/", ""),
	}
}


//router-------------------------------------------------------------------

type Router struct {
	sensitive    bool             //if case sensitive
	host         bool             //if host match
	App          *Application

	errorHandler ErrorHandler     //error handler

	interceptors []Handler        //global interceptor
	routes       map[string]*tree //routes
}

//privare method

func (r *Router) match(method, pth string) ([]Handler, Params) {

	method = strings.ToUpper(method)

	if tree, ok := r.routes[method]; ok == true {

		handlers, params := tree.match(pth)

		ilength := len(r.interceptors)
		hlength := len(handlers)

		tlength := ilength + hlength

		results := make([]Handler, tlength, tlength)

		index := 0

		for i := 0; i < ilength; i++ {
			results[index] = r.interceptors[i]
			index++
		}

		for i := 0; i < hlength; i++ {
			results[index] = handlers[i]
			index++
		}

		return results, params
	}

	return r.interceptors, nil
}

//public method

func (r *Router) register(method, pth string, handlers ...Handler) {

	method = strings.ToUpper(method)

	if _, ok := r.routes[method]; ok == false {
		r.routes[method] = newTree()
	}

	for _, handler := range handlers {
		r.routes[method].add(pth, handler)
	}
}

func (r *Router) Use(handlers ...Handler) {
	r.interceptors = append(r.interceptors, handlers...)
}

func (r *Router) All(path string, handlers ...Handler) {
	r.register(http.MethodGet, path, handlers...)
	r.register(http.MethodPost, path, handlers...)
	r.register(http.MethodPut, path, handlers...)
	r.register(http.MethodDelete, path, handlers...)
	r.register(http.MethodHead, path, handlers...)
	r.register(http.MethodOptions, path, handlers...)
	r.register(http.MethodPatch, path, handlers...)
}

func (r *Router) Get(path string, handlers ...Handler) {
	r.register(http.MethodGet, path, handlers...)
}

func (r *Router) Head(path string, handlers ...Handler) {
	r.register(http.MethodHead, path, handlers...)
}

func (r *Router) Options(path string, handlers ...Handler) {
	r.register(http.MethodOptions, path, handlers...)
}

func (r *Router) Post(path string, handlers ...Handler) {
	r.register(http.MethodPost, path, handlers...)
}

func (r *Router) Put(path string, handlers ...Handler) {
	r.register(http.MethodPut, path, handlers...)
}

func (r *Router) Patch(path string, handlers ...Handler) {
	r.register(http.MethodPatch, path, handlers...)
}

func (r *Router) Delete(path string, handlers ...Handler) {
	r.register(http.MethodDelete, path, handlers...)
}

func (r *Router) Error(handler ErrorHandler) {
	r.errorHandler = handler
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	context := NewContext(NewRequest(req), NewResponse(writer), r)

	context.handle(nil)
}


// NewRouter: create a new router
func NewRouter() *Router {
	return &Router{
		sensitive: true,
		host: true,
		errorHandler: DefaultErrorHandler,
		routes: make(map[string]*tree),
	}
}