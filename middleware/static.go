package middleware

import (
	"synder.me/goe"
	"net/http"
	"strings"
	"mime"
	"path"
	"os"
	"errors"
)

func Static(base, root string) goe.Handler {

	base = path.Clean(base)
	root = path.Clean(root)

	_, err := os.Stat(root)

	if os.IsNotExist(err) {
		panic(errors.New("root path is not exist"))
	}

	fileServer := http.StripPrefix(base, http.FileServer(http.Dir(root)))

	return func(context *goe.Context) {

		if context.Request.Method != http.MethodGet && context.Request.Method != http.MethodHead {
			context.Next(nil)
			return
		}

		if strings.HasPrefix(context.Request.URL.Path, base) {
			ext := path.Ext(context.Request.URL.Path)
			context.Response.Set("Content-Type", mime.TypeByExtension(ext))
			fileServer.ServeHTTP(context.Response, context.Request.Request)
		}else {
			context.Next(nil)
		}
	}

}