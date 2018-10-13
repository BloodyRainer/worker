package handlers

import (
	"fmt"
	"net/http"
)

type RootHandler struct {}

func (rcv RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "fast")
}
