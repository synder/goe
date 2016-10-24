package goe

import (
	"io"
	"html/template"
	"path"
	"sync"
	"errors"
)

type Render interface {
	Execute(wr io.Writer, cache bool, tpth string, obj interface{}) (error)
}

type DefaultRender struct {
	sync.RWMutex
	caches map[string]*template.Template
}


func (render *DefaultRender) setCache(key string, t *template.Template) {
	render.Lock()
	render.caches[key] = t
	render.Unlock()
}

func (render *DefaultRender) getCache(key string) *template.Template  {
	render.RLock()
	temp := render.caches[key]
	render.RUnlock()
	return temp
}


func (render *DefaultRender) Execute(wr io.Writer, cache bool, tpth string, obj interface{}) (error) {

	tpth = path.Clean(tpth)

	var tmpl *template.Template
	var err error

	if cache == true {
		tmpl = render.getCache(tpth)
	}

	if tmpl == nil {
		tmpl, err = template.ParseFiles(tpth)

		if err != nil {
			return err
		}

		render.setCache(tpth, tmpl)
	}

	if tmpl == nil {
		return errors.New("template not found in " + tpth)

	}

	err = tmpl.Execute(wr, obj)

	return err
}

func NewDefaultRender() Render {
	return &DefaultRender{
		caches: make(map[string]*template.Template),
	}
}