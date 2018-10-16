package handlers

import (
	"fastworker/bottleneck"
	"net/http"
	"sync"
	"testing"
)

type responseWriterMock struct {
	wg *sync.WaitGroup
}

func (rcv responseWriterMock) Header() http.Header {
	return nil
}

func (rcv responseWriterMock) Write(p []byte) (int, error) {
	rcv.wg.Done()
	return len(p), nil
}

func (rcv responseWriterMock) WriteHeader(int) {}

func TestRootHandler_ServeHTTP(t *testing.T) {

	bottleneck.Init(5)

	const n = 5

}
