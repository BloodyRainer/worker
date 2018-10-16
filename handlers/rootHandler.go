package handlers

import (
	"bytes"
	"fastworker/work"
	"io"
	"net/http"
	"os"
	"time"
)

type RootHandler struct{}

func (rcv RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	do := func(rw http.ResponseWriter, req *http.Request) error {
		buf := new(bytes.Buffer)

		d := 10 * time.Millisecond
		time.Sleep(d)
		//s := fmt.Sprintf("task will take %v \n", d)
		//_, err := buf.WriteString(s)
		//if err != nil {
		//	return fmt.Errorf("error writing string to buffer: %v", err)
		//}

		mw := io.MultiWriter(rw, os.Stdout)
		io.WriteString(mw, buf.String())

		return nil
	}

	wg := work.SubmitTask(w, r, do)
	wg.Wait()

}
