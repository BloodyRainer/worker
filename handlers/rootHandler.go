package handlers

import (
	"fastworker/work"
	"fmt"
	"net/http"
)

type RootHandler struct{}

func (rcv RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println("serve http")

	//mw := io.MultiWriter(w, os.Stdout)
	work.AddTask(w, r)

	fmt.Println("task added")

}
