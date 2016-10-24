package middleware

import (
	"fmt"
	"synder.me/goe"
)


func Recovery() goe.Handler {
	var handle = func(context *goe.Context) {
		if err := recover(); err != nil {
			context.Response.WriteHeader(500)
			context.Response.Write([]byte(fmt.Sprintf("%s", err)))
		}
	}

	return func(context *goe.Context) {
		defer handle(context)
		context.Next(nil)
	}

}
