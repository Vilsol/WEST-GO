package WEST

import (
	"fmt"
	"net/http"
	"testing"
)

func TestListen(t *testing.T) {
	Listen("0.0.0.0", 81, HelloWorldResponse())
}

func HelloWorldResponse() http.Handler {
	fmt.Println("Initialized")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Hello World</h1>"))
	})
}
