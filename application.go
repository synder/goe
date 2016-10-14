package goe

import (
	"net/http"
	"strconv"
)


type Application struct {
	*Router

	locals map[string]string
}

func (app *Application) SetValue(key, value string) {
	app.locals[key] = value
}

func (app *Application) GetValue(key string) string {
	return app.locals[key]
}


func (app *Application) Listen(port int, host string) (string, error) {
	addr := host + ":" + strconv.Itoa(port)
	err := http.ListenAndServe(addr, app)
	return addr, err
}

func NewApp() *Application  {
	locals := make(map[string]string)

	router := NewRouter()

	app := &Application{
		Router: router,
		locals: locals,
	}

	router.App = app

	return app
}
