package impl

import (
	"io/ioutil"
	"net/http"
)
var HttpChan=make(chan []byte,1)

func HttpHandle(w http.ResponseWriter, r *http.Request)  {
	var (
		body []byte
		err error
	)
	if body,err=ioutil.ReadAll(r.Body);err!=nil{
		HttpChan<-body
	}
	w.Write(body)
}