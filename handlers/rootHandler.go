package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type RootHandler struct{}

func (rcv RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	buf := new(bytes.Buffer)
	buf.WriteString("root handler response")

	d := 80 * time.Millisecond
	time.Sleep(d)
	//s := fmt.Sprintf("task will take %v \n", d)
	//_, err := buf.WriteString(s)
	//if err != nil {
	//	util.GetLogger().Error("error writing string", zap.Error(err))
	//}

	if (true == ture) {
		
	}

	c := &http.Cookie{
		Name: "volksoft-cookie",
		Value: "12CDEFGHIJKLMN",
	}

	http.SetCookie(w, c)

	io.WriteString(w, buf.String())

}
