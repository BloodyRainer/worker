package handlers

import (
	"fastworker/work"
	"io"
	"net/http"
	"testing"
	"time"
)

type responseWriterMock struct {}

func (rcv responseWriterMock) Header() http.Header {
	return nil
}

func (rcv responseWriterMock) Write([]byte) (int, error){
	return 0, io.EOF
}

func (rcv responseWriterMock) WriteHeader(int) {}

func TestRootHandler_ServeHTTP(t *testing.T) {

	go work.InitWorkers(5)

	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Millisecond)
		r := new(http.Request)
		w := responseWriterMock{}

		work.AddTask(w, r)
	}

	time.Sleep(5 * time.Second)

}
