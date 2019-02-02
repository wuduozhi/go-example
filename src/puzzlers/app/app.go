package hello

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/hello", helloHandler)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from the go app")
}
