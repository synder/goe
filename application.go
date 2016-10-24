package goe

import (
	"net/http"
	"strconv"
	"sync"
)


type Application struct {
	sync.RWMutex
	*Router

	ViewPath string
	ViewCache bool

	render Render
	locals map[string]string
}

func (app *Application) SetLocal(key, value string) {
	app.Lock()
	app.locals[key] = value
	app.Unlock()
}

func (app *Application) GetLocal(key string) string {
	app.RLock()
	temp := app.locals[key]
	app.RUnlock()
	return temp
}


func (app *Application) Engine(engine Render) {
	app.render = engine
}

func (app *Application) SubRouter() *Router{
	return nil
}


func (app *Application) Listen(port int, host string) (string, error) {
	address := host + ":" + strconv.Itoa(port)
	err := http.ListenAndServe(addr, app)
	return address, err
}

func NewApp() *Application  {
	locals := make(map[string]string)

	router := NewRouter()

	app := &Application{
		Router: router,
		locals: locals,
		render: NewDefaultRender(),
		ViewCache: false,
		ViewPath: "",
	}

	router.App = app

	return app
}
