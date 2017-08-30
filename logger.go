package WEST

import (
	"fmt"
	"net/http"
	"time"
)

func LoggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		end := time.Now()
		latency := end.Sub(start)

		clientIP := r.RemoteAddr
		method := r.Method

		statusCode := w.(*WESTWriter).Status()

		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)

		path := r.RequestURI

		fmt.Printf("[API] %v |%s %3d %s| %12v | %21s |%s %-6s %s %s\n",
			end.Format("2006/01/02 - 15:04:05"),
			statusColor, statusCode, reset,
			latency,
			clientIP,
			methodColor, method, reset,
			path,
		)
	})
}

var (
	green   = string([]byte{27, 91, 51, 48, 59, 52, 50, 109})
	yellow  = string([]byte{27, 91, 51, 48, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return magenta
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	default:
		return reset
	}
}
