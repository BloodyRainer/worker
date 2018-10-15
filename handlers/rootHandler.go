package handlers

import (
	"fastworker/work"
	"net/http"
)

type RootHandler struct{}

func (rcv RootHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	w := work.SubmitTask(rw, r)
	w.Wait()

}
