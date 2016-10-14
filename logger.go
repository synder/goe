package goe

import (
	"time"
	"fmt"
	"os"
)

var (
	GREEN string  = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	WHITE  string  = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	YELLOW string = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	RED  string   = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	BLUE  string  = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	MAGENTA string= string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	CYAN string   = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	RESET string  = string([]byte{27, 91, 48, 109})
)

const (
	LOG_FORMAT = "%s [%-7s] %v %s %3d %s %13v  %s %s  %s\n"
	DATE_FORMAT = "2006/01/02 - 15:04:05"
)


func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return GREEN
	case code >= 300 && code < 400:
		return WHITE
	case code >= 400 && code < 500:
		return YELLOW
	default:
		return RED
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return BLUE
	case "POST":
		return CYAN
	case "PUT":
		return YELLOW
	case "DELETE":
		return RED
	case "PATCH":
		return GREEN
	case "HEAD":
		return MAGENTA
	case "OPTIONS":
		return WHITE
	default:
		return RESET
	}
}

func Logger() Handler{
	return func(c *Context){
		start := time.Now()
		path := c.Request.URL.Path



		c.Next(nil)

		end := time.Now()
		latency := end.Sub(start)

		clientIP := "127.0.0.1"
		method := c.Request.Method
		statusCode := c.Response.Status(0)
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)

		fmt.Fprintf(os.Stdout,
			LOG_FORMAT,
			methodColor,
			method,
			end.Format(DATE_FORMAT),
			statusColor,
			statusCode,
			RESET,
			latency,
			clientIP,
			RESET,
			path,
		)
	}
}
