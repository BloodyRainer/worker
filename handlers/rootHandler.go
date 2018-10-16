package handlers

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"
)

type RootHandler struct{}

func (rcv RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	buf := new(bytes.Buffer)

	d := 200 * time.Millisecond
	time.Sleep(d)
	//s := fmt.Sprintf("task will take %v \n", d)
	//_, err := buf.WriteString(s)
	//if err != nil {
	//	util.GetLogger().Error("error writing string", zap.Error(err))
	//}

	mw := io.MultiWriter(w, os.Stdout)
	io.WriteString(mw, buf.String())

}
