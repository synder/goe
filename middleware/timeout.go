package middleware

import (
	"time"
	"synder.me/goe"
)

func Timeout(d time.Duration) goe.Handler{

	timeout := func(context *goe.Context, timer *time.Timer) {
		<- timer.C
		context.Status(408)
		context.Response.Write([]byte{})
	}

	return func(context *goe.Context) {
		timer := time.NewTimer(d)

		go timeout(context, timer)

		context.Next(nil)
	}

}