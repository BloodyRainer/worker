package handlers

import (
	"fastworker/work"
	"net/http"
	"sync"
	"testing"
	"time"
)

type responseWriterMock struct {
	wg *sync.WaitGroup
}

func (rcv responseWriterMock) Header() http.Header {
	return nil
}

func (rcv responseWriterMock) Write(p []byte) (int, error){
	rcv.wg.Done()
	return len(p), nil
}

func (rcv responseWriterMock) WriteHeader(int) {}

func TestRootHandler_ServeHTTP(t *testing.T) {

	const n = 5

	go work.InitWorkers(5)

	time.Sleep(100 * time.Millisecond)

	wg := &sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {

		time.Sleep(1 * time.Millisecond)

		r := new(http.Request)
		w := responseWriterMock{
			wg: wg,
		}

		work.SubmitTask(w, r)
	}

	wg.Wait()
}
