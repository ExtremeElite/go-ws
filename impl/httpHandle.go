package impl

import (
	"io/ioutil"
	"net/http"
)
var HttpChan chan []byte
func HttpHandle(w http.ResponseWriter, r *http.Request)  {
	var (
		body []byte
		err error
	)
	if body,err=ioutil.ReadAll(r.Body);err!=nil{
		return
	}
	HttpChan<-body
	w.Write(body)
}