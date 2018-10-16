package handlers

import (
	"bytes"
	"fastworker/work"
	"fmt"
	"io"
	"net/http"
	"os"
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

	do := func(rw http.ResponseWriter, req *http.Request) error {
		buf := new(bytes.Buffer)

		d := 10 * time.Millisecond
		time.Sleep(d)
		s := fmt.Sprintf("task will take %v \n", d)
		_, err := buf.WriteString(s)
		if err != nil {
			return fmt.Errorf("error writing string to buffer: %v", err)
		}

		mw := io.MultiWriter(rw, os.Stdout)
		io.WriteString(mw, buf.String())

		return nil
	}

	for i := 0; i < n; i++ {

		time.Sleep(1 * time.Millisecond)

		r := new(http.Request)
		w := responseWriterMock{
			wg: wg,
		}

		work.SubmitTask(w, r, do)
	}

	wg.Wait()
}
