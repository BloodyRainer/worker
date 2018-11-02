package handlers

import (
	"bytes"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type RootHandler struct{}

func (rcv RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	buf := new(bytes.Buffer)
	buf.WriteString("root handler response")

	d := time.Duration(rand.Intn(100)) * time.Millisecond
	time.Sleep(d)

	c := &http.Cookie{
		Name: "bloody-cookie",
		Value: "12CDEFGHIJKLMN",
	}

	http.SetCookie(w, c)

	io.WriteString(w, buf.String())

}
